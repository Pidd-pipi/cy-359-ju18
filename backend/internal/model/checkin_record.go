package model

import "time"

// CheckinRecord 打卡记录实体。
type CheckinRecord struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
	ActivityID   int64     `gorm:"index;not null" json:"activity_id"`
	CheckpointID int64     `gorm:"index;not null" json:"checkpoint_id"`
	TeamID       int64     `gorm:"index;not null" json:"team_id"`
	UserID       int64     `gorm:"index;not null" json:"user_id"`
	CheckinType  string    `gorm:"size:20;not null" json:"checkin_type"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	Answer       string    `gorm:"size:255" json:"answer"`
	PhotoURL     string    `gorm:"size:255" json:"photo_url"`
	Result       string    `gorm:"size:20;default:none;not null" json:"result"`
	PointsEarned int       `gorm:"default:0;not null" json:"points_earned"`
	CheckedInAt  time.Time `json:"checked_in_at"`
}
