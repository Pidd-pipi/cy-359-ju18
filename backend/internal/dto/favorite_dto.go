package dto

import "time"

// FavoriteRequest 收藏/取消收藏入参。
type FavoriteRequest struct {
	ActivityID int64 `json:"activity_id" binding:"required"`
}

// FavoriteView 收藏视图。
type FavoriteView struct {
	ID           int64        `json:"id"`
	ActivityID   int64        `json:"activity_id"`
	Activity     ActivityView `json:"activity"`
	CreatedAt    time.Time    `json:"created_at"`
}
