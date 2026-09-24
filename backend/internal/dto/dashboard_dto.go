package dto

import "time"

// DashboardStatsResponse 仪表盘统计返回。
type DashboardStatsResponse struct {
	StatusDistribution  map[string]int64 `json:"statusDistribution"`
	TopBorrows          []TopBorrowItem  `json:"topBorrows"`
	ExpiringWarranty    []ExpiringItem   `json:"expiringWarranty"`
	PendingBorrows      int64            `json:"pendingBorrows"`
	PendingReservations int64            `json:"pendingReservations"`
}

// TopBorrowItem 本月借用次数排行项。
type TopBorrowItem struct {
	EquipmentID uint   `json:"equipmentId"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Count       int64  `json:"count"`
}

// ExpiringItem 即将过保设备项。
type ExpiringItem struct {
	ID             uint       `json:"id"`
	Name           string     `json:"name"`
	Code           string     `json:"code"`
	WarrantyExpiry *time.Time `json:"warrantyExpiry"`
}
