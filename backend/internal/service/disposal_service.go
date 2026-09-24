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

// DisposalService 设备处置（报废）审批业务服务。
type DisposalService struct {
	repo          repository.DisposalRepository
	equipmentRepo repository.EquipmentRepository
	audit         *AuditService
	logger        *slog.Logger
}

// NewDisposalService 构造处置审批服务。
func NewDisposalService(
	repo repository.DisposalRepository,
	equipmentRepo repository.EquipmentRepository,
	audit *AuditService,
	logger *slog.Logger,
) *DisposalService {
	return &DisposalService{repo: repo, equipmentRepo: equipmentRepo, audit: audit, logger: logger}
}

// Submit 提交设备报废处置申请。提交后设备进入审批中，暂停新的借用和预约。
func (s *DisposalService) Submit(ctx context.Context, approval *model.DisposalApproval, actor Actor) (*model.DisposalApproval, error) {
	equipment, err := s.equipmentRepo.FindByID(ctx, approval.EquipmentID)
	if err != nil {
		return nil, s.mapEquipmentNotFound(err)
	}
	switch equipment.Status {
	case constants.AssetStatusRetired:
		return nil, apperrors.NewBusinessError(40900, 409, "设备已报废，无需重复提交")
	case constants.AssetStatusLost:
		return nil, apperrors.NewBusinessError(40900, 409, "已报失设备不能提交报废处置")
	}
	if _, err := s.repo.FindActiveByEquipment(ctx, approval.EquipmentID); err == nil {
		return nil, apperrors.NewBusinessError(40900, 409, "该设备已有进行中的处置审批，请勿重复提交")
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("find active disposal approval: %w", err)
	}

	approval.ApplicantID = actor.UserID
	approval.Status = constants.DisposalStatusPending
	approval.Active = true
	activeEquipmentID := approval.EquipmentID
	approval.ActiveEquipmentID = &activeEquipmentID
	approval.Blockers = model.DisposalBlockers{}
	if err := s.repo.Create(ctx, approval); err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return nil, apperrors.NewBusinessError(40900, 409, "该设备已有进行中的处置审批，请勿重复提交")
		}
		return nil, fmt.Errorf("submit disposal approval: %w", err)
	}
	if err := s.audit.Log(ctx, actor, "equipment.disposal.submit", "disposal_approval", approval.ID,
		fmt.Sprintf("提交设备 %s 报废处置审批", equipment.Code)); err != nil {
		return nil, err
	}
	return approval, nil
}

// Review 处置审批：无阻塞项则通过并报废设备；存在未归还借用或未结束预约则退回。
func (s *DisposalService) Review(ctx context.Context, id uint, actor Actor) (*model.DisposalApproval, *repository.DisposalReviewResult, error) {
	approval, result, err := s.repo.Review(ctx, id, actor.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, nil, apperrors.NewBusinessError(40400, 404, "处置审批单不存在")
		}
		if errors.Is(err, repository.ErrConflict) {
			return nil, nil, apperrors.NewBusinessError(40900, 409, "审批单已被处理，请勿重复审批")
		}
		return nil, nil, fmt.Errorf("review disposal approval: %w", err)
	}
	action := "equipment.disposal.approve"
	detail := fmt.Sprintf("审批通过设备报废处置单 %d，设备已退役", id)
	if !result.Approved {
		action = "equipment.disposal.return"
		detail = fmt.Sprintf("处置审批单 %d 存在 %d 项阻塞，已退回", id, len(result.Blockers))
	}
	if err := s.audit.Log(ctx, actor, action, "disposal_approval", id, detail); err != nil {
		return nil, nil, err
	}
	return approval, result, nil
}

// Get 获取处置审批详情。
func (s *DisposalService) Get(ctx context.Context, id uint) (*model.DisposalApproval, error) {
	approval, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, s.mapNotFound(err)
	}
	return approval, nil
}

// GetByEquipment 获取设备当前审批（优先进行中，否则最近一次）；没有则返回 nil。
func (s *DisposalService) GetByEquipment(ctx context.Context, equipmentID uint) (*model.DisposalApproval, error) {
	if approval, err := s.repo.FindActiveByEquipment(ctx, equipmentID); err == nil {
		return approval, nil
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("find active disposal approval: %w", err)
	}
	approval, err := s.repo.FindLatestByEquipment(ctx, equipmentID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("find latest disposal approval: %w", err)
	}
	return approval, nil
}

// LiveBlockers 获取设备实时阻塞项，供审批中页面展示。
func (s *DisposalService) LiveBlockers(ctx context.Context, equipmentID uint) (model.DisposalBlockers, error) {
	blockers, err := s.repo.ListBlockers(ctx, equipmentID)
	if err != nil {
		return nil, fmt.Errorf("list disposal blockers: %w", err)
	}
	return blockers, nil
}

// List 分页查询处置审批。
func (s *DisposalService) List(ctx context.Context, filter repository.DisposalFilter) ([]model.DisposalApproval, int64, error) {
	list, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("list disposal approvals: %w", err)
	}
	return list, total, nil
}

func (s *DisposalService) mapNotFound(err error) error {
	if errors.Is(err, repository.ErrNotFound) {
		return apperrors.NewBusinessError(40400, 404, "处置审批单不存在")
	}
	return fmt.Errorf("find disposal approval: %w", err)
}

func (s *DisposalService) mapEquipmentNotFound(err error) error {
	if errors.Is(err, repository.ErrNotFound) {
		return apperrors.NewBusinessError(40400, 404, "设备不存在")
	}
	return fmt.Errorf("find equipment: %w", err)
}
