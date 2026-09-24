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

// DisposalReviewResult 处置审批结果。
type DisposalReviewResult struct {
	Approved bool
	Blockers model.DisposalBlockers
}

// DisposalRepository 设备处置（报废）审批仓储接口。
type DisposalRepository interface {
	Create(ctx context.Context, approval *model.DisposalApproval) error
	FindByID(ctx context.Context, id uint) (*model.DisposalApproval, error)
	FindActiveByEquipment(ctx context.Context, equipmentID uint) (*model.DisposalApproval, error)
	FindLatestByEquipment(ctx context.Context, equipmentID uint) (*model.DisposalApproval, error)
	MapActiveByEquipmentIDs(ctx context.Context, equipmentIDs []uint) (map[uint]model.DisposalApproval, error)
	List(ctx context.Context, filter DisposalFilter) ([]model.DisposalApproval, int64, error)
	ListBlockers(ctx context.Context, equipmentID uint) (model.DisposalBlockers, error)
	// Review 原子审批：无阻塞项则通过并报废设备，存在阻塞项则退回并快照阻塞项。
	Review(ctx context.Context, id, reviewerID uint) (*model.DisposalApproval, *DisposalReviewResult, error)
}

// DisposalFilter 处置审批筛选条件。
type DisposalFilter struct {
	Status      constants.DisposalStatus
	EquipmentID uint
	Pagination
}

type disposalRepository struct {
	db *gorm.DB
}

// NewDisposalRepository 构造处置审批仓储。
func NewDisposalRepository(db *gorm.DB) DisposalRepository {
	return &disposalRepository{db: db}
}

func (r *disposalRepository) Create(ctx context.Context, approval *model.DisposalApproval) error {
	if err := r.db.WithContext(ctx).Create(approval).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrConflict
		}
		return fmt.Errorf("create disposal approval: %w", err)
	}
	return nil
}

func (r *disposalRepository) FindByID(ctx context.Context, id uint) (*model.DisposalApproval, error) {
	var approval model.DisposalApproval
	err := r.db.WithContext(ctx).
		Preload("Equipment.Category").
		Preload("Applicant.Role").
		Preload("Reviewer.Role").
		First(&approval, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find disposal approval by id: %w", err)
	}
	return &approval, nil
}

func (r *disposalRepository) FindActiveByEquipment(ctx context.Context, equipmentID uint) (*model.DisposalApproval, error) {
	var approval model.DisposalApproval
	err := r.db.WithContext(ctx).
		Preload("Equipment.Category").
		Preload("Applicant.Role").
		Preload("Reviewer.Role").
		Where("equipment_id = ? AND active = ?", equipmentID, true).
		First(&approval).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find active disposal approval: %w", err)
	}
	return &approval, nil
}

func (r *disposalRepository) FindLatestByEquipment(ctx context.Context, equipmentID uint) (*model.DisposalApproval, error) {
	var approval model.DisposalApproval
	err := r.db.WithContext(ctx).
		Preload("Equipment.Category").
		Preload("Applicant.Role").
		Preload("Reviewer.Role").
		Where("equipment_id = ?", equipmentID).
		Order("id DESC").
		First(&approval).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find latest disposal approval: %w", err)
	}
	return &approval, nil
}

func (r *disposalRepository) MapActiveByEquipmentIDs(ctx context.Context, equipmentIDs []uint) (map[uint]model.DisposalApproval, error) {
	result := make(map[uint]model.DisposalApproval)
	if len(equipmentIDs) == 0 {
		return result, nil
	}
	var approvals []model.DisposalApproval
	if err := r.db.WithContext(ctx).
		Where("equipment_id IN ? AND active = ?", equipmentIDs, true).
		Find(&approvals).Error; err != nil {
		return nil, fmt.Errorf("map active disposal approvals: %w", err)
	}
	for _, approval := range approvals {
		result[approval.EquipmentID] = approval
	}
	return result, nil
}

func (r *disposalRepository) List(ctx context.Context, filter DisposalFilter) ([]model.DisposalApproval, int64, error) {
	filter.Normalize()
	query := r.db.WithContext(ctx).Model(&model.DisposalApproval{}).
		Preload("Equipment.Category").
		Preload("Applicant.Role").
		Preload("Reviewer.Role")
	if filter.Status.Valid() {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.EquipmentID > 0 {
		query = query.Where("equipment_id = ?", filter.EquipmentID)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count disposal approvals: %w", err)
	}
	var list []model.DisposalApproval
	if err := query.Order("id DESC").Offset(filter.Offset()).Limit(filter.PageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list disposal approvals: %w", err)
	}
	return list, total, nil
}

func (r *disposalRepository) ListBlockers(ctx context.Context, equipmentID uint) (model.DisposalBlockers, error) {
	return r.queryBlockers(r.db.WithContext(ctx), equipmentID)
}

// Review 在单个事务内完成阻塞项检查、审批单状态流转与设备报废，
// 通过条件更新保证并发审批只有一个请求成功。
func (r *disposalRepository) Review(ctx context.Context, id, reviewerID uint) (*model.DisposalApproval, *DisposalReviewResult, error) {
	var approval *model.DisposalApproval
	result := &DisposalReviewResult{}
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var current model.DisposalApproval
		if err := tx.First(&current, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return fmt.Errorf("find disposal approval: %w", err)
		}
		if current.Status != constants.DisposalStatusPending || !current.Active {
			return ErrConflict
		}

		blockers, err := r.queryBlockers(tx, current.EquipmentID)
		if err != nil {
			return err
		}

		now := time.Now()
		updates := map[string]any{
			"reviewer_id":         reviewerID,
			"reviewed_at":         now,
			"active":              false,
			"active_equipment_id": nil,
			"blockers":            blockers,
		}

		if len(blockers) > 0 {
			result.Approved = false
			updates["status"] = constants.DisposalStatusReturned
		} else {
			result.Approved = true
			updates["status"] = constants.DisposalStatusApproved
		}
		result.Blockers = blockers

		query := tx.Model(&model.DisposalApproval{}).
			Where("id = ? AND status = ? AND active = ?", id, constants.DisposalStatusPending, true).
			Updates(updates)
		if query.Error != nil {
			return fmt.Errorf("review disposal approval: %w", query.Error)
		}
		if query.RowsAffected == 0 {
			return ErrConflict
		}

		if result.Approved {
			if err := tx.Model(&model.Equipment{}).
				Where("id = ?", current.EquipmentID).
				Update("status", constants.AssetStatusRetired).Error; err != nil {
				return fmt.Errorf("retire equipment: %w", err)
			}
		}

		loaded, err := r.findByIDWithTx(tx, id)
		if err != nil {
			return err
		}
		approval = loaded
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return approval, result, nil
}

func (r *disposalRepository) findByIDWithTx(tx *gorm.DB, id uint) (*model.DisposalApproval, error) {
	var approval model.DisposalApproval
	err := tx.Preload("Equipment.Category").
		Preload("Applicant.Role").
		Preload("Reviewer.Role").
		First(&approval, id).Error
	if err != nil {
		return nil, fmt.Errorf("reload disposal approval: %w", err)
	}
	return &approval, nil
}

// queryBlockers 汇总设备当前未归还的借用与未结束的预约。
func (r *disposalRepository) queryBlockers(tx *gorm.DB, equipmentID uint) (model.DisposalBlockers, error) {
	blockers := model.DisposalBlockers{}

	var borrows []model.BorrowRecord
	if err := tx.Preload("Borrower").
		Where("equipment_id = ?", equipmentID).
		Where("status IN ?", []constants.BorrowStatus{
			constants.BorrowStatusPending,
			constants.BorrowStatusApproved,
			constants.BorrowStatusOverdue,
		}).
		Find(&borrows).Error; err != nil {
		return nil, fmt.Errorf("list active borrow blockers: %w", err)
	}
	for _, borrow := range borrows {
		blocker := model.DisposalBlocker{
			Type:         "Borrow",
			RecordID:     borrow.ID,
			UserID:       borrow.BorrowerID,
			Status:       string(borrow.Status),
			Detail:       fmt.Sprintf("借用人 %s 尚未归还该设备", borrowerName(borrow)),
			ExpectedTime: borrow.ExpectedReturnDate.Format("2006-01-02 15:04"),
		}
		if borrow.Borrower != nil {
			blocker.UserName = borrow.Borrower.Name
		}
		blockers = append(blockers, blocker)
	}

	var reservations []model.Reservation
	if err := tx.Preload("User").
		Where("equipment_id = ?", equipmentID).
		Where("status IN ?", []constants.ReservationStatus{
			constants.ReservationStatusPending,
			constants.ReservationStatusApproved,
		}).
		Find(&reservations).Error; err != nil {
		return nil, fmt.Errorf("list active reservation blockers: %w", err)
	}
	for _, reservation := range reservations {
		blocker := model.DisposalBlocker{
			Type:         "Reservation",
			RecordID:     reservation.ID,
			UserID:       reservation.UserID,
			Status:       string(reservation.Status),
			Detail: fmt.Sprintf(
				"使用人 %s 已预约 %s 至 %s",
				reservationUserName(reservation),
				reservation.StartTime.Format("2006-01-02 15:04"),
				reservation.EndTime.Format("2006-01-02 15:04"),
			),
			ExpectedTime: reservation.StartTime.Format("2006-01-02 15:04"),
		}
		if reservation.User != nil {
			blocker.UserName = reservation.User.Name
		}
		blockers = append(blockers, blocker)
	}

	return blockers, nil
}

func borrowerName(borrow model.BorrowRecord) string {
	if borrow.Borrower != nil {
		return borrow.Borrower.Name
	}
	return fmt.Sprintf("用户#%d", borrow.BorrowerID)
}

func reservationUserName(reservation model.Reservation) string {
	if reservation.User != nil {
		return reservation.User.Name
	}
	return fmt.Sprintf("用户#%d", reservation.UserID)
}
