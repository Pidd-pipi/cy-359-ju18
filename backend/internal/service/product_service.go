package service

import (
	"fmt"
	"log/slog"

	"github.com/orienteering/platform/internal/dto"
	"github.com/orienteering/platform/internal/model"
	"github.com/orienteering/platform/internal/repository"
)

// ProductService 商城商品业务服务。
type ProductService struct {
	repo *repository.ProductRepository
}

// NewProductService 构造商品服务。
func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// Create 创建商品（管理员）。
func (s *ProductService) Create(req *dto.CreateProductRequest, logger *slog.Logger) (*model.Product, error) {
	product := &model.Product{
		Name: req.Name, Description: req.Description, Type: req.Type,
		PointsCost: req.PointsCost, Stock: req.Stock, Status: req.Status, ImageURL: req.ImageURL,
	}
	if err := s.repo.Create(product); err != nil {
		return nil, err
	}
	logger.Info("product create", "product_id", product.ID, "name", req.Name)
	return product, nil
}

// Update 编辑商品（管理员）。
func (s *ProductService) Update(productID int64, req *dto.CreateProductRequest, logger *slog.Logger) (*model.Product, error) {
	product, err := s.repo.GetByID(productID)
	if err != nil {
		return nil, err
	}
	product.Name = req.Name
	product.Description = req.Description
	product.Type = req.Type
	product.PointsCost = req.PointsCost
	product.Stock = req.Stock
	product.Status = req.Status
	product.ImageURL = req.ImageURL
	if err := s.repo.Update(product); err != nil {
		return nil, err
	}
	logger.Info("product update", "product_id", productID, "name", req.Name)
	return product, nil
}

// UpdateStatus 上下架商品（管理员）。
func (s *ProductService) UpdateStatus(productID int64, status string, logger *slog.Logger) (*model.Product, error) {
	product, err := s.repo.GetByID(productID)
	if err != nil {
		return nil, err
	}
	if err := s.repo.UpdateStatus(productID, status); err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf("product status change product_id=%d from=%s to=%s", productID, product.Status, status))
	product.Status = status
	return product, nil
}

// List 分页查询商品。
func (s *ProductService) List(page, pageSize, offset int, status, keyword string) ([]model.Product, int64, error) {
	return s.repo.List(page, pageSize, offset, status, keyword)
}

// Get 查询商品详情。
func (s *ProductService) Get(productID int64) (*model.Product, error) {
	return s.repo.GetByID(productID)
}

// ToProductView 商品转展示视图。
func ToProductView(p *model.Product) dto.ProductView {
	return dto.ProductView{
		ID: p.ID, Name: p.Name, Description: p.Description, Type: p.Type,
		PointsCost: p.PointsCost, Stock: p.Stock, Status: p.Status,
		ImageURL: p.ImageURL, CreatedAt: p.CreatedAt,
	}
}

