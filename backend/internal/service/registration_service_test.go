package service

import (
	"fmt"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/orienteering/platform/internal/constants"
	"github.com/orienteering/platform/internal/model"
	"github.com/orienteering/platform/internal/repository"
)

// newRegTestDB 构造单连接内存 SQLite：单连接下事务串行执行，
// 可复现"两个管理员/两支队伍同时操作"的并发原子性场景。
func newRegTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared&_busy_timeout=5000"), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get sql db: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(
		&model.User{}, &model.Activity{}, &model.Checkpoint{},
		&model.Team{}, &model.TeamMember{}, &model.Registration{},
		&model.CheckinRecord{}, &model.Product{}, &model.Redemption{},
		&model.Favorite{}, &model.AuditLog{},
	); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	return db
}

func newRegService(t *testing.T) (*gorm.DB, *RegistrationService, *repository.RegistrationRepository, *repository.ActivityRepository) {
	db := newRegTestDB(t)
	regRepo := repository.NewRegistrationRepository(db)
	actRepo := repository.NewActivityRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	lb := NewLeaderboardService(regRepo, repository.NewCheckinRepository(db), actRepo, nil)
	svc := NewRegistrationService(regRepo, actRepo, teamRepo, lb)
	return db, svc, regRepo, actRepo
}

var testLogger = slog.New(slog.NewTextHandler(io.Discard, nil))

// seedPublishedActivity 创建一个容量 maxTeams 的已发布活动。
func seedPublishedActivity(t *testing.T, db *gorm.DB, maxTeams int) *model.Activity {
	t.Helper()
	now := time.Now()
	act := &model.Activity{
		Title: "候补测试活动", Difficulty: "adult",
		Status: constants.ActivityStatusPublished, CreatorID: 99, MaxTeams: maxTeams,
		StartTime: now, EndTime: now.Add(3 * time.Hour),
	}
	if err := db.Create(act).Error; err != nil {
		t.Fatalf("create activity: %v", err)
	}
	return act
}

// seedTeam 创建 id=n 的队长为 n 的团队。
func seedTeam(t *testing.T, db *gorm.DB, id int64) {
	t.Helper()
	team := &model.Team{ID: id, Name: fmt.Sprintf("team-%d", id), CaptainID: id}
	if err := db.Create(team).Error; err != nil {
		t.Fatalf("create team %d: %v", id, err)
	}
}

// apply 封装报名并返回报名记录。
func apply(t *testing.T, svc *RegistrationService, teamID, activityID int64) *model.Registration {
	t.Helper()
	reg, err := svc.Apply(teamID, activityID, teamID, testLogger)
	if err != nil {
		t.Fatalf("team %d apply: %v", teamID, err)
	}
	return reg
}

// TestApply_PendingOccupiesSlot 名额未满：待审核也占用名额。
func TestApply_PendingOccupiesSlot(t *testing.T) {
	db, svc, _, actRepo := newRegService(t)
	act := seedPublishedActivity(t, db, 2)
	seedTeam(t, db, 1)
	seedTeam(t, db, 2)

	r1 := apply(t, svc, 1, act.ID)
	if r1.Status != constants.RegistrationStatusPending {
		t.Fatalf("first reg status = %s, want pending", r1.Status)
	}
	r2 := apply(t, svc, 2, act.ID)
	if r2.Status != constants.RegistrationStatusPending {
		t.Fatalf("second reg status = %s, want pending (pending must occupy slot)", r2.Status)
	}
	occupied, err := actRepo.CountRegistrations(act.ID)
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	if occupied != 2 {
		t.Fatalf("occupied = %d, want 2 (pending occupies slot)", occupied)
	}
}

// TestApply_OverflowGoesWaitlistInOrder 名额不足：按提交顺序进入候补，前面队数正确。
func TestApply_OverflowGoesWaitlistInOrder(t *testing.T) {
	db, svc, _, _ := newRegService(t)
	act := seedPublishedActivity(t, db, 1)
	for _, id := range []int64{1, 2, 3} {
		seedTeam(t, db, id)
	}

	r1 := apply(t, svc, 1, act.ID)
	r2 := apply(t, svc, 2, act.ID)
	r3 := apply(t, svc, 3, act.ID)
	if r1.Status != constants.RegistrationStatusPending {
		t.Fatalf("r1 = %s, want pending", r1.Status)
	}
	if r2.Status != constants.RegistrationStatusWaitlist || r3.Status != constants.RegistrationStatusWaitlist {
		t.Fatalf("overflow regs = %s,%s, want waitlist,waitlist", r2.Status, r3.Status)
	}

	// 候补视图前面队数：r2 前面 0 队，r3 前面 1 队。
	regs, err := svc.ListByActivity(act.ID)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	views, err := svc.ToViews(regs)
	if err != nil {
		t.Fatalf("to views: %v", err)
	}
	ahead := map[int64]int{}
	status := map[int64]string{}
	for _, v := range views {
		ahead[v.TeamID] = v.WaitlistAhead
		status[v.TeamID] = v.Status
	}
	if status[1] != constants.RegistrationStatusPending ||
		status[2] != constants.RegistrationStatusWaitlist ||
		status[3] != constants.RegistrationStatusWaitlist {
		t.Fatalf("unexpected statuses: %+v", status)
	}
	if ahead[2] != 0 {
		t.Fatalf("team2 waitlist_ahead = %d, want 0", ahead[2])
	}
	if ahead[3] != 1 {
		t.Fatalf("team3 waitlist_ahead = %d, want 1", ahead[3])
	}
}

// TestReject_PendingPromotesEarliestWaitlist 拒绝待审核：最早候补自动补成待审核。
func TestReject_PendingPromotesEarliestWaitlist(t *testing.T) {
	db, svc, regRepo, actRepo := newRegService(t)
	act := seedPublishedActivity(t, db, 1)
	for _, id := range []int64{1, 2, 3, 4} {
		seedTeam(t, db, id)
	}
	r1 := apply(t, svc, 1, act.ID)
	r2 := apply(t, svc, 2, act.ID)
	r3 := apply(t, svc, 3, act.ID)
	r4 := apply(t, svc, 4, act.ID)
	_ = r4

	// 拒绝唯一的待审核 r1，最早候补 r2 应自动补位。
	if _, err := svc.Reject(r1.ID, testLogger); err != nil {
		t.Fatalf("reject: %v", err)
	}
	got2, err := regRepo.GetByID(r2.ID)
	if err != nil {
		t.Fatalf("get r2: %v", err)
	}
	if got2.Status != constants.RegistrationStatusPending {
		t.Fatalf("earliest waitlist status = %s, want promoted to pending", got2.Status)
	}
	got3, _ := regRepo.GetByID(r3.ID)
	if got3.Status != constants.RegistrationStatusWaitlist {
		t.Fatalf("team3 status = %s, want still waitlist", got3.Status)
	}
	// 候补 r3 前面现在 0 队。
	views, err := svc.ToViews([]model.Registration{*got3})
	if err != nil {
		t.Fatalf("views: %v", err)
	}
	if views[0].WaitlistAhead != 0 {
		t.Fatalf("team3 waitlist_ahead = %d, want 0 after promotion", views[0].WaitlistAhead)
	}
	// 名额占用始终为 1，拒绝+补位后不超额也不空缺。
	occupied, err := actRepo.CountRegistrations(act.ID)
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	if occupied != 1 {
		t.Fatalf("occupied after reject+promote = %d, want 1", occupied)
	}

	// 再拒绝补位的 r2，r3 顶上；拒绝 r3，r4 顶上：链式补位且名额恒为 1。
	if _, err := svc.Reject(r2.ID, testLogger); err != nil {
		t.Fatalf("reject r2: %v", err)
	}
	got3, _ = regRepo.GetByID(r3.ID)
	if got3.Status != constants.RegistrationStatusPending {
		t.Fatalf("team3 after chain reject = %s, want pending", got3.Status)
	}
	if _, err := svc.Reject(r3.ID, testLogger); err != nil {
		t.Fatalf("reject r3: %v", err)
	}
	got4, _ := regRepo.GetByID(r4.ID)
	if got4.Status != constants.RegistrationStatusPending {
		t.Fatalf("team4 after chain reject = %s, want pending", got4.Status)
	}
	occupied, _ = actRepo.CountRegistrations(act.ID)
	if occupied != 1 {
		t.Fatalf("occupied after chain = %d, want 1", occupied)
	}
}

// TestReject_WaitlistDoesNotPromote 拒绝候补不触发补位，名额不变。
func TestReject_WaitlistDoesNotPromote(t *testing.T) {
	db, svc, regRepo, actRepo := newRegService(t)
	act := seedPublishedActivity(t, db, 1)
	for _, id := range []int64{1, 2, 3} {
		seedTeam(t, db, id)
	}
	r1 := apply(t, svc, 1, act.ID)
	r2 := apply(t, svc, 2, act.ID)
	r3 := apply(t, svc, 3, act.ID)

	if _, err := svc.Reject(r3.ID, testLogger); err != nil {
		t.Fatalf("reject waitlist: %v", err)
	}
	got1, _ := regRepo.GetByID(r1.ID)
	got2, _ := regRepo.GetByID(r2.ID)
	if got1.Status != constants.RegistrationStatusPending {
		t.Fatalf("pending changed to %s", got1.Status)
	}
	if got2.Status != constants.RegistrationStatusWaitlist {
		t.Fatalf("earlier waitlist changed to %s, want waitlist", got2.Status)
	}
	occupied, _ := actRepo.CountRegistrations(act.ID)
	if occupied != 1 {
		t.Fatalf("occupied = %d, want 1 after rejecting a waitlist team", occupied)
	}
}

// TestApply_DuplicateTeamRejected 同一队重复提交不能产生两条记录/多占名额。
func TestApply_DuplicateTeamRejected(t *testing.T) {
	db, svc, regRepo, actRepo := newRegService(t)
	act := seedPublishedActivity(t, db, 5)
	seedTeam(t, db, 1)

	r1 := apply(t, svc, 1, act.ID)
	if _, err := svc.Apply(1, act.ID, 1, testLogger); err == nil {
		t.Fatal("duplicate apply should fail")
	}
	var n int64
	if err := db.Model(&model.Registration{}).
		Where("team_id = 1 AND activity_id = ?", act.ID).Count(&n).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	if n != 1 {
		t.Fatalf("registration rows for same team = %d, want 1", n)
	}
	occupied, _ := actRepo.CountRegistrations(act.ID)
	if occupied != 1 {
		t.Fatalf("occupied = %d, want 1 after duplicate apply", occupied)
	}
	// 被拒绝后同队仍不能重复提交（唯一索引拦截）。
	if _, err := svc.Reject(r1.ID, testLogger); err != nil {
		t.Fatalf("reject: %v", err)
	}
	if _, err := svc.Apply(1, act.ID, 1, testLogger); err == nil {
		t.Fatal("re-apply with existing rejected row should be rejected by unique index")
	}
	_ = regRepo
}

// TestConcurrentApprove_OnlyOneSucceeds 两个管理员同时审核同一条报名：只有一个成功。
func TestConcurrentApprove_OnlyOneSucceeds(t *testing.T) {
	db, svc, regRepo, _ := newRegService(t)
	act := seedPublishedActivity(t, db, 2)
	for _, id := range []int64{1, 2} {
		seedTeam(t, db, id)
	}
	r1 := apply(t, svc, 1, act.ID)
	r2 := apply(t, svc, 2, act.ID)
	_ = r2

	// 单连接串行执行两个 Approve（模拟两管理员同时点通过）。
	_, err1 := svc.Approve(r1.ID, 1, testLogger)
	_, err2 := svc.Approve(r1.ID, 2, testLogger)
	if err1 != nil {
		t.Fatalf("first approve should succeed, got %v", err1)
	}
	if err2 == nil {
		t.Fatal("second concurrent approve should fail with conflict")
	}
	got, _ := regRepo.GetByID(r1.ID)
	if got.Status != constants.RegistrationStatusApproved {
		t.Fatalf("status = %s, want approved", got.Status)
	}
}

// TestConcurrentRejectAndApprove_OnlyOneWins 一个管理员通过、另一个同时拒绝：状态唯一，不补位错误。
func TestConcurrentRejectAndApprove_OnlyOneWins(t *testing.T) {
	db, svc, regRepo, actRepo := newRegService(t)
	act := seedPublishedActivity(t, db, 1)
	for _, id := range []int64{1, 2} {
		seedTeam(t, db, id)
	}
	r1 := apply(t, svc, 1, act.ID)
	r2 := apply(t, svc, 2, act.ID)

	// 先拒绝 r1（r2 补位 pending），再对 r1 调 Approve 必须失败。
	if _, err := svc.Reject(r1.ID, testLogger); err != nil {
		t.Fatalf("reject: %v", err)
	}
	if _, err := svc.Approve(r1.ID, 1, testLogger); err == nil {
		t.Fatal("approve after reject should fail")
	}
	// r2 已补位为待审核；通过 r2 后名额仍恰好 1，不会超额。
	if _, err := svc.Approve(r2.ID, 1, testLogger); err != nil {
		t.Fatalf("approve promoted waitlist r2: %v", err)
	}
	occupied, _ := actRepo.CountRegistrations(act.ID)
	if occupied != 1 {
		t.Fatalf("occupied = %d, want exactly maxTeams 1", occupied)
	}
	got1, _ := regRepo.GetByID(r1.ID)
	got2, _ := regRepo.GetByID(r2.ID)
	if got1.Status != constants.RegistrationStatusRejected || got2.Status != constants.RegistrationStatusApproved {
		t.Fatalf("statuses = %s,%s, want rejected,approved", got1.Status, got2.Status)
	}
}

// TestApprove_WaitlistNotApprovable 候补不能直接通过，必须先排队补位。
func TestApprove_WaitlistNotApprovable(t *testing.T) {
	db, svc, _, _ := newRegService(t)
	act := seedPublishedActivity(t, db, 1)
	for _, id := range []int64{1, 2} {
		seedTeam(t, db, id)
	}
	_ = apply(t, svc, 1, act.ID)
	r2 := apply(t, svc, 2, act.ID)
	if _, err := svc.Approve(r2.ID, 1, testLogger); err == nil {
		t.Fatal("approving a waitlist registration directly should fail")
	}
}

// TestCapacityNeverExceedsUnderSerialContention 串行竞争下名额绝不超出上限。
func TestCapacityNeverExceedsUnderSerialContention(t *testing.T) {
	db, svc, _, actRepo := newRegService(t)
	const maxTeams = 3
	act := seedPublishedActivity(t, db, maxTeams)
	const teamN = 8
	for id := int64(1); id <= teamN; id++ {
		seedTeam(t, db, id)
	}
	for id := int64(1); id <= teamN; id++ {
		apply(t, svc, id, act.ID)
	}
	occupied, _ := actRepo.CountRegistrations(act.ID)
	if occupied != maxTeams {
		t.Fatalf("occupied = %d, want exactly %d", occupied, maxTeams)
	}

	// 通过全部待审核（模拟管理员把所有待审核都点通过），正式队伍数仍恰好等于上限。
	regs, _ := svc.ListByActivity(act.ID)
	pendingN := 0
	for _, r := range regs {
		if r.Status == constants.RegistrationStatusPending {
			pendingN++
			if _, err := svc.Approve(r.ID, 1, testLogger); err != nil {
				t.Fatalf("approve pending %d: %v", r.ID, err)
			}
		}
	}
	if pendingN != maxTeams {
		t.Fatalf("pending count = %d, want %d", pendingN, maxTeams)
	}
	occupied, _ = actRepo.CountRegistrations(act.ID)
	if occupied != maxTeams {
		t.Fatalf("occupied after approving all pending = %d, want %d (must not exceed cap)", occupied, maxTeams)
	}
}
