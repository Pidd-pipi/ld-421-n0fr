package service

import (
	"context"
	"testing"
	"time"

	"github.com/labequipment/lab-equipment/internal/model"
)

func TestReservationService_CreateAndConflict(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	actor := Actor{UserID: env.ownerID, Username: "admin", Role: "Student"}
	start := time.Now().Add(24 * time.Hour)
	end := start.Add(2 * time.Hour)

	first, err := env.reservationService.Create(ctx, &model.Reservation{
		EquipmentID: 1,
		StartTime:   start,
		EndTime:     end,
		Purpose:     "实验",
	}, actor)
	if err != nil {
		t.Fatalf("create first: %v", err)
	}
	if first.ID == 0 {
		t.Fatalf("expected created id > 0")
	}

	_, err = env.reservationService.Create(ctx, &model.Reservation{
		EquipmentID: 1,
		StartTime:   start.Add(30 * time.Minute),
		EndTime:     end.Add(30 * time.Minute),
		Purpose:     "冲突预约",
	}, actor)
	if err == nil {
		t.Fatalf("expected conflict error")
	}
}
