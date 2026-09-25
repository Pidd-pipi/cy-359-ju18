package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/orienteering/platform/internal/handler"
	"github.com/orienteering/platform/internal/middleware"
)

// registerUserRoutes 用户路由。
func registerUserRoutes(g *gin.RouterGroup, h *handler.UserHandler, authMW gin.HandlerFunc, logger *slog.Logger) {
	users := g.Group("/users")
	users.POST("/register", h.Register)
	users.POST("/login", h.Login)
	auth := users.Group("", authMW)
	auth.GET("/profile", h.Profile)
	auth.PUT("/profile", h.UpdateProfile)
	admin := g.Group("/users", authMW, middleware.RequireRole(logger, "admin"))
	admin.GET("", h.List)
	admin.PATCH("/:id/disabled", h.SetDisabled)
}
