package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/labequipment/lab-equipment/internal/model"
	"gorm.io/gorm"
)

// UserRepository 用户仓储接口。
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id uint) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	List(ctx context.Context, filter UserFilter) ([]model.User, int64, error)
}

// UserFilter 用户列表筛选条件。
type UserFilter struct {
	Keyword string
	RoleID  uint
	Pagination
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 构造用户仓储。
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrConflict
		}
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Preload("Role").First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Preload("Role").Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by username: %w", err)
	}
	return &user, nil
}

func (r *userRepository) List(ctx context.Context, filter UserFilter) ([]model.User, int64, error) {
	filter.Normalize()
	query := r.db.WithContext(ctx).Model(&model.User{}).Preload("Role")
	if filter.Keyword != "" {
		like := "%" + filter.Keyword + "%"
		query = query.Where("username LIKE ? OR name LIKE ?", like, like)
	}
	if filter.RoleID > 0 {
		query = query.Where("role_id = ?", filter.RoleID)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}
	var users []model.User
	if err := query.Order("id ASC").Offset(filter.Offset()).Limit(filter.PageSize).Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	return users, total, nil
}
