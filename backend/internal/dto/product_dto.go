package dto

import "time"

// CreateProductRequest 创建/编辑商品入参。
type CreateProductRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=120"`
	Description string `json:"description"`
	Type        string `json:"type" binding:"required,oneof=equipment coupon badge"`
	PointsCost  int    `json:"points_cost" binding:"required,min=1"`
	Stock       int    `json:"stock" binding:"required,min=0"`
	Status      string `json:"status" binding:"required,oneof=on off"`
	ImageURL    string `json:"image_url"`
}

// ProductView 商品展示视图。
type ProductView struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	PointsCost  int       `json:"points_cost"`
	Stock       int       `json:"stock"`
	Status      string    `json:"status"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
}

// RedeemRequest 兑换入参。
type RedeemRequest struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int   `json:"quantity" binding:"required,min=1,max=99"`
}

// RedemptionView 兑换订单视图。
type RedemptionView struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	Username    string     `json:"username"`
	ProductID   int64      `json:"product_id"`
	ProductName string     `json:"product_name"`
	Quantity    int        `json:"quantity"`
	PointsCost  int        `json:"points_cost"`
	TotalPoints int        `json:"total_points"`
	Status      string     `json:"status"`
	RedeemedAt  time.Time  `json:"redeemed_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

// RedemptionStatusRequest 兑换订单状态流转入参。
type RedemptionStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending completed cancelled"`
}
