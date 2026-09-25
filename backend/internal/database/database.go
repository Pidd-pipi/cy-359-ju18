package database

import (
	"fmt"
	"io"
	"log/slog"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/orienteering/platform/internal/model"
)

// New 创建并初始化数据库连接，执行自动迁移。
func New(dsn string, logger *slog.Logger) (*gorm.DB, error) {
	gormLogger := gormlogger.New(
		&slogWriter{logger: logger},
		gormlogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormlogger.Warn,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
		},
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLogger})
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migrate database: %w", err)
	}
	return db, nil
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Activity{},
		&model.Checkpoint{},
		&model.Team{},
		&model.TeamMember{},
		&model.Registration{},
		&model.CheckinRecord{},
		&model.Product{},
		&model.Redemption{},
		&model.Favorite{},
		&model.AuditLog{},
	)
}

// slogWriter 将 GORM 日志桥接到 slog。
type slogWriter struct {
	logger *slog.Logger
}

func (w *slogWriter) Write(p []byte) (int, error) {
	w.logger.Warn("gorm: " + string(p))
	return len(p), nil
}

// Printf 实现 gorm logger.Writer 接口。
func (w *slogWriter) Printf(format string, args ...any) {
	w.logger.Warn(fmt.Sprintf("gorm: "+format, args...))
}

var _ io.Writer = (*slogWriter)(nil)
