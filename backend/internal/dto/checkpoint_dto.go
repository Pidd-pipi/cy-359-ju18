package dto

// CreateCheckpointRequest 创建/编辑打卡点入参。
type CreateCheckpointRequest struct {
	Name           string  `json:"name" binding:"required,min=1,max=80"`
	Sequence       int     `json:"sequence" binding:"required,min=1"`
	Lat            float64 `json:"lat" binding:"required"`
	Lng            float64 `json:"lng" binding:"required"`
	Clue           string  `json:"clue"`
	TaskType       string  `json:"task_type" binding:"required,oneof=none quiz photo"`
	TaskContent    string  `json:"task_content"`
	ExpectedAnswer string  `json:"expected_answer"`
	RadiusMeters   int     `json:"radius_meters" binding:"required,min=50,max=2000"`
	QRCode         string  `json:"qr_code" binding:"required,min=6,max=64"`
}

// CheckpointView 打卡点展示视图。
type CheckpointView struct {
	ID          int64   `json:"id"`
	ActivityID  int64   `json:"activity_id"`
	Name        string  `json:"name"`
	Sequence    int     `json:"sequence"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	Clue        string  `json:"clue"`
	TaskType    string  `json:"task_type"`
	TaskContent string  `json:"task_content"`
	RadiusMeters int    `json:"radius_meters"`
	QRCode      string  `json:"qr_code"`
}

// CheckinRequest 打卡入参。
type CheckinRequest struct {
	CheckpointID int64   `json:"checkpoint_id" binding:"required"`
	CheckinType  string  `json:"checkin_type" binding:"required,oneof=gps qrcode"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Answer       string  `json:"answer"`
	PhotoURL     string  `json:"photo_url"`
}
