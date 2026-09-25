package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/orienteering/platform/internal/config"
	"github.com/orienteering/platform/internal/constants"
	"github.com/orienteering/platform/internal/database"
	"github.com/orienteering/platform/internal/model"
	"gorm.io/gorm"
	"github.com/orienteering/platform/internal/repository"
	"github.com/orienteering/platform/internal/router"
	"github.com/orienteering/platform/internal/service"
	"github.com/orienteering/platform/internal/util"
	"github.com/orienteering/platform/pkg/ws"
)

func main() {
	logger := util.NewLogger(os.Getenv("APP_ENV"))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("load config failed", "err", err)
		os.Exit(1)
	}
	logger.Info(fmt.Sprintf(constants.LogConfigLoaded, cfg.AppName, cfg.AppEnv, cfg.Port))

	db, err := database.New(cfg.DSN(), logger)
	if err != nil {
		logger.Error(fmt.Sprintf(constants.LogDBNotConnected, err))
		os.Exit(1)
	}
	logger.Info(fmt.Sprintf(constants.LogDBConnected, cfg.DBHost, cfg.DBName))

	rdb := database.NewRedis(cfg, logger)
	hub := ws.NewHub(logger)

	repos := repository.New(db)
	services := service.New(cfg, repos, rdb, hub)
	services.Leaderboard.SetLogger(logger)

	if err := seedAdmin(db, cfg, logger); err != nil {
		logger.Error("seed admin failed", "err", err)
		os.Exit(1)
	}

	engine := router.Setup(cfg, logger, repos, services, hub)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		logger.Info(fmt.Sprintf(constants.LogServerStarted, cfg.Port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(fmt.Sprintf(constants.LogServerStopped, err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server shutdown error", "err", err)
	}
	logger.Info("server stopped gracefully")
}

// seedAdmin 初始化默认管理员账号。
func seedAdmin(db *gorm.DB, cfg *config.Config, logger *slog.Logger) error {
	var count int64
	if err := db.Model(&model.User{}).Where("role = ?", constants.RoleAdmin).Count(&count).Error; err != nil {
		return fmt.Errorf("count admin: %w", err)
	}
	if count > 0 {
		return nil
	}
	hash, err := util.HashPassword("admin123456")
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}
	admin := &model.User{
		Username:     "admin",
		PasswordHash: hash,
		Nickname:     "系统管理员",
		Email:        "admin@orienteering.local",
		Role:         constants.RoleAdmin,
		Points:       0,
	}
	if err := db.Create(admin).Error; err != nil {
		return fmt.Errorf("create admin: %w", err)
	}
	logger.Info("seed admin created", "username", admin.Username)
	return nil
}
