package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/orienteering/platform/internal/handler"
)

// registerCheckpointRoutes 打卡点路由。
func registerCheckpointRoutes(g *gin.RouterGroup, h *handler.CheckpointHandler, authMW gin.HandlerFunc, adminMW gin.HandlerFunc, logger *slog.Logger) {
	checkpoints := g.Group("/checkpoints")
	admin := checkpoints.Group("", authMW, adminMW)
	admin.POST("/activities/:id", h.Create)
	admin.PUT("/:id", h.Update)
	admin.DELETE("/:id", h.Delete)
	_ = logger
}
