package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/labequipment/lab-equipment/internal/constants"
	apperrors "github.com/labequipment/lab-equipment/internal/errors"
	"github.com/labequipment/lab-equipment/internal/model"
	"github.com/labequipment/lab-equipment/internal/repository"
)

// DisposalItem 报废申请单及其阻塞项，供审批结果与页面展示使用。
type DisposalItem struct {
	Request  *model.DisposalRequest
	Blockers []model.DisposalBlocker
}

// DisposalService 设备报废处置审批业务服务。
type DisposalService struct {
	repo            repository.DisposalRepository
	equipmentRepo   repository.EquipmentRepository
	borrowRepo      repository.BorrowRepository
	reservationRepo repository.ReservationRepository
	audit           *AuditService
	logger          *slog.Logger
}

// NewDisposalService 构造报废处置审批服务。
func NewDisposalService(
	repo repository.DisposalRepository,
	equipmentRepo repository.EquipmentRepository,
	borrowRepo repository.BorrowRepository,
	reservationRepo repository.ReservationRepository,
	audit *AuditService,
	logger *slog.Logger,
) *DisposalService {
	return &DisposalService{
		repo:            repo,
		equipmentRepo:   equipmentRepo,
		borrowRepo:      borrowRepo,
		reservationRepo: reservationRepo,
		audit:           audit,
		logger:          logger,
	}
}

// Submit 提交报废申请。提交后设备暂停新的借用与预约，但状态不变，待审批通过才转为 Retired。
func (s *DisposalService) Submit(ctx context.Context, equipmentID uint, reason string, actor Actor) (*DisposalItem, error) {
	equipment, err := s.equipmentRepo.FindByID(ctx, equipmentID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperrors.NewBusinessError(40400, 404, "设备不存在")
		}
		return nil, fmt.Errorf("find equipment: %w", err)
	}
	if equipment.Status == constants.AssetStatusRetired {
		return nil, apperrors.NewBusinessError(40900, 409, "设备已报废，请勿重复提交")
	}
	pending, err := s.repo.HasPendingForEquipment(ctx, equipmentID)
	if err != nil {
		return nil, fmt.Errorf("check pending disposal: %w", err)
	}
	if pending {
		return nil, apperrors.NewBusinessError(40900, 409, "已存在待审批的报废申请")
	}
	request := &model.DisposalRequest{
		EquipmentID: equipmentID,
		Reason:      reason,
		Status:      constants.DisposalStatusPending,
		ApplicantID: actor.UserID,
	}
	if err := s.repo.CreateIfNoPending(ctx, request); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperrors.NewBusinessError(40400, 404, "设备不存在")
		}
		if errors.Is(err, repository.ErrConflict) {
			return nil, apperrors.NewBusinessError(40900, 409, "已存在待审批的报废申请或设备已报废")
		}
		return nil, fmt.Errorf("submit disposal request: %w", err)
	}
	if err := s.audit.Log(ctx, actor, "disposal.submit", "disposal", request.ID, fmt.Sprintf("提交设备 %s 报废申请", equipment.Code)); err != nil {
		return nil, err
	}
	return &DisposalItem{Request: request}, nil
}

// Approve 审批报废申请。存在未归还借用或未结束预约时列出阻塞项并退回；
// 无阻塞项时审批通过并将设备转为 Retired，原借用与预约记录保留。
func (s *DisposalService) Approve(ctx context.Context, equipmentID, disposalID uint, actor Actor) (*DisposalItem, error) {
	request, err := s.repo.FindByID(ctx, disposalID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperrors.NewBusinessError(40400, 404, "报废申请不存在")
		}
		return nil, fmt.Errorf("find disposal request: %w", err)
	}
	if request.EquipmentID != equipmentID {
		return nil, apperrors.NewBusinessError(40400, 404, "报废申请不存在")
	}
	if request.Status != constants.DisposalStatusPending {
		return nil, apperrors.NewBusinessError(40900, 409, "该申请已处理，请勿重复审批")
	}
	blockers, err := s.collectBlockers(ctx, request.EquipmentID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	if len(blockers) > 0 {
		snapshot, err := json.Marshal(blockers)
		if err != nil {
			return nil, fmt.Errorf("marshal blockers: %w", err)
		}
		updated, err := s.repo.RejectIfPending(ctx, disposalID, actor.UserID, string(snapshot), now)
		if err != nil {
			return nil, fmt.Errorf("return disposal request: %w", err)
		}
		if !updated {
			return nil, apperrors.NewBusinessError(40900, 409, "该申请已被处理，请勿重复审批")
		}
		if err := s.audit.Log(ctx, actor, "disposal.return", "disposal", disposalID, fmt.Sprintf("退回报废申请 %d，阻塞项 %d 项", disposalID, len(blockers))); err != nil {
			return nil, err
		}
		request.Status = constants.DisposalStatusRejected
		request.ApproverID = &actor.UserID
		request.Blockers = string(snapshot)
		request.ProcessedAt = &now
		return &DisposalItem{Request: request, Blockers: blockers}, nil
	}
	updated, err := s.repo.ApproveAndRetire(ctx, disposalID, request.EquipmentID, actor.UserID, now)
	if err != nil {
		return nil, fmt.Errorf("approve disposal request: %w", err)
	}
	if !updated {
		return nil, apperrors.NewBusinessError(40900, 409, "该申请已被处理，请勿重复审批")
	}
	if err := s.audit.Log(ctx, actor, "disposal.approve", "disposal", disposalID, fmt.Sprintf("审批通过报废申请 %d，设备转为报废", disposalID)); err != nil {
		return nil, err
	}
	request.Status = constants.DisposalStatusApproved
	request.ApproverID = &actor.UserID
	request.ProcessedAt = &now
	return &DisposalItem{Request: request}, nil
}

// ListByEquipment 查询设备的报废申请（最新在前），附带阻塞项：待审批单取实时阻塞项，退回单取退回时快照。
func (s *DisposalService) ListByEquipment(ctx context.Context, equipmentID uint) ([]DisposalItem, error) {
	if _, err := s.equipmentRepo.FindByID(ctx, equipmentID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperrors.NewBusinessError(40400, 404, "设备不存在")
		}
		return nil, fmt.Errorf("find equipment: %w", err)
	}
	requests, err := s.repo.ListByEquipment(ctx, equipmentID, 20)
	if err != nil {
		return nil, fmt.Errorf("list disposal requests: %w", err)
	}
	items := make([]DisposalItem, 0, len(requests))
	for i := range requests {
		request := requests[i]
		item := DisposalItem{Request: &request}
		switch request.Status {
		case constants.DisposalStatusPending:
			blockers, err := s.collectBlockers(ctx, request.EquipmentID)
			if err != nil {
				return nil, err
			}
			item.Blockers = blockers
		case constants.DisposalStatusRejected:
			blockers, err := decodeBlockers(request.Blockers)
			if err != nil {
				return nil, err
			}
			item.Blockers = blockers
		}
		items = append(items, item)
	}
	return items, nil
}

// collectBlockers 汇总设备当前未归还借用与未结束预约。
func (s *DisposalService) collectBlockers(ctx context.Context, equipmentID uint) ([]model.DisposalBlocker, error) {
	borrows, err := s.borrowRepo.ListUnreturnedByEquipment(ctx, equipmentID)
	if err != nil {
		return nil, fmt.Errorf("list unreturned borrows: %w", err)
	}
	reservations, err := s.reservationRepo.ListUnfinishedByEquipment(ctx, equipmentID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("list unfinished reservations: %w", err)
	}
	blockers := make([]model.DisposalBlocker, 0, len(borrows)+len(reservations))
	for _, record := range borrows {
		userName := ""
		if record.Borrower != nil {
			userName = record.Borrower.Name
		}
		blockers = append(blockers, model.DisposalBlocker{
			Type:     "borrow",
			RecordID: record.ID,
			Status:   string(record.Status),
			UserName: userName,
			Detail:   fmt.Sprintf("借用单 #%d，预计归还 %s", record.ID, record.ExpectedReturnDate.Format("2006-01-02")),
		})
	}
	for _, reservation := range reservations {
		userName := ""
		if reservation.User != nil {
			userName = reservation.User.Name
		}
		blockers = append(blockers, model.DisposalBlocker{
			Type:     "reservation",
			RecordID: reservation.ID,
			Status:   string(reservation.Status),
			UserName: userName,
			Detail: fmt.Sprintf("预约单 #%d，%s 至 %s", reservation.ID,
				reservation.StartTime.Format("2006-01-02 15:04"),
				reservation.EndTime.Format("2006-01-02 15:04")),
		})
	}
	return blockers, nil
}

func decodeBlockers(snapshot string) ([]model.DisposalBlocker, error) {
	if snapshot == "" {
		return nil, nil
	}
	var blockers []model.DisposalBlocker
	if err := json.Unmarshal([]byte(snapshot), &blockers); err != nil {
		return nil, fmt.Errorf("decode blockers snapshot: %w", err)
	}
	return blockers, nil
}
