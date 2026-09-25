package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/orienteering/platform/internal/constants"
)

// FormatDateTime 将时间格式化为 "2006-01-02 15:04:05"。
func FormatDateTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// FormatDuration 将秒数格式化为 "HH:MM:SS"。
func FormatDuration(totalSeconds int) string {
	if totalSeconds <= 0 {
		return "--:--:--"
	}
	h := totalSeconds / 3600
	m := (totalSeconds % 3600) / 60
	s := totalSeconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

// ActivityStatusText 活动状态文本（状态机在 formatters 的镜像）。
func ActivityStatusText(status string) string {
	switch status {
	case constants.ActivityStatusDraft:
		return "草稿"
	case constants.ActivityStatusPublished:
		return "已发布"
	case constants.ActivityStatusOngoing:
		return "进行中"
	case constants.ActivityStatusFinished:
		return "已结束"
	case constants.ActivityStatusCancelled:
		return "已取消"
	default:
		return "未知"
	}
}

// DifficultyText 活动难度文本。
func DifficultyText(difficulty string) string {
	switch difficulty {
	case constants.DifficultyFamily:
		return "亲子"
	case constants.DifficultyAdult:
		return "成人"
	case constants.DifficultyPro:
		return "专业"
	default:
		return "未知"
	}
}

// TaskTypeText 打卡任务类型文本。
func TaskTypeText(taskType string) string {
	switch taskType {
	case constants.TaskTypeNone:
		return "无任务"
	case constants.TaskTypeQuiz:
		return "答题"
	case constants.TaskTypePhoto:
		return "拍照"
	default:
		return "未知"
	}
}

// RegistrationStatusText 报名状态文本。
func RegistrationStatusText(status string) string {
	switch status {
	case constants.RegistrationStatusPending:
		return "待审核"
	case constants.RegistrationStatusApproved:
		return "已通过"
	case constants.RegistrationStatusRejected:
		return "已拒绝"
	case constants.RegistrationStatusFinished:
		return "已完成"
	default:
		return "未知"
	}
}

// RedemptionStatusText 兑换订单状态文本。
func RedemptionStatusText(status string) string {
	switch status {
	case constants.RedemptionStatusPending:
		return "待处理"
	case constants.RedemptionStatusCompleted:
		return "已完成"
	case constants.RedemptionStatusCancelled:
		return "已取消"
	default:
		return "未知"
	}
}

// ProductTypeText 商品类型文本。
func ProductTypeText(productType string) string {
	switch productType {
	case constants.ProductTypeEquipment:
		return "户外装备"
	case constants.ProductTypeCoupon:
		return "活动优惠券"
	case constants.ProductTypeBadge:
		return "虚拟勋章"
	default:
		return "未知"
	}
}

// RoleText 角色文本。
func RoleText(role string) string {
	switch role {
	case constants.RoleAdmin:
		return "管理员"
	case constants.RoleUser:
		return "普通用户"
	default:
		return "未知"
	}
}

// StatusColor 返回状态对应的前端语义色。
func StatusColor(status string) string {
	// 与前端 constants/statusConfig 保持一致的色值映射
	switch status {
	case constants.ActivityStatusDraft, constants.RegistrationStatusPending:
		return "default"
	case constants.ActivityStatusPublished, constants.RegistrationStatusApproved, constants.ProductStatusOn:
		return "processing"
	case constants.ActivityStatusOngoing, constants.CheckinResultCorrect:
		return "success"
	case constants.ActivityStatusFinished, constants.RedemptionStatusCompleted, constants.ProductTypeBadge:
		return "cyan"
	case constants.ActivityStatusCancelled, constants.RegistrationStatusRejected, constants.CheckinResultWrong:
		return "error"
	default:
		return "default"
	}
}

// MaskPhone 脱敏手机号。
func MaskPhone(phone string) string {
	if len(phone) != 11 {
		return phone
	}
	return phone[:3] + "****" + phone[7:]
}

// SanitizeKeyword 清理关键字首尾空格。
func SanitizeKeyword(keyword string) string {
	return strings.TrimSpace(keyword)
}
