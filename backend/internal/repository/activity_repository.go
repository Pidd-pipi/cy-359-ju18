package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/orienteering/platform/internal/model"
)

// ActivityRepository 活动仓储。
type ActivityRepository struct {
	db *gorm.DB
}

// NewActivityRepository 构造活动仓储。
func NewActivityRepository(db *gorm.DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}

func (r *ActivityRepository) Create(activity *model.Activity) error {
	if err := r.db.Create(activity).Error; err != nil {
		return fmt.Errorf("create activity: %w", err)
	}
	return nil
}

func (r *ActivityRepository) GetByID(id int64) (*model.Activity, error) {
	var activity model.Activity
	if err := r.db.Preload("Checkpoints").First(&activity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get activity by id: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("get activity by id: %w", err)
	}
	return &activity, nil
}

func (r *ActivityRepository) GetByIDForUpdate(tx *gorm.DB, id int64) (*model.Activity, error) {
	var activity model.Activity
	if err := tx.Clauses(lockedClause()).First(&activity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get activity for update: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("get activity for update: %w", err)
	}
	return &activity, nil
}

func (r *ActivityRepository) Update(activity *model.Activity) error {
	if err := r.db.Save(activity).Error; err != nil {
		return fmt.Errorf("update activity: %w", err)
	}
	return nil
}

func (r *ActivityRepository) UpdateStatus(activityID int64, status string) error {
	if err := r.db.Model(&model.Activity{}).Where("id = ?", activityID).
		Update("status", status).Error; err != nil {
		return fmt.Errorf("update activity status: %w", err)
	}
	return nil
}

func (r *ActivityRepository) List(page, pageSize, offset int, status, difficulty, keyword string) ([]model.Activity, int64, error) {
	query := r.db.Model(&model.Activity{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("title LIKE ? OR address LIKE ?", like, like)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count activities: %w", err)
	}
	var activities []model.Activity
	if err := query.Order("id DESC").Limit(pageSize).Offset(offset).Find(&activities).Error; err != nil {
		return nil, 0, fmt.Errorf("list activities: %w", err)
	}
	return activities, total, nil
}

func (r *ActivityRepository) CountCheckpoints(activityID int64) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Checkpoint{}).Where("activity_id = ?", activityID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count checkpoints: %w", err)
	}
	return count, nil
}

func (r *ActivityRepository) CountRegistrations(activityID int64) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Registration{}).
		Where("activity_id = ? AND status IN ?", activityID,
			[]string{"approved", "finished"}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count registrations: %w", err)
	}
	return count, nil
}

func (r *ActivityRepository) Delete(id int64) error {
	if err := r.db.Delete(&model.Activity{}, id).Error; err != nil {
		return fmt.Errorf("delete activity: %w", err)
	}
	return nil
}
