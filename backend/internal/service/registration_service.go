package service

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"

	"github.com/orienteering/platform/internal/constants"
	"github.com/orienteering/platform/internal/dto"
	"github.com/orienteering/platform/internal/model"
	"github.com/orienteering/platform/internal/repository"
	"github.com/orienteering/platform/internal/util"
)

// RegistrationService 团队报名业务服务。
type RegistrationService struct {
	repo        *repository.RegistrationRepository
	activity    *repository.ActivityRepository
	team        *repository.TeamRepository
	leaderboard *LeaderboardService
}

// NewRegistrationService 构造报名服务。
func NewRegistrationService(repo *repository.RegistrationRepository, activity *repository.ActivityRepository, team *repository.TeamRepository, leaderboard *LeaderboardService) *RegistrationService {
	return &RegistrationService{repo: repo, activity: activity, team: team, leaderboard: leaderboard}
}

// Apply 团队报名活动（事务：校验活动已发布、团队存在、操作者是队长、未重复报名）。
func (s *RegistrationService) Apply(teamID, activityID, operatorID int64, logger *slog.Logger) (*model.Registration, error) {
	var reg *model.Registration
	err := s.repo.WithTx(func(tx *gorm.DB) error {
		activity, err := s.activity.GetByIDForUpdate(tx, activityID)
		if err != nil {
			return err
		}
		if activity.Status != constants.ActivityStatusPublished {
			return util.NewAppError(constants.CodeActivityClosed,
				fmt.Sprintf(constants.MsgInvalidStatus, activity.Title, activity.Status, "apply"), nil)
		}
		teamCount, err := s.activity.CountRegistrations(activityID)
		if err != nil {
			return err
		}
		if int(teamCount) >= activity.MaxTeams {
			return util.NewAppError(constants.CodeConflict,
				fmt.Sprintf("活动[%s]报名团队已满[%d/%d]", activity.Title, teamCount, activity.MaxTeams), nil)
		}
		team, err := s.team.GetByIDForUpdate(tx, teamID)
		if err != nil {
			return err
		}
		if team.CaptainID != operatorID {
			return util.NewAppError(constants.CodeForbidden,
				fmt.Sprintf("用户[%d]不是团队[%s]的队长，不能报名", operatorID, team.Name), nil)
		}
		existing, err := s.repo.GetByTeamAndActivity(teamID, activityID)
		if err == nil && existing != nil {
			return util.NewAppError(constants.CodeAlreadyApplied,
				fmt.Sprintf(constants.MsgAlreadyApplied, "", ""), nil)
		}
		if err != nil && !errors.Is(err, repository.ErrNotFound) {
			return err
		}
		reg = &model.Registration{
			TeamID: teamID, ActivityID: activityID,
			Status: constants.RegistrationStatusPending, RegisteredAt: time.Now(),
		}
		return s.repo.Create(tx, reg)
	})
	if err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf(constants.LogTeamApply, teamID, activityID))
	return reg, nil
}

// Approve 审核通过报名（管理员）。
func (s *RegistrationService) Approve(regID int64, operatorID int64, logger *slog.Logger) (*model.Registration, error) {
	reg, err := s.repo.GetByID(regID)
	if err != nil {
		return nil, err
	}
	if reg.Status != constants.RegistrationStatusPending {
		return nil, util.NewAppError(constants.CodeConflict,
			fmt.Sprintf("报名记录[%d]当前状态[%s]不允许审核", regID, reg.Status), nil)
	}
	err = s.repo.WithTx(func(tx *gorm.DB) error {
		activity, err := s.activity.GetByIDForUpdate(tx, reg.ActivityID)
		if err != nil {
			return err
		}
		if activity.Status != constants.ActivityStatusPublished {
			return util.NewAppError(constants.CodeActivityClosed,
				fmt.Sprintf(constants.MsgInvalidStatus, activity.Title, activity.Status, "approve"), nil)
		}
		return s.repo.UpdateStatus(tx, regID, constants.RegistrationStatusApproved)
	})
	if err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf(constants.LogTeamApprove, regID, reg.TeamID))
	_ = operatorID
	_ = s.leaderboard.Invalidate(reg.ActivityID)
	return s.repo.GetByID(regID)
}

// Reject 拒绝报名（管理员）。
func (s *RegistrationService) Reject(regID int64, logger *slog.Logger) (*model.Registration, error) {
	reg, err := s.repo.GetByID(regID)
	if err != nil {
		return nil, err
	}
	if reg.Status != constants.RegistrationStatusPending {
		return nil, util.NewAppError(constants.CodeConflict,
			fmt.Sprintf("报名记录[%d]当前状态[%s]不允许拒绝", regID, reg.Status), nil)
	}
	if err := s.repo.WithTx(func(tx *gorm.DB) error {
		return s.repo.UpdateStatus(tx, regID, constants.RegistrationStatusRejected)
	}); err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf(constants.LogTeamApprove, regID, reg.TeamID))
	return s.repo.GetByID(regID)
}

// Start 活动开始时为已通过团队记录开始时间（活动服务状态机调用）。
func (s *RegistrationService) Start(activityID int64, logger *slog.Logger) error {
	err := s.repo.WithTx(func(tx *gorm.DB) error {
		regs, err := s.repo.ListByActivityForUpdate(tx, activityID)
		if err != nil {
			return err
		}
		now := time.Now()
		for i := range regs {
			if regs[i].Status == constants.RegistrationStatusApproved && regs[i].StartTime == nil {
				regs[i].StartTime = &now
				if err := s.repo.Update(tx, &regs[i]); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf(constants.LogTeamFinish, activityID, 0))
	return nil
}

// ListByActivity 查询活动的报名列表。
func (s *RegistrationService) ListByActivity(activityID int64) ([]model.Registration, error) {
	return s.repo.ListByActivity(activityID)
}

// ListMine 查询我所在团队的报名记录。
func (s *RegistrationService) ListMine(userID int64) ([]model.Registration, error) {
	return s.repo.ListByUserTeams(userID)
}

// Get 查询报名详情。
func (s *RegistrationService) Get(regID int64) (*model.Registration, error) {
	return s.repo.GetByID(regID)
}

// ToRegistrationView 报名记录转展示视图。
func ToRegistrationView(r *model.Registration, teamName, activityTitle string) dto.RegistrationView {
	duration := ""
	if r.TotalSeconds > 0 {
		duration = util.FormatDuration(r.TotalSeconds)
	}
	return dto.RegistrationView{
		ID: r.ID, TeamID: r.TeamID, TeamName: teamName,
		ActivityID: r.ActivityID, ActivityTitle: activityTitle,
		Status: r.Status, StartTime: r.StartTime, FinishTime: r.FinishTime,
		TotalSeconds: r.TotalSeconds, Duration: duration, RegisteredAt: r.RegisteredAt,
	}
}
