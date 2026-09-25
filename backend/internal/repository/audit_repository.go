package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/orienteering/platform/internal/model"
)

// AuditRepository 审计日志仓储。
type AuditRepository struct {
	db *gorm.DB
}

// NewAuditRepository 构造审计日志仓储。
func NewAuditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Create(log *model.AuditLog) error {
	if err := r.db.Create(log).Error; err != nil {
		return fmt.Errorf("create audit log: %w", err)
	}
	return nil
}

func (r *AuditRepository) List(page, pageSize, offset int, action, resourceType, keyword string) ([]model.AuditLog, int64, error) {
	query := r.db.Model(&model.AuditLog{})
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("username LIKE ? OR detail LIKE ?", like, like)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count audit logs: %w", err)
	}
	var logs []model.AuditLog
	if err := query.Order("id DESC").Limit(pageSize).Offset(offset).Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("list audit logs: %w", err)
	}
	return logs, total, nil
}
