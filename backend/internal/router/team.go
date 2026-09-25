package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/orienteering/platform/internal/handler"
)

// registerTeamRoutes 团队与报名路由。
func registerTeamRoutes(
	g *gin.RouterGroup,
	h *handler.TeamHandler,
	reg *handler.RegistrationHandler,
	checkin *handler.CheckinHandler,
	authMW gin.HandlerFunc,
	adminMW gin.HandlerFunc,
	logger *slog.Logger,
) {
	teams := g.Group("/teams")
	auth := teams.Group("", authMW)
	auth.POST("", h.Create)
	auth.GET("/mine", h.ListMine)
	auth.GET("/:id", h.Get)
	auth.POST("/:id/join", h.Join)
	auth.DELETE("/:id/leave", h.Leave)
	auth.GET("/:id/checkins", checkin.ListByTeam)
	auth.POST("/:id/checkin", checkin.Checkin)

	admin := teams.Group("", authMW, adminMW)
	admin.GET("", h.ListAll)

	regGroup := g.Group("/registrations")
	regAuth := regGroup.Group("", authMW)
	regAuth.POST("", reg.Apply)
	regAuth.GET("/mine", reg.ListMine)
	regAdmin := regGroup.Group("", authMW, adminMW)
	regAdmin.POST("/:id/approve", reg.Approve)
	regAdmin.POST("/:id/reject", reg.Reject)
	_ = logger
}
