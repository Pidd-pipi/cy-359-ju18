package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Body 统一响应体。
type Body struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// OK 返回成功响应。
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Body{Code: 0, Message: "ok", Data: data})
}

// OKMessage 返回带自定义 message 的成功响应。
func OKMessage(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, Body{Code: 0, Message: message, Data: data})
}

// Fail 返回带 HTTP 状态码与业务码的错误响应。
func Fail(c *gin.Context, httpStatus, code int, message string) {
	c.JSON(httpStatus, Body{Code: code, Message: message})
}
