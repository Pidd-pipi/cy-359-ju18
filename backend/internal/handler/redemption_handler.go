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

// RedemptionHandler 积分兑换 HTTP 处理器。
type RedemptionHandler struct {
	svc    *service.RedemptionService
	logger *slog.Logger
}

// NewRedemptionHandler 构造兑换处理器。
func NewRedemptionHandler(svc *service.RedemptionService, logger *slog.Logger) *RedemptionHandler {
	return &RedemptionHandler{svc: svc, logger: logger}
}

// Redeem 兑换商品。
func (h *RedemptionHandler) Redeem(c *gin.Context) {
	var req dto.RedeemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "参数错误："+err.Error())
		return
	}
	redemption, err := h.svc.Redeem(middleware.UserID(c), &req, h.logger)
	if err != nil {
		respondError(c, err)
		return
	}
	util.OKMessage(c, constants.MsgRedeemOK, redemption)
}

// ListMine 我的兑换记录。
func (h *RedemptionHandler) ListMine(c *gin.Context) {
	pq := util.ParsePageQuery(c)
	redemptions, total, err := h.svc.ListMine(middleware.UserID(c), pq.Page, pq.PageSize, pq.Offset)
	if err != nil {
		respondError(c, err)
		return
	}
	views := make([]dto.RedemptionView, 0, len(redemptions))
	for i := range redemptions {
		views = append(views, toRedemptionView(&redemptions[i], "", ""))
	}
	util.OK(c, util.PageResult{List: views, Total: total, Page: pq.Page, PageSize: pq.PageSize})
}

// ListAll 管理员查询全部兑换记录。
func (h *RedemptionHandler) ListAll(c *gin.Context) {
	pq := util.ParsePageQuery(c)
	status := c.Query("status")
	if status != "" && !constants.Contains(constants.AllRedemptionStatuses, status) {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "status 参数不合法")
		return
	}
	redemptions, total, err := h.svc.ListAll(pq.Page, pq.PageSize, pq.Offset, status)
	if err != nil {
		respondError(c, err)
		return
	}
	views := make([]dto.RedemptionView, 0, len(redemptions))
	for i := range redemptions {
		views = append(views, toRedemptionView(&redemptions[i], "", ""))
	}
	util.OK(c, util.PageResult{List: views, Total: total, Page: pq.Page, PageSize: pq.PageSize})
}

// UpdateStatus 管理员处理兑换订单。
func (h *RedemptionHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "兑换 id 参数错误")
		return
	}
	var req dto.RedemptionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "参数错误："+err.Error())
		return
	}
	redemption, err := h.svc.UpdateStatus(id, req.Status, h.logger)
	if err != nil {
		respondError(c, err)
		return
	}
	util.OKMessage(c, constants.MsgUpdated, redemption)
}

func toRedemptionView(r *model.Redemption, username, productName string) dto.RedemptionView {
	return dto.RedemptionView{
		ID: r.ID, UserID: r.UserID, Username: username,
		ProductID: r.ProductID, ProductName: productName,
		Quantity: r.Quantity, PointsCost: r.PointsCost, TotalPoints: r.TotalPoints,
		Status: r.Status, RedeemedAt: r.RedeemedAt, CompletedAt: r.CompletedAt,
	}
}
