package dto

import "time"

// CreateEquipmentRequest 新增设备请求。
type CreateEquipmentRequest struct {
	Name           string  `json:"name" binding:"required,max=128"`
	Code           string  `json:"code" binding:"required,max=64"`
	CategoryID     uint    `json:"categoryId" binding:"required"`
	BrandModel     string  `json:"brandModel" binding:"omitempty,max=128"`
	SerialNumber   string  `json:"serialNumber" binding:"omitempty,max=128"`
	PurchaseDate   string  `json:"purchaseDate" binding:"omitempty"`
	PurchasePrice  float64 `json:"purchasePrice" binding:"omitempty,gte=0"`
	Location       string  `json:"location" binding:"omitempty,max=255"`
	OwnerID        uint    `json:"ownerId" binding:"required"`
	Supplier       string  `json:"supplier" binding:"omitempty,max=128"`
	WarrantyExpiry string  `json:"warrantyExpiry" binding:"omitempty"`
	ImageURL       string  `json:"imageUrl" binding:"omitempty,max=512"`
}

// UpdateEquipmentRequest 更新设备请求。
type UpdateEquipmentRequest struct {
	Name           string  `json:"name" binding:"required,max=128"`
	Code           string  `json:"code" binding:"required,max=64"`
	CategoryID     uint    `json:"categoryId" binding:"required"`
	BrandModel     string  `json:"brandModel" binding:"omitempty,max=128"`
	SerialNumber   string  `json:"serialNumber" binding:"omitempty,max=128"`
	PurchaseDate   string  `json:"purchaseDate" binding:"omitempty"`
	PurchasePrice  float64 `json:"purchasePrice" binding:"omitempty,gte=0"`
	Location       string  `json:"location" binding:"omitempty,max=255"`
	OwnerID        uint    `json:"ownerId" binding:"required"`
	Supplier       string  `json:"supplier" binding:"omitempty,max=128"`
	WarrantyExpiry string  `json:"warrantyExpiry" binding:"omitempty"`
	ImageURL       string  `json:"imageUrl" binding:"omitempty,max=512"`
}

// UpdateEquipmentStatusRequest 修改设备资产状态。
type UpdateEquipmentStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// TransferOwnerRequest 转移责任人。
type TransferOwnerRequest struct {
	OwnerID uint `json:"ownerId" binding:"required"`
}

// EquipmentResponse 设备信息返回。
type EquipmentResponse struct {
	ID             uint       `json:"id"`
	Name           string     `json:"name"`
	Code           string     `json:"code"`
	CategoryID     uint       `json:"categoryId"`
	CategoryName   string     `json:"categoryName,omitempty"`
	BrandModel     string     `json:"brandModel"`
	SerialNumber   string     `json:"serialNumber"`
	PurchaseDate   *time.Time `json:"purchaseDate"`
	PurchasePrice  float64    `json:"purchasePrice"`
	Location       string     `json:"location"`
	Status         string     `json:"status"`
	OwnerID        uint       `json:"ownerId"`
	OwnerName      string     `json:"ownerName,omitempty"`
	Supplier       string     `json:"supplier"`
	WarrantyExpiry *time.Time `json:"warrantyExpiry"`
	ImageURL       string     `json:"imageUrl"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}
