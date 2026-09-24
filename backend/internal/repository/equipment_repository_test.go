package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/labequipment/lab-equipment/internal/constants"
	"github.com/labequipment/lab-equipment/internal/model"
)

func TestEquipmentRepository_CreateAndFindByID(t *testing.T) {
	db := newTestDB(t)
	repo := NewEquipmentRepository(db)
	ctx := context.Background()

	equipment := &model.Equipment{Name: "分光光度计", Code: "EQ-100", CategoryID: 1, OwnerID: 1, Status: constants.AssetStatusAvailable}
	if err := repo.Create(ctx, equipment); err != nil {
		t.Fatalf("create equipment: %v", err)
	}
	got, err := repo.FindByID(ctx, equipment.ID)
	if err != nil {
		t.Fatalf("find by id: %v", err)
	}
	if got.Code != "EQ-100" {
		t.Fatalf("expected code EQ-100, got %s", got.Code)
	}
}

func TestEquipmentRepository_FindByID_NotFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewEquipmentRepository(db)
	_, err := repo.FindByID(context.Background(), 999)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestEquipmentRepository_List(t *testing.T) {
	db := newTestDB(t)
	repo := NewEquipmentRepository(db)
	ctx := context.Background()
	items := []model.Equipment{
		{Name: "光谱仪", Code: "EQ-1", CategoryID: 1, OwnerID: 1, Status: constants.AssetStatusAvailable},
		{Name: "色谱仪", Code: "EQ-2", CategoryID: 1, OwnerID: 1, Status: constants.AssetStatusMaintenance},
	}
	for i := range items {
		if err := repo.Create(ctx, &items[i]); err != nil {
			t.Fatalf("create item %d: %v", i, err)
		}
	}
	tests := []struct {
		name   string
		filter EquipmentFilter
		want   int
	}{
		{name: "all", filter: EquipmentFilter{Pagination: Pagination{Page: 1, PageSize: 10}}, want: 2},
		{name: "keyword", filter: EquipmentFilter{Keyword: "光谱", Pagination: Pagination{Page: 1, PageSize: 10}}, want: 1},
		{name: "status", filter: EquipmentFilter{Status: constants.AssetStatusMaintenance, Pagination: Pagination{Page: 1, PageSize: 10}}, want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, total, err := repo.List(ctx, tt.filter)
			if err != nil {
				t.Fatalf("list: %v", err)
			}
			if total != int64(tt.want) {
				t.Fatalf("expected total %d, got %d", tt.want, total)
			}
			if len(list) != tt.want {
				t.Fatalf("expected len %d, got %d", tt.want, len(list))
			}
		})
	}
}
