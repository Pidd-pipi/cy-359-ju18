package service

import (
	"fmt"
	"log/slog"

	"github.com/orienteering/platform/internal/constants"
	"github.com/orienteering/platform/internal/dto"
	"github.com/orienteering/platform/internal/model"
	"github.com/orienteering/platform/internal/repository"
)

// AuditService 操作审计日志业务服务。
type AuditService struct {
	repo *repository.AuditRepository
}

// NewAuditService 构造审计服务。
func NewAuditService(repo *repository.AuditRepository) *AuditService {
	return &AuditService{repo: repo}
}

// Record 记录一条审计日志。
func (s *AuditService) Record(log *model.AuditLog, logger *slog.Logger) error {
	if err := s.repo.Create(log); err != nil {
		logger.Error("audit record failed", "err", err)
		return err
	}
	logger.Info(fmt.Sprintf(constants.LogAuditWrite, log.UserID, log.Action, log.ResourceType))
	return nil
}

// List 分页查询审计日志（管理员）。
func (s *AuditService) List(page, pageSize, offset int, action, resourceType, keyword string) ([]model.AuditLog, int64, error) {
	return s.repo.List(page, pageSize, offset, action, resourceType, keyword)
}

// ToAuditLogView 审计日志转展示视图。
func ToAuditLogView(l *model.AuditLog) dto.AuditLogView {
	return dto.AuditLogView{
		ID: l.ID, UserID: l.UserID, Username: l.Username, Action: l.Action,
		ResourceType: l.ResourceType, ResourceID: l.ResourceID, Detail: l.Detail,
		IP: l.IP, RequestID: l.RequestID, CreatedAt: l.CreatedAt,
	}
}
