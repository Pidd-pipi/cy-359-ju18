package util

import (
	"errors"
	"fmt"
)

// AppError 业务错误，包含响应码与人类可读 message。
type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code=%d message=%s cause=%v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code=%d message=%s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error { return e.Err }

// NewAppError 构造业务错误。
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

// Wrap 在业务错误上继续包装 message（屎山要求：handler 再次包装 service 错误）。
func Wrap(err error, code int, message string) error {
	var ae *AppError
	if errors.As(err, &ae) {
		return &AppError{Code: ae.Code, Message: message + "：" + ae.Message, Err: ae}
	}
	return &AppError{Code: code, Message: message, Err: err}
}
