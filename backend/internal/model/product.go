package model

import "time"

// Product 积分商城商品实体。
type Product struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:120;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Type        string    `gorm:"size:20;default:equipment;not null" json:"type"`
	PointsCost  int       `gorm:"not null" json:"points_cost"`
	Stock       int       `gorm:"not null" json:"stock"`
	Status      string    `gorm:"size:20;default:on;not null" json:"status"`
	ImageURL    string    `gorm:"size:255" json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
