package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/orienteering/platform/internal/constants"
	"github.com/orienteering/platform/internal/util"
)

type bucket struct {
	tokens float64
	last   time.Time
}

// RateLimit 基于令牌桶的限流中间件（内存实现）。
func RateLimit(rate float64, capacity float64, logger *slog.Logger) gin.HandlerFunc {
	var mu sync.Mutex
	buckets := make(map[string]*bucket)
	return func(c *gin.Context) {
		key := c.ClientIP() + c.Request.URL.Path
		now := time.Now()
		mu.Lock()
		b, ok := buckets[key]
		if !ok {
			b = &bucket{tokens: capacity, last: now}
			buckets[key] = b
		}
		elapsed := now.Sub(b.last).Seconds()
		b.tokens += elapsed * rate
		if b.tokens > capacity {
			b.tokens = capacity
		}
		b.last = now
		if b.tokens < 1 {
			mu.Unlock()
			logger.Warn(fmt.Sprintf(constants.LogRateLimited, c.GetString("request_id"), c.Request.URL.Path))
			util.Fail(c, http.StatusTooManyRequests, constants.CodeBadRequest, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}
		b.tokens--
		mu.Unlock()
		c.Next()
	}
}
