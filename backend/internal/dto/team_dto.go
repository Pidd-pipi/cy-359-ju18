package dto

import "time"

// CreateTeamRequest 创建团队入参。
type CreateTeamRequest struct {
	Name   string `json:"name" binding:"required,min=2,max=80"`
	Slogan string `json:"slogan" binding:"omitempty,max=255"`
}

// TeamView 团队展示视图。
type TeamView struct {
	ID        int64          `json:"id"`
	Name      string         `json:"name"`
	Slogan    string         `json:"slogan"`
	CaptainID int64          `json:"captain_id"`
	MemberCount int          `json:"member_count"`
	Members   []TeamMemberView `json:"members,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

// TeamMemberView 团队成员视图。
type TeamMemberView struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
}

// ApplyRequest 团队报名入参。
type ApplyRequest struct {
	TeamID     int64 `json:"team_id" binding:"required"`
	ActivityID int64 `json:"activity_id" binding:"required"`
}

// RegistrationView 报名记录视图。
type RegistrationView struct {
	ID           int64      `json:"id"`
	TeamID       int64      `json:"team_id"`
	TeamName     string     `json:"team_name"`
	ActivityID   int64      `json:"activity_id"`
	ActivityTitle string    `json:"activity_title"`
	Status       string     `json:"status"`
	// WaitlistAhead 候补队伍前面还有多少队（仅 status=waitlisted 时有意义，由 id 提交顺序计算）。
	WaitlistAhead int       `json:"waitlist_ahead"`
	StartTime    *time.Time `json:"start_time"`
	FinishTime   *time.Time `json:"finish_time"`
	TotalSeconds int        `json:"total_seconds"`
	Duration     string     `json:"duration"`
	RegisteredAt time.Time  `json:"registered_at"`
}

// JoinTeamRequest 加入团队入参。
type JoinTeamRequest struct {
	TeamID int64 `json:"team_id" binding:"required"`
}
