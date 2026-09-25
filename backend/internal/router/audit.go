package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/orienteering/platform/internal/handler"
)

// registerAuditRoutes 审计日志路由（仅管理员）。
func registerAuditRoutes(g *gin.RouterGroup, h *handler.AuditHandler, authMW gin.HandlerFunc, adminMW gin.HandlerFunc, logger *slog.Logger) {
	audit := g.Group("/audit-logs")
	admin := audit.Group("", authMW, adminMW)
	admin.GET("", h.List)
	_ = logger
}
