package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/orienteering/platform/internal/model"
)

// UserRepository 用户仓储。
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 构造用户仓储。
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetByID(id int64) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get user by id: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get user by username: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("get user by username: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	if err := r.db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, fmt.Errorf("exists user by username: %w", err)
	}
	return count > 0, nil
}

func (r *UserRepository) Update(user *model.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (r *UserRepository) AddPoints(userID int64, delta int) error {
	if err := r.db.Model(&model.User{}).Where("id = ?", userID).
		UpdateColumn("points", gorm.Expr("points + ?", delta)).Error; err != nil {
		return fmt.Errorf("add points: %w", err)
	}
	return nil
}

func (r *UserRepository) List(page, pageSize, offset int, keyword string) ([]model.User, int64, error) {
	query := r.db.Model(&model.User{})
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("username LIKE ? OR nickname LIKE ?", like, like)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}
	var users []model.User
	if err := query.Order("id DESC").Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	return users, total, nil
}

func (r *UserRepository) SetDisabled(userID int64, disabled bool) error {
	if err := r.db.Model(&model.User{}).Where("id = ?", userID).
		Update("disabled", disabled).Error; err != nil {
		return fmt.Errorf("set user disabled: %w", err)
	}
	return nil
}
