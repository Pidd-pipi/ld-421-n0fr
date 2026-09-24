package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/labequipment/lab-equipment/internal/model"
	"github.com/labequipment/lab-equipment/internal/repository"
)

// DashboardService 仪表盘统计服务。
type DashboardService struct {
	equipmentRepo   repository.EquipmentRepository
	borrowRepo      repository.BorrowRepository
	reservationRepo repository.ReservationRepository
	logger          *slog.Logger
}

// NewDashboardService 构造仪表盘服务。
func NewDashboardService(
	equipmentRepo repository.EquipmentRepository,
	borrowRepo repository.BorrowRepository,
	reservationRepo repository.ReservationRepository,
	logger *slog.Logger,
) *DashboardService {
	return &DashboardService{equipmentRepo: equipmentRepo, borrowRepo: borrowRepo, reservationRepo: reservationRepo, logger: logger}
}

// Stats 聚合仪表盘所需统计。
func (s *DashboardService) Stats(ctx context.Context) (*DashboardStats, error) {
	statusMap, err := s.equipmentRepo.CountByStatus(ctx)
	if err != nil {
		return nil, fmt.Errorf("dashboard status distribution: %w", err)
	}
	topBorrows, err := s.borrowRepo.TopEquipmentThisMonth(ctx, 10)
	if err != nil {
		return nil, fmt.Errorf("dashboard top borrows: %w", err)
	}
	expiring, err := s.equipmentRepo.ListExpiringWarranty(ctx, 30)
	if err != nil {
		return nil, fmt.Errorf("dashboard expiring warranty: %w", err)
	}
	pendingBorrows, err := s.borrowRepo.CountPending(ctx)
	if err != nil {
		return nil, fmt.Errorf("dashboard pending borrows: %w", err)
	}
	pendingReservations, err := s.reservationRepo.CountPending(ctx)
	if err != nil {
		return nil, fmt.Errorf("dashboard pending reservations: %w", err)
	}
	statusDistribution := make(map[string]int64, len(statusMap))
	for k, v := range statusMap {
		statusDistribution[string(k)] = v
	}
	return &DashboardStats{
		StatusDistribution:  statusDistribution,
		TopBorrows:          topBorrows,
		ExpiringWarranty:    expiring,
		PendingBorrows:      pendingBorrows,
		PendingReservations: pendingReservations,
	}, nil
}

// DashboardStats 仪表盘统计结果。
type DashboardStats struct {
	StatusDistribution  map[string]int64
	TopBorrows          []repository.BorrowTopStat
	ExpiringWarranty    []model.Equipment
	PendingBorrows      int64
	PendingReservations int64
}
