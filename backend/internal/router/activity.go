package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/orienteering/platform/internal/handler"
)

// registerActivityRoutes 活动路由。
func registerActivityRoutes(
	g *gin.RouterGroup,
	h *handler.ActivityHandler,
	reg *handler.RegistrationHandler,
	checkin *handler.CheckinHandler,
	authMW gin.HandlerFunc,
	adminMW gin.HandlerFunc,
	logger *slog.Logger,
) {
	activities := g.Group("/activities")
	activities.GET("", h.List)
	activities.GET("/:id", h.Get)

	auth := activities.Group("", authMW)
	auth.POST("", h.Create)
	auth.PUT("/:id", h.Update)
	auth.POST("/:id/transition", h.Transition)
	auth.GET("/:id/leaderboard", checkin.Leaderboard)
	auth.GET("/:id/checkins", checkin.ListByActivity)
	auth.GET("/:id/registrations", reg.ListByActivity)

	admin := activities.Group("", authMW, adminMW)
	admin.DELETE("/:id", h.Delete)
	_ = logger
}
