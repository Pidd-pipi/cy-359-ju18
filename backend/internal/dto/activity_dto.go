package dto

import "time"

// CreateActivityRequest 创建/编辑活动入参。
type CreateActivityRequest struct {
	Title                string    `json:"title" binding:"required,min=2,max=120"`
	Description          string    `json:"description"`
	Difficulty           string    `json:"difficulty" binding:"required,oneof=family adult pro"`
	DurationMinutes      int       `json:"duration_minutes" binding:"required,min=15,max=1440"`
	EquipmentRequirement string    `json:"equipment_requirement"`
	StartTime            time.Time `json:"start_time" binding:"required"`
	EndTime              time.Time `json:"end_time" binding:"required"`
	StartLat             float64   `json:"start_lat" binding:"required"`
	StartLng             float64   `json:"start_lng" binding:"required"`
	EndLat               float64   `json:"end_lat" binding:"required"`
	EndLng               float64   `json:"end_lng" binding:"required"`
	Address              string    `json:"address" binding:"required,max=255"`
	MaxTeams             int       `json:"max_teams" binding:"required,min=1,max=500"`
}

// ActivityView 活动展示视图。
type ActivityView struct {
	ID                   int64          `json:"id"`
	Title                string         `json:"title"`
	Description          string         `json:"description"`
	Difficulty           string         `json:"difficulty"`
	DurationMinutes      int            `json:"duration_minutes"`
	EquipmentRequirement string         `json:"equipment_requirement"`
	StartTime            time.Time      `json:"start_time"`
	EndTime              time.Time      `json:"end_time"`
	Status               string         `json:"status"`
	CreatorID            int64          `json:"creator_id"`
	StartLat             float64        `json:"start_lat"`
	StartLng             float64        `json:"start_lng"`
	EndLat               float64        `json:"end_lat"`
	EndLng               float64        `json:"end_lng"`
	Address              string         `json:"address"`
	MaxTeams             int            `json:"max_teams"`
	CheckpointCount      int            `json:"checkpoint_count"`
	TeamCount            int            `json:"team_count"`
	Checkpoints          []CheckpointView `json:"checkpoints,omitempty"`
	CreatedAt            time.Time      `json:"created_at"`
}

// ActivityStatusRequest 活动状态流转入参。
type ActivityStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=draft published ongoing finished cancelled"`
}

// LeaderboardRow 排行榜行。
type LeaderboardRow struct {
	Rank         int    `json:"rank"`
	TeamID       int64  `json:"team_id"`
	TeamName     string `json:"team_name"`
	TotalSeconds int    `json:"total_seconds"`
	Duration     string `json:"duration"`
	CheckpointCount int `json:"checkpoint_count"`
	Status       string `json:"status"`
}
