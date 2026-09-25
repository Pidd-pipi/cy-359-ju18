package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/orienteering/platform/internal/model"
)

// CheckpointRepository 打卡点仓储。
type CheckpointRepository struct {
	db *gorm.DB
}

// NewCheckpointRepository 构造打卡点仓储。
func NewCheckpointRepository(db *gorm.DB) *CheckpointRepository {
	return &CheckpointRepository{db: db}
}

func (r *CheckpointRepository) Create(cp *model.Checkpoint) error {
	if err := r.db.Create(cp).Error; err != nil {
		return fmt.Errorf("create checkpoint: %w", err)
	}
	return nil
}

func (r *CheckpointRepository) GetByID(id int64) (*model.Checkpoint, error) {
	var cp model.Checkpoint
	if err := r.db.First(&cp, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get checkpoint by id: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("get checkpoint by id: %w", err)
	}
	return &cp, nil
}

func (r *CheckpointRepository) ListByActivity(activityID int64) ([]model.Checkpoint, error) {
	var cps []model.Checkpoint
	if err := r.db.Where("activity_id = ?", activityID).Order("sequence ASC").Find(&cps).Error; err != nil {
		return nil, fmt.Errorf("list checkpoints: %w", err)
	}
	return cps, nil
}

func (r *CheckpointRepository) Update(cp *model.Checkpoint) error {
	if err := r.db.Save(cp).Error; err != nil {
		return fmt.Errorf("update checkpoint: %w", err)
	}
	return nil
}

func (r *CheckpointRepository) Delete(id int64) error {
	if err := r.db.Delete(&model.Checkpoint{}, id).Error; err != nil {
		return fmt.Errorf("delete checkpoint: %w", err)
	}
	return nil
}

func (r *CheckpointRepository) CountByActivity(activityID int64) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Checkpoint{}).Where("activity_id = ?", activityID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count checkpoints: %w", err)
	}
	return count, nil
}
