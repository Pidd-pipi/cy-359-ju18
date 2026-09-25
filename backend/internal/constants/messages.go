package constants

// 接口返回文案集中定义，供 handler / service 手动拼接使用。
const (
	MsgOK          = "ok"
	MsgCreated     = "创建成功"
	MsgUpdated     = "更新成功"
	MsgDeleted     = "删除成功"
	MsgRegistered  = "注册成功"
	MsgLoggedIn    = "登录成功"
	MsgCheckinOK   = "打卡成功"
	MsgRedeemOK    = "兑换成功"
	MsgFavoriteOK  = "收藏成功"
	MsgUnfavoriteOK = "取消收藏成功"

	MsgUserNotFound     = "用户[%s]不存在"
	MsgActivityNotFound = "活动[%s]不存在"
	MsgCheckpointNotFound = "打卡点[%s]不存在"
	MsgTeamNotFound     = "团队[%s]不存在"
	MsgProductNotFound  = "商品[%s]不存在"
	MsgRegistrationNotFound = "报名记录[%s]不存在"
	MsgPermissionDenied = "角色[%s]无权执行该操作"

	MsgInvalidStatus = "活动[%s]当前状态[%s]不允许执行[%s]"
	MsgInvalidRole   = "用户[%s]角色[%s]不合法"
	MsgWrongAnswer   = "打卡点[%s]答案[%s]错误"
	MsgAlreadyJoined = "用户[%s]已加入团队[%s]"
	MsgTeamFull      = "团队[%d]人数已满"
	MsgAlreadyApplied = "团队[%s]已报名活动[%s]"
	MsgPointsNotEnough = "用户[%s]积分[%d]不足，兑换需要[%d]"
	MsgStockNotEnough = "商品[%s]库存[%d]不足"
	MsgNotRegistered  = "团队[%d]未报名活动[%s]"
	MsgCheckpointDone = "打卡点[%s]已被团队[%d]打卡"
	MsgActivityClosed = "活动[%s]当前状态[%s]不可打卡"
)
