package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/labequipment/lab-equipment/internal/dto"
	"github.com/labequipment/lab-equipment/internal/service"
)

// DashboardHandler 仪表盘处理器。
type DashboardHandler struct {
	dashboardService *service.DashboardService
}

// NewDashboardHandler 构造仪表盘处理器。
func NewDashboardHandler(dashboardService *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

// Stats 获取仪表盘统计。
func (h *DashboardHandler) Stats(c *gin.Context) {
	stats, err := h.dashboardService.Stats(c.Request.Context())
	if err != nil {
		failError(c, err)
		return
	}
	topBorrows := make([]dto.TopBorrowItem, 0, len(stats.TopBorrows))
	for _, item := range stats.TopBorrows {
		topBorrows = append(topBorrows, dto.TopBorrowItem{
			EquipmentID: item.EquipmentID,
			Name:        item.Name,
			Code:        item.Code,
			Count:       item.Count,
		})
	}
	expiring := make([]dto.ExpiringItem, 0, len(stats.ExpiringWarranty))
	for _, equipment := range stats.ExpiringWarranty {
		expiring = append(expiring, dto.ExpiringItem{
			ID:             equipment.ID,
			Name:           equipment.Name,
			Code:           equipment.Code,
			WarrantyExpiry: equipment.WarrantyExpiry,
		})
	}
	ok(c, dto.DashboardStatsResponse{
		StatusDistribution:  stats.StatusDistribution,
		TopBorrows:          topBorrows,
		ExpiringWarranty:    expiring,
		PendingBorrows:      stats.PendingBorrows,
		PendingReservations: stats.PendingReservations,
	})
}
