package middleware

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/orienteering/platform/internal/model"
	"github.com/orienteering/platform/internal/repository"
)

// Audit 操作审计中间件：为写操作（POST/PUT/PATCH/DELETE）记录审计日志。
func Audit(repo *repository.AuditRepository, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		method := c.Request.Method
		if method == "GET" || method == "HEAD" || c.Writer.Status() >= 400 {
			return
		}
		claims, ok := CurrentUser(c)
		userID := int64(0)
		username := "anonymous"
		if ok {
			userID = claims.UserID
			username = claims.Username
		}
		resourceType := strings.TrimPrefix(c.FullPath(), "/api/v1/")
		if idx := strings.Index(resourceType, "/"); idx > 0 {
			resourceType = resourceType[:idx]
		}
		log := &model.AuditLog{
			UserID:       userID,
			Username:     username,
			Action:       method + " " + c.FullPath(),
			ResourceType: resourceType,
			ResourceID:   c.Param("id"),
			Detail:       fmt.Sprintf("%s %s", method, c.Request.URL.RequestURI()),
			IP:           c.ClientIP(),
			RequestID:    c.GetString("request_id"),
		}
		if err := repo.Create(log); err != nil {
			logger.Error("audit middleware create failed", "err", err)
		}
	}
}
