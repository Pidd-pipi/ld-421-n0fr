package model

import (
	"time"

	"github.com/labequipment/lab-equipment/internal/constants"
)

// Reservation 设备使用预约实体。
type Reservation struct {
	Base
	EquipmentID uint                        `gorm:"index;not null" json:"equipmentId"`
	UserID      uint                        `gorm:"index;not null" json:"userId"`
	StartTime   time.Time                   `json:"startTime"`
	EndTime     time.Time                   `json:"endTime"`
	Purpose     string                      `gorm:"size:512" json:"purpose"`
	Status      constants.ReservationStatus `gorm:"size:32;index;not null;default:Pending" json:"status"`
	ApproverID  *uint                       `json:"approverId"`
	Equipment   *Equipment                  `gorm:"foreignKey:EquipmentID" json:"equipment,omitempty"`
	User        *User                       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Approver    *User                       `gorm:"foreignKey:ApproverID" json:"approver,omitempty"`
}

func (Reservation) TableName() string { return "reservations" }
