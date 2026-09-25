package service

import (
	"fmt"
	"log/slog"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/orienteering/platform/internal/constants"
	"github.com/orienteering/platform/internal/dto"
	"github.com/orienteering/platform/internal/model"
	"github.com/orienteering/platform/internal/repository"
	"github.com/orienteering/platform/internal/util"
	"github.com/orienteering/platform/pkg/ws"
)

// CheckinService 打卡业务服务。
type CheckinService struct {
	checkin     *repository.CheckinRepository
	checkpoint  *repository.CheckpointRepository
	registration *repository.RegistrationRepository
	activity    *repository.ActivityRepository
	user        *repository.UserRepository
	leaderboard *LeaderboardService
	hub         *ws.Hub
}

// NewCheckinService 构造打卡服务。
func NewCheckinService(checkin *repository.CheckinRepository, checkpoint *repository.CheckpointRepository, registration *repository.RegistrationRepository, activity *repository.ActivityRepository, user *repository.UserRepository, leaderboard *LeaderboardService, hub *ws.Hub) *CheckinService {
	return &CheckinService{checkin: checkin, checkpoint: checkpoint, registration: registration, activity: activity, user: user, leaderboard: leaderboard, hub: hub}
}

// Checkin 打卡（事务）：校验活动进行中、报名通过、未重复打卡；GPS 距离校验 + 答题校验；记录打卡、加分、可能完赛。
func (s *CheckinService) Checkin(teamID, userID int64, req *dto.CheckinRequest, logger *slog.Logger) (*model.CheckinRecord, error) {
	cp, err := s.checkpoint.GetByID(req.CheckpointID)
	if err != nil {
		return nil, err
	}
	var record *model.CheckinRecord
	err = s.checkin.WithTx(func(tx *gorm.DB) error {
		activity, err := s.activity.GetByIDForUpdate(tx, cp.ActivityID)
		if err != nil {
			return err
		}
		if activity.Status != constants.ActivityStatusOngoing {
			return util.NewAppError(constants.CodeActivityClosed,
				fmt.Sprintf(constants.MsgActivityClosed, activity.Title, activity.Status), nil)
		}
		reg, err := s.registration.GetByTeamAndActivity(teamID, cp.ActivityID)
		if err != nil {
			if isNotFound(err) {
				return util.NewAppError(constants.CodeNotRegistered,
					fmt.Sprintf(constants.MsgNotRegistered, teamID, activity.Title), nil)
			}
			return err
		}
		if reg.Status != constants.RegistrationStatusApproved {
			return util.NewAppError(constants.CodeConflict,
				fmt.Sprintf("报名记录[%d]状态[%s]不可打卡", reg.ID, reg.Status), nil)
		}
		done, err := s.checkin.Exists(tx, req.CheckpointID, teamID)
		if err != nil {
			return err
		}
		if done {
			return util.NewAppError(constants.CodeCheckpointDone,
				fmt.Sprintf(constants.MsgCheckpointDone, cp.Name, teamID), nil)
		}
		// GPS 距离校验
		if req.CheckinType == constants.CheckinTypeGPS {
			distance := haversine(req.Latitude, req.Longitude, cp.Lat, cp.Lng)
			if distance > float64(cp.RadiusMeters) {
				return util.NewAppError(constants.CodeBadRequest,
					fmt.Sprintf("打卡点[%s]距离[%.0f米]超出范围[%d米]", cp.Name, distance, cp.RadiusMeters), nil)
			}
		}
		// 答题任务校验
		result := constants.CheckinResultNone
		pointsEarned := 10
		if cp.TaskType == constants.TaskTypeQuiz {
			if strings.TrimSpace(req.Answer) != strings.TrimSpace(cp.ExpectedAnswer) {
				return util.NewAppError(constants.CodeWrongAnswer,
					fmt.Sprintf(constants.MsgWrongAnswer, cp.Name, req.Answer), nil)
			}
			result = constants.CheckinResultCorrect
			pointsEarned = 20
		}
		record = &model.CheckinRecord{
			ActivityID: cp.ActivityID, CheckpointID: cp.ID, TeamID: teamID,
			UserID: userID, CheckinType: req.CheckinType, Latitude: req.Latitude,
			Longitude: req.Longitude, Answer: req.Answer, PhotoURL: req.PhotoURL,
			Result: result, PointsEarned: pointsEarned, CheckedInAt: time.Now(),
		}
		if err := s.checkin.Create(tx, record); err != nil {
			return err
		}
		if err := s.user.AddPoints(userID, pointsEarned); err != nil {
			return err
		}
		// 检查是否完成全部打卡点 -> 完赛
		total, err := s.checkpoint.CountByActivity(cp.ActivityID)
		if err != nil {
			return err
		}
		checked, err := s.checkin.CountByTeam(tx, cp.ActivityID, teamID)
		if err != nil {
			return err
		}
		if checked >= total && reg.StartTime != nil {
			now := time.Now()
			reg.Status = constants.RegistrationStatusFinished
			reg.FinishTime = &now
			reg.TotalSeconds = int(now.Sub(*reg.StartTime).Seconds())
			if err := s.registration.Update(tx, reg); err != nil {
				return err
			}
			logger.Info(fmt.Sprintf(constants.LogTeamFinish, reg.ID, reg.TotalSeconds))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf(constants.LogCheckin, cp.ID, teamID, req.CheckinType, record.Result, record.PointsEarned))
	// 失效排行榜缓存并推送实时更新
	_ = s.leaderboard.Invalidate(cp.ActivityID)
	if s.hub != nil {
		s.hub.Broadcast(cp.ActivityID, ws.Message{Type: "leaderboard_update", ActivityID: cp.ActivityID})
	}
	return record, nil
}

// ListByTeam 查询团队打卡记录。
func (s *CheckinService) ListByTeam(teamID, activityID int64) ([]model.CheckinRecord, error) {
	return s.checkin.ListByTeam(teamID, activityID)
}

// ListByActivity 查询活动全部打卡记录（管理员）。
func (s *CheckinService) ListByActivity(activityID int64) ([]model.CheckinRecord, error) {
	return s.checkin.ListByActivity(activityID)
}

func isNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), "not found")
}

// haversine 计算两个经纬度点的大圆距离（米）。
func haversine(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371000.0
	toRad := math.Pi / 180.0
	dLat := (lat2 - lat1) * toRad
	dLng := (lng2 - lng1) * toRad
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*toRad)*math.Cos(lat2*toRad)*math.Sin(dLng/2)*math.Sin(dLng/2)
	return 2 * earthRadius * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}
