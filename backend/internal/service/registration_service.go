package service

import (
	"errors"
	"fmt"
	"log/slog"
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

// ApplyResult 报名结果：返回报名记录与是否进入候补。
type ApplyResult struct {
	Registration *model.Registration
	Waitlisted   bool
	// WaitlistAhead 进入候补时，前面还有多少队。
	WaitlistAhead int64
}

// Apply 团队报名活动（事务：校验活动已发布、操作者是队长、未重复报名；
// 待审核同样占用名额，名额不足时按提交顺序进入候补）。
func (s *RegistrationService) Apply(teamID, activityID, operatorID int64, logger *slog.Logger) (*ApplyResult, error) {
	result := &ApplyResult{}
	err := s.repo.WithTx(func(tx *gorm.DB) error {
		// 先锁活动行：同一活动的报名/审核/拒绝全部经此串行化，名额计数不会出现并发幻读。
		activity, err := s.activity.GetByIDForUpdate(tx, activityID)
		if err != nil {
			return err
		}
		if activity.Status != constants.ActivityStatusPublished {
			return util.NewAppError(constants.CodeActivityClosed,
				fmt.Sprintf(constants.MsgInvalidStatus, activity.Title, activity.Status, "apply"), nil)
		}
		// 锁团队行（所有报名事务都遵循 活动行 → 团队行 的加锁顺序，避免死锁），并校验队长身份。
		team, err := s.team.GetByIDForUpdate(tx, teamID)
		if err != nil {
			return err
		}
		if team.CaptainID != operatorID {
			return util.NewAppError(constants.CodeForbidden,
				fmt.Sprintf("用户[%d]不是团队[%s]的队长，不能报名", operatorID, team.Name), nil)
		}
		// 同一队重复提交：唯一索引 idx_team_activity 兜底，行锁内先做一次友好判定。
		existing, err := s.repo.GetByTeamAndActivityForUpdate(tx, teamID, activityID)
		if err == nil && existing != nil {
			return util.NewAppError(constants.CodeAlreadyApplied,
				fmt.Sprintf(constants.MsgAlreadyApplied, team.Name, activity.Title), nil)
		}
		if err != nil && !errors.Is(err, repository.ErrNotFound) {
			return err
		}
		// 占用名额计数（pending/approved/finished）必须在活动行锁内读取，两个并发报名不会同时通过判定。
		occupied, err := s.activity.CountOccupied(tx, activityID)
		if err != nil {
			return err
		}
		status := constants.RegistrationStatusPending
		if int(occupied) >= activity.MaxTeams {
			status = constants.RegistrationStatusWaitlisted
		}
		reg := &model.Registration{
			TeamID: teamID, ActivityID: activityID,
			Status: status, RegisteredAt: time.Now(),
		}
		if err := s.repo.Create(tx, reg); err != nil {
			// 并发下唯一索引兜底：两个相同 (team_id, activity_id) 的插入只有一条成功。
			if isDuplicateKeyErr(err) {
				return util.NewAppError(constants.CodeAlreadyApplied,
					fmt.Sprintf(constants.MsgAlreadyApplied, team.Name, activity.Title), nil)
			}
			return err
		}
		result.Registration = reg
		if status == constants.RegistrationStatusWaitlisted {
			result.Waitlisted = true
			result.WaitlistAhead, err = s.repo.CountWaitlistedBefore(tx, activityID, reg.ID)
			if err != nil {
				return err
			}
		}
		logger.Info(fmt.Sprintf(constants.LogTeamApply, teamID, activityID, status, occupied, activity.MaxTeams))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Approve 审核通过报名（管理员）。候补队伍不能直接通过，必须先按顺序递补为待审核。
func (s *RegistrationService) Approve(regID int64, operatorID int64, logger *slog.Logger) (*model.Registration, error) {
	var approved *model.Registration
	err := s.repo.WithTx(func(tx *gorm.DB) error {
		// 事务内先不加锁取活动 id，随后统一按“活动行 → 报名行”顺序加锁，与 Apply/Reject 一致，杜绝死锁。
		reg, err := s.repo.GetByIDTx(tx, regID)
		if err != nil {
			return err
		}
		activity, err := s.activity.GetByIDForUpdate(tx, reg.ActivityID)
		if err != nil {
			return err
		}
		reg, err = s.repo.GetByIDForUpdate(tx, regID)
		if err != nil {
			return err
		}
		if reg.Status != constants.RegistrationStatusPending {
			return util.NewAppError(constants.CodeConflict,
				fmt.Sprintf("报名记录[%d]当前状态[%s]不允许审核", regID, reg.Status), nil)
		}
		if activity.Status != constants.ActivityStatusPublished {
			return util.NewAppError(constants.CodeActivityClosed,
				fmt.Sprintf(constants.MsgInvalidStatus, activity.Title, activity.Status, "approve"), nil)
		}
		// 待审核本就占用名额，pending→approved 不改变占用计数；CAS 防止两个管理员同时操作同一记录。
		affected, err := s.repo.UpdateStatusCAS(tx, regID,
			constants.RegistrationStatusApproved, constants.RegistrationStatusPending)
		if err != nil {
			return err
		}
		if affected == 0 {
			return util.NewAppError(constants.CodeConflict,
				fmt.Sprintf("报名记录[%d]状态已被其他管理员变更，请刷新后重试", regID), nil)
		}
		reg.Status = constants.RegistrationStatusApproved
		approved = reg
		return nil
	})
	if err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf(constants.LogTeamApprove, regID, approved.TeamID))
	_ = operatorID
	_ = s.leaderboard.Invalidate(approved.ActivityID)
	return s.repo.GetByID(regID)
}

// RejectResult 拒绝报名结果：被拒绝的记录与自动递补的候补记录（可能为 nil）。
type RejectResult struct {
	Registration *model.Registration
	// Promoted 拒绝待审核释放名额后，自动递补为待审核的最早候补记录。
	Promoted *model.Registration
}

// Reject 拒绝报名（管理员）。拒绝待审核会释放一个名额，
// 同事务内把候补队列中提交最早的一队自动递补为待审核。
func (s *RegistrationService) Reject(regID int64, logger *slog.Logger) (*RejectResult, error) {
	result := &RejectResult{}
	err := s.repo.WithTx(func(tx *gorm.DB) error {
		// 统一按“活动行 → 报名行”顺序加锁，与 Apply/Approve 一致，杜绝死锁。
		pre, err := s.repo.GetByIDTx(tx, regID)
		if err != nil {
			return err
		}
		if _, err := s.activity.GetByIDForUpdate(tx, pre.ActivityID); err != nil {
			return err
		}
		reg, err := s.repo.GetByIDForUpdate(tx, regID)
		if err != nil {
			return err
		}
		if reg.Status != constants.RegistrationStatusPending && reg.Status != constants.RegistrationStatusWaitlisted {
			return util.NewAppError(constants.CodeConflict,
				fmt.Sprintf("报名记录[%d]当前状态[%s]不允许拒绝", regID, reg.Status), nil)
		}
		oldStatus := reg.Status
		// CAS：两个管理员同时拒绝/审核同一记录时只有一人成功，名额不会异常变化。
		affected, err := s.repo.UpdateStatusCAS(tx, regID,
			constants.RegistrationStatusRejected, oldStatus)
		if err != nil {
			return err
		}
		if affected == 0 {
			return util.NewAppError(constants.CodeConflict,
				fmt.Sprintf("报名记录[%d]状态已被其他管理员变更，请刷新后重试", regID), nil)
		}
		reg.Status = constants.RegistrationStatusRejected
		result.Registration = reg

		// 只有拒绝“占用名额”的待审核记录才释放名额；候补本身不占位，拒绝候补不触发递补。
		if oldStatus == constants.RegistrationStatusPending {
			first, err := s.repo.FirstWaitlistedForUpdate(tx, reg.ActivityID)
			if err == nil {
				promoted, perr := s.repo.UpdateStatusCAS(tx, first.ID,
					constants.RegistrationStatusPending, constants.RegistrationStatusWaitlisted)
				if perr != nil {
					return perr
				}
				if promoted > 0 {
					first.Status = constants.RegistrationStatusPending
					result.Promoted = first
					logger.Info(fmt.Sprintf(constants.LogTeamWaitlistPromote, reg.ActivityID, first.ID, first.TeamID))
				}
			} else if !errors.Is(err, repository.ErrNotFound) {
				return err
			}
		}
		logger.Info(fmt.Sprintf(constants.LogTeamReject, regID, reg.TeamID, promotedID(result.Promoted)))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
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

// ListByActivity 查询活动的报名列表（管理员审核页，按提交顺序）。
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

// ToViews 批量把报名记录转为展示视图：补全团队名/活动标题，并为候补记录计算前面还有几队。
func (s *RegistrationService) ToViews(regs []model.Registration) ([]dto.RegistrationView, error) {
	teamIDs := make([]int64, 0, len(regs))
	activityIDs := make([]int64, 0, len(regs))
	seenTeam := map[int64]struct{}{}
	seenActivity := map[int64]struct{}{}
	waitlistRegIDs := map[int64][]int64{} // 活动 -> 视图内该活动的候补报名 id
	for i := range regs {
		if _, ok := seenTeam[regs[i].TeamID]; !ok {
			seenTeam[regs[i].TeamID] = struct{}{}
			teamIDs = append(teamIDs, regs[i].TeamID)
		}
		if _, ok := seenActivity[regs[i].ActivityID]; !ok {
			seenActivity[regs[i].ActivityID] = struct{}{}
			activityIDs = append(activityIDs, regs[i].ActivityID)
		}
		if regs[i].Status == constants.RegistrationStatusWaitlisted {
			waitlistRegIDs[regs[i].ActivityID] = append(waitlistRegIDs[regs[i].ActivityID], regs[i].ID)
		}
	}
	teamNames, err := s.team.NamesByIDs(teamIDs)
	if err != nil {
		return nil, err
	}
	activityTitles, err := s.activity.TitlesByIDs(activityIDs)
	if err != nil {
		return nil, err
	}
	// 排位必须按全局候补队列计算（mine 列表经过滤，不能用列表内相对位置）；每个活动一次查询。
	aheadByID := map[int64]int64{}
	for activityID, ids := range waitlistRegIDs {
		rows := []struct {
			ID    int64
			Ahead int64
		}{}
		if err := s.repo.CountWaitlistedBeforeIDs(nil, activityID, ids, &rows); err != nil {
			return nil, err
		}
		for _, row := range rows {
			aheadByID[row.ID] = row.Ahead
		}
	}
	views := make([]dto.RegistrationView, 0, len(regs))
	for i := range regs {
		r := &regs[i]
		view := ToRegistrationView(r, teamNames[r.TeamID], activityTitles[r.ActivityID])
		if r.Status == constants.RegistrationStatusWaitlisted {
			view.WaitlistAhead = int(aheadByID[r.ID])
		}
		views = append(views, view)
	}
	return views, nil
}

// ToView 单条报名记录转展示视图（候补排位一并计算）。
func (s *RegistrationService) ToView(r *model.Registration) (dto.RegistrationView, error) {
	views, err := s.ToViews([]model.Registration{*r})
	if err != nil {
		return dto.RegistrationView{}, err
	}
	return views[0], nil
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

// promotedID 安全取递补记录 id（无递补时为 0）。
func promotedID(r *model.Registration) int64 {
	if r == nil {
		return 0
	}
	return r.ID
}

// isDuplicateKeyErr 判断是否为唯一约束冲突（PostgreSQL 23505 / SQLite UNIQUE constraint）。
func isDuplicateKeyErr(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "23505") ||
		strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "unique constraint")
}
