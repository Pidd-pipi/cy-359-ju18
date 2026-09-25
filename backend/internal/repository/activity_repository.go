package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/orienteering/platform/internal/constants"
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

// CountRegistrations 统计活动已占用名额的报名数（pending/approved/finished）。
// 待审核同样占用名额，避免管理员全部通过后正式队伍超过 MaxTeams。
func (r *ActivityRepository) CountRegistrations(activityID int64) (int64, error) {
	return r.CountOccupied(r.db, activityID)
}

// CountOccupied 在给定事务/连接内统计占用名额的报名数（与 SELECT ... FOR UPDATE 同事务，保证并发下名额判定准确）。
func (r *ActivityRepository) CountOccupied(tx *gorm.DB, activityID int64) (int64, error) {
	var count int64
	if err := tx.Model(&model.Registration{}).
		Where("activity_id = ? AND status IN ?", activityID,
			constants.RegistrationStatusSlotOccupied).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count occupied registrations: %w", err)
	}
	return count, nil
}

// CountWaitlisted 统计活动候补队列中的队伍数（候补不占名额）。
func (r *ActivityRepository) CountWaitlisted(activityID int64) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Registration{}).
		Where("activity_id = ? AND status = ?", activityID,
			constants.RegistrationStatusWaitlisted).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count waitlisted registrations: %w", err)
	}
	return count, nil
}

func (r *ActivityRepository) Delete(id int64) error {
	if err := r.db.Delete(&model.Activity{}, id).Error; err != nil {
		return fmt.Errorf("delete activity: %w", err)
	}
	return nil
}

// TitlesByIDs 按 id 批量查询活动标题（报名列表展示复用）。
func (r *ActivityRepository) TitlesByIDs(ids []int64) (map[int64]string, error) {
	titles := make(map[int64]string, len(ids))
	if len(ids) == 0 {
		return titles, nil
	}
	var activities []model.Activity
	if err := r.db.Select("id", "title").Where("id IN ?", ids).Find(&activities).Error; err != nil {
		return nil, fmt.Errorf("list activity titles: %w", err)
	}
	for i := range activities {
		titles[activities[i].ID] = activities[i].Title
	}
	return titles, nil
}
