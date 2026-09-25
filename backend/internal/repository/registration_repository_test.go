package repository

import (
	"errors"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/orienteering/platform/internal/constants"
	"github.com/orienteering/platform/internal/model"
)

func newPublishedActivity(t *testing.T, db *gorm.DB, maxTeams int) *model.Activity {
	t.Helper()
	act := &model.Activity{
		Title: "名额统计", Status: constants.ActivityStatusPublished, CreatorID: 1,
		StartTime: time.Now(), EndTime: time.Now().Add(time.Hour), MaxTeams: maxTeams,
	}
	if err := db.Create(act).Error; err != nil {
		t.Fatalf("create activity: %v", err)
	}
	return act
}

func TestRegistrationRepository_CountOccupiedAndWaitlisted(t *testing.T) {
	db := newTestDB(t)
	act := newPublishedActivity(t, db, 3)
	seedStatuses := []string{
		constants.RegistrationStatusPending,
		constants.RegistrationStatusApproved,
		constants.RegistrationStatusWaitlisted,
		constants.RegistrationStatusRejected,
		constants.RegistrationStatusFinished,
		constants.RegistrationStatusWaitlisted,
	}
	for i, st := range seedStatuses {
		if err := db.Create(&model.Registration{
			TeamID: int64(200 + i), ActivityID: act.ID, Status: st, RegisteredAt: time.Now(),
		}).Error; err != nil {
			t.Fatal(err)
		}
	}

	actRepo := NewActivityRepository(db)
	regRepo := NewRegistrationRepository(db)

	occupied, err := actRepo.CountOccupied(db, act.ID)
	if err != nil {
		t.Fatalf("count occupied: %v", err)
	}
	if occupied != 3 {
		t.Fatalf("占位应为 pending+approved+finished=3，got %d", occupied)
	}
	waitlist, err := actRepo.CountWaitlisted(act.ID)
	if err != nil {
		t.Fatalf("count waitlisted: %v", err)
	}
	if waitlist != 2 {
		t.Fatalf("候补应为 2，got %d", waitlist)
	}
	// 对外 CountRegistrations 与事务内 CountOccupied 口径一致。
	external, err := actRepo.CountRegistrations(act.ID)
	if err != nil {
		t.Fatal(err)
	}
	if external != occupied {
		t.Fatalf("CountRegistrations=%d 应等于 CountOccupied=%d", external, occupied)
	}

	// 候补排位：第二条候补前面应有 1 队。
	var regs []model.Registration
	db.Where("activity_id = ? AND status = ?", act.ID, constants.RegistrationStatusWaitlisted).
		Order("id ASC").Find(&regs)
	if len(regs) != 2 {
		t.Fatalf("候补记录应为 2，got %d", len(regs))
	}
	ahead, err := regRepo.CountWaitlistedBefore(nil, act.ID, regs[1].ID)
	if err != nil {
		t.Fatal(err)
	}
	if ahead != 1 {
		t.Fatalf("第二候补前面应为 1，got %d", ahead)
	}

	// 批量排位：只查第二条候补（模拟过滤后的子集），结果仍应基于全局队列得到 1。
	type row struct {
		ID    int64
		Ahead int64
	}
	var rows []row
	if err := regRepo.CountWaitlistedBeforeIDs(nil, act.ID, []int64{regs[1].ID}, &rows); err != nil {
		t.Fatalf("batch ahead: %v", err)
	}
	if len(rows) != 1 || rows[0].Ahead != 1 {
		t.Fatalf("批量排位应为 [{id:%d ahead:1}]，got %+v", regs[1].ID, rows)
	}
}

func TestRegistrationRepository_FirstWaitlistedAndCAS(t *testing.T) {
	db := newTestDB(t)
	act := newPublishedActivity(t, db, 1)
	ids := make([]int64, 3)
	statuses := []string{
		constants.RegistrationStatusPending,
		constants.RegistrationStatusWaitlisted,
		constants.RegistrationStatusWaitlisted,
	}
	for i, st := range statuses {
		reg := &model.Registration{TeamID: int64(300 + i), ActivityID: act.ID, Status: st, RegisteredAt: time.Now()}
		if err := db.Create(reg).Error; err != nil {
			t.Fatal(err)
		}
		ids[i] = reg.ID
	}
	regRepo := NewRegistrationRepository(db)

	// 事务内取最早候补并 CAS 递补。
	err := db.Transaction(func(tx *gorm.DB) error {
		first, ferr := regRepo.FirstWaitlistedForUpdate(tx, act.ID)
		if ferr != nil {
			return ferr
		}
		if first.ID != ids[1] {
			t.Fatalf("最早候补应为 %d，got %d", ids[1], first.ID)
		}
		n, uerr := regRepo.UpdateStatusCAS(tx, first.ID,
			constants.RegistrationStatusPending, constants.RegistrationStatusWaitlisted)
		if uerr != nil {
			return uerr
		}
		if n != 1 {
			t.Fatalf("CAS 应影响 1 行，got %d", n)
		}
		// 再次以 waitlisted 为前置条件更新已变为 pending 的记录，应影响 0 行。
		n2, uerr := regRepo.UpdateStatusCAS(tx, first.ID,
			constants.RegistrationStatusRejected, constants.RegistrationStatusWaitlisted)
		if uerr != nil {
			return uerr
		}
		if n2 != 0 {
			t.Fatalf("状态不符的 CAS 应影响 0 行，got %d", n2)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("tx: %v", err)
	}

	got, _ := regRepo.GetByID(ids[1])
	if got.Status != constants.RegistrationStatusPending {
		t.Fatalf("最早候补应已递补为待审核，got %s", got.Status)
	}

	// 候补取空时应返回 ErrNotFound。
	err = db.Transaction(func(tx *gorm.DB) error {
		// 先拒绝剩余候补，使队列清空。
		if _, err := regRepo.UpdateStatusCAS(tx, ids[2],
			constants.RegistrationStatusRejected, constants.RegistrationStatusWaitlisted); err != nil {
			return err
		}
		_, ferr := regRepo.FirstWaitlistedForUpdate(tx, act.ID)
		return ferr
	})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("候补为空应返回 ErrNotFound，got %v", err)
	}
}
