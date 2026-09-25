package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/orienteering/platform/internal/constants"
)

// RequestID 请求追踪中间件：生成/透传 X-Request-ID。
func RequestID(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader("X-Request-ID")
		if rid == "" {
			rid = newRequestID()
		}
		c.Set("request_id", rid)
		c.Header("X-Request-ID", rid)
		logger.Info(constants.LogRequestStart, "request_id", rid, "method", c.Request.Method, "path", c.Request.URL.Path)
		c.Next()
		logger.Info(constants.LogRequestEnd, "request_id", rid, "method", c.Request.Method, "path", c.Request.URL.Path, "status", c.Writer.Status())
	}
}

func newRequestID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(b)
}
