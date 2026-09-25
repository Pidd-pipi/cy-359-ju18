package model

import "time"

// AuditLog 操作审计日志实体。
type AuditLog struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
	UserID       int64     `gorm:"index;not null" json:"user_id"`
	Username     string    `gorm:"size:50" json:"username"`
	Action       string    `gorm:"size:50;not null" json:"action"`
	ResourceType string    `gorm:"size:50;not null" json:"resource_type"`
	ResourceID   string    `gorm:"size:64" json:"resource_id"`
	Detail       string    `gorm:"type:text" json:"detail"`
	IP           string    `gorm:"size:50" json:"ip"`
	RequestID    string    `gorm:"size:64;index" json:"request_id"`
	CreatedAt    time.Time `json:"created_at"`
}
