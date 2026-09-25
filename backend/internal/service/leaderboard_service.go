package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/orienteering/platform/internal/constants"
	"github.com/orienteering/platform/internal/dto"
	"github.com/orienteering/platform/internal/repository"
	"github.com/orienteering/platform/internal/util"
)

// LeaderboardService 排行榜服务，Redis 缓存 + DB 兜底。
type LeaderboardService struct {
	registration *repository.RegistrationRepository
	checkin      *repository.CheckinRepository
	activity     *repository.ActivityRepository
	rdb          *redis.Client
	logger       *slog.Logger
}

// NewLeaderboardService 构造排行榜服务。
func NewLeaderboardService(registration *repository.RegistrationRepository, checkin *repository.CheckinRepository, activity *repository.ActivityRepository, rdb *redis.Client) *LeaderboardService {
	return &LeaderboardService{registration: registration, checkin: checkin, activity: activity, rdb: rdb, logger: slog.Default()}
}

// SetLogger 注入日志器。
func (s *LeaderboardService) SetLogger(logger *slog.Logger) {
	s.logger = logger
}

// Leaderboard 计算活动排行榜：已完成团队按用时升序，未完成按打卡数降序。
func (s *LeaderboardService) Leaderboard(activityID int64) ([]dto.LeaderboardRow, error) {
	if cached, err := s.fromCache(activityID); err == nil && cached != nil {
		s.logger.Info(fmt.Sprintf(constants.LogLeaderboard, activityID, true))
		return cached, nil
	}
	if _, err := s.activity.GetByID(activityID); err != nil {
		return nil, err
	}
	regs, err := s.registration.ListByActivity(activityID)
	if err != nil {
		return nil, err
	}
	rows := make([]dto.LeaderboardRow, 0, len(regs))
	finished := make([]dto.LeaderboardRow, 0)
	inProgress := make([]dto.LeaderboardRow, 0)
	for _, reg := range regs {
		if reg.Status != constants.RegistrationStatusApproved && reg.Status != constants.RegistrationStatusFinished {
			continue
		}
		count, err := s.checkin.CountByTeam(nil, activityID, reg.TeamID)
		if err != nil {
			return nil, err
		}
		row := dto.LeaderboardRow{
			TeamID: reg.TeamID, TotalSeconds: reg.TotalSeconds,
			Duration: util.FormatDuration(reg.TotalSeconds),
			CheckpointCount: int(count), Status: reg.Status,
		}
		if reg.Status == constants.RegistrationStatusFinished {
			finished = append(finished, row)
		} else {
			inProgress = append(inProgress, row)
		}
	}
	// 稳定排序：已完成按用时升序；进行中按打卡数降序、团队 ID 升序
	finished = sortFinished(finished)
	inProgress = sortInProgress(inProgress)
	rows = append(finished, inProgress...)
	for i := range rows {
		rows[i].Rank = i + 1
	}
	_ = s.cache(activityID, rows)
	s.logger.Info(fmt.Sprintf(constants.LogLeaderboard, activityID, false))
	return rows, nil
}

// Invalidate 清除活动排行榜缓存。
func (s *LeaderboardService) Invalidate(activityID int64) error {
	if s.rdb == nil {
		return nil
	}
	return s.rdb.Del(context.Background(), cacheKey(activityID)).Err()
}

func (s *LeaderboardService) fromCache(activityID int64) ([]dto.LeaderboardRow, error) {
	if s.rdb == nil {
		return nil, nil
	}
	raw, err := s.rdb.Get(context.Background(), cacheKey(activityID)).Bytes()
	if err != nil {
		return nil, err
	}
	var rows []dto.LeaderboardRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *LeaderboardService) cache(activityID int64, rows []dto.LeaderboardRow) error {
	if s.rdb == nil || len(rows) == 0 {
		return nil
	}
	raw, err := json.Marshal(rows)
	if err != nil {
		return err
	}
	return s.rdb.Set(context.Background(), cacheKey(activityID), raw, 5*time.Minute).Err()
}

func cacheKey(activityID int64) string {
	return fmt.Sprintf("orienteering:leaderboard:%d", activityID)
}

func sortFinished(rows []dto.LeaderboardRow) []dto.LeaderboardRow {
	for i := 0; i < len(rows); i++ {
		for j := i + 1; j < len(rows); j++ {
			if rows[j].TotalSeconds < rows[i].TotalSeconds ||
				(rows[j].TotalSeconds == rows[i].TotalSeconds && rows[j].TeamID < rows[i].TeamID) {
				rows[i], rows[j] = rows[j], rows[i]
			}
		}
	}
	return rows
}

func sortInProgress(rows []dto.LeaderboardRow) []dto.LeaderboardRow {
	for i := 0; i < len(rows); i++ {
		for j := i + 1; j < len(rows); j++ {
			if rows[j].CheckpointCount > rows[i].CheckpointCount ||
				(rows[j].CheckpointCount == rows[i].CheckpointCount && rows[j].TeamID < rows[i].TeamID) {
				rows[i], rows[j] = rows[j], rows[i]
			}
		}
	}
	return rows
}

