package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/orienteering/platform/internal/constants"
	"github.com/orienteering/platform/internal/dto"
	"github.com/orienteering/platform/internal/service"
	"github.com/orienteering/platform/internal/util"
)

// CheckpointHandler 打卡点 HTTP 处理器。
type CheckpointHandler struct {
	svc    *service.CheckpointService
	logger *slog.Logger
}

// NewCheckpointHandler 构造打卡点处理器。
func NewCheckpointHandler(svc *service.CheckpointService, logger *slog.Logger) *CheckpointHandler {
	return &CheckpointHandler{svc: svc, logger: logger}
}

// Create 创建打卡点。
func (h *CheckpointHandler) Create(c *gin.Context) {
	activityID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "活动 id 参数错误")
		return
	}
	var req dto.CreateCheckpointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "参数错误："+err.Error())
		return
	}
	cp, err := h.svc.Create(activityID, &req, h.logger)
	if err != nil {
		respondError(c, err)
		return
	}
	util.OKMessage(c, constants.MsgCreated, cp)
}

// Update 编辑打卡点。
func (h *CheckpointHandler) Update(c *gin.Context) {
	cpID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "打卡点 id 参数错误")
		return
	}
	var req dto.CreateCheckpointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "参数错误："+err.Error())
		return
	}
	cp, err := h.svc.Update(cpID, &req, h.logger)
	if err != nil {
		respondError(c, err)
		return
	}
	util.OKMessage(c, constants.MsgUpdated, cp)
}

// Delete 删除打卡点。
func (h *CheckpointHandler) Delete(c *gin.Context) {
	cpID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "打卡点 id 参数错误")
		return
	}
	if err := h.svc.Delete(cpID, h.logger); err != nil {
		respondError(c, err)
		return
	}
	util.OKMessage(c, constants.MsgDeleted, gin.H{"id": cpID})
}
