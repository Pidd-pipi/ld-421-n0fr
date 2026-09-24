package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/labequipment/lab-equipment/internal/constants"
	"github.com/labequipment/lab-equipment/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DisposalRepository 报废处置审批仓储接口。
type DisposalRepository interface {
	CreateIfNoPending(ctx context.Context, request *model.DisposalRequest) error
	FindByID(ctx context.Context, id uint) (*model.DisposalRequest, error)
	ListByEquipment(ctx context.Context, equipmentID uint, limit int) ([]model.DisposalRequest, error)
	HasPendingForEquipment(ctx context.Context, equipmentID uint) (bool, error)
	ApproveAndRetire(ctx context.Context, id, equipmentID, approverID uint, processedAt time.Time) (bool, error)
	RejectIfPending(ctx context.Context, id, approverID uint, blockers string, processedAt time.Time) (bool, error)
}

type disposalRepository struct {
	db *gorm.DB
}

// NewDisposalRepository 构造报废处置审批仓储。
func NewDisposalRepository(db *gorm.DB) DisposalRepository {
	return &disposalRepository{db: db}
}

// lockForUpdate 仅在支持行锁的方言上追加 FOR UPDATE，SQLite（测试）下自动跳过。
func lockForUpdate(tx *gorm.DB) *gorm.DB {
	if tx.Dialector.Name() == "mysql" {
		return tx.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	return tx
}

// CreateIfNoPending 在事务中提交报废申请：锁定设备行后校验，并发重复提交时后到的返回 ErrConflict。
func (r *disposalRepository) CreateIfNoPending(ctx context.Context, request *model.DisposalRequest) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var equipment model.Equipment
		err := lockForUpdate(tx).First(&equipment, request.EquipmentID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		if err != nil {
			return fmt.Errorf("lock equipment: %w", err)
		}
		if equipment.Status == constants.AssetStatusRetired {
			return ErrConflict
		}
		var count int64
		if err := tx.Model(&model.DisposalRequest{}).
			Where("equipment_id = ? AND status = ?", request.EquipmentID, constants.DisposalStatusPending).
			Count(&count).Error; err != nil {
			return fmt.Errorf("count pending disposal requests: %w", err)
		}
		if count > 0 {
			return ErrConflict
		}
		if err := tx.Create(request).Error; err != nil {
			return fmt.Errorf("create disposal request: %w", err)
		}
		return nil
	})
}

func (r *disposalRepository) FindByID(ctx context.Context, id uint) (*model.DisposalRequest, error) {
	var request model.DisposalRequest
	err := r.db.WithContext(ctx).
		Preload("Equipment").
		Preload("Applicant").
		Preload("Approver").
		First(&request, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find disposal request by id: %w", err)
	}
	return &request, nil
}

func (r *disposalRepository) ListByEquipment(ctx context.Context, equipmentID uint, limit int) ([]model.DisposalRequest, error) {
	if limit <= 0 {
		limit = 20
	}
	var list []model.DisposalRequest
	err := r.db.WithContext(ctx).
		Preload("Applicant").
		Preload("Approver").
		Where("equipment_id = ?", equipmentID).
		Order("id DESC").
		Limit(limit).
		Find(&list).Error
	if err != nil {
		return nil, fmt.Errorf("list disposal requests: %w", err)
	}
	return list, nil
}

func (r *disposalRepository) HasPendingForEquipment(ctx context.Context, equipmentID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.DisposalRequest{}).
		Where("equipment_id = ? AND status = ?", equipmentID, constants.DisposalStatusPending).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("count pending disposal requests: %w", err)
	}
	return count > 0, nil
}

// ApproveAndRetire 条件更新：仅当申请仍处于 Pending 时置为 Approved，并在同一事务内将设备转为 Retired。
// 返回 false 表示申请已被其他人处理（并发审批冲突）。
func (r *disposalRepository) ApproveAndRetire(ctx context.Context, id, equipmentID, approverID uint, processedAt time.Time) (bool, error) {
	approved := false
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&model.DisposalRequest{}).
			Where("id = ? AND status = ?", id, constants.DisposalStatusPending).
			Updates(map[string]any{
				"status":       constants.DisposalStatusApproved,
				"approver_id":  approverID,
				"processed_at": processedAt,
			})
		if result.Error != nil {
			return fmt.Errorf("approve disposal request: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return nil
		}
		if err := tx.Model(&model.Equipment{}).
			Where("id = ?", equipmentID).
			Update("status", constants.AssetStatusRetired).Error; err != nil {
			return fmt.Errorf("retire equipment: %w", err)
		}
		approved = true
		return nil
	})
	if err != nil {
		return false, err
	}
	return approved, nil
}

// RejectIfPending 条件更新：仅当申请仍处于 Pending 时置为 Rejected（退回）并记录阻塞项快照。
// 返回 false 表示申请已被其他人处理（并发审批冲突）。
func (r *disposalRepository) RejectIfPending(ctx context.Context, id, approverID uint, blockers string, processedAt time.Time) (bool, error) {
	result := r.db.WithContext(ctx).Model(&model.DisposalRequest{}).
		Where("id = ? AND status = ?", id, constants.DisposalStatusPending).
		Updates(map[string]any{
			"status":       constants.DisposalStatusRejected,
			"approver_id":  approverID,
			"blockers":     blockers,
			"processed_at": processedAt,
		})
	if result.Error != nil {
		return false, fmt.Errorf("reject disposal request: %w", result.Error)
	}
	return result.RowsAffected > 0, nil
}
