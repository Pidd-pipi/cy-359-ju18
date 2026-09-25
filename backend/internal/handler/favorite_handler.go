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

// FavoriteHandler 收藏 HTTP 处理器。
type FavoriteHandler struct {
	svc      *service.FavoriteService
	activity *service.ActivityService
	logger   *slog.Logger
}

// NewFavoriteHandler 构造收藏处理器。
func NewFavoriteHandler(svc *service.FavoriteService, activity *service.ActivityService, logger *slog.Logger) *FavoriteHandler {
	return &FavoriteHandler{svc: svc, activity: activity, logger: logger}
}

// Add 收藏线路。
func (h *FavoriteHandler) Add(c *gin.Context) {
	var req dto.FavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "参数错误："+err.Error())
		return
	}
	activity, err := h.activity.Get(req.ActivityID)
	if err != nil {
		respondError(c, err)
		return
	}
	fav, err := h.svc.Add(middleware.UserID(c), req.ActivityID, activity.Status, h.logger)
	if err != nil {
		respondError(c, err)
		return
	}
	util.OKMessage(c, constants.MsgFavoriteOK, fav)
}

// Remove 取消收藏。
func (h *FavoriteHandler) Remove(c *gin.Context) {
	activityID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "活动 id 参数错误")
		return
	}
	if err := h.svc.Remove(middleware.UserID(c), activityID, h.logger); err != nil {
		respondError(c, err)
		return
	}
	util.OKMessage(c, constants.MsgUnfavoriteOK, nil)
}

// List 我的收藏。
func (h *FavoriteHandler) List(c *gin.Context) {
	pq := util.ParsePageQuery(c)
	favorites, total, err := h.svc.List(middleware.UserID(c), pq.Page, pq.PageSize, pq.Offset)
	if err != nil {
		respondError(c, err)
		return
	}
	views := make([]dto.FavoriteView, 0, len(favorites))
	for i := range favorites {
		f := &favorites[i]
		cpCount, _ := h.activity.CountCheckpoints(f.ActivityID)
		teamCount, _ := h.activity.CountRegistrations(f.ActivityID)
		av := service.ToActivityView(&f.Activity, cpCount, teamCount)
		views = append(views, service.ToFavoriteView(f, av))
	}
	util.OK(c, util.PageResult{List: views, Total: total, Page: pq.Page, PageSize: pq.PageSize})
}
