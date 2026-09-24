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

// BorrowRepository 借用记录仓储接口。
type BorrowRepository interface {
	Create(ctx context.Context, record *model.BorrowRecord) error
	Update(ctx context.Context, record *model.BorrowRecord) error
	FindByID(ctx context.Context, id uint) (*model.BorrowRecord, error)
	List(ctx context.Context, filter BorrowFilter) ([]model.BorrowRecord, int64, error)
	CountPending(ctx context.Context) (int64, error)
	TopEquipmentThisMonth(ctx context.Context, limit int) ([]BorrowTopStat, error)
}

// BorrowFilter 借用记录筛选条件。
type BorrowFilter struct {
	Status      constants.BorrowStatus
	EquipmentID uint
	BorrowerID  uint
	Pagination
}

// BorrowTopStat 本月借用次数统计结果。
type BorrowTopStat struct {
	EquipmentID uint
	Name        string
	Code        string
	Count       int64
}

type borrowRepository struct {
	db *gorm.DB
}

// NewBorrowRepository 构造借用记录仓储。
func NewBorrowRepository(db *gorm.DB) BorrowRepository {
	return &borrowRepository{db: db}
}

func (r *borrowRepository) Create(ctx context.Context, record *model.BorrowRecord) error {
	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return fmt.Errorf("create borrow record: %w", err)
	}
	return nil
}

func (r *borrowRepository) Update(ctx context.Context, record *model.BorrowRecord) error {
	if err := r.db.WithContext(ctx).Save(record).Error; err != nil {
		return fmt.Errorf("update borrow record: %w", err)
	}
	return nil
}

func (r *borrowRepository) FindByID(ctx context.Context, id uint) (*model.BorrowRecord, error) {
	var record model.BorrowRecord
	err := r.db.WithContext(ctx).Preload("Equipment.Category").Preload("Borrower.Role").Preload("Approver.Role").First(&record, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find borrow record by id: %w", err)
	}
	return &record, nil
}

func (r *borrowRepository) List(ctx context.Context, filter BorrowFilter) ([]model.BorrowRecord, int64, error) {
	filter.Normalize()
	query := r.db.WithContext(ctx).Model(&model.BorrowRecord{}).Preload("Equipment.Category").Preload("Borrower.Role").Preload("Approver.Role")
	if filter.Status.Valid() {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.EquipmentID > 0 {
		query = query.Where("equipment_id = ?", filter.EquipmentID)
	}
	if filter.BorrowerID > 0 {
		query = query.Where("borrower_id = ?", filter.BorrowerID)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count borrow records: %w", err)
	}
	var list []model.BorrowRecord
	if err := query.Order("id DESC").Offset(filter.Offset()).Limit(filter.PageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list borrow records: %w", err)
	}
	return list, total, nil
}

func (r *borrowRepository) CountPending(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.BorrowRecord{}).
		Where("status = ?", constants.BorrowStatusPending).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count pending borrow records: %w", err)
	}
	return count, nil
}

func (r *borrowRepository) TopEquipmentThisMonth(ctx context.Context, limit int) ([]BorrowTopStat, error) {
	if limit <= 0 {
		limit = 10
	}
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	rows := []BorrowTopStat{}
	err := r.db.WithContext(ctx).Model(&model.BorrowRecord{}).
		Select("borrow_records.equipment_id, equipment.name, equipment.code, count(*) as count").
		Joins("LEFT JOIN equipment ON equipment.id = borrow_records.equipment_id").
		Where("borrow_records.borrow_date >= ?", monthStart).
		Group("borrow_records.equipment_id, equipment.name, equipment.code").
		Order("count DESC").
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("top equipment this month: %w", err)
	}
	return rows, nil
}
