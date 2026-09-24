package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/labequipment/lab-equipment/internal/constants"
	"github.com/labequipment/lab-equipment/internal/model"
	"gorm.io/gorm"
)

// EquipmentRepository 设备仓储接口。
type EquipmentRepository interface {
	Create(ctx context.Context, equipment *model.Equipment) error
	Update(ctx context.Context, equipment *model.Equipment) error
	FindByID(ctx context.Context, id uint) (*model.Equipment, error)
	List(ctx context.Context, filter EquipmentFilter) ([]model.Equipment, int64, error)
	CountByStatus(ctx context.Context) (map[constants.AssetStatus]int64, error)
	ListExpiringWarranty(ctx context.Context, days int) ([]model.Equipment, error)
	UpdateOwner(ctx context.Context, id, ownerID uint) error
	UpdateStatus(ctx context.Context, id uint, status constants.AssetStatus) error
}

// EquipmentFilter 设备列表筛选条件。
type EquipmentFilter struct {
	Keyword    string
	CategoryID uint
	Status     constants.AssetStatus
	Pagination
}

type equipmentRepository struct {
	db *gorm.DB
}

// NewEquipmentRepository 构造设备仓储。
func NewEquipmentRepository(db *gorm.DB) EquipmentRepository {
	return &equipmentRepository{db: db}
}

func (r *equipmentRepository) Create(ctx context.Context, equipment *model.Equipment) error {
	if err := r.db.WithContext(ctx).Create(equipment).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrConflict
		}
		return fmt.Errorf("create equipment: %w", err)
	}
	return nil
}

func (r *equipmentRepository) Update(ctx context.Context, equipment *model.Equipment) error {
	if err := r.db.WithContext(ctx).Save(equipment).Error; err != nil {
		return fmt.Errorf("update equipment: %w", err)
	}
	return nil
}

func (r *equipmentRepository) FindByID(ctx context.Context, id uint) (*model.Equipment, error) {
	var equipment model.Equipment
	err := r.db.WithContext(ctx).Preload("Category").Preload("Owner.Role").First(&equipment, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find equipment by id: %w", err)
	}
	return &equipment, nil
}

func (r *equipmentRepository) List(ctx context.Context, filter EquipmentFilter) ([]model.Equipment, int64, error) {
	filter.Normalize()
	query := r.db.WithContext(ctx).Model(&model.Equipment{}).Preload("Category").Preload("Owner.Role")
	if filter.Keyword != "" {
		like := "%" + filter.Keyword + "%"
		query = query.Where("name LIKE ? OR code LIKE ? OR brand_model LIKE ?", like, like, like)
	}
	if filter.CategoryID > 0 {
		query = query.Where("category_id = ?", filter.CategoryID)
	}
	if filter.Status.Valid() {
		query = query.Where("status = ?", filter.Status)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count equipment: %w", err)
	}
	var list []model.Equipment
	if err := query.Order("id DESC").Offset(filter.Offset()).Limit(filter.PageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list equipment: %w", err)
	}
	return list, total, nil
}

func (r *equipmentRepository) CountByStatus(ctx context.Context) (map[constants.AssetStatus]int64, error) {
	rows := []struct {
		Status constants.AssetStatus
		Count  int64
	}{}
	if err := r.db.WithContext(ctx).Model(&model.Equipment{}).
		Select("status, count(*) as count").
		Group("status").
		Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("count equipment by status: %w", err)
	}
	result := make(map[constants.AssetStatus]int64, len(rows))
	for _, row := range rows {
		result[row.Status] = row.Count
	}
	return result, nil
}

func (r *equipmentRepository) ListExpiringWarranty(ctx context.Context, days int) ([]model.Equipment, error) {
	if days <= 0 {
		days = 30
	}
	now := time.Now()
	deadline := now.AddDate(0, 0, days)
	var list []model.Equipment
	err := r.db.WithContext(ctx).
		Preload("Category").
		Where("warranty_expiry IS NOT NULL").
		Where("warranty_expiry BETWEEN ? AND ?", now, deadline).
		Where("status NOT IN ?", []constants.AssetStatus{constants.AssetStatusRetired, constants.AssetStatusLost}).
		Order("warranty_expiry ASC").
		Find(&list).Error
	if err != nil {
		return nil, fmt.Errorf("list expiring warranty: %w", err)
	}
	return list, nil
}

func (r *equipmentRepository) UpdateOwner(ctx context.Context, id, ownerID uint) error {
	if err := r.db.WithContext(ctx).Model(&model.Equipment{}).Where("id = ?", id).Update("owner_id", ownerID).Error; err != nil {
		return fmt.Errorf("update equipment owner: %w", err)
	}
	return nil
}

func (r *equipmentRepository) UpdateStatus(ctx context.Context, id uint, status constants.AssetStatus) error {
	if err := r.db.WithContext(ctx).Model(&model.Equipment{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return fmt.Errorf("update equipment status: %w", err)
	}
	return nil
}
