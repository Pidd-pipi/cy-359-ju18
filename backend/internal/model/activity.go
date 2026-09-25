package model

import "time"

// Activity 城市定向越野活动（线路）实体。
type Activity struct {
	ID                 int64     `gorm:"primaryKey" json:"id"`
	Title              string    `gorm:"size:120;not null" json:"title"`
	Description        string    `gorm:"type:text" json:"description"`
	Difficulty         string    `gorm:"size:20;default:adult;not null" json:"difficulty"`
	DurationMinutes    int       `gorm:"default:120;not null" json:"duration_minutes"`
	EquipmentRequirement string  `gorm:"type:text" json:"equipment_requirement"`
	StartTime          time.Time `json:"start_time"`
	EndTime            time.Time `json:"end_time"`
	Status             string    `gorm:"size:20;default:draft;not null;index" json:"status"`
	CreatorID          int64     `gorm:"index;not null" json:"creator_id"`
	StartLat           float64   `json:"start_lat"`
	StartLng           float64   `json:"start_lng"`
	EndLat             float64   `json:"end_lat"`
	EndLng             float64   `json:"end_lng"`
	Address            string    `gorm:"size:255" json:"address"`
	MaxTeams           int       `gorm:"default:50;not null" json:"max_teams"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`

	Checkpoints []Checkpoint `gorm:"foreignKey:ActivityID" json:"checkpoints,omitempty"`
}
