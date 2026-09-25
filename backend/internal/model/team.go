package model

import "time"

// Team 参赛团队实体。
type Team struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:80;uniqueIndex;not null" json:"name"`
	Slogan    string    `gorm:"size:255" json:"slogan"`
	CaptainID int64     `gorm:"index;not null" json:"captain_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Members []TeamMember `gorm:"foreignKey:TeamID" json:"members,omitempty"`
}

// TeamMember 团队成员关联实体。
type TeamMember struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	TeamID    int64     `gorm:"index;not null;uniqueIndex:idx_team_user" json:"team_id"`
	UserID    int64     `gorm:"index;not null;uniqueIndex:idx_team_user" json:"user_id"`
	Role      string    `gorm:"size:20;default:member;not null" json:"role"`
	JoinedAt  time.Time `json:"joined_at"`
}

// Registration 团队报名活动实体。
type Registration struct {
	ID           int64      `gorm:"primaryKey" json:"id"`
	TeamID       int64      `gorm:"index;not null;uniqueIndex:idx_team_activity" json:"team_id"`
	ActivityID   int64      `gorm:"index;not null;uniqueIndex:idx_team_activity" json:"activity_id"`
	Status       string     `gorm:"size:20;default:pending;not null;index" json:"status"`
	StartTime    *time.Time `json:"start_time"`
	FinishTime   *time.Time `json:"finish_time"`
	TotalSeconds int        `gorm:"default:0;not null" json:"total_seconds"`
	RegisteredAt time.Time  `json:"registered_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
