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

func disposalActor(env *testEnv) Actor {
	return Actor{UserID: env.ownerID, Username: "admin", Role: "LabManager"}
}

func expectBusinessCode(t *testing.T, err error, code int) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected business error code %d, got nil", code)
	}
	var be *apperrors.BusinessError
	if !errors.As(err, &be) || be.Code != code {
		t.Fatalf("expected business error code %d, got %v", code, err)
	}
}

// 提交报废申请后：设备状态不变，但新的借用与预约被暂停；重复提交返回冲突。
func TestDisposalService_SubmitPausesBorrowAndReservation(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	actor := disposalActor(env)

	item, err := env.disposalService.Submit(ctx, 1, "设备老化，精度下降", actor)
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	if item.Request.Status != constants.DisposalStatusPending {
		t.Fatalf("expected Pending, got %s", item.Request.Status)
	}

	equipment, err := env.equipmentRepo.FindByID(ctx, 1)
	if err != nil {
		t.Fatalf("find equipment: %v", err)
	}
	if equipment.Status != constants.AssetStatusAvailable {
		t.Fatalf("equipment status should stay Available before approval, got %s", equipment.Status)
	}

	_, err = env.disposalService.Submit(ctx, 1, "重复提交", actor)
	expectBusinessCode(t, err, 40900)

	_, err = env.borrowService.Create(ctx, &model.BorrowRecord{
		EquipmentID:        1,
		BorrowDate:         time.Now(),
		ExpectedReturnDate: time.Now().AddDate(0, 0, 7),
	}, actor)
	expectBusinessCode(t, err, 40900)

	_, err = env.reservationService.Create(ctx, &model.Reservation{
		EquipmentID: 1,
		StartTime:   time.Now().Add(24 * time.Hour),
		EndTime:     time.Now().Add(26 * time.Hour),
	}, actor)
	expectBusinessCode(t, err, 40900)
}

// 审批时存在未归还借用：列出阻塞项并退回，设备不转为 Retired；
// 归还后重新提交并审批通过，设备才转为 Retired，原借用记录保留。
func TestDisposalService_ApproveWithUnreturnedBorrow(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	actor := disposalActor(env)

	borrow, err := env.borrowService.Create(ctx, &model.BorrowRecord{
		EquipmentID:        1,
		BorrowDate:         time.Now(),
		ExpectedReturnDate: time.Now().AddDate(0, 0, 7),
	}, actor)
	if err != nil {
		t.Fatalf("create borrow: %v", err)
	}
	if err := env.borrowService.Approve(ctx, borrow.ID, actor); err != nil {
		t.Fatalf("approve borrow: %v", err)
	}

	submitted, err := env.disposalService.Submit(ctx, 1, "设备老化", actor)
	if err != nil {
		t.Fatalf("submit disposal: %v", err)
	}

	returned, err := env.disposalService.Approve(ctx, 1, submitted.Request.ID, actor)
	if err != nil {
		t.Fatalf("approve disposal: %v", err)
	}
	if returned.Request.Status != constants.DisposalStatusRejected {
		t.Fatalf("expected Rejected（退回）, got %s", returned.Request.Status)
	}
	if len(returned.Blockers) != 1 || returned.Blockers[0].Type != "borrow" || returned.Blockers[0].RecordID != borrow.ID {
		t.Fatalf("expected one borrow blocker, got %+v", returned.Blockers)
	}
	equipment, err := env.equipmentRepo.FindByID(ctx, 1)
	if err != nil {
		t.Fatalf("find equipment: %v", err)
	}
	if equipment.Status == constants.AssetStatusRetired {
		t.Fatalf("equipment must not be Retired while blocked")
	}

	// 退回快照可在列表中查到。
	items, err := env.disposalService.ListByEquipment(ctx, 1)
	if err != nil {
		t.Fatalf("list disposals: %v", err)
	}
	if len(items) != 1 || len(items[0].Blockers) != 1 {
		t.Fatalf("expected stored blocker snapshot, got %+v", items)
	}

	// 归还借用后重新提交，审批通过，设备转为 Retired，借用记录保留。
	if err := env.borrowService.Return(ctx, borrow.ID, time.Now(), constants.ReturnConditionGood, actor); err != nil {
		t.Fatalf("return borrow: %v", err)
	}
	resubmitted, err := env.disposalService.Submit(ctx, 1, "设备老化，再次申请", actor)
	if err != nil {
		t.Fatalf("resubmit disposal: %v", err)
	}
	approved, err := env.disposalService.Approve(ctx, 1, resubmitted.Request.ID, actor)
	if err != nil {
		t.Fatalf("approve disposal: %v", err)
	}
	if approved.Request.Status != constants.DisposalStatusApproved {
		t.Fatalf("expected Approved, got %s", approved.Request.Status)
	}
	equipment, err = env.equipmentRepo.FindByID(ctx, 1)
	if err != nil {
		t.Fatalf("find equipment: %v", err)
	}
	if equipment.Status != constants.AssetStatusRetired {
		t.Fatalf("expected Retired, got %s", equipment.Status)
	}
	record, err := env.borrowService.Get(ctx, borrow.ID)
	if err != nil {
		t.Fatalf("borrow record should be kept: %v", err)
	}
	if record.Status != constants.BorrowStatusReturned {
		t.Fatalf("expected borrow record Returned, got %s", record.Status)
	}

	// 已报废设备不允许再次提交报废申请。
	_, err = env.disposalService.Submit(ctx, 1, "重复报废", actor)
	expectBusinessCode(t, err, 40900)
}

// 审批时存在未结束预约：列出阻塞项并退回；预约结束后方可审批通过。
func TestDisposalService_ApproveWithUnfinishedReservation(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	actor := disposalActor(env)

	reservation, err := env.reservationService.Create(ctx, &model.Reservation{
		EquipmentID: 1,
		StartTime:   time.Now().Add(24 * time.Hour),
		EndTime:     time.Now().Add(26 * time.Hour),
		Purpose:     "实验",
	}, actor)
	if err != nil {
		t.Fatalf("create reservation: %v", err)
	}

	submitted, err := env.disposalService.Submit(ctx, 1, "设备淘汰", actor)
	if err != nil {
		t.Fatalf("submit disposal: %v", err)
	}
	returned, err := env.disposalService.Approve(ctx, 1, submitted.Request.ID, actor)
	if err != nil {
		t.Fatalf("approve disposal: %v", err)
	}
	if returned.Request.Status != constants.DisposalStatusRejected {
		t.Fatalf("expected Rejected（退回）, got %s", returned.Request.Status)
	}
	if len(returned.Blockers) != 1 || returned.Blockers[0].Type != "reservation" || returned.Blockers[0].RecordID != reservation.ID {
		t.Fatalf("expected one reservation blocker, got %+v", returned.Blockers)
	}

	if err := env.reservationService.Cancel(ctx, reservation.ID, actor); err != nil {
		t.Fatalf("cancel reservation: %v", err)
	}
	resubmitted, err := env.disposalService.Submit(ctx, 1, "设备淘汰，再次申请", actor)
	if err != nil {
		t.Fatalf("resubmit disposal: %v", err)
	}
	approved, err := env.disposalService.Approve(ctx, 1, resubmitted.Request.ID, actor)
	if err != nil {
		t.Fatalf("approve disposal: %v", err)
	}
	if approved.Request.Status != constants.DisposalStatusApproved {
		t.Fatalf("expected Approved, got %s", approved.Request.Status)
	}
	kept, err := env.reservationService.Get(ctx, reservation.ID)
	if err != nil {
		t.Fatalf("reservation record should be kept: %v", err)
	}
	if kept.Status != constants.ReservationStatusCancelled {
		t.Fatalf("expected reservation Cancelled, got %s", kept.Status)
	}
}

// 两人同时审批同一申请：后到操作返回冲突。
func TestDisposalService_DoubleApproveConflict(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	actor := disposalActor(env)

	submitted, err := env.disposalService.Submit(ctx, 1, "设备老化", actor)
	if err != nil {
		t.Fatalf("submit disposal: %v", err)
	}
	if _, err := env.disposalService.Approve(ctx, 1, submitted.Request.ID, actor); err != nil {
		t.Fatalf("first approve: %v", err)
	}
	_, err = env.disposalService.Approve(ctx, 1, submitted.Request.ID, actor)
	expectBusinessCode(t, err, 40900)
}
