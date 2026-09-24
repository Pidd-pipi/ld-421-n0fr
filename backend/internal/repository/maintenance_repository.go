package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/labequipment/lab-equipment/internal/constants"
	"github.com/labequipment/lab-equipment/internal/model"
	"gorm.io/gorm"
)

// MaintenanceRepository 维护记录仓储接口。
type MaintenanceRepository interface {
	Create(ctx context.Context, record *model.MaintenanceRecord) error
	Update(ctx context.Context, record *model.MaintenanceRecord) error
	FindByID(ctx context.Context, id uint) (*model.MaintenanceRecord, error)
	List(ctx context.Context, filter MaintenanceFilter) ([]model.MaintenanceRecord, int64, error)
	Stats(ctx context.Context) (*MaintenanceStats, error)
}

// MaintenanceFilter 维护记录筛选条件。
type MaintenanceFilter struct {
	EquipmentID uint
	Pagination
}

// MaintenanceStats 维护统计结果。
type MaintenanceStats struct {
	TotalRecords int64
	TotalCost    float64
	ResultCount  map[constants.MaintenanceResult]int64
}

type maintenanceRepository struct {
	db *gorm.DB
}

// NewMaintenanceRepository 构造维护记录仓储。
func NewMaintenanceRepository(db *gorm.DB) MaintenanceRepository {
	return &maintenanceRepository{db: db}
}

func (r *maintenanceRepository) Create(ctx context.Context, record *model.MaintenanceRecord) error {
	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return fmt.Errorf("create maintenance record: %w", err)
	}
	return nil
}

func (r *maintenanceRepository) Update(ctx context.Context, record *model.MaintenanceRecord) error {
	if err := r.db.WithContext(ctx).Save(record).Error; err != nil {
		return fmt.Errorf("update maintenance record: %w", err)
	}
	return nil
}

func (r *maintenanceRepository) FindByID(ctx context.Context, id uint) (*model.MaintenanceRecord, error) {
	var record model.MaintenanceRecord
	err := r.db.WithContext(ctx).Preload("Equipment.Category").Preload("Maintainer.Role").First(&record, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find maintenance record by id: %w", err)
	}
	return &record, nil
}

func (r *maintenanceRepository) List(ctx context.Context, filter MaintenanceFilter) ([]model.MaintenanceRecord, int64, error) {
	filter.Normalize()
	query := r.db.WithContext(ctx).Model(&model.MaintenanceRecord{}).Preload("Equipment.Category").Preload("Maintainer.Role")
	if filter.EquipmentID > 0 {
		query = query.Where("equipment_id = ?", filter.EquipmentID)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count maintenance records: %w", err)
	}
	var list []model.MaintenanceRecord
	if err := query.Order("id DESC").Offset(filter.Offset()).Limit(filter.PageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list maintenance records: %w", err)
	}
	return list, total, nil
}

func (r *maintenanceRepository) Stats(ctx context.Context) (*MaintenanceStats, error) {
	stats := &MaintenanceStats{ResultCount: map[constants.MaintenanceResult]int64{}}
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceRecord{}).Count(&stats.TotalRecords).Error; err != nil {
		return nil, fmt.Errorf("count maintenance records: %w", err)
	}
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceRecord{}).
		Select("COALESCE(SUM(cost), 0)").
		Scan(&stats.TotalCost).Error; err != nil {
		return nil, fmt.Errorf("sum maintenance cost: %w", err)
	}
	rows := []struct {
		Result constants.MaintenanceResult
		Count  int64
	}{}
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceRecord{}).
		Select("result, count(*) as count").
		Group("result").
		Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("count maintenance by result: %w", err)
	}
	for _, row := range rows {
		stats.ResultCount[row.Result] = row.Count
	}
	return stats, nil
}
