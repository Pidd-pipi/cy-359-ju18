package service

import (
	"fmt"
	"log/slog"

	"github.com/orienteering/platform/internal/constants"
	"github.com/orienteering/platform/internal/dto"
	"github.com/orienteering/platform/internal/model"
	"github.com/orienteering/platform/internal/repository"
	"github.com/orienteering/platform/internal/util"
)

// CheckpointService 打卡点业务服务。
type CheckpointService struct {
	repo     *repository.CheckpointRepository
	activity *repository.ActivityRepository
}

// NewCheckpointService 构造打卡点服务。
func NewCheckpointService(repo *repository.CheckpointRepository, activity *repository.ActivityRepository) *CheckpointService {
	return &CheckpointService{repo: repo, activity: activity}
}

// Create 创建打卡点（仅活动草稿状态允许）。
func (s *CheckpointService) Create(activityID int64, req *dto.CreateCheckpointRequest, logger *slog.Logger) (*model.Checkpoint, error) {
	activity, err := s.activity.GetByID(activityID)
	if err != nil {
		return nil, err
	}
	if activity.Status != constants.ActivityStatusDraft {
		return nil, util.NewAppError(constants.CodeActivityClosed,
			fmt.Sprintf(constants.MsgInvalidStatus, activity.Title, activity.Status, "create_checkpoint"), nil)
	}
	cp := &model.Checkpoint{
		ActivityID:     activityID,
		Name:           req.Name,
		Sequence:       req.Sequence,
		Lat:            req.Lat,
		Lng:            req.Lng,
		Clue:           req.Clue,
		TaskType:       req.TaskType,
		TaskContent:    req.TaskContent,
		ExpectedAnswer: req.ExpectedAnswer,
		RadiusMeters:   req.RadiusMeters,
		QRCode:         req.QRCode,
	}
	if err := s.repo.Create(cp); err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf(constants.LogCheckpointCreate, activityID, req.Name, req.Sequence))
	return cp, nil
}

// Update 编辑打卡点（仅活动草稿状态允许）。
func (s *CheckpointService) Update(cpID int64, req *dto.CreateCheckpointRequest, logger *slog.Logger) (*model.Checkpoint, error) {
	cp, err := s.repo.GetByID(cpID)
	if err != nil {
		return nil, err
	}
	activity, err := s.activity.GetByID(cp.ActivityID)
	if err != nil {
		return nil, err
	}
	if activity.Status != constants.ActivityStatusDraft {
		return nil, util.NewAppError(constants.CodeActivityClosed,
			fmt.Sprintf(constants.MsgInvalidStatus, activity.Title, activity.Status, "update_checkpoint"), nil)
	}
	cp.Name = req.Name
	cp.Sequence = req.Sequence
	cp.Lat = req.Lat
	cp.Lng = req.Lng
	cp.Clue = req.Clue
	cp.TaskType = req.TaskType
	cp.TaskContent = req.TaskContent
	cp.ExpectedAnswer = req.ExpectedAnswer
	cp.RadiusMeters = req.RadiusMeters
	cp.QRCode = req.QRCode
	if err := s.repo.Update(cp); err != nil {
		return nil, err
	}
	return cp, nil
}

// Delete 删除打卡点（仅活动草稿状态允许）。
func (s *CheckpointService) Delete(cpID int64, logger *slog.Logger) error {
	cp, err := s.repo.GetByID(cpID)
	if err != nil {
		return err
	}
	activity, err := s.activity.GetByID(cp.ActivityID)
	if err != nil {
		return err
	}
	if activity.Status != constants.ActivityStatusDraft {
		return util.NewAppError(constants.CodeActivityClosed,
			fmt.Sprintf(constants.MsgInvalidStatus, activity.Title, activity.Status, "delete_checkpoint"), nil)
	}
	if err := s.repo.Delete(cpID); err != nil {
		return err
	}
	logger.Info(fmt.Sprintf(constants.LogCheckpointDelete, cpID, cp.ActivityID))
	return nil
}

// ListByActivity 查询活动的全部打卡点（按序号排序）。
func (s *CheckpointService) ListByActivity(activityID int64) ([]model.Checkpoint, error) {
	return s.repo.ListByActivity(activityID)
}

// CountByActivity 复用：统计活动打卡点数量。
func (s *CheckpointService) CountByActivity(activityID int64) (int64, error) {
	return s.repo.CountByActivity(activityID)
}

// ToCheckpointView 打卡点转展示视图。
func ToCheckpointView(cp *model.Checkpoint) dto.CheckpointView {
	return dto.CheckpointView{
		ID:           cp.ID,
		ActivityID:   cp.ActivityID,
		Name:         cp.Name,
		Sequence:     cp.Sequence,
		Lat:          cp.Lat,
		Lng:          cp.Lng,
		Clue:         cp.Clue,
		TaskType:     cp.TaskType,
		TaskContent:  cp.TaskContent,
		RadiusMeters: cp.RadiusMeters,
		QRCode:       cp.QRCode,
	}
}
