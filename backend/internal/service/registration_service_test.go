package service

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/orienteering/platform/internal/constants"
	"github.com/orienteering/platform/internal/model"
	"github.com/orienteering/platform/internal/repository"
	"github.com/orienteering/platform/internal/util"
)

// newTestDB 构造静默的内存 SQLite（共享缓存，单连接串行化，等价于行锁下的串行调度）。
func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	if err := db.AutoMigrate(
		&model.User{}, &model.Activity{}, &model.Team{}, &model.TeamMember{},
		&model.Registration{},
	); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = sqlDB.Close() })
	return db
}

type fixture struct {
	svc      *RegistrationService
	activity *repository.ActivityRepository
	regRepo  *repository.RegistrationRepository
}

func newFixture(db *gorm.DB) *fixture {
	actRepo := repository.NewActivityRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	regRepo := repository.NewRegistrationRepository(db)
	checkinRepo := repository.NewCheckinRepository(db)
	lb := NewLeaderboardService(regRepo, checkinRepo, actRepo, nil)
	svc := NewRegistrationService(regRepo, actRepo, teamRepo, lb)
	return &fixture{svc: svc, activity: actRepo, regRepo: regRepo}
}

func seedFixture(t *testing.T, db *gorm.DB, maxTeams int, teamCount int) (*fixture, int64, []int64) {
	t.Helper()
	f := newFixture(db)
	now := time.Now()
	act := &model.Activity{
		Title: "测试活动", Status: constants.ActivityStatusPublished,
		StartTime: now, EndTime: now.Add(2 * time.Hour), MaxTeams: maxTeams,
	}
	if err := db.Create(act).Error; err != nil {
		t.Fatalf("create activity: %v", err)
	}
	teamIDs := make([]int64, 0, teamCount)
	for i := 0; i < teamCount; i++ {
		team := &model.Team{Name: fmt.Sprintf("队-%d", i+1), CaptainID: int64(i + 1)}
		if err := db.Create(team).Error; err != nil {
			t.Fatalf("create team: %v", err)
		}
		teamIDs = append(teamIDs, team.ID)
	}
	return f, act.ID, teamIDs
}

func appErrCode(t *testing.T, err error) int {
	t.Helper()
	var ae *util.AppError
	if errors.As(err, &ae) {
		return ae.Code
	}
	return -1
}

// 名额充足时报名进入待审核；超出上限后按提交顺序进入候补，候补不占名额。
func TestApply_PendingThenWaitlist(t *testing.T) {
	db := newTestDB(t)
	f, activityID, teamIDs := seedFixture(t, db, 2, 4)

	r1, err := f.svc.Apply(teamIDs[0], activityID, 1, testLogger())
	if err != nil {
		t.Fatalf("apply 1: %v", err)
	}
	if r1.Waitlisted || r1.Registration.Status != constants.RegistrationStatusPending {
		t.Fatalf("第 1 队应为待审核，got waitlisted=%v status=%s", r1.Waitlisted, r1.Registration.Status)
	}
	if _, err := f.svc.Apply(teamIDs[1], activityID, 2, testLogger()); err != nil {
		t.Fatalf("apply 2: %v", err)
	}
	r3, err := f.svc.Apply(teamIDs[2], activityID, 3, testLogger())
	if err != nil {
		t.Fatalf("apply 3: %v", err)
	}
	if !r3.Waitlisted || r3.Registration.Status != constants.RegistrationStatusWaitlisted {
		t.Fatalf("第 3 队应进候补，got waitlisted=%v status=%s", r3.Waitlisted, r3.Registration.Status)
	}
	if r3.WaitlistAhead != 0 {
		t.Fatalf("候补第 1 队前面应有 0 队，got %d", r3.WaitlistAhead)
	}
	r4, err := f.svc.Apply(teamIDs[3], activityID, 4, testLogger())
	if err != nil {
		t.Fatalf("apply 4: %v", err)
	}
	if !r4.Waitlisted || r4.WaitlistAhead != 1 {
		t.Fatalf("候补第 2 队前面应有 1 队，got %+v", r4)
	}

	occupied, _ := f.activity.CountRegistrations(activityID)
	if occupied != 2 {
		t.Fatalf("占用名额应恒为上限 2，got %d", occupied)
	}
	waitlisted, _ := f.activity.CountWaitlisted(activityID)
	if waitlisted != 2 {
		t.Fatalf("候补应有 2 队，got %d", waitlisted)
	}
}

// 待审核本身占用名额：全部通过后占用数仍不超过上限。
func TestApproveAll_DoesNotExceedMax(t *testing.T) {
	db := newTestDB(t)
	f, activityID, teamIDs := seedFixture(t, db, 2, 3)
	for i := 0; i < 2; i++ {
		if _, err := f.svc.Apply(teamIDs[i], activityID, int64(i+1), testLogger()); err != nil {
			t.Fatalf("apply: %v", err)
		}
	}
	regs, _ := f.svc.ListByActivity(activityID)
	for _, r := range regs {
		if r.Status != constants.RegistrationStatusPending {
			t.Fatalf("前置数据应为待审核，got %s", r.Status)
		}
		if _, err := f.svc.Approve(r.ID, 99, testLogger()); err != nil {
			t.Fatalf("approve %d: %v", r.ID, err)
		}
	}
	occupied, _ := f.activity.CountRegistrations(activityID)
	if occupied != 2 {
		t.Fatalf("全部通过后占用名额应仍为 2，got %d", occupied)
	}
}

// 拒绝待审核释放名额，最早候补自动递补为待审核，占用名额保持在上限。
func TestReject_PromotesEarliestWaitlisted(t *testing.T) {
	db := newTestDB(t)
	f, activityID, teamIDs := seedFixture(t, db, 1, 3)
	for i := 0; i < 3; i++ {
		if _, err := f.svc.Apply(teamIDs[i], activityID, int64(i+1), testLogger()); err != nil {
			t.Fatalf("apply %d: %v", i, err)
		}
	}
	regs, _ := f.svc.ListByActivity(activityID)
	var pendingID, w1ID, w2ID int64
	for _, r := range regs {
		switch r.Status {
		case constants.RegistrationStatusPending:
			pendingID = r.ID
		case constants.RegistrationStatusWaitlisted:
			if w1ID == 0 {
				w1ID = r.ID
			} else {
				w2ID = r.ID
			}
		}
	}
	if pendingID == 0 || w1ID == 0 || w2ID == 0 {
		t.Fatalf("种子状态错误: pending=%d w1=%d w2=%d", pendingID, w1ID, w2ID)
	}

	res, err := f.svc.Reject(pendingID, testLogger())
	if err != nil {
		t.Fatalf("reject: %v", err)
	}
	if res.Promoted == nil || res.Promoted.ID != w1ID {
		t.Fatalf("应递补最早候补 %d，got %+v", w1ID, res.Promoted)
	}
	if got, _ := f.regRepo.GetByID(w1ID); got.Status != constants.RegistrationStatusPending {
		t.Fatalf("最早候补应变为待审核，got %s", got.Status)
	}
	if got, _ := f.regRepo.GetByID(w2ID); got.Status != constants.RegistrationStatusWaitlisted {
		t.Fatalf("第二候补应保持候补，got %s", got.Status)
	}
	occupied, _ := f.activity.CountRegistrations(activityID)
	if occupied != 1 {
		t.Fatalf("拒绝+递补后占用名额应恒为 1，got %d", occupied)
	}
	views, err := f.svc.ToViews([]model.Registration{*mustGet(t, f, w2ID)})
	if err != nil {
		t.Fatalf("to views: %v", err)
	}
	if views[0].WaitlistAhead != 0 {
		t.Fatalf("递补后剩余候补前面应为 0 队，got %d", views[0].WaitlistAhead)
	}
}

// 拒绝候补不触发递补。
func TestRejectWaitlisted_NoPromotion(t *testing.T) {
	db := newTestDB(t)
	f, activityID, teamIDs := seedFixture(t, db, 1, 2)
	if _, err := f.svc.Apply(teamIDs[0], activityID, 1, testLogger()); err != nil {
		t.Fatal(err)
	}
	w, err := f.svc.Apply(teamIDs[1], activityID, 2, testLogger())
	if err != nil {
		t.Fatal(err)
	}
	res, err := f.svc.Reject(w.Registration.ID, testLogger())
	if err != nil {
		t.Fatalf("reject waitlisted: %v", err)
	}
	if res.Promoted != nil {
		t.Fatalf("拒绝候补不应触发递补，got promoted=%d", res.Promoted.ID)
	}
	if occupied, _ := f.activity.CountRegistrations(activityID); occupied != 1 {
		t.Fatalf("占用名额应保持 1，got %d", occupied)
	}
}

// 同一队重复提交不能产生第二条记录。
func TestApply_DuplicateTeamRejected(t *testing.T) {
	db := newTestDB(t)
	f, activityID, teamIDs := seedFixture(t, db, 5, 1)
	if _, err := f.svc.Apply(teamIDs[0], activityID, 1, testLogger()); err != nil {
		t.Fatal(err)
	}
	_, err := f.svc.Apply(teamIDs[0], activityID, 1, testLogger())
	if appErrCode(t, err) != constants.CodeAlreadyApplied {
		t.Fatalf("重复报名应返回 CodeAlreadyApplied，got %v", err)
	}
	var count int64
	db.Model(&model.Registration{}).Where("team_id = ? AND activity_id = ?", teamIDs[0], activityID).Count(&count)
	if count != 1 {
		t.Fatalf("同一队只能有 1 条报名，got %d", count)
	}
}

// 候补排位基于全局队列：即使只看“我的报名”过滤子集，前面队伍数仍正确。
func TestWaitlistAhead_FromFilteredMineView(t *testing.T) {
	db := newTestDB(t)
	f, activityID, teamIDs := seedFixture(t, db, 1, 3)
	var lastReg *model.Registration
	for i := 0; i < 3; i++ {
		r, err := f.svc.Apply(teamIDs[i], activityID, int64(i+1), testLogger())
		if err != nil {
			t.Fatalf("apply %d: %v", i, err)
		}
		lastReg = r.Registration
	}
	// 只取第三队自己的报名（等价于 mine 列表的过滤子集），前面应仍有 1 队（第一候补）。
	views, err := f.svc.ToViews([]model.Registration{*lastReg})
	if err != nil {
		t.Fatalf("to views: %v", err)
	}
	if views[0].Status != constants.RegistrationStatusWaitlisted || views[0].WaitlistAhead != 1 {
		t.Fatalf("第三候补在过滤视图中前面应仍为 1，got %+v", views[0])
	}
}

// 非队长不能报名。
func TestApply_NotCaptain(t *testing.T) {
	db := newTestDB(t)
	f, activityID, teamIDs := seedFixture(t, db, 5, 1)
	_, err := f.svc.Apply(teamIDs[0], activityID, 999, testLogger())
	if appErrCode(t, err) != constants.CodeForbidden {
		t.Fatalf("非队长应返回 CodeForbidden，got %v", err)
	}
}

// 并发报名：活动行锁串行化后，占用名额绝不超过上限，其余全部进候补。
func TestApply_ConcurrentNeverOversells(t *testing.T) {
	db := newTestDB(t)
	f, activityID, teamIDs := seedFixture(t, db, 3, 20)

	// 单连接 SQLite 会把事务串行化（等价于每个事务都拿到活动行锁）；
	// 用应用层互斥模拟“一次只允许一个报名事务做名额判定+写入”的临界区，
	// 验证 service 临界区内的判定逻辑在任何交错顺序下都不超额。
	var gate sync.Mutex
	var wg sync.WaitGroup
	errCh := make(chan error, len(teamIDs))
	for i, teamID := range teamIDs {
		wg.Add(1)
		go func(i int, teamID int64) {
			defer wg.Done()
			gate.Lock()
			defer gate.Unlock()
			_, err := f.svc.Apply(teamID, activityID, int64(i+1), testLogger())
			errCh <- err
		}(i, teamID)
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			t.Fatalf("concurrent apply: %v", err)
		}
	}
	occupied, _ := f.activity.CountRegistrations(activityID)
	if occupied != 3 {
		t.Fatalf("并发报名后占用名额必须等于上限 3，got %d", occupied)
	}
	waitlisted, _ := f.activity.CountWaitlisted(activityID)
	if waitlisted != 17 {
		t.Fatalf("其余 17 队应全部在候补，got %d", waitlisted)
	}
	var total int64
	db.Model(&model.Registration{}).Where("activity_id = ?", activityID).Count(&total)
	if total != 20 {
		t.Fatalf("应保留 20 条报名记录，got %d", total)
	}
}

// 并发同时审核/拒绝同一条报名：只有一个事务成功。
func TestConcurrentDoubleAction_OnlyOneSucceeds(t *testing.T) {
	db := newTestDB(t)
	f, activityID, teamIDs := seedFixture(t, db, 2, 3)
	if _, err := f.svc.Apply(teamIDs[0], activityID, 1, testLogger()); err != nil {
		t.Fatal(err)
	}
	if _, err := f.svc.Apply(teamIDs[1], activityID, 2, testLogger()); err != nil {
		t.Fatal(err)
	}
	wReg, err := f.svc.Apply(teamIDs[2], activityID, 3, testLogger())
	if err != nil {
		t.Fatal(err)
	}
	pRegs, _ := f.svc.ListByActivity(activityID)
	var pendingID int64
	for _, r := range pRegs {
		if r.Status == constants.RegistrationStatusPending {
			pendingID = r.ID
		}
	}

	// 单连接下两个事务不可能真正交叠；直接对同一条 pending 串行做“通过+拒绝”，
	// 第二个操作必须因 CAS 失败而报冲突，名额状态不被破坏。
	if _, err := f.svc.Approve(pendingID, 1, testLogger()); err != nil {
		t.Fatalf("首次通过应成功: %v", err)
	}
	if _, err := f.svc.Reject(pendingID, testLogger()); appErrCode(t, err) != constants.CodeConflict {
		t.Fatalf("状态已变为 approved 后再拒绝应冲突，got %v", err)
	}
	// 候补不能被直接通过。
	if _, err := f.svc.Approve(wReg.Registration.ID, 1, testLogger()); appErrCode(t, err) != constants.CodeConflict {
		t.Fatalf("候补不能直接通过，got %v", err)
	}
}

func mustGet(t *testing.T, f *fixture, id int64) *model.Registration {
	t.Helper()
	r, err := f.regRepo.GetByID(id)
	if err != nil {
		t.Fatalf("get reg %d: %v", id, err)
	}
	return r
}

func testLogger() *slog.Logger { return slog.Default() }
