package repository

import (
	"context"
	"testing"
	"time"

	"github.com/labequipment/lab-equipment/internal/constants"
	"github.com/labequipment/lab-equipment/internal/model"
)

func TestDisposalRepository_ReviewWithBlockers(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	role := model.Role{Code: "Student", Name: "学生"}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}
	user := model.User{Username: "stu", PasswordHash: "x", Name: "学生甲", RoleID: role.ID, Active: true}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	equipment := model.Equipment{Name: "设备A", Code: "EQ-A", OwnerID: user.ID, Status: constants.AssetStatusAvailable}
	if err := db.Create(&equipment).Error; err != nil {
		t.Fatalf("create equipment: %v", err)
	}

	// 一笔未归还借用
	borrow := model.BorrowRecord{
		EquipmentID:        equipment.ID,
		BorrowerID:         user.ID,
		BorrowDate:         time.Now(),
		ExpectedReturnDate: time.Now().AddDate(0, 0, 3),
		Status:             constants.BorrowStatusApproved,
	}
	if err := db.Create(&borrow).Error; err != nil {
		t.Fatalf("create borrow: %v", err)
	}
	// 一笔未结束预约
	reservation := model.Reservation{
		EquipmentID: equipment.ID,
		UserID:      user.ID,
		StartTime:   time.Now().AddDate(0, 0, 5),
		EndTime:     time.Now().AddDate(0, 0, 5).Add(time.Hour),
		Status:      constants.ReservationStatusPending,
	}
	if err := db.Create(&reservation).Error; err != nil {
		t.Fatalf("create reservation: %v", err)
	}

	repo := NewDisposalRepository(db)
	activeEquipmentID := equipment.ID
	approval := &model.DisposalApproval{
		EquipmentID:       equipment.ID,
		Reason:            "损坏无法维修",
		Status:            constants.DisposalStatusPending,
		ApplicantID:       user.ID,
		Active:            true,
		ActiveEquipmentID: &activeEquipmentID,
	}
	if err := repo.Create(ctx, approval); err != nil {
		t.Fatalf("create approval: %v", err)
	}

	reviewed, result, err := repo.Review(ctx, approval.ID, user.ID)
	if err != nil {
		t.Fatalf("review: %v", err)
	}
	if result.Approved {
		t.Fatalf("expected review returned due to blockers")
	}
	if reviewed.Status != constants.DisposalStatusReturned {
		t.Fatalf("expected Returned, got %s", reviewed.Status)
	}
	if len(result.Blockers) != 2 {
		t.Fatalf("expected 2 blockers, got %d: %+v", len(result.Blockers), result.Blockers)
	}
	if reviewed.Active || reviewed.ActiveEquipmentID != nil {
		t.Fatalf("approval should be inactive after review")
	}
	// 存在阻塞项时设备状态保持不变。
	var updated model.Equipment
	if err := db.First(&updated, equipment.ID).Error; err != nil {
		t.Fatalf("reload equipment: %v", err)
	}
	if updated.Status != constants.AssetStatusAvailable {
		t.Fatalf("equipment status should remain Available, got %s", updated.Status)
	}
}

func TestDisposalRepository_ReviewApprovedRetiresEquipment(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	user := model.User{Username: "admin", PasswordHash: "x", Name: "管理员", Active: true}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	equipment := model.Equipment{Name: "设备B", Code: "EQ-B", OwnerID: user.ID, Status: constants.AssetStatusAvailable}
	if err := db.Create(&equipment).Error; err != nil {
		t.Fatalf("create equipment: %v", err)
	}

	repo := NewDisposalRepository(db)
	activeEquipmentID := equipment.ID
	approval := &model.DisposalApproval{
		EquipmentID:       equipment.ID,
		Reason:            "到龄退役",
		Status:            constants.DisposalStatusPending,
		ApplicantID:       user.ID,
		Active:            true,
		ActiveEquipmentID: &activeEquipmentID,
	}
	if err := repo.Create(ctx, approval); err != nil {
		t.Fatalf("create approval: %v", err)
	}

	_, result, err := repo.Review(ctx, approval.ID, user.ID)
	if err != nil {
		t.Fatalf("review: %v", err)
	}
	if !result.Approved || len(result.Blockers) != 0 {
		t.Fatalf("expected clean approval, got %+v", result)
	}
	var updated model.Equipment
	if err := db.First(&updated, equipment.ID).Error; err != nil {
		t.Fatalf("reload equipment: %v", err)
	}
	if updated.Status != constants.AssetStatusRetired {
		t.Fatalf("expected Retired, got %s", updated.Status)
	}

	// 重复审批冲突。
	if _, _, err := repo.Review(ctx, approval.ID, user.ID); err != ErrConflict {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}
