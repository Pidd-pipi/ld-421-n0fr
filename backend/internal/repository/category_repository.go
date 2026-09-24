package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/labequipment/lab-equipment/internal/model"
	"gorm.io/gorm"
)

// CategoryRepository 设备分类仓储接口。
type CategoryRepository interface {
	Create(ctx context.Context, category *model.EquipmentCategory) error
	Update(ctx context.Context, category *model.EquipmentCategory) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*model.EquipmentCategory, error)
	ListAll(ctx context.Context) ([]model.EquipmentCategory, error)
}

type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 构造分类仓储。
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *model.EquipmentCategory) error {
	if err := r.db.WithContext(ctx).Create(category).Error; err != nil {
		return fmt.Errorf("create category: %w", err)
	}
	return nil
}

func (r *categoryRepository) Update(ctx context.Context, category *model.EquipmentCategory) error {
	if err := r.db.WithContext(ctx).Save(category).Error; err != nil {
		return fmt.Errorf("update category: %w", err)
	}
	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.EquipmentCategory{}, id).Error; err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	return nil
}

func (r *categoryRepository) FindByID(ctx context.Context, id uint) (*model.EquipmentCategory, error) {
	var category model.EquipmentCategory
	err := r.db.WithContext(ctx).First(&category, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find category by id: %w", err)
	}
	return &category, nil
}

func (r *categoryRepository) ListAll(ctx context.Context) ([]model.EquipmentCategory, error) {
	var categories []model.EquipmentCategory
	if err := r.db.WithContext(ctx).Order("id ASC").Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	return categories, nil
}
