package constants

// 统一响应码：0 表示成功，非 0 表示业务/系统错误。
const (
	CodeOK             = 0    // 成功
	CodeBadRequest     = 40000 // 参数错误
	CodeUnauthorized   = 40100 // 未认证/凭证失效
	CodeForbidden      = 40300 // 无权限
	CodeNotFound       = 40400 // 资源不存在
	CodeConflict       = 40900 // 状态冲突
	CodeInternalError  = 50000 // 系统内部错误
	CodeUserExists     = 40001 // 用户名已存在
	CodeLoginFailed    = 40101 // 用户名或密码错误
	CodeUserDisabled   = 40301 // 用户被禁用
	CodeInvalidRole    = 40302 // 角色不合法
	CodeTeamFull       = 40901 // 团队人数已满
	CodeAlreadyJoined  = 40902 // 已加入该团队
	CodeAlreadyApplied = 40903 // 已报名该活动
	CodePointsNotEnough = 40904 // 积分不足
	CodeStockNotEnough = 40905 // 库存不足
	CodeNotRegistered  = 40906 // 未报名该活动
	CodeCheckpointDone = 40907 // 该打卡点已打卡
	CodeWrongAnswer    = 42200 // 答题错误
	CodeActivityClosed = 40908 // 活动不在可操作状态
	CodeAlreadyFavored = 40909 // 已收藏该线路
	CodeAlreadyMember  = 40910 // 已是该团队队员
)
