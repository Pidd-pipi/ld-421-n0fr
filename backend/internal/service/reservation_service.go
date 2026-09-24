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

// ReservationService 预约记录业务服务。
type ReservationService struct {
	repo          repository.ReservationRepository
	equipmentRepo repository.EquipmentRepository
	audit         *AuditService
	logger        *slog.Logger
}

// NewReservationService 构造预约服务。
func NewReservationService(
	repo repository.ReservationRepository,
	equipmentRepo repository.EquipmentRepository,
	audit *AuditService,
	logger *slog.Logger,
) *ReservationService {
	return &ReservationService{repo: repo, equipmentRepo: equipmentRepo, audit: audit, logger: logger}
}

// Create 创建预约。
func (s *ReservationService) Create(ctx context.Context, reservation *model.Reservation, actor Actor) (*model.Reservation, error) {
	if !reservation.EndTime.After(reservation.StartTime) {
		return nil, apperrors.NewBusinessError(40000, 400, "结束时间必须晚于开始时间")
	}
	if _, err := s.equipmentRepo.FindByID(ctx, reservation.EquipmentID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperrors.NewBusinessError(40400, 404, "设备不存在")
		}
		return nil, fmt.Errorf("find equipment: %w", err)
	}
	conflict, err := s.repo.HasConflict(ctx, reservation.EquipmentID, reservation.StartTime, reservation.EndTime, 0)
	if err != nil {
		return nil, fmt.Errorf("check conflict: %w", err)
	}
	if conflict {
		return nil, apperrors.NewBusinessError(40900, 409, "该时间段已被预约")
	}
	reservation.UserID = actor.UserID
	reservation.Status = constants.ReservationStatusPending
	if err := s.repo.Create(ctx, reservation); err != nil {
		return nil, fmt.Errorf("create reservation: %w", err)
	}
	if err := s.audit.Log(ctx, actor, "reservation.create", "reservation", reservation.ID, fmt.Sprintf("创建预约 %d", reservation.ID)); err != nil {
		return nil, err
	}
	return reservation, nil
}

// Approve 审批通过预约。
func (s *ReservationService) Approve(ctx context.Context, id uint, actor Actor) error {
	reservation, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return s.mapNotFound(err)
	}
	if reservation.Status != constants.ReservationStatusPending {
		return apperrors.NewBusinessError(40900, 409, "仅待审批预约可审批")
	}
	approverID := actor.UserID
	reservation.Status = constants.ReservationStatusApproved
	reservation.ApproverID = &approverID
	if err := s.repo.Update(ctx, reservation); err != nil {
		return fmt.Errorf("approve reservation: %w", err)
	}
	if err := s.audit.Log(ctx, actor, "reservation.approve", "reservation", id, fmt.Sprintf("审批通过预约 %d", id)); err != nil {
		return err
	}
	return nil
}

// Reject 驳回预约。
func (s *ReservationService) Reject(ctx context.Context, id uint, actor Actor) error {
	reservation, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return s.mapNotFound(err)
	}
	if reservation.Status != constants.ReservationStatusPending {
		return apperrors.NewBusinessError(40900, 409, "仅待审批预约可驳回")
	}
	approverID := actor.UserID
	reservation.Status = constants.ReservationStatusRejected
	reservation.ApproverID = &approverID
	if err := s.repo.Update(ctx, reservation); err != nil {
		return fmt.Errorf("reject reservation: %w", err)
	}
	if err := s.audit.Log(ctx, actor, "reservation.reject", "reservation", id, fmt.Sprintf("驳回预约 %d", id)); err != nil {
		return err
	}
	return nil
}

// Cancel 取消预约。
func (s *ReservationService) Cancel(ctx context.Context, id uint, actor Actor) error {
	reservation, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return s.mapNotFound(err)
	}
	if reservation.Status == constants.ReservationStatusCancelled {
		return apperrors.NewBusinessError(40900, 409, "预约已取消")
	}
	reservation.Status = constants.ReservationStatusCancelled
	if err := s.repo.Update(ctx, reservation); err != nil {
		return fmt.Errorf("cancel reservation: %w", err)
	}
	if err := s.audit.Log(ctx, actor, "reservation.cancel", "reservation", id, fmt.Sprintf("取消预约 %d", id)); err != nil {
		return err
	}
	return nil
}

// Get 获取预约详情。
func (s *ReservationService) Get(ctx context.Context, id uint) (*model.Reservation, error) {
	reservation, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, s.mapNotFound(err)
	}
	return reservation, nil
}

// List 分页查询预约。
func (s *ReservationService) List(ctx context.Context, filter repository.ReservationFilter) ([]model.Reservation, int64, error) {
	list, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("list reservations: %w", err)
	}
	return list, total, nil
}

func (s *ReservationService) mapNotFound(err error) error {
	if errors.Is(err, repository.ErrNotFound) {
		return apperrors.NewBusinessError(40400, 404, "预约记录不存在")
	}
	return fmt.Errorf("find reservation: %w", err)
}
