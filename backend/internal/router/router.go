package router

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/orienteering/platform/internal/config"
	"github.com/orienteering/platform/internal/constants"
	"github.com/orienteering/platform/internal/handler"
	"github.com/orienteering/platform/internal/middleware"
	"github.com/orienteering/platform/internal/repository"
	"github.com/orienteering/platform/internal/service"
	"github.com/orienteering/platform/pkg/ws"
)

// Setup 装配全部路由。
func Setup(
	cfg *config.Config,
	logger *slog.Logger,
	repos *repository.Repositories,
	services *service.Services,
	hub *ws.Hub,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// 全局中间件：请求 ID -> 恢复/错误处理 -> 请求日志 -> 限流 -> 审计
	r.Use(middleware.RequestID(logger))
	r.Use(middleware.ErrorHandler(logger))
	r.Use(middleware.RequestLog(logger))
	r.Use(middleware.RateLimit(20, 60, logger))
	r.Use(middleware.Audit(repos.Audit, logger))

	// 健康检查
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now().Format(time.RFC3339)})
	})

	api := r.Group("/api/v1")

	// 认证 / 授权中间件实例
	authMW := middleware.Auth(cfg.JWTSecret, logger)
	adminMW := middleware.RequireRole(logger, constants.RoleAdmin)

	// 处理器
	userHandler := handler.NewUserHandler(services.User, logger)
	activityHandler := handler.NewActivityHandler(services.Activity, services.Registration, logger)
	checkpointHandler := handler.NewCheckpointHandler(services.Checkpoint, logger)
	teamHandler := handler.NewTeamHandler(services.Team, logger)
	registrationHandler := handler.NewRegistrationHandler(services.Registration, logger)
	checkinHandler := handler.NewCheckinHandler(services.Checkin, services.Leaderboard, hub, cfg.JWTSecret, logger)
	// WebSocket 实时排行榜（依赖 checkinHandler）
	r.GET("/ws/leaderboard", checkinHandler.LeaderboardWS)
	productHandler := handler.NewProductHandler(services.Product, logger)
	redemptionHandler := handler.NewRedemptionHandler(services.Redemption, logger)
	favoriteHandler := handler.NewFavoriteHandler(services.Favorite, services.Activity, logger)
	auditHandler := handler.NewAuditHandler(services.Audit, logger)

	// 各实体路由注册
	registerUserRoutes(api, userHandler, authMW, logger)
	registerActivityRoutes(api, activityHandler, registrationHandler, checkinHandler, authMW, adminMW, logger)
	registerTeamRoutes(api, teamHandler, registrationHandler, checkinHandler, authMW, adminMW, logger)
	registerCheckpointRoutes(api, checkpointHandler, authMW, adminMW, logger)
	registerProductRoutes(api, productHandler, redemptionHandler, authMW, adminMW, logger)
	registerFavoriteRoutes(api, favoriteHandler, authMW, logger)
	registerAuditRoutes(api, auditHandler, authMW, adminMW, logger)

	return r
}
