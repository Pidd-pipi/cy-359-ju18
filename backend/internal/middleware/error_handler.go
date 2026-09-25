package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/orienteering/platform/internal/constants"
	"github.com/orienteering/platform/internal/repository"
	"github.com/orienteering/platform/internal/util"
)

// ErrorHandler 全局错误处理与恢复中间件。
func ErrorHandler(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Error(fmt.Sprintf(constants.LogPanicRecover, c.GetString("request_id"), rec))
				util.Fail(c, http.StatusInternalServerError, constants.CodeInternalError, "系统内部错误")
				c.Abort()
			}
		}()
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			handleError(c, logger, err)
		}
	}
}

func handleError(c *gin.Context, logger *slog.Logger, err error) {
	logger.Error("request error",
		"request_id", c.GetString("request_id"),
		"path", c.FullPath(),
		"error", err.Error())
	var appErr *util.AppError
	if errors.As(err, &appErr) {
		util.Fail(c, httpStatusForCode(appErr.Code), appErr.Code, appErr.Message)
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
