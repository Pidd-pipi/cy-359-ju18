package service

import (
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"

	"github.com/orienteering/platform/internal/constants"
	"github.com/orienteering/platform/internal/dto"
	"github.com/orienteering/platform/internal/model"
	"github.com/orienteering/platform/internal/repository"
	"github.com/orienteering/platform/internal/util"
)

// RedemptionService 积分兑换业务服务。
type RedemptionService struct {
	repo    *repository.RedemptionRepository
	product *repository.ProductRepository
	user    *repository.UserRepository
}

// NewRedemptionService 构造兑换服务。
func NewRedemptionService(repo *repository.RedemptionRepository, product *repository.ProductRepository, user *repository.UserRepository) *RedemptionService {
	return &RedemptionService{repo: repo, product: product, user: user}
}

// Redeem 兑换商品（事务：行锁商品、校验库存与积分、扣库存、扣积分、创建订单）。
func (s *RedemptionService) Redeem(userID int64, req *dto.RedeemRequest, logger *slog.Logger) (*model.Redemption, error) {
	var redemption *model.Redemption
	err := s.repo.WithTx(func(tx *gorm.DB) error {
		product, err := s.product.GetByIDForUpdate(tx, req.ProductID)
		if err != nil {
			return err
		}
		if product.Status != constants.ProductStatusOn {
			return util.NewAppError(constants.CodeConflict,
				fmt.Sprintf("商品[%s]已下架", product.Name), nil)
		}
		total := product.PointsCost * req.Quantity
		user, err := s.user.GetByID(userID)
		if err != nil {
			return err
		}
		if user.Points < total {
			return util.NewAppError(constants.CodePointsNotEnough,
				fmt.Sprintf(constants.MsgPointsNotEnough, user.Username, user.Points, total), nil)
		}
		if err := s.product.DecrementStock(tx, req.ProductID, req.Quantity); err != nil {
			if isNotFound(err) {
				return util.NewAppError(constants.CodeStockNotEnough,
					fmt.Sprintf(constants.MsgStockNotEnough, product.Name, product.Stock), nil)
			}
			return err
		}
		if err := s.user.AddPoints(userID, -total); err != nil {
			return err
		}
		redemption = &model.Redemption{
			UserID: userID, ProductID: req.ProductID, Quantity: req.Quantity,
			PointsCost: product.PointsCost, TotalPoints: total,
			Status: constants.RedemptionStatusPending, RedeemedAt: time.Now(),
		}
		return s.repo.Create(tx, redemption)
	})
	if err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf(constants.LogRedeem, userID, req.ProductID, req.Quantity, redemption.TotalPoints))
	return redemption, nil
}

// ListMine 查询我的兑换记录。
func (s *RedemptionService) ListMine(userID int64, page, pageSize, offset int) ([]model.Redemption, int64, error) {
	return s.repo.ListByUser(userID, page, pageSize, offset)
}

// ListAll 分页查询全部兑换记录（管理员）。
func (s *RedemptionService) ListAll(page, pageSize, offset int, status string) ([]model.Redemption, int64, error) {
	return s.repo.ListAll(page, pageSize, offset, status)
}

// UpdateStatus 兑换订单状态流转（管理员）：pending -> completed / cancelled。
func (s *RedemptionService) UpdateStatus(redemptionID int64, status string, logger *slog.Logger) (*model.Redemption, error) {
	redemption, err := s.repo.GetByID(redemptionID)
	if err != nil {
		return nil, err
	}
	from := redemption.Status
	valid := false
	switch status {
	case constants.RedemptionStatusCompleted, constants.RedemptionStatusCancelled:
		valid = from == constants.RedemptionStatusPending
	}
	if !valid {
		return nil, util.NewAppError(constants.CodeConflict,
			fmt.Sprintf("兑换订单[%d]当前状态[%s]不允许流转到[%s]", redemptionID, from, status), nil)
	}
	err = s.repo.WithTx(func(tx *gorm.DB) error {
		if err := s.repo.UpdateStatus(tx, redemptionID, status); err != nil {
			return err
		}
		if status == constants.RedemptionStatusCompleted {
			return tx.Model(&model.Redemption{}).Where("id = ?", redemptionID).
				Update("completed_at", time.Now()).Error
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf(constants.LogRedeemStatus, redemptionID, from, status))
	return s.repo.GetByID(redemptionID)
}

// Get 查询兑换详情。
func (s *RedemptionService) Get(redemptionID int64) (*model.Redemption, error) {
	return s.repo.GetByID(redemptionID)
}
