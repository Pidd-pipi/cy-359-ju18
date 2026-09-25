package service

import (
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"strings"
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

// Apply 团队报名活动。
// 事务：活动行锁 + 行内计数。待审核也占用名额；名额不足时按提交顺序进入候补（waitlist）。
func (s *RegistrationService) Apply(teamID, activityID, operatorID int64, logger *slog.Logger) (*model.Registration, error) {
	var reg *model.Registration
	err := s.repo.WithTx(func(tx *gorm.DB) error {
		// 1. 锁定活动行：两个并发报名在此串行化，计数-判定-插入是原子的。
		activity, err := s.activity.GetByIDForUpdate(tx, activityID)
		if err != nil {
			return err
		}
		if activity.Status != constants.ActivityStatusPublished {
			return util.NewAppError(constants.CodeActivityClosed,
				fmt.Sprintf(constants.MsgInvalidStatus, activity.Title, activity.Status, "apply"), nil)
		}
		// 2. 锁定团队行并校验队长。
		team, err := s.team.GetByIDForUpdate(tx, teamID)
		if err != nil {
			return err
		}
		if team.CaptainID != operatorID {
			return util.NewAppError(constants.CodeForbidden,
				fmt.Sprintf("用户[%d]不是团队[%s]的队长，不能报名", operatorID, team.Name), nil)
		}
		// 3. 同队重复提交：唯一索引 idx_team_activity 兜底，先给出友好错误。
		existing, err := s.repo.GetByTeamAndActivityTx(tx, teamID, activityID)
		if err == nil && existing != nil {
			return util.NewAppError(constants.CodeAlreadyApplied,
				fmt.Sprintf(constants.MsgAlreadyApplied, team.Name, activity.Title), nil)
		}
		if err != nil && !errors.Is(err, repository.ErrNotFound) {
			return err
		}
		// 4. 行内统计占用名额（待审核 + 已通过 + 已完成），决定待审核还是候补。
		occupied, err := s.repo.CountOccupied(tx, activityID)
		if err != nil {
			return err
		}
		status := constants.RegistrationStatusPending
		if int(occupied) >= activity.MaxTeams {
			status = constants.RegistrationStatusWaitlist
		}
		reg = &model.Registration{
			TeamID: teamID, ActivityID: activityID,
			Status: status, RegisteredAt: time.Now(),
		}
		if err := s.repo.Create(tx, reg); err != nil {
			// 并发下同队重复提交由唯一索引拦截，转成业务错误，避免名额/状态错乱。
			if isDuplicateKeyErr(err) {
				return util.NewAppError(constants.CodeAlreadyApplied,
					fmt.Sprintf(constants.MsgAlreadyApplied, team.Name, activity.Title), nil)
			}
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf(constants.LogTeamApply, teamID, activityID) + fmt.Sprintf(" status=%s", reg.Status))
	return reg, nil
}

// Approve 审核通过报名（管理员）。
// 仅待审核可通过；待审核已占用名额，故通过不会让正式队伍超过上限。
// 条件更新（CAS）保证两个管理员同时操作时只有一个成功。
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
		// 锁定报名行并复查状态，防止另一管理员同时审核/拒绝。
		locked, err := s.repo.GetByIDForUpdate(tx, regID)
		if err != nil {
			return err
		}
		if locked.Status != constants.RegistrationStatusPending {
			return util.NewAppError(constants.CodeConflict,
				fmt.Sprintf("报名记录[%d]当前状态[%s]不允许审核", regID, locked.Status), nil)
		}
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
// 可拒绝待审核或候补；拒绝待审核会释放一个名额，最早候补队伍（按 registered_at, id）
// 在同一事务内自动补成待审核，名额不会空缺也不会超额。CAS + 行锁保证两个管理员并发操作安全。
func (s *RegistrationService) Reject(regID int64, logger *slog.Logger) (*model.Registration, error) {
	reg, err := s.repo.GetByID(regID)
	if err != nil {
		return nil, err
	}
	if reg.Status != constants.RegistrationStatusPending && reg.Status != constants.RegistrationStatusWaitlist {
		return nil, util.NewAppError(constants.CodeConflict,
			fmt.Sprintf("报名记录[%d]当前状态[%s]不允许拒绝", regID, reg.Status), nil)
	}
	var promotedID int64
	err = s.repo.WithTx(func(tx *gorm.DB) error {
		// 1. 锁定报名行并复查状态，防止与另一管理员的审核/拒绝并发冲突。
		locked, err := s.repo.GetByIDForUpdate(tx, regID)
		if err != nil {
			return err
		}
		if locked.Status != constants.RegistrationStatusPending && locked.Status != constants.RegistrationStatusWaitlist {
			return util.NewAppError(constants.CodeConflict,
				fmt.Sprintf("报名记录[%d]当前状态[%s]不允许拒绝", regID, locked.Status), nil)
		}
		wasPending := locked.Status == constants.RegistrationStatusPending
		// 2. CAS 条件更新：只有状态未变才置为已拒绝。
		affected, err := s.repo.UpdateStatusIf(tx, regID, locked.Status, constants.RegistrationStatusRejected)
		if err != nil {
			return err
		}
		if affected == 0 {
			return util.NewAppError(constants.CodeConflict,
				fmt.Sprintf("报名记录[%d]状态已变更，请刷新后重试", regID), nil)
		}
		// 3. 仅当拒绝的是待审核（释放了名额）时，最早候补自动补位为待审核。
		if wasPending {
			next, err := s.repo.GetEarliestWaitlist(tx, reg.ActivityID)
			if err != nil {
				return err
			}
			if next != nil {
				rows, err := s.repo.UpdateStatusIf(tx, next.ID,
					constants.RegistrationStatusWaitlist, constants.RegistrationStatusPending)
				if err != nil {
					return err
				}
				if rows > 0 {
					promotedID = next.ID
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf(constants.LogTeamApprove, regID, reg.TeamID) +
		fmt.Sprintf(" action=reject promoted=%d", promotedID))
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

// ListByActivity 查询活动的报名列表（含候补排位）。
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

// ToViews 将报名记录批量转展示视图，并为候补记录计算前面还有几队。
// 候补排位按 (registered_at ASC, id ASC) 计数；同活动内一次批量计算，调用方传入的切片顺序无关。
func (s *RegistrationService) ToViews(regs []model.Registration) ([]dto.RegistrationView, error) {
	views := make([]dto.RegistrationView, 0, len(regs))
	if len(regs) == 0 {
		return views, nil
	}
	// 收集各活动候补排队顺序（先按提交时间、id 升序稳定排序后计数）。
	waitlistOrder := make([]*model.Registration, 0, len(regs))
	for i := range regs {
		if regs[i].Status == constants.RegistrationStatusWaitlist {
			waitlistOrder = append(waitlistOrder, &regs[i])
		}
	}
	sort.Slice(waitlistOrder, func(i, j int) bool {
		if !waitlistOrder[i].RegisteredAt.Equal(waitlistOrder[j].RegisteredAt) {
			return waitlistOrder[i].RegisteredAt.Before(waitlistOrder[j].RegisteredAt)
		}
		return waitlistOrder[i].ID < waitlistOrder[j].ID
	})
	aheadMap := make(map[int64]int, len(waitlistOrder))
	counters := make(map[int64]int)
	for _, r := range waitlistOrder {
		aheadMap[r.ID] = counters[r.ActivityID]
		counters[r.ActivityID]++
	}
	for i := range regs {
		v := ToRegistrationView(&regs[i], "", "")
		if regs[i].Status == constants.RegistrationStatusWaitlist {
			v.WaitlistAhead = aheadMap[regs[i].ID]
		}
		views = append(views, v)
	}
	return views, nil
}

// ToView 单条报名记录转展示视图（候补记录实时查询前面还有几队）。
func (s *RegistrationService) ToView(r *model.Registration) (dto.RegistrationView, error) {
	v := ToRegistrationView(r, "", "")
	if r.Status == constants.RegistrationStatusWaitlist {
		ahead, err := s.repo.CountWaitlistAheadDB(r.ActivityID, r.RegisteredAt, r.ID)
		if err != nil {
			return v, err
		}
		v.WaitlistAhead = int(ahead)
	}
	return v, nil
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

// isDuplicateKeyErr 判断唯一约束冲突（PostgreSQL: duplicate key；SQLite: UNIQUE constraint）。
func isDuplicateKeyErr(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate key") || strings.Contains(msg, "unique constraint")
}
