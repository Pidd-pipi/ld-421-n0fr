package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/labequipment/lab-equipment/internal/model"
	"gorm.io/gorm"
)

// RoleRepository 角色仓储接口。
type RoleRepository interface {
	Create(ctx context.Context, role *model.Role) error
	FindByCode(ctx context.Context, code string) (*model.Role, error)
	List(ctx context.Context) ([]model.Role, error)
}

type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository 构造角色仓储。
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *model.Role) error {
	if err := r.db.WithContext(ctx).Create(role).Error; err != nil {
		return fmt.Errorf("create role: %w", err)
	}
	return nil
}

func (r *roleRepository) FindByCode(ctx context.Context, code string) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find role by code: %w", err)
	}
	return &role, nil
}

func (r *roleRepository) List(ctx context.Context) ([]model.Role, error) {
	var roles []model.Role
	if err := r.db.WithContext(ctx).Order("id ASC").Find(&roles).Error; err != nil {
		return nil, fmt.Errorf("list roles: %w", err)
	}
	return roles, nil
}
