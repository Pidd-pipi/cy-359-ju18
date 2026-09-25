package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/orienteering/platform/internal/handler"
)

// registerProductRoutes 商城商品与兑换路由。
func registerProductRoutes(g *gin.RouterGroup, h *handler.ProductHandler, redemption *handler.RedemptionHandler, authMW gin.HandlerFunc, adminMW gin.HandlerFunc, logger *slog.Logger) {
	products := g.Group("/products")
	products.GET("", h.List)
	products.GET("/:id", h.Get)

	admin := g.Group("", authMW, adminMW)
	admin.POST("/products", h.Create)
	admin.PUT("/products/:id", h.Update)
	admin.PATCH("/products/:id/status", h.UpdateStatus)

	auth := g.Group("", authMW)
	auth.POST("/redemptions", redemption.Redeem)
	auth.GET("/redemptions/mine", redemption.ListMine)
	adminRed := g.Group("", authMW, adminMW)
	adminRed.GET("/redemptions", redemption.ListAll)
	adminRed.PUT("/redemptions/:id/status", redemption.UpdateStatus)
	_ = logger
}
