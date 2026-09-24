package dto

import "time"

// CreateReservationRequest 创建设备预约。
type CreateReservationRequest struct {
	EquipmentID uint   `json:"equipmentId" binding:"required"`
	StartTime   string `json:"startTime" binding:"required"`
	EndTime     string `json:"endTime" binding:"required"`
	Purpose     string `json:"purpose" binding:"omitempty,max=512"`
}

// ReservationResponse 预约记录返回。
type ReservationResponse struct {
	ID            uint      `json:"id"`
	EquipmentID   uint      `json:"equipmentId"`
	EquipmentName string    `json:"equipmentName,omitempty"`
	UserID        uint      `json:"userId"`
	UserName      string    `json:"userName,omitempty"`
	StartTime     time.Time `json:"startTime"`
	EndTime       time.Time `json:"endTime"`
	Purpose       string    `json:"purpose"`
	Status        string    `json:"status"`
	ApproverID    *uint     `json:"approverId"`
	ApproverName  string    `json:"approverName,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
}
