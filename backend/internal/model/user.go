package model

import "time"

// User 用户实体。
type User struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"size:50;uniqueIndex;not null" json:"username"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Nickname     string    `gorm:"size:50" json:"nickname"`
	Email        string    `gorm:"size:100" json:"email"`
	Phone        string    `gorm:"size:20" json:"phone"`
	Role         string    `gorm:"size:20;default:user;not null" json:"role"`
	Points       int       `gorm:"default:0;not null" json:"points"`
	Disabled     bool      `gorm:"default:false;not null" json:"disabled"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
