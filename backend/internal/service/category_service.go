package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	apperrors "github.com/labequipment/lab-equipment/internal/errors"
	"github.com/labequipment/lab-equipment/internal/model"
	"github.com/labequipment/lab-equipment/internal/repository"
)

// CategoryService 设备分类服务。
type CategoryService struct {
	repo   repository.CategoryRepository
	logger *slog.Logger
}

// NewCategoryService 构造分类服务。
func NewCategoryService(repo repository.CategoryRepository, logger *slog.Logger) *CategoryService {
	return &CategoryService{repo: repo, logger: logger}
}

// Create 新增分类。
func (s *CategoryService) Create(ctx context.Context, category *model.EquipmentCategory) (*model.EquipmentCategory, error) {
	if category.ParentID != nil && *category.ParentID > 0 {
		if _, err := s.repo.FindByID(ctx, *category.ParentID); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, apperrors.NewBusinessError(40400, 404, "父分类不存在")
			}
			return nil, fmt.Errorf("find parent category: %w", err)
		}
	}
	if err := s.repo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("create category: %w", err)
	}
	return category, nil
}

// Update 更新分类。
func (s *CategoryService) Update(ctx context.Context, id uint, category *model.EquipmentCategory) (*model.EquipmentCategory, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperrors.NewBusinessError(40400, 404, "分类不存在")
		}
		return nil, fmt.Errorf("find category: %w", err)
	}
	existing.Name = category.Name
	existing.ParentID = category.ParentID
	existing.Description = category.Description
	existing.Icon = category.Icon
	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("update category: %w", err)
	}
	return existing, nil
}

// Delete 删除分类。
func (s *CategoryService) Delete(ctx context.Context, id uint) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return apperrors.NewBusinessError(40400, 404, "分类不存在")
		}
		return fmt.Errorf("find category: %w", err)
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	return nil
}

// ListTree 返回树形分类。
func (s *CategoryService) ListTree(ctx context.Context) ([]model.EquipmentCategory, error) {
	all, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	return buildCategoryTree(all), nil
}

func buildCategoryTree(categories []model.EquipmentCategory) []model.EquipmentCategory {
	nodeMap := make(map[uint]*model.EquipmentCategory, len(categories))
	for i := range categories {
		item := categories[i]
		item.Children = []model.EquipmentCategory{}
		nodeMap[item.ID] = &item
	}
	roots := make([]*model.EquipmentCategory, 0)
	for i := range categories {
		item := categories[i]
		node := nodeMap[item.ID]
		if item.ParentID != nil && *item.ParentID > 0 {
			if parent, ok := nodeMap[*item.ParentID]; ok {
				parent.Children = append(parent.Children, *node)
				continue
			}
		}
		roots = append(roots, node)
	}
	result := make([]model.EquipmentCategory, 0, len(roots))
	for _, root := range roots {
		result = append(result, *root)
	}
	return result
}
