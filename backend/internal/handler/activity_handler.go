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

// ActivityHandler 活动 HTTP 处理器。
type ActivityHandler struct {
	svc    *service.ActivityService
	regSvc *service.RegistrationService
	logger *slog.Logger
}

// NewActivityHandler 构造活动处理器。
func NewActivityHandler(svc *service.ActivityService, regSvc *service.RegistrationService, logger *slog.Logger) *ActivityHandler {
	return &ActivityHandler{svc: svc, regSvc: regSvc, logger: logger}
}

// Create 创建活动。
func (h *ActivityHandler) Create(c *gin.Context) {
	var req dto.CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "参数错误："+err.Error())
		return
	}
	creatorID := middleware.UserID(c)
	activity, err := h.svc.Create(creatorID, &req, h.logger)
	if err != nil {
		respondError(c, err)
		return
	}
	util.OKMessage(c, constants.MsgCreated, activity)
}

// Update 编辑活动。
func (h *ActivityHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "活动 id 参数错误")
		return
	}
	var req dto.CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "参数错误："+err.Error())
		return
	}
	activity, err := h.svc.Update(id, &req, h.logger)
	if err != nil {
		respondError(c, err)
		return
	}
	util.OKMessage(c, constants.MsgUpdated, activity)
}

// Transition 活动状态流转（发布/开始/结束/取消）。
func (h *ActivityHandler) Transition(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "活动 id 参数错误")
		return
	}
	var req dto.ActivityStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "参数错误："+err.Error())
		return
	}
	activity, err := h.svc.Transition(id, req.Status, middleware.UserID(c), h.logger)
	if err != nil {
		respondError(c, err)
		return
	}
	if req.Status == constants.ActivityStatusOngoing {
		// 复用 RegistrationService.Start：活动开始时为已通过团队记录开始时间
		if err := h.regSvc.Start(id, h.logger); err != nil {
			respondError(c, err)
			return
		}
	}
	util.OKMessage(c, "状态已更新为 "+util.ActivityStatusText(req.Status), activity)
}

// List 分页查询活动。
func (h *ActivityHandler) List(c *gin.Context) {
	pq := util.ParsePageQuery(c)
	status := c.Query("status")
	difficulty := c.Query("difficulty")
	if status != "" && !constants.Contains(constants.AllActivityStatuses, status) {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "活动 status 参数不合法")
		return
	}
	if difficulty != "" && !constants.Contains(constants.AllDifficulties, difficulty) {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "活动 difficulty 参数不合法")
		return
	}
	activities, total, err := h.svc.List(pq.Page, pq.PageSize, pq.Offset, status, difficulty, c.Query("keyword"))
	if err != nil {
		respondError(c, err)
		return
	}
	h.logger.Info(fmt.Sprintf(constants.LogActivityList, pq.Page, pq.PageSize, status))
	views := make([]dto.ActivityView, 0, len(activities))
	for i := range activities {
		cpCount, err1 := h.svc.CountCheckpoints(activities[i].ID)
		if err1 != nil {
			respondError(c, err1)
			return
		}
		teamCount, err2 := h.svc.CountRegistrations(activities[i].ID)
		if err2 != nil {
			respondError(c, err2)
			return
		}
		waitlistCount, err3 := h.svc.CountWaitlisted(activities[i].ID)
		if err3 != nil {
			respondError(c, err3)
			return
		}
		views = append(views, service.ToActivityView(&activities[i], cpCount, teamCount, waitlistCount))
	}
	util.OK(c, util.PageResult{List: views, Total: total, Page: pq.Page, PageSize: pq.PageSize})
}

// Get 活动详情（含打卡点、报名数）。
func (h *ActivityHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "活动 id 参数错误")
		return
	}
	activity, err := h.svc.Get(id)
	if err != nil {
		respondError(c, err)
		return
	}
	cpCount, err := h.svc.CountCheckpoints(id)
	if err != nil {
		respondError(c, err)
		return
	}
	teamCount, err := h.svc.CountRegistrations(id)
	if err != nil {
		respondError(c, err)
		return
	}
	waitlistCount, err := h.svc.CountWaitlisted(id)
	if err != nil {
		respondError(c, err)
		return
	}
	util.OK(c, service.ToActivityView(activity, cpCount, teamCount, waitlistCount))
}

// Delete 删除活动（管理员）。
func (h *ActivityHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "活动 id 参数错误")
		return
	}
	if err := h.svc.Delete(id); err != nil {
		respondError(c, err)
		return
	}
	util.OKMessage(c, constants.MsgDeleted, gin.H{"id": id})
}
