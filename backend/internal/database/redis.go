package database

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/orienteering/platform/internal/config"
)

// NewRedis 创建 Redis 客户端（连接失败不阻塞启动，仅记录日志，供排行榜缓存降级）。
func NewRedis(cfg *config.Config, logger *slog.Logger) *redis.Client {
	if cfg.RedisAddr == "" {
		return nil
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       cfg.RedisDB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Warn("redis ping failed, leaderboard cache disabled", "addr", cfg.RedisAddr, "err", err)
		return rdb
	}
	logger.Info("redis connected", "addr", cfg.RedisAddr)
	return rdb
}
