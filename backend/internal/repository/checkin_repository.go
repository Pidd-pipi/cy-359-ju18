package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/orienteering/platform/internal/model"
)

// CheckinRepository 打卡记录仓储。
type CheckinRepository struct {
	db *gorm.DB
}

// NewCheckinRepository 构造打卡记录仓储。
func NewCheckinRepository(db *gorm.DB) *CheckinRepository {
	return &CheckinRepository{db: db}
}

func (r *CheckinRepository) Create(tx *gorm.DB, record *model.CheckinRecord) error {
	if err := tx.Create(record).Error; err != nil {
		return fmt.Errorf("create checkin record: %w", err)
	}
	return nil
}

func (r *CheckinRepository) Exists(tx *gorm.DB, checkpointID, teamID int64) (bool, error) {
	var count int64
	if err := tx.Model(&model.CheckinRecord{}).
		Where("checkpoint_id = ? AND team_id = ?", checkpointID, teamID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("check checkin exists: %w", err)
	}
	return count > 0, nil
}

func (r *CheckinRepository) CountByTeam(tx *gorm.DB, activityID, teamID int64) (int64, error) {
	db := r.db
	if tx != nil {
		db = tx
	}
	var count int64
	if err := db.Model(&model.CheckinRecord{}).
		Where("activity_id = ? AND team_id = ?", activityID, teamID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count checkins by team: %w", err)
	}
	return count, nil
}

func (r *CheckinRepository) ListByTeam(teamID, activityID int64) ([]model.CheckinRecord, error) {
	var records []model.CheckinRecord
	if err := r.db.Where("team_id = ? AND activity_id = ?", teamID, activityID).
		Order("checked_in_at ASC").Find(&records).Error; err != nil {
		return nil, fmt.Errorf("list checkins by team: %w", err)
	}
	return records, nil
}

func (r *CheckinRepository) ListByActivity(activityID int64) ([]model.CheckinRecord, error) {
	var records []model.CheckinRecord
	if err := r.db.Where("activity_id = ?", activityID).
		Order("checked_in_at ASC").Find(&records).Error; err != nil {
		return nil, fmt.Errorf("list checkins by activity: %w", err)
	}
	return records, nil
}

// WithTx 在事务中执行业务函数。
func (r *CheckinRepository) WithTx(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
