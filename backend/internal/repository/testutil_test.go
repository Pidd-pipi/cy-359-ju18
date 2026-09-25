package repository

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/orienteering/platform/internal/model"
)

// newTestDB 创建内存 sqlite 数据库并迁移测试表。
func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	if err := db.AutoMigrate(
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
	); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}
	return db
}
