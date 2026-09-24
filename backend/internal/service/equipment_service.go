package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/labequipment/lab-equipment/internal/constants"
	apperrors "github.com/labequipment/lab-equipment/internal/errors"
	"github.com/labequipment/lab-equipment/internal/model"
	"github.com/labequipment/lab-equipment/internal/repository"
)

// EquipmentService 设备业务服务。
type EquipmentService struct {
	repo         repository.EquipmentRepository
	categoryRepo repository.CategoryRepository
	userRepo     repository.UserRepository
	audit        *AuditService
	logger       *slog.Logger
}

// NewEquipmentService 构造设备服务。
func NewEquipmentService(
	repo repository.EquipmentRepository,
	categoryRepo repository.CategoryRepository,
	userRepo repository.UserRepository,
	audit *AuditService,
	logger *slog.Logger,
) *EquipmentService {
	return &EquipmentService{repo: repo, categoryRepo: categoryRepo, userRepo: userRepo, audit: audit, logger: logger}
}

// Create 登记入库新设备。
func (s *EquipmentService) Create(ctx context.Context, equipment *model.Equipment, actor Actor) (*model.Equipment, error) {
	if err := s.validateRelations(ctx, equipment); err != nil {
		return nil, err
	}
	if !equipment.Status.Valid() {
		equipment.Status = constants.AssetStatusAvailable
	}
	if err := s.repo.Create(ctx, equipment); err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return nil, apperrors.NewBusinessError(40900, 409, "设备编号已存在")
		}
		return nil, fmt.Errorf("create equipment: %w", err)
	}
	if err := s.audit.Log(ctx, actor, "equipment.register", "equipment", equipment.ID, fmt.Sprintf("登记设备 %s", equipment.Code)); err != nil {
		return nil, err
	}
	return equipment, nil
}

// Update 更新设备基础信息。
func (s *EquipmentService) Update(ctx context.Context, id uint, equipment *model.Equipment, actor Actor) (*model.Equipment, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, s.mapNotFound(err, "设备不存在")
	}
	if err := s.validateRelations(ctx, equipment); err != nil {
		return nil, err
	}
	existing.Name = equipment.Name
	existing.Code = equipment.Code
	existing.CategoryID = equipment.CategoryID
	existing.BrandModel = equipment.BrandModel
	existing.SerialNumber = equipment.SerialNumber
	existing.PurchaseDate = equipment.PurchaseDate
	existing.PurchasePrice = equipment.PurchasePrice
	existing.Location = equipment.Location
	existing.OwnerID = equipment.OwnerID
	existing.Supplier = equipment.Supplier
	existing.WarrantyExpiry = equipment.WarrantyExpiry
	existing.ImageURL = equipment.ImageURL
	if err := s.repo.Update(ctx, existing); err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return nil, apperrors.NewBusinessError(40900, 409, "设备编号已存在")
		}
		return nil, fmt.Errorf("update equipment: %w", err)
	}
	if err := s.audit.Log(ctx, actor, "equipment.update", "equipment", id, fmt.Sprintf("更新设备 %s", existing.Code)); err != nil {
		return nil, err
	}
	return existing, nil
}

// Get 获取设备详情。
func (s *EquipmentService) Get(ctx context.Context, id uint) (*model.Equipment, error) {
	equipment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, s.mapNotFound(err, "设备不存在")
	}
	return equipment, nil
}

// List 分页查询设备。
func (s *EquipmentService) List(ctx context.Context, filter repository.EquipmentFilter) ([]model.Equipment, int64, error) {
	list, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("list equipment: %w", err)
	}
	return list, total, nil
}

// Retire 报废设备。
func (s *EquipmentService) Retire(ctx context.Context, id uint, actor Actor) error {
	equipment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return s.mapNotFound(err, "设备不存在")
	}
	if err := s.repo.UpdateStatus(ctx, id, constants.AssetStatusRetired); err != nil {
		return fmt.Errorf("retire equipment: %w", err)
	}
	if err := s.audit.Log(ctx, actor, "equipment.retire", "equipment", id, fmt.Sprintf("报废设备 %s", equipment.Code)); err != nil {
		return err
	}
	return nil
}

// TransferOwner 转移责任人。
func (s *EquipmentService) TransferOwner(ctx context.Context, id, newOwnerID uint, actor Actor) error {
	equipment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return s.mapNotFound(err, "设备不存在")
	}
	if _, err := s.userRepo.FindByID(ctx, newOwnerID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return apperrors.NewBusinessError(40400, 404, "责任人不存在")
		}
		return fmt.Errorf("find owner: %w", err)
	}
	if err := s.repo.UpdateOwner(ctx, id, newOwnerID); err != nil {
		return fmt.Errorf("transfer equipment owner: %w", err)
	}
	if err := s.audit.Log(ctx, actor, "equipment.transfer_owner", "equipment", id, fmt.Sprintf("设备 %s 责任人变更为 %d", equipment.Code, newOwnerID)); err != nil {
		return err
	}
	return nil
}

func (s *EquipmentService) validateRelations(ctx context.Context, equipment *model.Equipment) error {
	if _, err := s.categoryRepo.FindByID(ctx, equipment.CategoryID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return apperrors.NewBusinessError(40400, 404, "设备分类不存在")
		}
		return fmt.Errorf("find category: %w", err)
	}
	if _, err := s.userRepo.FindByID(ctx, equipment.OwnerID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return apperrors.NewBusinessError(40400, 404, "责任人不存在")
		}
		return fmt.Errorf("find owner: %w", err)
	}
	return nil
}

func (s *EquipmentService) mapNotFound(err error, message string) error {
	if errors.Is(err, repository.ErrNotFound) {
		return apperrors.NewBusinessError(40400, 404, message)
	}
	return fmt.Errorf("find equipment: %w", err)
}
