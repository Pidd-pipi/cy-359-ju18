package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/orienteering/platform/internal/constants"
	"github.com/orienteering/platform/internal/dto"
	"github.com/orienteering/platform/internal/middleware"
	"github.com/orienteering/platform/internal/service"
	"github.com/orienteering/platform/internal/util"
	"github.com/orienteering/platform/pkg/ws"
)

// CheckinHandler 打卡与排行榜 HTTP 处理器。
type CheckinHandler struct {
	svc    *service.CheckinService
	lb     *service.LeaderboardService
	hub    *ws.Hub
	secret string
	logger *slog.Logger
}

// NewCheckinHandler 构造打卡处理器。
func NewCheckinHandler(svc *service.CheckinService, lb *service.LeaderboardService, hub *ws.Hub, secret string, logger *slog.Logger) *CheckinHandler {
	return &CheckinHandler{svc: svc, lb: lb, hub: hub, secret: secret, logger: logger}
}

// Checkin 打卡。
func (h *CheckinHandler) Checkin(c *gin.Context) {
	teamID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "团队 id 参数错误")
		return
	}
	var req dto.CheckinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "参数错误："+err.Error())
		return
	}
	record, err := h.svc.Checkin(teamID, middleware.UserID(c), &req, h.logger)
	if err != nil {
		respondError(c, err)
		return
	}
	util.OKMessage(c, constants.MsgCheckinOK, record)
}

// ListByTeam 团队打卡记录。
func (h *CheckinHandler) ListByTeam(c *gin.Context) {
	teamID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "团队 id 参数错误")
		return
	}
	activityID, _ := strconv.ParseInt(c.Query("activity_id"), 10, 64)
	records, err := h.svc.ListByTeam(teamID, activityID)
	if err != nil {
		respondError(c, err)
		return
	}
	util.OK(c, records)
}

// ListByActivity 活动全部打卡记录（管理员）。
func (h *CheckinHandler) ListByActivity(c *gin.Context) {
	activityID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "活动 id 参数错误")
		return
	}
	records, err := h.svc.ListByActivity(activityID)
	if err != nil {
		respondError(c, err)
		return
	}
	util.OK(c, records)
}

// Leaderboard 实时排行榜。
func (h *CheckinHandler) Leaderboard(c *gin.Context) {
	activityID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "活动 id 参数错误")
		return
	}
	rows, err := h.lb.Leaderboard(activityID)
	if err != nil {
		respondError(c, err)
		return
	}
	util.OK(c, rows)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// LeaderboardWS 排行榜 WebSocket 实时推送。
func (h *CheckinHandler) LeaderboardWS(c *gin.Context) {
	activityID, err := strconv.ParseInt(c.Query("activity_id"), 10, 64)
	if err != nil || activityID <= 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "activity_id 参数错误")
		return
	}
	// 通过 query token 校验登录
	tokenStr := c.Query("token")
	if tokenStr == "" {
		tokenStr = c.GetHeader("Sec-WebSocket-Protocol")
	}
	claims, err := util.ParseToken(h.secret, tokenStr)
	if err != nil {
		util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, "登录凭证无效")
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Warn("ws upgrade failed", "err", err)
		return
	}
	ch := h.hub.Subscribe(activityID)
	defer func() {
		h.hub.Unsubscribe(activityID, ch)
		_ = conn.Close()
	}()
	h.logger.Info(fmt.Sprintf("ws connected activity_id=%d user_id=%d", activityID, claims.UserID))
	// 先推送当前榜单
	if rows, err := h.lb.Leaderboard(activityID); err == nil {
		_ = conn.WriteJSON(ws.Message{Type: "leaderboard", ActivityID: activityID, Data: rows})
	}
	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()
	for payload := range ch {
		if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
			return
		}
	}
}
