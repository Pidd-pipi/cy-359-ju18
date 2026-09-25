package service

import (
	"fmt"
	"log/slog"

	"gorm.io/gorm"

	"github.com/orienteering/platform/internal/constants"
	"github.com/orienteering/platform/internal/dto"
	"github.com/orienteering/platform/internal/model"
	"github.com/orienteering/platform/internal/repository"
	"github.com/orienteering/platform/internal/util"
)

// ActivityService 活动业务服务。
type ActivityService struct {
	repo *repository.ActivityRepository
}

// NewActivityService 构造活动服务。
func NewActivityService(repo *repository.ActivityRepository) *ActivityService {
	return &ActivityService{repo: repo}
}

// Create 创建活动（默认草稿状态）。
func (s *ActivityService) Create(creatorID int64, req *dto.CreateActivityRequest, logger *slog.Logger) (*model.Activity, error) {
	if req.EndTime.Before(req.StartTime) {
		return nil, util.NewAppError(constants.CodeBadRequest,
			"活动[%s]结束时间早于开始时间", nil)
	}
	activity := &model.Activity{
		Title:                req.Title,
		Description:          req.Description,
		Difficulty:           req.Difficulty,
		DurationMinutes:      req.DurationMinutes,
		EquipmentRequirement: req.EquipmentRequirement,
		StartTime:            req.StartTime,
		EndTime:              req.EndTime,
		Status:               constants.ActivityStatusDraft,
		CreatorID:            creatorID,
		StartLat:             req.StartLat,
		StartLng:             req.StartLng,
		EndLat:               req.EndLat,
		EndLng:               req.EndLng,
		Address:              req.Address,
		MaxTeams:             req.MaxTeams,
	}
	if err := s.repo.Create(activity); err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf(constants.LogActivityCreate, req.Title, req.Difficulty, creatorID))
	return activity, nil
}

// Update 编辑活动（仅草稿状态允许）。
func (s *ActivityService) Update(activityID int64, req *dto.CreateActivityRequest, logger *slog.Logger) (*model.Activity, error) {
	activity, err := s.repo.GetByID(activityID)
	if err != nil {
		return nil, err
	}
	if activity.Status != constants.ActivityStatusDraft {
		return nil, util.NewAppError(constants.CodeActivityClosed,
			fmt.Sprintf(constants.MsgInvalidStatus, activity.Title, activity.Status, "update"), nil)
	}
	activity.Title = req.Title
	activity.Description = req.Description
	activity.Difficulty = req.Difficulty
	activity.DurationMinutes = req.DurationMinutes
	activity.EquipmentRequirement = req.EquipmentRequirement
	activity.StartTime = req.StartTime
	activity.EndTime = req.EndTime
	activity.StartLat = req.StartLat
	activity.StartLng = req.StartLng
	activity.EndLat = req.EndLat
	activity.EndLng = req.EndLng
	activity.Address = req.Address
	activity.MaxTeams = req.MaxTeams
	if err := s.repo.Update(activity); err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf(constants.LogActivityStatus, activityID, "draft", activity.Status, activity.CreatorID))
	return activity, nil
}

// Transition 活动状态机：draft→published→ongoing→finished；draft/ongoing 可取消。
func (s *ActivityService) Transition(activityID int64, target string, operatorID int64, logger *slog.Logger) (*model.Activity, error) {
	activity, err := s.repo.GetByID(activityID)
	if err != nil {
		return nil, err
	}
	from := activity.Status
	valid := false
	switch target {
	case constants.ActivityStatusPublished:
		valid = from == constants.ActivityStatusDraft
	case constants.ActivityStatusOngoing:
		valid = from == constants.ActivityStatusPublished
	case constants.ActivityStatusFinished:
		valid = from == constants.ActivityStatusOngoing
	case constants.ActivityStatusCancelled:
		valid = from == constants.ActivityStatusDraft || from == constants.ActivityStatusPublished
	}
	if !valid {
		return nil, util.NewAppError(constants.CodeActivityClosed,
			fmt.Sprintf(constants.MsgInvalidStatus, activity.Title, from, target), nil)
	}
	activity.Status = target
	if err := s.repo.Update(activity); err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf(constants.LogActivityStatus, activityID, from, target, operatorID))
	return activity, nil
}

// List 分页查询活动。
func (s *ActivityService) List(page, pageSize, offset int, status, difficulty, keyword string) ([]model.Activity, int64, error) {
	return s.repo.List(page, pageSize, offset, status, difficulty, keyword)
}

// Get 获取活动详情（含打卡点）。
func (s *ActivityService) Get(activityID int64) (*model.Activity, error) {
	return s.repo.GetByID(activityID)
}

// Delete 删除活动（管理员，仅草稿/已取消）。
func (s *ActivityService) Delete(activityID int64) error {
	activity, err := s.repo.GetByID(activityID)
	if err != nil {
		return err
	}
	if activity.Status != constants.ActivityStatusDraft && activity.Status != constants.ActivityStatusCancelled {
		return util.NewAppError(constants.CodeActivityClosed,
			fmt.Sprintf(constants.MsgInvalidStatus, activity.Title, activity.Status, "delete"), nil)
	}
	return s.repo.Delete(activityID)
}

// CountCheckpoints 复用：统计活动打卡点数量。
func (s *ActivityService) CountCheckpoints(activityID int64) (int64, error) {
	return s.repo.CountCheckpoints(activityID)
}

// CountRegistrations 复用：统计活动已占用名额的团队数量（待审核 + 已通过 + 已完成）。
func (s *ActivityService) CountRegistrations(activityID int64) (int64, error) {
	return s.repo.CountRegistrations(activityID)
}

// CountWaitlisted 复用：统计活动候补队列中的团队数量。
func (s *ActivityService) CountWaitlisted(activityID int64) (int64, error) {
	return s.repo.CountWaitlisted(activityID)
}

// ToView 转换为展示视图。
func ToActivityView(a *model.Activity, checkpointCount, teamCount, waitlistCount int64) dto.ActivityView {
	cps := make([]dto.CheckpointView, 0, len(a.Checkpoints))
	for _, cp := range a.Checkpoints {
		cps = append(cps, ToCheckpointView(&cp))
	}
	return dto.ActivityView{
		ID:                   a.ID,
		Title:                a.Title,
		Description:          a.Description,
		Difficulty:           a.Difficulty,
		DurationMinutes:      a.DurationMinutes,
		EquipmentRequirement: a.EquipmentRequirement,
		StartTime:            a.StartTime,
		EndTime:              a.EndTime,
		Status:               a.Status,
		CreatorID:            a.CreatorID,
		StartLat:             a.StartLat,
		StartLng:             a.StartLng,
		EndLat:               a.EndLat,
		EndLng:               a.EndLng,
		Address:              a.Address,
		MaxTeams:             a.MaxTeams,
		CheckpointCount:      int(checkpointCount),
		TeamCount:            int(teamCount),
		WaitlistCount:        int(waitlistCount),
		Checkpoints:          cps,
		CreatedAt:            a.CreatedAt,
	}
}

// GetTx 供其他服务在事务中读取活动（并发锁）。
func (s *ActivityService) GetTx(tx *gorm.DB, activityID int64) (*model.Activity, error) {
	return s.repo.GetByIDForUpdate(tx, activityID)
}

