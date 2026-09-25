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

// ProductHandler 商城商品 HTTP 处理器。
type ProductHandler struct {
	svc    *service.ProductService
	logger *slog.Logger
}

// NewProductHandler 构造商品处理器。
func NewProductHandler(svc *service.ProductService, logger *slog.Logger) *ProductHandler {
	return &ProductHandler{svc: svc, logger: logger}
}

// Create 创建商品（管理员）。
func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "参数错误："+err.Error())
		return
	}
	product, err := h.svc.Create(&req, h.logger)
	if err != nil {
		respondError(c, err)
		return
	}
	util.OKMessage(c, constants.MsgCreated, product)
}

// Update 编辑商品（管理员）。
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "商品 id 参数错误")
		return
	}
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "参数错误："+err.Error())
		return
	}
	product, err := h.svc.Update(id, &req, h.logger)
	if err != nil {
		respondError(c, err)
		return
	}
	util.OKMessage(c, constants.MsgUpdated, product)
}

// UpdateStatus 上下架（管理员）。
func (h *ProductHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "商品 id 参数错误")
		return
	}
	status := c.Query("status")
	if status != constants.ProductStatusOn && status != constants.ProductStatusOff {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "status 参数不合法")
		return
	}
	product, err := h.svc.UpdateStatus(id, status, h.logger)
	if err != nil {
		respondError(c, err)
		return
	}
	util.OKMessage(c, constants.MsgUpdated, product)
}

// List 分页查询商品。
func (h *ProductHandler) List(c *gin.Context) {
	pq := util.ParsePageQuery(c)
	status := c.Query("status")
	if status != "" && status != constants.ProductStatusOn && status != constants.ProductStatusOff {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "status 参数不合法")
		return
	}
	products, total, err := h.svc.List(pq.Page, pq.PageSize, pq.Offset, status, c.Query("keyword"))
	if err != nil {
		respondError(c, err)
		return
	}
	views := make([]dto.ProductView, 0, len(products))
	for i := range products {
		views = append(views, service.ToProductView(&products[i]))
	}
	util.OK(c, util.PageResult{List: views, Total: total, Page: pq.Page, PageSize: pq.PageSize})
}

// Get 商品详情。
func (h *ProductHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "商品 id 参数错误")
		return
	}
	product, err := h.svc.Get(id)
	if err != nil {
		respondError(c, err)
		return
	}
	util.OK(c, service.ToProductView(product))
}
