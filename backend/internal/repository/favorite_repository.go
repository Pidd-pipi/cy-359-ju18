package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/orienteering/platform/internal/model"
)

// FavoriteRepository 收藏仓储。
type FavoriteRepository struct {
	db *gorm.DB
}

// NewFavoriteRepository 构造收藏仓储。
func NewFavoriteRepository(db *gorm.DB) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

func (r *FavoriteRepository) Create(fav *model.Favorite) error {
	if err := r.db.Create(fav).Error; err != nil {
		return fmt.Errorf("create favorite: %w", err)
	}
	return nil
}

func (r *FavoriteRepository) Exists(userID, activityID int64) (bool, error) {
	var count int64
	if err := r.db.Model(&model.Favorite{}).
		Where("user_id = ? AND activity_id = ?", userID, activityID).Count(&count).Error; err != nil {
		return false, fmt.Errorf("check favorite: %w", err)
	}
	return count > 0, nil
}

func (r *FavoriteRepository) Delete(userID, activityID int64) error {
	if err := r.db.Where("user_id = ? AND activity_id = ?", userID, activityID).
		Delete(&model.Favorite{}).Error; err != nil {
		return fmt.Errorf("delete favorite: %w", err)
	}
	return nil
}

func (r *FavoriteRepository) ListByUser(userID int64, page, pageSize, offset int) ([]model.Favorite, int64, error) {
	query := r.db.Model(&model.Favorite{}).Where("user_id = ?", userID)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count favorites: %w", err)
	}
	var favorites []model.Favorite
	if err := query.Preload("Activity").
		Order("id DESC").Limit(pageSize).Offset(offset).Find(&favorites).Error; err != nil {
		return nil, 0, fmt.Errorf("list favorites: %w", err)
	}
	return favorites, total, nil
}

func (r *FavoriteRepository) Get(userID, activityID int64) (*model.Favorite, error) {
	var fav model.Favorite
	if err := r.db.Where("user_id = ? AND activity_id = ?", userID, activityID).First(&fav).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get favorite: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("get favorite: %w", err)
	}
	return &fav, nil
}
