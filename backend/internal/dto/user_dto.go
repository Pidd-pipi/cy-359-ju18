package dto

import "time"

// RegisterRequest 注册入参。
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=64"`
	Nickname string `json:"nickname" binding:"required,min=1,max=50"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty,len=11"`
}

// LoginRequest 登录入参。
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录返回。
type LoginResponse struct {
	Token string    `json:"token"`
	User  UserBrief `json:"user"`
}

// UserBrief 用户简要信息。
type UserBrief struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
	Points   int    `json:"points"`
}

// UpdateProfileRequest 更新个人资料入参。
type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"omitempty,min=1,max=50"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty,len=11"`
}

// UserView 用户展示视图。
type UserView struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"`
	Points    int       `json:"points"`
	Disabled  bool      `json:"disabled"`
	CreatedAt time.Time `json:"created_at"`
}
