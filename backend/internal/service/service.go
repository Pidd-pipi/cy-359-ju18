package service

import (
	"github.com/redis/go-redis/v9"

	"github.com/orienteering/platform/internal/config"
	"github.com/orienteering/platform/internal/repository"
	"github.com/orienteering/platform/pkg/ws"
)

// Services 聚合所有业务服务，handler 通过构造器注入。
type Services struct {
	User         *UserService
	Activity     *ActivityService
	Checkpoint   *CheckpointService
	Team         *TeamService
	Registration *RegistrationService
	Checkin      *CheckinService
	Product      *ProductService
	Redemption   *RedemptionService
	Favorite     *FavoriteService
	Audit        *AuditService
	Leaderboard  *LeaderboardService
}

// New 装配全部服务。
func New(
	cfg *config.Config,
	repos *repository.Repositories,
	rdb *redis.Client,
	hub *ws.Hub,
) *Services {
	leaderboard := NewLeaderboardService(repos.Registration, repos.Checkin, repos.Activity, rdb)
	return &Services{
		User:         NewUserService(repos.User, cfg),
		Activity:     NewActivityService(repos.Activity),
		Checkpoint:   NewCheckpointService(repos.Checkpoint, repos.Activity),
		Team:         NewTeamService(repos.Team, repos.User),
		Registration: NewRegistrationService(repos.Registration, repos.Activity, repos.Team, leaderboard),
		Checkin:      NewCheckinService(repos.Checkin, repos.Checkpoint, repos.Registration, repos.Activity, repos.User, leaderboard, hub),
		Product:      NewProductService(repos.Product),
		Redemption:   NewRedemptionService(repos.Redemption, repos.Product, repos.User),
		Favorite:     NewFavoriteService(repos.Favorite),
		Audit:        NewAuditService(repos.Audit),
		Leaderboard:  leaderboard,
	}
}
