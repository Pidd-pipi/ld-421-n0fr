package service

import (
	"context"
	"testing"
	"time"

	"github.com/labequipment/lab-equipment/internal/constants"
	"github.com/labequipment/lab-equipment/internal/model"
)

func TestBorrowService_FullFlow(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	actor := Actor{UserID: env.ownerID, Username: "admin", Role: "LabManager"}
	equipment, err := env.equipmentRepo.FindByID(ctx, 1)
	if err != nil {
		t.Fatalf("find equipment: %v", err)
	}

	record, err := env.borrowService.Create(ctx, &model.BorrowRecord{
		EquipmentID:        equipment.ID,
		BorrowDate:         time.Now(),
		ExpectedReturnDate: time.Now().AddDate(0, 0, 7),
	}, actor)
	if err != nil {
		t.Fatalf("create borrow: %v", err)
	}
	if record.Status != constants.BorrowStatusPending {
		t.Fatalf("expected pending, got %s", record.Status)
	}

	if err := env.borrowService.Approve(ctx, record.ID, actor); err != nil {
		t.Fatalf("approve: %v", err)
	}
	got, _ := env.borrowService.Get(ctx, record.ID)
	if got.Status != constants.BorrowStatusApproved {
		t.Fatalf("expected approved, got %s", got.Status)
	}

	if err := env.borrowService.Return(ctx, record.ID, time.Now(), constants.ReturnConditionGood, actor); err != nil {
		t.Fatalf("return: %v", err)
	}
	got, _ = env.borrowService.Get(ctx, record.ID)
	if got.Status != constants.BorrowStatusReturned {
		t.Fatalf("expected returned, got %s", got.Status)
	}
}

func TestBorrowService_CreateUnavailable(t *testing.T) {
	env := newTestEnv(t)
	actor := Actor{UserID: env.ownerID, Username: "admin", Role: "Student"}
	equipment, _ := env.equipmentRepo.FindByID(context.Background(), 1)
	equipment.Status = constants.AssetStatusMaintenance
	_ = env.equipmentRepo.Update(context.Background(), equipment)

	_, err := env.borrowService.Create(context.Background(), &model.BorrowRecord{
		EquipmentID:        equipment.ID,
		BorrowDate:         time.Now(),
		ExpectedReturnDate: time.Now().AddDate(0, 0, 7),
	}, actor)
	if err == nil {
		t.Fatalf("expected error for unavailable equipment")
	}
}
