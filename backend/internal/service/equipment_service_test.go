package service

import (
	"context"
	"errors"
	"testing"

	"github.com/labequipment/lab-equipment/internal/constants"
	apperrors "github.com/labequipment/lab-equipment/internal/errors"
	"github.com/labequipment/lab-equipment/internal/model"
)

func TestEquipmentService_Create(t *testing.T) {
	env := newTestEnv(t)
	actor := Actor{UserID: env.ownerID, Username: "admin", Role: "Admin"}
	tests := []struct {
		name      string
		equipment model.Equipment
		wantErr   bool
		wantCode  int
	}{
		{
			name:      "valid",
			equipment: model.Equipment{Name: "新设备", Code: "EQ-NEW", CategoryID: env.categoryID, OwnerID: env.ownerID},
			wantErr:   false,
		},
		{
			name:      "invalid category",
			equipment: model.Equipment{Name: "新设备2", Code: "EQ-NEW2", CategoryID: 999, OwnerID: env.ownerID},
			wantErr:   true,
			wantCode:  40400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			created, err := env.equipmentService.Create(context.Background(), &tt.equipment, actor)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				var be *apperrors.BusinessError
				if !errors.As(err, &be) || be.Code != tt.wantCode {
					t.Fatalf("expected business error code %d, got %v", tt.wantCode, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if created.ID == 0 {
				t.Fatalf("expected created id > 0")
			}
		})
	}
}

func TestEquipmentService_Retire(t *testing.T) {
	env := newTestEnv(t)
	actor := Actor{UserID: env.ownerID, Username: "admin", Role: "Admin"}
	created, err := env.equipmentService.Create(context.Background(), &model.Equipment{
		Name: "待报废", Code: "EQ-RET", CategoryID: env.categoryID, OwnerID: env.ownerID,
	}, actor)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	approval, err := env.disposalService.Submit(context.Background(), &model.DisposalApproval{
		EquipmentID: created.ID,
		Reason:      "设备老化严重",
	}, actor)
	if err != nil {
		t.Fatalf("submit disposal: %v", err)
	}
	reviewed, result, err := env.disposalService.Review(context.Background(), approval.ID, actor)
	if err != nil {
		t.Fatalf("review disposal: %v", err)
	}
	if !result.Approved {
		t.Fatalf("expected approval without blockers, got blockers: %v", result.Blockers)
	}
	if reviewed.Status != constants.DisposalStatusApproved {
		t.Fatalf("expected Approved, got %s", reviewed.Status)
	}
	got, err := env.equipmentRepo.FindByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if got.Status != constants.AssetStatusRetired {
		t.Fatalf("expected Retired, got %s", got.Status)
	}
}
