package model

import "time"

// Checkpoint 线索打卡点（CP 点）实体。
type Checkpoint struct {
	ID             int64     `gorm:"primaryKey" json:"id"`
	ActivityID     int64     `gorm:"index;not null" json:"activity_id"`
	Name           string    `gorm:"size:80;not null" json:"name"`
	Sequence       int       `gorm:"not null" json:"sequence"`
	Lat            float64   `json:"lat"`
	Lng            float64   `json:"lng"`
	Clue           string    `gorm:"size:255" json:"clue"`
	TaskType       string    `gorm:"size:20;default:none;not null" json:"task_type"`
	TaskContent    string    `gorm:"type:text" json:"task_content"`
	ExpectedAnswer string    `gorm:"size:255" json:"-"`
	RadiusMeters   int       `gorm:"default:200;not null" json:"radius_meters"`
	QRCode         string    `gorm:"size:64;uniqueIndex" json:"qr_code"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
