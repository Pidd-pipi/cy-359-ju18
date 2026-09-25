package model

import "time"

// Favorite 历史线路收藏实体。
type Favorite struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	UserID     int64     `gorm:"index;not null;uniqueIndex:idx_user_activity" json:"user_id"`
	ActivityID int64     `gorm:"index;not null;uniqueIndex:idx_user_activity" json:"activity_id"`
	CreatedAt  time.Time `json:"created_at"`

	Activity Activity `gorm:"foreignKey:ActivityID" json:"activity,omitempty"`
}
