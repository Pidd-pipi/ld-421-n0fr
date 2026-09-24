package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/labequipment/lab-equipment/internal/constants"
	apperrors "github.com/labequipment/lab-equipment/internal/errors"
	"github.com/labequipment/lab-equipment/internal/model"
)

// submitDisposalForEquipment 提交并返回一个待审批的处置申请。
func submitDisposalForEquipment(t *testing.T, env *testEnv, equipmentID uint, actor Actor) *model.DisposalApproval {
	t.Helper()
	approval, err := env.disposalService.Submit(context.Background(), &model.DisposalApproval{
		EquipmentID: equipmentID,
		Reason:      "设备无法校准，申请报废",
	}, actor)
	if err != nil {
		t.Fatalf("submit disposal: %v", err)
	}
	return approval
}

func TestDisposalService_SubmitBlocksNewBorrowAndReservation(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	actor := Actor{UserID: env.ownerID, Username: "admin", Role: "LabManager"}
	equipment, err := env.equipmentRepo.FindByID(ctx, 1)
	if err != nil {
		t.Fatalf("find equipment: %v", err)
	}

	approval := submitDisposalForEquipment(t, env, equipment.ID, actor)
	if approval.Status != constants.DisposalStatusPending {
		t.Fatalf("expected pending, got %s", approval.Status)
	}

	// 审批中设备状态不变，不应立即停用。
	equipment, _ = env.equipmentRepo.FindByID(ctx, equipment.ID)
	if equipment.Status == constants.AssetStatusRetired {
		t.Fatalf("equipment must not be retired while approval pending")
	}

	_, err = env.borrowService.Create(ctx, &model.BorrowRecord{
		EquipmentID:        equipment.ID,
		BorrowDate:         time.Now(),
		ExpectedReturnDate: time.Now().AddDate(0, 0, 7),
	}, actor)
	if err == nil {
		t.Fatalf("expected borrow blocked during disposal approval")
	}
	var be *apperrors.BusinessError
	if !errors.As(err, &be) || be.HTTPStatus != 409 {
		t.Fatalf("expected 409 business error, got %v", err)
	}

	_, err = env.reservationService.Create(ctx, &model.Reservation{
		EquipmentID: equipment.ID,
		StartTime:   time.Now().AddDate(0, 0, 1),
		EndTime:     time.Now().AddDate(0, 0, 1).Add(2 * time.Hour),
		Purpose:     "实验",
	}, actor)
	if err == nil {
		t.Fatalf("expected reservation blocked during disposal approval")
	}
}

func TestDisposalService_DuplicateSubmitConflict(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	actor := Actor{UserID: env.ownerID, Username: "admin", Role: "Admin"}
	submitDisposalForEquipment(t, env, 1, actor)

	_, err := env.disposalService.Submit(ctx, &model.DisposalApproval{EquipmentID: 1, Reason: "重复提交"}, actor)
	var be *apperrors.BusinessError
	if !errors.As(err, &be) || be.HTTPStatus != 409 {
		t.Fatalf("expected 409 conflict on duplicate submit, got %v", err)
	}
}

func TestDisposalService_ReviewReturnsWithBorrowBlocker(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	manager := Actor{UserID: env.ownerID, Username: "admin", Role: "LabManager"}
	equipment, _ := env.equipmentRepo.FindByID(ctx, 1)

	// 先产生一笔已审批（未归还）的借用。
	record, err := env.borrowService.Create(ctx, &model.BorrowRecord{
		EquipmentID:        equipment.ID,
		BorrowDate:         time.Now(),
		ExpectedReturnDate: time.Now().AddDate(0, 0, 7),
	}, manager)
	if err != nil {
		t.Fatalf("create borrow: %v", err)
	}
	if err := env.borrowService.Approve(ctx, record.ID, manager); err != nil {
		t.Fatalf("approve borrow: %v", err)
	}

	approval := submitDisposalForEquipment(t, env, equipment.ID, manager)
	reviewed, result, err := env.disposalService.Review(ctx, approval.ID, manager)
	if err != nil {
		t.Fatalf("review: %v", err)
	}
	if result.Approved {
		t.Fatalf("expected return due to blockers")
	}
	if reviewed.Status != constants.DisposalStatusReturned {
		t.Fatalf("expected returned status, got %s", reviewed.Status)
	}
	if len(reviewed.Blockers) != 1 || reviewed.Blockers[0].Type != "Borrow" {
		t.Fatalf("expected one borrow blocker snapshot, got %+v", reviewed.Blockers)
	}
	// 退回后设备不能变成 Retired。
	equipment, _ = env.equipmentRepo.FindByID(ctx, equipment.ID)
	if equipment.Status == constants.AssetStatusRetired {
		t.Fatalf("equipment must not be retired when review returned")
	}
	// 原借用记录保留。
	borrow, err := env.borrowService.Get(ctx, record.ID)
	if err != nil {
		t.Fatalf("borrow record must be kept: %v", err)
	}
	if borrow.Status != constants.BorrowStatusApproved {
		t.Fatalf("borrow status must remain Approved, got %s", borrow.Status)
	}
}

func TestDisposalService_ReviewReturnsWithReservationBlocker(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	manager := Actor{UserID: env.ownerID, Username: "admin", Role: "LabManager"}

	reservation, err := env.reservationService.Create(ctx, &model.Reservation{
		EquipmentID: 1,
		StartTime:   time.Now().AddDate(0, 0, 2),
		EndTime:     time.Now().AddDate(0, 0, 2).Add(2 * time.Hour),
		Purpose:     "预约实验",
	}, manager)
	if err != nil {
		t.Fatalf("create reservation: %v", err)
	}
	if err := env.reservationService.Approve(ctx, reservation.ID, manager); err != nil {
		t.Fatalf("approve reservation: %v", err)
	}

	approval := submitDisposalForEquipment(t, env, 1, manager)
	_, result, err := env.disposalService.Review(ctx, approval.ID, manager)
	if err != nil {
		t.Fatalf("review: %v", err)
	}
	if result.Approved || len(result.Blockers) != 1 || result.Blockers[0].Type != "Reservation" {
		t.Fatalf("expected one reservation blocker, got approved=%v blockers=%+v", result.Approved, result.Blockers)
	}
}

func TestDisposalService_ConcurrentReviewConflict(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	manager := Actor{UserID: env.ownerID, Username: "admin", Role: "LabManager"}
	approval := submitDisposalForEquipment(t, env, 1, manager)

	if _, _, err := env.disposalService.Review(ctx, approval.ID, manager); err != nil {
		t.Fatalf("first review: %v", err)
	}
	_, _, err := env.disposalService.Review(ctx, approval.ID, manager)
	var be *apperrors.BusinessError
	if !errors.As(err, &be) || be.HTTPStatus != 409 {
		t.Fatalf("expected 409 on second review, got %v", err)
	}
}

func TestDisposalService_ReviewApprovedClearsActiveThenRetires(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	manager := Actor{UserID: env.ownerID, Username: "admin", Role: "LabManager"}

	first := submitDisposalForEquipment(t, env, 1, manager)
	_, result, err := env.disposalService.Review(ctx, first.ID, manager)
	if err != nil || !result.Approved {
		t.Fatalf("first review should be approved: %v blockers=%v", err, result)
	}
	// 审批结束后活动标记释放，但设备已报废无法再次提交。
	_, err = env.disposalService.Submit(ctx, &model.DisposalApproval{EquipmentID: 1, Reason: "再次提交"}, manager)
	var be *apperrors.BusinessError
	if !errors.As(err, &be) || be.HTTPStatus != 409 {
		t.Fatalf("retired equipment cannot be submitted again, got %v", err)
	}
}
