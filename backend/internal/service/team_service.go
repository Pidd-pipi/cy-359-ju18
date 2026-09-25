package service

import (
	"errors"
	"fmt"
	"log/slog"

	"gorm.io/gorm"

	"github.com/orienteering/platform/internal/constants"
	"github.com/orienteering/platform/internal/dto"
	"github.com/orienteering/platform/internal/model"
	"github.com/orienteering/platform/internal/repository"
	"github.com/orienteering/platform/internal/util"
)

// TeamService 团队业务服务。
type TeamService struct {
	repo *repository.TeamRepository
	user *repository.UserRepository
}

// NewTeamService 构造团队服务。
func NewTeamService(repo *repository.TeamRepository, user *repository.UserRepository) *TeamService {
	return &TeamService{repo: repo, user: user}
}

// Create 创建团队，事务内创建团队并写入队长成员关系。
func (s *TeamService) Create(captainID int64, req *dto.CreateTeamRequest, logger *slog.Logger) (*model.Team, error) {
	exists, err := s.repo.ExistsByName(req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, util.NewAppError(constants.CodeConflict,
			fmt.Sprintf("团队[%s]名称已存在", req.Name), nil)
	}
	var team *model.Team
	err = s.repo.WithTx(func(tx *gorm.DB) error {
		team = &model.Team{Name: req.Name, Slogan: req.Slogan, CaptainID: captainID}
		if err := s.repo.Create(tx, team); err != nil {
			return err
		}
		member := &model.TeamMember{TeamID: team.ID, UserID: captainID, Role: constants.TeamRoleCaptain}
		return s.repo.AddMember(tx, member)
	})
	if err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf(constants.LogTeamCreate, req.Name, captainID))
	return team, nil
}

// Join 加入团队（并发保护：行锁 + 成员数限制）。
func (s *TeamService) Join(teamID, userID int64, logger *slog.Logger) error {
	err := s.repo.WithTx(func(tx *gorm.DB) error {
		if _, err := s.repo.GetByIDForUpdate(tx, teamID); err != nil {
			return err
		}
		already, err := s.repo.IsMember(teamID, userID)
		if err != nil {
			return err
		}
		if already {
			return util.NewAppError(constants.CodeAlreadyMember,
				fmt.Sprintf(constants.MsgAlreadyJoined, "", ""), nil)
		}
		count, err := s.repo.CountMembers(teamID)
		if err != nil {
			return err
		}
		if count >= 8 {
			return util.NewAppError(constants.CodeTeamFull,
				fmt.Sprintf(constants.MsgTeamFull, teamID), nil)
		}
		member := &model.TeamMember{TeamID: teamID, UserID: userID, Role: constants.TeamRoleMember}
		return s.repo.AddMember(tx, member)
	})
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf(constants.LogTeamJoin, teamID, userID, constants.TeamRoleMember))
	return nil
}

// Leave 退出团队（队长不可退出）。
func (s *TeamService) Leave(teamID, userID int64, logger *slog.Logger) error {
	team, err := s.repo.GetByID(teamID)
	if err != nil {
		return err
	}
	if team.CaptainID == userID {
		return util.NewAppError(constants.CodeForbidden,
			"队长[%d]不能退出团队，请先转让", nil)
	}
	isMember, err := s.repo.IsMember(teamID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return util.NewAppError(constants.CodeNotFound,
			fmt.Sprintf("用户[%d]不在团队[%s]中", userID, team.Name), nil)
	}
	if err := s.repo.Leave(teamID, userID); err != nil {
		return err
	}
	logger.Info(fmt.Sprintf(constants.LogTeamJoin, teamID, userID, "leave"))
	return nil
}

// ListMine 查询我加入的团队。
func (s *TeamService) ListMine(userID int64) ([]model.Team, error) {
	return s.repo.ListByUser(userID)
}

// ListAll 分页查询全部团队（管理员）。
func (s *TeamService) ListAll(page, pageSize, offset int) ([]model.Team, int64, error) {
	return s.repo.ListAll(page, pageSize, offset)
}

// Get 获取团队详情。
func (s *TeamService) Get(teamID int64) (*model.Team, error) {
	return s.repo.GetByID(teamID)
}

// ToTeamView 团队转展示视图。
func ToTeamView(t *model.Team, users map[int64]string) dto.TeamView {
	members := make([]dto.TeamMemberView, 0, len(t.Members))
	for _, m := range t.Members {
		username := users[m.UserID]
		if username == "" {
			username = fmt.Sprintf("user_%d", m.UserID)
		}
		members = append(members, dto.TeamMemberView{
			UserID: m.UserID, Username: username, Nickname: username, Role: m.Role,
		})
	}
	return dto.TeamView{
		ID: t.ID, Name: t.Name, Slogan: t.Slogan, CaptainID: t.CaptainID,
		MemberCount: len(members), Members: members, CreatedAt: t.CreatedAt,
	}
}

var _ = errors.Is
