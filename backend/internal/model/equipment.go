package model

import (
	"time"

	"github.com/labequipment/lab-equipment/internal/constants"
)

// Equipment 设备资产实体。
type Equipment struct {
	Base
	Name           string                `gorm:"size:128;not null" json:"name"`
	Code           string                `gorm:"size:64;uniqueIndex;not null" json:"code"`
	CategoryID     uint                  `gorm:"index;not null" json:"categoryId"`
	BrandModel     string                `gorm:"size:128" json:"brandModel"`
	SerialNumber   string                `gorm:"size:128" json:"serialNumber"`
	PurchaseDate   *time.Time            `json:"purchaseDate"`
	PurchasePrice  float64               `gorm:"type:decimal(12,2)" json:"purchasePrice"`
	Location       string                `gorm:"size:255" json:"location"`
	Status         constants.AssetStatus `gorm:"size:32;index;not null;default:Available" json:"status"`
	OwnerID        uint                  `gorm:"index;not null" json:"ownerId"`
	Supplier       string                `gorm:"size:128" json:"supplier"`
	WarrantyExpiry *time.Time            `json:"warrantyExpiry"`
	ImageURL       string                `gorm:"size:512" json:"imageUrl"`
	Category       *EquipmentCategory    `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Owner          *User                 `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
}

func (Equipment) TableName() string { return "equipment" }
