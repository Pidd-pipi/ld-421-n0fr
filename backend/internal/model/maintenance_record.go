package model

import (
	"time"

	"github.com/labequipment/lab-equipment/internal/constants"
)

// MaintenanceRecord 设备维护记录实体。
type MaintenanceRecord struct {
	Base
	EquipmentID         uint                        `gorm:"index;not null" json:"equipmentId"`
	Type                constants.MaintenanceType   `gorm:"size:32;not null" json:"type"`
	Content             string                      `gorm:"size:1024" json:"content"`
	MaintenanceDate     time.Time                   `json:"maintenanceDate"`
	NextMaintenanceDate *time.Time                  `json:"nextMaintenanceDate"`
	Cost                float64                     `gorm:"type:decimal(12,2)" json:"cost"`
	MaintainerID        uint                        `gorm:"index;not null" json:"maintainerId"`
	Result              constants.MaintenanceResult `gorm:"size:32;not null;default:Pass" json:"result"`
	Equipment           *Equipment                  `gorm:"foreignKey:EquipmentID" json:"equipment,omitempty"`
	Maintainer          *User                       `gorm:"foreignKey:MaintainerID" json:"maintainer,omitempty"`
}

func (MaintenanceRecord) TableName() string { return "maintenance_records" }
