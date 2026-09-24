package repository

import (
	"context"
	"testing"
	"time"

	"github.com/labequipment/lab-equipment/internal/constants"
	"github.com/labequipment/lab-equipment/internal/model"
)

func TestBorrowRepository_ListAndCountPending(t *testing.T) {
	db := newTestDB(t)
	repo := NewBorrowRepository(db)
	ctx := context.Background()

	record := &model.BorrowRecord{
		EquipmentID:        1,
		BorrowerID:         1,
		BorrowDate:         time.Now(),
		ExpectedReturnDate: time.Now().AddDate(0, 0, 7),
		Status:             constants.BorrowStatusPending,
	}
	if err := repo.Create(ctx, record); err != nil {
		t.Fatalf("create borrow: %v", err)
	}
	pending, err := repo.CountPending(ctx)
	if err != nil {
		t.Fatalf("count pending: %v", err)
	}
	if pending != 1 {
		t.Fatalf("expected 1 pending, got %d", pending)
	}
	list, total, err := repo.List(ctx, BorrowFilter{Pagination: Pagination{Page: 1, PageSize: 10}})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 1 || len(list) != 1 {
		t.Fatalf("expected 1 record, got total=%d len=%d", total, len(list))
	}
}
