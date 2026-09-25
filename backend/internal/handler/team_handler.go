package handler

import (
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

// TeamHandler 团队 HTTP 处理器。
type TeamHandler struct {
	svc    *service.TeamService
	logger *slog.Logger
}

// NewTeamHandler 构造团队处理器。
func NewTeamHandler(svc *service.TeamService, logger *slog.Logger) *TeamHandler {
	return &TeamHandler{svc: svc, logger: logger}
}

// Create 创建团队。
func (h *TeamHandler) Create(c *gin.Context) {
	var req dto.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "参数错误："+err.Error())
		return
	}
	team, err := h.svc.Create(middleware.UserID(c), &req, h.logger)
	if err != nil {
		respondError(c, err)
		return
	}
	util.OKMessage(c, constants.MsgCreated, team)
}

// Join 加入团队。
func (h *TeamHandler) Join(c *gin.Context) {
	var req dto.JoinTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "参数错误："+err.Error())
		return
	}
	if err := h.svc.Join(req.TeamID, middleware.UserID(c), h.logger); err != nil {
		respondError(c, err)
		return
	}
	util.OKMessage(c, "加入成功", nil)
}

// Leave 退出团队。
func (h *TeamHandler) Leave(c *gin.Context) {
	teamID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "团队 id 参数错误")
		return
	}
	if err := h.svc.Leave(teamID, middleware.UserID(c), h.logger); err != nil {
		respondError(c, err)
		return
	}
	util.OKMessage(c, "已退出团队", nil)
}

// ListMine 我的团队。
func (h *TeamHandler) ListMine(c *gin.Context) {
	teams, err := h.svc.ListMine(middleware.UserID(c))
	if err != nil {
		respondError(c, err)
		return
	}
	views := make([]dto.TeamView, 0, len(teams))
	for i := range teams {
		views = append(views, service.ToTeamView(&teams[i], nil))
	}
	util.OK(c, views)
}

// ListAll 管理员分页查询团队。
func (h *TeamHandler) ListAll(c *gin.Context) {
	pq := util.ParsePageQuery(c)
	teams, total, err := h.svc.ListAll(pq.Page, pq.PageSize, pq.Offset)
	if err != nil {
		respondError(c, err)
		return
	}
	views := make([]dto.TeamView, 0, len(teams))
	for i := range teams {
		views = append(views, service.ToTeamView(&teams[i], nil))
	}
	util.OK(c, util.PageResult{List: views, Total: total, Page: pq.Page, PageSize: pq.PageSize})
}

// Get 团队详情。
func (h *TeamHandler) Get(c *gin.Context) {
	teamID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "团队 id 参数错误")
		return
	}
	team, err := h.svc.Get(teamID)
	if err != nil {
		respondError(c, err)
		return
	}
	util.OK(c, service.ToTeamView(team, nil))
}
