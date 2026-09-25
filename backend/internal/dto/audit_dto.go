package dto

import "time"

// AuditLogView 审计日志视图。
type AuditLogView struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Username     string    `json:"username"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resource_type"`
	ResourceID   string    `json:"resource_id"`
	Detail       string    `json:"detail"`
	IP           string    `json:"ip"`
	RequestID    string    `json:"request_id"`
	CreatedAt    time.Time `json:"created_at"`
}
