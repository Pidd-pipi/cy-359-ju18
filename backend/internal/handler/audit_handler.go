package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/orienteering/platform/internal/dto"
	"github.com/orienteering/platform/internal/service"
	"github.com/orienteering/platform/internal/util"
)

// AuditHandler 审计日志 HTTP 处理器。
type AuditHandler struct {
	svc    *service.AuditService
	logger *slog.Logger
}

// NewAuditHandler 构造审计处理器。
func NewAuditHandler(svc *service.AuditService, logger *slog.Logger) *AuditHandler {
	return &AuditHandler{svc: svc, logger: logger}
}

// List 管理员分页查询审计日志。
func (h *AuditHandler) List(c *gin.Context) {
	pq := util.ParsePageQuery(c)
	logs, total, err := h.svc.List(pq.Page, pq.PageSize, pq.Offset,
		c.Query("action"), c.Query("resource_type"), c.Query("keyword"))
	if err != nil {
		respondError(c, err)
		return
	}
	views := make([]dto.AuditLogView, 0, len(logs))
	for i := range logs {
		views = append(views, service.ToAuditLogView(&logs[i]))
	}
	util.OK(c, util.PageResult{List: views, Total: total, Page: pq.Page, PageSize: pq.PageSize})
}
