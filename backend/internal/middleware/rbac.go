package middleware

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/orienteering/platform/internal/constants"
	"github.com/orienteering/platform/internal/util"
)

// RequireRole RBAC 权限中间件：仅允许指定角色访问。
func RequireRole(logger *slog.Logger, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := CurrentUser(c)
		if !ok {
			util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, "请先登录")
			c.Abort()
			return
		}
		allowed := false
		for _, role := range roles {
			if claims.Role == role {
				allowed = true
				break
			}
		}
		if !allowed {
			logger.Warn("rbac forbidden",
				"request_id", c.GetString("request_id"),
				"path", c.FullPath(),
				"role", claims.Role)
			util.Fail(c, http.StatusForbidden, constants.CodeForbidden,
				fmt.Sprintf(constants.MsgPermissionDenied, claims.Role))
			c.Abort()
			return
		}
		c.Next()
	}
}
