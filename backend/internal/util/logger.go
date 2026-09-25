package util

import (
	"log/slog"
	"os"
)

// NewLogger 创建 slog 结构化日志器，统一 JSON 输出。
func NewLogger(env string) *slog.Logger {
	level := slog.LevelInfo
	if env == "development" {
		level = slog.LevelDebug
	}
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}
