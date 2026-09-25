package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/orienteering/platform/internal/constants"
	"github.com/orienteering/platform/internal/util"
)

// ContextUserKey 当前用户上下文键。
const ContextUserKey = "current_user"

// AuthClaims 注入到上下文的用户信息。
type AuthClaims struct {
	UserID   int64
	Username string
	Role     string
}

// Auth 认证中间件：解析 Bearer Token 并注入用户上下文。
func Auth(secret string, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized,
				"未登录或缺少 Authorization 请求头，请携带 Bearer token")
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := util.ParseToken(secret, tokenStr)
		if err != nil {
			logger.Warn("auth token invalid", "request_id", c.GetString("request_id"), "err", err)
			util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized,
				"登录凭证无效或已过期，请重新登录")
			c.Abort()
			return
		}
		c.Set(ContextUserKey, &AuthClaims{
			UserID: claims.UserID, Username: claims.Username, Role: claims.Role,
		})
		c.Next()
	}
}

// CurrentUser 从上下文读取当前用户。
func CurrentUser(c *gin.Context) (*AuthClaims, bool) {
	val, ok := c.Get(ContextUserKey)
	if !ok {
		return nil, false
	}
	claims, ok := val.(*AuthClaims)
	return claims, ok
}

// UserID 读取当前用户 ID。
func UserID(c *gin.Context) int64 {
	if claims, ok := CurrentUser(c); ok {
		return claims.UserID
	}
	return 0
}

// Role 读取当前用户角色。
func Role(c *gin.Context) string {
	if claims, ok := CurrentUser(c); ok {
		return claims.Role
	}
	return ""
}

var _ = fmt.Sprintf
