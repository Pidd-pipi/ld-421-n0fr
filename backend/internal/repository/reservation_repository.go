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

// ReservationRepository 预约记录仓储接口。
type ReservationRepository interface {
	Create(ctx context.Context, reservation *model.Reservation) error
	Update(ctx context.Context, reservation *model.Reservation) error
	FindByID(ctx context.Context, id uint) (*model.Reservation, error)
	List(ctx context.Context, filter ReservationFilter) ([]model.Reservation, int64, error)
	CountPending(ctx context.Context) (int64, error)
	HasConflict(ctx context.Context, equipmentID uint, start, end time.Time, excludeID uint) (bool, error)
}

// ReservationFilter 预约记录筛选条件。
type ReservationFilter struct {
	EquipmentID uint
	UserID      uint
	Status      constants.ReservationStatus
	Pagination
}

type reservationRepository struct {
	db *gorm.DB
}

// NewReservationRepository 构造预约记录仓储。
func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) Create(ctx context.Context, reservation *model.Reservation) error {
	if err := r.db.WithContext(ctx).Create(reservation).Error; err != nil {
		return fmt.Errorf("create reservation: %w", err)
	}
	return nil
}

func (r *reservationRepository) Update(ctx context.Context, reservation *model.Reservation) error {
	if err := r.db.WithContext(ctx).Save(reservation).Error; err != nil {
		return fmt.Errorf("update reservation: %w", err)
	}
	return nil
}

func (r *reservationRepository) FindByID(ctx context.Context, id uint) (*model.Reservation, error) {
	var reservation model.Reservation
	err := r.db.WithContext(ctx).Preload("Equipment.Category").Preload("User.Role").Preload("Approver.Role").First(&reservation, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find reservation by id: %w", err)
	}
	return &reservation, nil
}

func (r *reservationRepository) List(ctx context.Context, filter ReservationFilter) ([]model.Reservation, int64, error) {
	filter.Normalize()
	query := r.db.WithContext(ctx).Model(&model.Reservation{}).Preload("Equipment.Category").Preload("User.Role").Preload("Approver.Role")
	if filter.EquipmentID > 0 {
		query = query.Where("equipment_id = ?", filter.EquipmentID)
	}
	if filter.UserID > 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.Status.Valid() {
		query = query.Where("status = ?", filter.Status)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count reservations: %w", err)
	}
	var list []model.Reservation
	if err := query.Order("start_time ASC").Offset(filter.Offset()).Limit(filter.PageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list reservations: %w", err)
	}
	return list, total, nil
}

func (r *reservationRepository) CountPending(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Reservation{}).
		Where("status = ?", constants.ReservationStatusPending).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count pending reservations: %w", err)
	}
	return count, nil
}

func (r *reservationRepository) HasConflict(ctx context.Context, equipmentID uint, start, end time.Time, excludeID uint) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&model.Reservation{}).
		Where("equipment_id = ?", equipmentID).
		Where("status IN ?", []constants.ReservationStatus{constants.ReservationStatusPending, constants.ReservationStatusApproved}).
		Where("start_time < ? AND end_time > ?", end, start)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("check reservation conflict: %w", err)
	}
	return count > 0, nil
}
