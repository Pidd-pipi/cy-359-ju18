package constants

// 角色类型枚举（数据库 user.role 字段）。
const (
	RoleUser  = "user"  // 普通用户
	RoleAdmin = "admin" // 管理员
)

// 活动状态枚举（数据库 activities.status 字段）。
const (
	ActivityStatusDraft     = "draft"     // 草稿
	ActivityStatusPublished = "published" // 已发布（可报名）
	ActivityStatusOngoing   = "ongoing"   // 进行中（已开始计时）
	ActivityStatusFinished  = "finished"  // 已结束（可收藏/查看成绩）
	ActivityStatusCancelled = "cancelled" // 已取消
)

// 活动难度枚举（数据库 activities.difficulty 字段）。
const (
	DifficultyFamily = "family" // 亲子
	DifficultyAdult  = "adult"  // 成人
	DifficultyPro    = "pro"    // 专业
)

// 打卡点任务类型枚举（数据库 checkpoints.task_type 字段）。
const (
	TaskTypeNone  = "none"  // 无任务
	TaskTypeQuiz  = "quiz"  // 答题
	TaskTypePhoto = "photo" // 拍照
)

// 打卡方式枚举（数据库 checkin_records.checkin_type 字段）。
const (
	CheckinTypeGPS  = "gps"  // GPS 定位打卡
	CheckinTypeQR   = "qrcode" // 二维码打卡
)

// 报名状态枚举（数据库 registrations.status 字段）。
const (
	RegistrationStatusPending  = "pending"  // 待审核
	RegistrationStatusWaitlist = "waitlist" // 候补（名额已满，按提交顺序排队）
	RegistrationStatusApproved = "approved" // 已通过
	RegistrationStatusRejected = "rejected" // 已拒绝
	RegistrationStatusFinished = "finished" // 已完成
)

// 商城商品类型枚举（数据库 products.type 字段）。
const (
	ProductTypeEquipment = "equipment" // 户外装备
	ProductTypeCoupon    = "coupon"    // 活动优惠券
	ProductTypeBadge     = "badge"     // 虚拟勋章
)

// 商城商品上下架状态枚举（数据库 products.status 字段）。
const (
	ProductStatusOn  = "on"  // 上架
	ProductStatusOff = "off" // 下架
)

// 兑换订单状态枚举（数据库 redemptions.status 字段）。
const (
	RedemptionStatusPending   = "pending"   // 待处理
	RedemptionStatusCompleted = "completed" // 已完成
	RedemptionStatusCancelled = "cancelled" // 已取消
)

// 团队成员角色枚举（数据库 team_members.role 字段）。
const (
	TeamRoleCaptain = "captain" // 队长
	TeamRoleMember  = "member"  // 队员
)

// 答题打卡结果枚举（数据库 checkin_records.result 字段）。
const (
	CheckinResultCorrect = "correct" // 回答正确
	CheckinResultWrong   = "wrong"   // 回答错误
	CheckinResultNone    = "none"    // 无答题任务
)

// AllActivityStatuses 活动状态完整列表，用于状态机校验与前端筛选。
var AllActivityStatuses = []string{
	ActivityStatusDraft, ActivityStatusPublished, ActivityStatusOngoing,
	ActivityStatusFinished, ActivityStatusCancelled,
}

// AllDifficulties 活动难度完整列表。
var AllDifficulties = []string{DifficultyFamily, DifficultyAdult, DifficultyPro}

// AllRegistrationStatuses 报名状态完整列表。
var AllRegistrationStatuses = []string{
	RegistrationStatusPending, RegistrationStatusWaitlist, RegistrationStatusApproved,
	RegistrationStatusRejected, RegistrationStatusFinished,
}

// OccupiedRegistrationStatuses 占用活动名额的报名状态：待审核与已通过（含已完成）均占位，
// 候补不占位，被拒绝后释放名额。
var OccupiedRegistrationStatuses = []string{
	RegistrationStatusPending, RegistrationStatusApproved, RegistrationStatusFinished,
}

// AllRedemptionStatuses 兑换状态完整列表。
var AllRedemptionStatuses = []string{
	RedemptionStatusPending, RedemptionStatusCompleted, RedemptionStatusCancelled,
}

// Contains 判断切片是否包含目标字符串。
func Contains(list []string, target string) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}
