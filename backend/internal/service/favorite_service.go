package service

import (
	"fmt"
	"log/slog"

	"github.com/orienteering/platform/internal/constants"
	"github.com/orienteering/platform/internal/dto"
	"github.com/orienteering/platform/internal/model"
	"github.com/orienteering/platform/internal/repository"
	"github.com/orienteering/platform/internal/util"
)

// FavoriteService 历史线路收藏业务服务。
type FavoriteService struct {
	repo *repository.FavoriteRepository
}

// NewFavoriteService 构造收藏服务。
func NewFavoriteService(repo *repository.FavoriteRepository) *FavoriteService {
	return &FavoriteService{repo: repo}
}

// Add 收藏线路（仅已结束活动可收藏）。
func (s *FavoriteService) Add(userID, activityID int64, activityStatus string, logger *slog.Logger) (*model.Favorite, error) {
	if activityStatus != constants.ActivityStatusFinished {
		return nil, util.NewAppError(constants.CodeActivityClosed,
			fmt.Sprintf("活动[%d]状态[%s]不可收藏，仅已结束线路可收藏", activityID, activityStatus), nil)
	}
	exists, err := s.repo.Exists(userID, activityID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, util.NewAppError(constants.CodeAlreadyFavored,
			fmt.Sprintf("用户[%d]已收藏活动[%d]", userID, activityID), nil)
	}
	fav := &model.Favorite{UserID: userID, ActivityID: activityID}
	if err := s.repo.Create(fav); err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf(constants.LogFavorite, userID, activityID, "add"))
	return fav, nil
}

// Remove 取消收藏。
func (s *FavoriteService) Remove(userID, activityID int64, logger *slog.Logger) error {
	exists, err := s.repo.Exists(userID, activityID)
	if err != nil {
		return err
	}
	if !exists {
		return util.NewAppError(constants.CodeNotFound,
			fmt.Sprintf("用户[%d]未收藏活动[%d]", userID, activityID), nil)
	}
	if err := s.repo.Delete(userID, activityID); err != nil {
		return err
	}
	logger.Info(fmt.Sprintf(constants.LogFavorite, userID, activityID, "remove"))
	return nil
}

// List 分页查询我的收藏。
func (s *FavoriteService) List(userID int64, page, pageSize, offset int) ([]model.Favorite, int64, error) {
	return s.repo.ListByUser(userID, page, pageSize, offset)
}

// IsFavored 查询是否已收藏（复用 repo.Exists）。
func (s *FavoriteService) IsFavored(userID, activityID int64) (bool, error) {
	return s.repo.Exists(userID, activityID)
}

// ToFavoriteView 收藏转展示视图。
func ToFavoriteView(f *model.Favorite, activity dto.ActivityView) dto.FavoriteView {
	return dto.FavoriteView{ID: f.ID, ActivityID: f.ActivityID, Activity: activity, CreatedAt: f.CreatedAt}
}
