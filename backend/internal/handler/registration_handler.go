package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/orienteering/platform/internal/constants"
	"github.com/orienteering/platform/internal/dto"
	"github.com/orienteering/platform/internal/middleware"
	"github.com/orienteering/platform/internal/service"
	"github.com/orienteering/platform/internal/util"
)

// RegistrationHandler 报名 HTTP 处理器。
type RegistrationHandler struct {
	svc    *service.RegistrationService
	logger *slog.Logger
}

// NewRegistrationHandler 构造报名处理器。
func NewRegistrationHandler(svc *service.RegistrationService, logger *slog.Logger) *RegistrationHandler {
	return &RegistrationHandler{svc: svc, logger: logger}
}

// Apply 团队报名活动。
func (h *RegistrationHandler) Apply(c *gin.Context) {
	var req dto.ApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "参数错误："+err.Error())
		return
	}
	reg, err := h.svc.Apply(req.TeamID, req.ActivityID, middleware.UserID(c), h.logger)
	if err != nil {
		respondError(c, err)
		return
	}
	view, err := h.svc.ToView(reg)
	if err != nil {
		respondError(c, err)
		return
	}
	msg := constants.MsgRegistered
	if reg.Status == constants.RegistrationStatusWaitlist {
		msg = fmt.Sprintf("活动名额已满，已进入候补，前面还有 %d 队", view.WaitlistAhead)
	}
	util.OKMessage(c, msg, view)
}

// Approve 管理员通过报名。
func (h *RegistrationHandler) Approve(c *gin.Context) {
	regID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "报名 id 参数错误")
		return
	}
	reg, err := h.svc.Approve(regID, middleware.UserID(c), h.logger)
	if err != nil {
		respondError(c, err)
		return
	}
	util.OKMessage(c, "报名已通过", reg)
}

// Reject 管理员拒绝报名。
func (h *RegistrationHandler) Reject(c *gin.Context) {
	regID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "报名 id 参数错误")
		return
	}
	reg, err := h.svc.Reject(regID, h.logger)
	if err != nil {
		respondError(c, err)
		return
	}
	util.OKMessage(c, "报名已拒绝", reg)
}

// ListByActivity 活动的报名列表。
func (h *RegistrationHandler) ListByActivity(c *gin.Context) {
	activityID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "活动 id 参数错误")
		return
	}
	regs, err := h.svc.ListByActivity(activityID)
	if err != nil {
		respondError(c, err)
		return
	}
	views, err := h.svc.ToViews(regs)
	if err != nil {
		respondError(c, err)
		return
	}
	util.OK(c, views)
}

// ListMine 我的报名记录。
func (h *RegistrationHandler) ListMine(c *gin.Context) {
	regs, err := h.svc.ListMine(middleware.UserID(c))
	if err != nil {
		respondError(c, err)
		return
	}
	views, err := h.svc.ToViews(regs)
	if err != nil {
		respondError(c, err)
		return
	}
	util.OK(c, views)
}
