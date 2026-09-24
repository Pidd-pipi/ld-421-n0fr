package model

// EquipmentCategory 设备分类实体，支持树形层级。
type EquipmentCategory struct {
	Base
	Name        string              `gorm:"size:128;not null" json:"name"`
	ParentID    *uint               `gorm:"index" json:"parentId"`
	Description string              `gorm:"size:512" json:"description"`
	Icon        string              `gorm:"size:128" json:"icon"`
	Children    []EquipmentCategory `gorm:"-" json:"children,omitempty"`
}

func (EquipmentCategory) TableName() string { return "equipment_categories" }
