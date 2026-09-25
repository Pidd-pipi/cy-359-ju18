package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/orienteering/platform/internal/model"
)

// RedemptionRepository 兑换订单仓储。
type RedemptionRepository struct {
	db *gorm.DB
}

// NewRedemptionRepository 构造兑换订单仓储。
func NewRedemptionRepository(db *gorm.DB) *RedemptionRepository {
	return &RedemptionRepository{db: db}
}

func (r *RedemptionRepository) Create(tx *gorm.DB, redemption *model.Redemption) error {
	if err := tx.Create(redemption).Error; err != nil {
		return fmt.Errorf("create redemption: %w", err)
	}
	return nil
}

func (r *RedemptionRepository) GetByID(id int64) (*model.Redemption, error) {
	var redemption model.Redemption
	if err := r.db.First(&redemption, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get redemption by id: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("get redemption by id: %w", err)
	}
	return &redemption, nil
}

func (r *RedemptionRepository) ListByUser(userID int64, page, pageSize, offset int) ([]model.Redemption, int64, error) {
	query := r.db.Model(&model.Redemption{}).Where("user_id = ?", userID)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count redemptions: %w", err)
	}
	var redemptions []model.Redemption
	if err := query.Order("id DESC").Limit(pageSize).Offset(offset).Find(&redemptions).Error; err != nil {
		return nil, 0, fmt.Errorf("list redemptions: %w", err)
	}
	return redemptions, total, nil
}

func (r *RedemptionRepository) ListAll(page, pageSize, offset int, status string) ([]model.Redemption, int64, error) {
	query := r.db.Model(&model.Redemption{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count all redemptions: %w", err)
	}
	var redemptions []model.Redemption
	if err := query.Order("id DESC").Limit(pageSize).Offset(offset).Find(&redemptions).Error; err != nil {
		return nil, 0, fmt.Errorf("list all redemptions: %w", err)
	}
	return redemptions, total, nil
}

func (r *RedemptionRepository) UpdateStatus(tx *gorm.DB, id int64, status string) error {
	if err := tx.Model(&model.Redemption{}).Where("id = ?", id).
		Update("status", status).Error; err != nil {
		return fmt.Errorf("update redemption status: %w", err)
	}
	return nil
}

// WithTx 在事务中执行业务函数。
func (r *RedemptionRepository) WithTx(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
