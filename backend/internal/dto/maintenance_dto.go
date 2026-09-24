package dto

import "time"

// CreateMaintenanceRequest 创建维护记录。
type CreateMaintenanceRequest struct {
	EquipmentID         uint    `json:"equipmentId" binding:"required"`
	Type                string  `json:"type" binding:"required"`
	Content             string  `json:"content" binding:"omitempty,max=1024"`
	MaintenanceDate     string  `json:"maintenanceDate" binding:"required"`
	NextMaintenanceDate string  `json:"nextMaintenanceDate" binding:"omitempty"`
	Cost                float64 `json:"cost" binding:"omitempty,gte=0"`
	MaintainerID        uint    `json:"maintainerId" binding:"required"`
}

// ExecuteMaintenanceRequest 执行维护并记录结果。
type ExecuteMaintenanceRequest struct {
	Result string `json:"result" binding:"required"`
}

// MaintenanceResponse 维护记录返回。
type MaintenanceResponse struct {
	ID                  uint       `json:"id"`
	EquipmentID         uint       `json:"equipmentId"`
	EquipmentName       string     `json:"equipmentName,omitempty"`
	Type                string     `json:"type"`
	Content             string     `json:"content"`
	MaintenanceDate     time.Time  `json:"maintenanceDate"`
	NextMaintenanceDate *time.Time `json:"nextMaintenanceDate"`
	Cost                float64    `json:"cost"`
	MaintainerID        uint       `json:"maintainerId"`
	MaintainerName      string     `json:"maintainerName,omitempty"`
	Result              string     `json:"result"`
	CreatedAt           time.Time  `json:"createdAt"`
}

// MaintenanceStatsResponse 维护统计返回。
type MaintenanceStatsResponse struct {
	TotalRecords int64            `json:"totalRecords"`
	TotalCost    float64          `json:"totalCost"`
	ResultCount  map[string]int64 `json:"resultCount"`
}
