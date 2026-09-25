package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/orienteering/platform/internal/constants"
	"github.com/orienteering/platform/internal/dto"
	"github.com/orienteering/platform/internal/middleware"
	"github.com/orienteering/platform/internal/model"
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

// Apply 团队报名活动。名额充足进入待审核；名额不足按提交顺序进入候补。
func (h *RegistrationHandler) Apply(c *gin.Context) {
	var req dto.ApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "参数错误："+err.Error())
		return
	}
	result, err := h.svc.Apply(req.TeamID, req.ActivityID, middleware.UserID(c), h.logger)
	if err != nil {
		respondError(c, err)
		return
	}
	view, err := h.svc.ToView(result.Registration)
	if err != nil {
		respondError(c, err)
		return
	}
	if result.Waitlisted {
		// 候补是正常业务结果，仍以 code=0 返回，前端据此提示排队位置。
		util.OKMessage(c, "名额已满，已进入候补队列，前面还有"+
			strconv.FormatInt(result.WaitlistAhead, 10)+"队", view)
		return
	}
	util.OKMessage(c, "报名成功，等待管理员审核", view)
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

// Reject 管理员拒绝报名；待审核被拒绝后最早候补自动递补为待审核。
func (h *RegistrationHandler) Reject(c *gin.Context) {
	regID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "报名 id 参数错误")
		return
	}
	result, err := h.svc.Reject(regID, h.logger)
	if err != nil {
		respondError(c, err)
		return
	}
	msg := "报名已拒绝"
	if result.Promoted != nil {
		msg = "报名已拒绝，候补第 1 队（报名记录 " +
			strconv.FormatInt(result.Promoted.ID, 10) + "）已自动递补为待审核"
	}
	util.OKMessage(c, msg, gin.H{
		"registration_id": result.Registration.ID,
		"promoted_id":     promotedViewID(result.Promoted),
	})
}

// ListByActivity 活动的报名列表（管理员）。
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

// promotedViewID 安全取递补记录 id。
func promotedViewID(r *model.Registration) int64 {
	if r == nil {
		return 0
	}
	return r.ID
}
