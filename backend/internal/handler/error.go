package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/orienteering/platform/internal/constants"
	"github.com/orienteering/platform/internal/repository"
	"github.com/orienteering/platform/internal/util"
)

// respondError 将 service 返回的错误转换为统一响应（handler 再次包装）。
func respondError(c *gin.Context, err error) {
	var appErr *util.AppError
	if errors.As(err, &appErr) {
		// handler 再次包装 service 错误，保留实体/字段/角色上下文
		wrapped := util.Wrap(err, appErr.Code, "handler:"+c.FullPath())
		var w *util.AppError
		_ = errors.As(wrapped, &w)
		util.Fail(c, httpStatusForCode(w.Code), w.Code, w.Message)
		return
	}
	if errors.Is(err, repository.ErrNotFound) {
		util.Fail(c, http.StatusNotFound, constants.CodeNotFound, "资源不存在")
		return
	}
	util.Fail(c, http.StatusInternalServerError, constants.CodeInternalError, "系统内部错误")
}

func httpStatusForCode(code int) int {
	switch {
	case code >= 40000 && code < 40100:
		return http.StatusBadRequest
	case code >= 40100 && code < 40200:
		return http.StatusUnauthorized
	case code >= 40300 && code < 40400:
		return http.StatusForbidden
	case code >= 40400 && code < 40500:
		return http.StatusNotFound
	case code >= 40900 && code < 41000:
		return http.StatusConflict
	case code >= 42200 && code < 42300:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
