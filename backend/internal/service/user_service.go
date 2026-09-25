package service

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/orienteering/platform/internal/config"
	"github.com/orienteering/platform/internal/constants"
	"github.com/orienteering/platform/internal/dto"
	"github.com/orienteering/platform/internal/model"
	"github.com/orienteering/platform/internal/repository"
	"github.com/orienteering/platform/internal/util"
)

// UserService 用户业务服务。
type UserService struct {
	repo *repository.UserRepository
	cfg  *config.Config
}

// NewUserService 构造用户服务。
func NewUserService(repo *repository.UserRepository, cfg *config.Config) *UserService {
	return &UserService{repo: repo, cfg: cfg}
}

// Register 注册用户，校验用户名唯一并加密密码。
func (s *UserService) Register(req *dto.RegisterRequest) (*model.User, error) {
	exists, err := s.repo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("register user: %w", err)
	}
	if exists {
		return nil, util.NewAppError(constants.CodeUserExists,
			fmt.Sprintf("用户名[%s]已存在", req.Username), nil)
	}
	hash, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("register user: hash password: %w", err)
	}
	user := &model.User{
		Username:     req.Username,
		PasswordHash: hash,
		Nickname:     req.Nickname,
		Email:        req.Email,
		Phone:        req.Phone,
		Role:         constants.RoleUser,
		Points:       100,
	}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// Login 校验凭证并签发 JWT。
func (s *UserService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeLoginFailed,
				"用户名或密码错误", nil)
		}
		return nil, fmt.Errorf("login user: %w", err)
	}
	if !util.CheckPassword(user.PasswordHash, req.Password) {
		return nil, util.NewAppError(constants.CodeLoginFailed,
			"用户名或密码错误", nil)
	}
	if user.Disabled {
		return nil, util.NewAppError(constants.CodeUserDisabled,
			fmt.Sprintf("用户[%s]已被禁用", user.Username), nil)
	}
	token, err := util.GenerateToken(s.cfg.JWTSecret, s.cfg.JWTExpire, user.ID, user.Username, user.Role)
	if err != nil {
		return nil, fmt.Errorf("login user: generate token: %w", err)
	}
	return &dto.LoginResponse{
		Token: token,
		User: dto.UserBrief{
			ID: user.ID, Username: user.Username, Nickname: user.Nickname,
			Role: user.Role, Points: user.Points,
		},
	}, nil
}

// GetProfile 获取用户信息。
func (s *UserService) GetProfile(userID int64) (*model.User, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateProfile 更新个人资料。
func (s *UserService) UpdateProfile(userID int64, req *dto.UpdateProfileRequest) (*model.User, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

// List 分页查询用户（管理员）。
func (s *UserService) List(page, pageSize, offset int, keyword string) ([]model.User, int64, error) {
	return s.repo.List(page, pageSize, offset, keyword)
}

// SetDisabled 禁用/启用用户（管理员）。
func (s *UserService) SetDisabled(userID int64, disabled bool) error {
	if _, err := s.repo.GetByID(userID); err != nil {
		return err
	}
	return s.repo.SetDisabled(userID, disabled)
}

// AddPoints 增加用户积分（service 内部复用，供打卡/兑换服务调用）。
func (s *UserService) AddPoints(userID int64, delta int, reason string, logger *slog.Logger) (int, error) {
	if delta == 0 {
		return 0, nil
	}
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return 0, err
	}
	newPoints := user.Points + delta
	if newPoints < 0 {
		return 0, util.NewAppError(constants.CodePointsNotEnough,
			fmt.Sprintf("用户[%s]积分[%d]不足，操作需要[%d]", user.Username, user.Points, -delta), nil)
	}
	if err := s.repo.AddPoints(userID, delta); err != nil {
		return 0, err
	}
	logger.Info(fmt.Sprintf(constants.LogUserPointsChange, userID, delta, newPoints, reason))
	return newPoints, nil
}
