package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/orienteering/platform/internal/model"
)

// ProductRepository 商品仓储。
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository 构造商品仓储。
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(product *model.Product) error {
	if err := r.db.Create(product).Error; err != nil {
		return fmt.Errorf("create product: %w", err)
	}
	return nil
}

func (r *ProductRepository) GetByID(id int64) (*model.Product, error) {
	var product model.Product
	if err := r.db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get product by id: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("get product by id: %w", err)
	}
	return &product, nil
}

func (r *ProductRepository) GetByIDForUpdate(tx *gorm.DB, id int64) (*model.Product, error) {
	var product model.Product
	if err := tx.Clauses(lockedClause()).First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get product for update: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("get product for update: %w", err)
	}
	return &product, nil
}

func (r *ProductRepository) Update(product *model.Product) error {
	if err := r.db.Save(product).Error; err != nil {
		return fmt.Errorf("update product: %w", err)
	}
	return nil
}

func (r *ProductRepository) UpdateStatus(productID int64, status string) error {
	if err := r.db.Model(&model.Product{}).Where("id = ?", productID).
		Update("status", status).Error; err != nil {
		return fmt.Errorf("update product status: %w", err)
	}
	return nil
}

func (r *ProductRepository) List(page, pageSize, offset int, status, keyword string) ([]model.Product, int64, error) {
	query := r.db.Model(&model.Product{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", like, like)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count products: %w", err)
	}
	var products []model.Product
	if err := query.Order("id DESC").Limit(pageSize).Offset(offset).Find(&products).Error; err != nil {
		return nil, 0, fmt.Errorf("list products: %w", err)
	}
	return products, total, nil
}

func (r *ProductRepository) DecrementStock(tx *gorm.DB, productID int64, quantity int) error {
	res := tx.Model(&model.Product{}).
		Where("id = ? AND stock >= ?", productID, quantity).
		UpdateColumn("stock", gorm.Expr("stock - ?", quantity))
	if res.Error != nil {
		return fmt.Errorf("decrement stock: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotEnoughStock
	}
	return nil
}
