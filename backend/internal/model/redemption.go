package model

import "time"

// Redemption 积分兑换订单实体。
type Redemption struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	UserID      int64     `gorm:"index;not null" json:"user_id"`
	ProductID   int64     `gorm:"index;not null" json:"product_id"`
	Quantity    int       `gorm:"default:1;not null" json:"quantity"`
	PointsCost  int       `gorm:"not null" json:"points_cost"`
	TotalPoints int       `gorm:"not null" json:"total_points"`
	Status      string    `gorm:"size:20;default:pending;not null;index" json:"status"`
	RedeemedAt  time.Time `json:"redeemed_at"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
