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

// UserHandler 用户 HTTP 处理器。
type UserHandler struct {
	svc    *service.UserService
	logger *slog.Logger
}

// NewUserHandler 构造用户处理器。
func NewUserHandler(svc *service.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{svc: svc, logger: logger}
}

// Register 注册。
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "参数错误："+err.Error())
		return
	}
	user, err := h.svc.Register(&req)
	if err != nil {
		respondError(c, err)
		return
	}
	h.logger.Info(fmt.Sprintf(constants.LogUserRegister, user.Username, user.Role))
	util.OKMessage(c, constants.MsgRegistered, user)
}

// Login 登录。
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "参数错误："+err.Error())
		return
	}
	resp, err := h.svc.Login(&req)
	if err != nil {
		respondError(c, err)
		return
	}
	h.logger.Info(fmt.Sprintf(constants.LogUserLogin, resp.User.Username, resp.User.Role))
	util.OKMessage(c, constants.MsgLoggedIn, resp)
}

// Profile 查询我的信息。
func (h *UserHandler) Profile(c *gin.Context) {
	userID := middleware.UserID(c)
	user, err := h.svc.GetProfile(userID)
	if err != nil {
		respondError(c, err)
		return
	}
	h.logger.Info(fmt.Sprintf(constants.LogUserProfile, userID))
	util.OK(c, user)
}

// UpdateProfile 更新我的资料。
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "参数错误："+err.Error())
		return
	}
	userID := middleware.UserID(c)
	user, err := h.svc.UpdateProfile(userID, &req)
	if err != nil {
		respondError(c, err)
		return
	}
	h.logger.Info(fmt.Sprintf(constants.LogUserUpdate, userID, "profile"))
	util.OKMessage(c, constants.MsgUpdated, user)
}

// List 管理员分页查询用户。
func (h *UserHandler) List(c *gin.Context) {
	pq := util.ParsePageQuery(c)
	users, total, err := h.svc.List(pq.Page, pq.PageSize, pq.Offset, c.Query("keyword"))
	if err != nil {
		respondError(c, err)
		return
	}
	util.OK(c, util.PageResult{List: users, Total: total, Page: pq.Page, PageSize: pq.PageSize})
}

// SetDisabled 管理员禁用/启用用户。
func (h *UserHandler) SetDisabled(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "用户 id 参数错误")
		return
	}
	disabled, _ := strconv.ParseBool(c.DefaultQuery("disabled", "false"))
	if err := h.svc.SetDisabled(id, disabled); err != nil {
		respondError(c, err)
		return
	}
	util.OKMessage(c, constants.MsgUpdated, gin.H{"id": id, "disabled": disabled})
}
