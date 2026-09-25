package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/orienteering/platform/internal/handler"
)

// registerFavoriteRoutes 收藏路由。
func registerFavoriteRoutes(g *gin.RouterGroup, h *handler.FavoriteHandler, authMW gin.HandlerFunc, logger *slog.Logger) {
	favorites := g.Group("/favorites")
	auth := favorites.Group("", authMW)
	auth.GET("", h.List)
	auth.POST("", h.Add)
	auth.DELETE("/:id", h.Remove)
	_ = logger
}
