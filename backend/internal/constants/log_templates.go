package constants

// 日志模板集中定义（屎山要求：>=25 条，全部 handler/service/middleware 引用）。
const (
	LogConfigLoaded      = "config loaded app_name=%s env=%s port=%s"
	LogDBConnected       = "database connected host=%s db=%s"
	LogDBNotConnected    = "database connect failed err=%v"
	LogRedisConnected    = "redis connected addr=%s"
	LogRedisNotConnected = "redis connect failed addr=%s err=%v"
	LogServerStarted     = "server started addr=:%s"
	LogServerStopped     = "server stopped err=%v"

	LogUserRegister     = "user register username=%s role=%s"
	LogUserLogin        = "user login username=%s role=%s"
	LogUserProfile      = "user profile query user_id=%d"
	LogUserUpdate       = "user update user_id=%d fields=%s"
	LogUserPointsChange = "user points change user_id=%d delta=%d balance=%d reason=%s"

	LogActivityCreate = "activity create title=%s difficulty=%s creator_id=%d"
	LogActivityStatus = "activity status change activity_id=%d from=%s to=%s operator_id=%d"
	LogActivityList   = "activity list page=%d page_size=%d status=%s"

	LogCheckpointCreate = "checkpoint create activity_id=%d name=%s seq=%d"
	LogCheckpointDelete = "checkpoint delete checkpoint_id=%d activity_id=%d"

	LogTeamCreate    = "team create name=%s captain_id=%d"
	LogTeamJoin      = "team join team_id=%d user_id=%d role=%s"
	LogTeamApply     = "team apply team_id=%d activity_id=%d status=%s occupied=%d max=%d"
	LogTeamApprove   = "team approve registration_id=%d team_id=%d"
	LogTeamReject    = "team reject registration_id=%d team_id=%d promoted_registration_id=%d"
	LogTeamWaitlistPromote = "team waitlist promote activity_id=%d registration_id=%d team_id=%d"
	LogTeamFinish    = "team finish registration_id=%d total_seconds=%d"

	LogCheckin      = "checkin checkpoint_id=%d team_id=%d type=%s result=%s points=%d"
	LogLeaderboard  = "leaderboard query activity_id=%d cached=%t"
	LogRedeem       = "redeem user_id=%d product_id=%d quantity=%d points=%d"
	LogRedeemStatus = "redeem status change redemption_id=%d from=%s to=%s"
	LogFavorite     = "favorite user_id=%d activity_id=%d op=%s"

	LogAuditWrite   = "audit write user_id=%d action=%s resource=%s"
	LogRequestStart = "request start request_id=%s method=%s path=%s"
	LogRequestEnd   = "request end request_id=%s method=%s path=%s status=%d latency_ms=%d"
	LogPanicRecover = "panic recovered request_id=%s err=%v"
	LogRateLimited  = "rate limited request_id=%s path=%s"
)
