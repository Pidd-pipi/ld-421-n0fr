package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/labequipment/lab-equipment/internal/constants"
	"github.com/labequipment/lab-equipment/internal/dto"
	"github.com/labequipment/lab-equipment/internal/middleware"
	"github.com/labequipment/lab-equipment/internal/model"
	"github.com/labequipment/lab-equipment/internal/repository"
	"github.com/labequipment/lab-equipment/internal/service"
)

// EquipmentHandler 设备处理器。
type EquipmentHandler struct {
	equipmentService *service.EquipmentService
}

// NewEquipmentHandler 构造设备处理器。
func NewEquipmentHandler(equipmentService *service.EquipmentService) *EquipmentHandler {
	return &EquipmentHandler{equipmentService: equipmentService}
}

// List 分页查询设备。
func (h *EquipmentHandler) List(c *gin.Context) {
	page, pageSize := parsePage(c)
	categoryID, _ := strconv.ParseUint(c.Query("category_id"), 10, 64)
	filter := repository.EquipmentFilter{
		Keyword:    c.Query("keyword"),
		CategoryID: uint(categoryID),
		Status:     constants.AssetStatus(c.Query("status")),
		Pagination: repository.Pagination{Page: page, PageSize: pageSize},
	}
	list, total, err := h.equipmentService.List(c.Request.Context(), filter)
	if err != nil {
		failError(c, err)
		return
	}
	items := make([]dto.EquipmentResponse, 0, len(list))
	for _, equipment := range list {
		items = append(items, mapEquipment(equipment))
	}
	ok(c, dto.PageResult{List: items, Total: total, Page: page, PageSize: pageSize})
}

// Get 获取设备详情。
func (h *EquipmentHandler) Get(c *gin.Context) {
	id, valid := parseID(c)
	if !valid {
		return
	}
	equipment, err := h.equipmentService.Get(c.Request.Context(), id)
	if err != nil {
		failError(c, err)
		return
	}
	ok(c, mapEquipment(*equipment))
}

// Create 新增设备。
func (h *EquipmentHandler) Create(c *gin.Context) {
	var req dto.CreateEquipmentRequest
	if !bindJSON(c, &req) {
		return
	}
	purchaseDate, err := dto.ParseDate(req.PurchaseDate)
	if err != nil {
		fail(c, 400, 40000, "购买日期格式错误")
		return
	}
	warrantyExpiry, err := dto.ParseDate(req.WarrantyExpiry)
	if err != nil {
		fail(c, 400, 40000, "保修到期日格式错误")
		return
	}
	var purchase *time.Time
	if !purchaseDate.IsZero() {
		purchase = &purchaseDate
	}
	var warranty *time.Time
	if !warrantyExpiry.IsZero() {
		warranty = &warrantyExpiry
	}
	equipment := &model.Equipment{
		Name:           req.Name,
		Code:           req.Code,
		CategoryID:     req.CategoryID,
		BrandModel:     req.BrandModel,
		SerialNumber:   req.SerialNumber,
		PurchaseDate:   purchase,
		PurchasePrice:  req.PurchasePrice,
		Location:       req.Location,
		OwnerID:        req.OwnerID,
		Supplier:       req.Supplier,
		WarrantyExpiry: warranty,
		ImageURL:       req.ImageURL,
	}
	created, err := h.equipmentService.Create(c.Request.Context(), equipment, middleware.GetActor(c))
	if err != nil {
		failError(c, err)
		return
	}
	ok(c, mapEquipment(*created))
}

// Update 更新设备。
func (h *EquipmentHandler) Update(c *gin.Context) {
	id, valid := parseID(c)
	if !valid {
		return
	}
	var req dto.UpdateEquipmentRequest
	if !bindJSON(c, &req) {
		return
	}
	purchaseDate, err := dto.ParseDate(req.PurchaseDate)
	if err != nil {
		fail(c, 400, 40000, "购买日期格式错误")
		return
	}
	warrantyExpiry, err := dto.ParseDate(req.WarrantyExpiry)
	if err != nil {
		fail(c, 400, 40000, "保修到期日格式错误")
		return
	}
	var purchase *time.Time
	if !purchaseDate.IsZero() {
		purchase = &purchaseDate
	}
	var warranty *time.Time
	if !warrantyExpiry.IsZero() {
		warranty = &warrantyExpiry
	}
	equipment := &model.Equipment{
		Name:           req.Name,
		Code:           req.Code,
		CategoryID:     req.CategoryID,
		BrandModel:     req.BrandModel,
		SerialNumber:   req.SerialNumber,
		PurchaseDate:   purchase,
		PurchasePrice:  req.PurchasePrice,
		Location:       req.Location,
		OwnerID:        req.OwnerID,
		Supplier:       req.Supplier,
		WarrantyExpiry: warranty,
		ImageURL:       req.ImageURL,
	}
	updated, err := h.equipmentService.Update(c.Request.Context(), id, equipment, middleware.GetActor(c))
	if err != nil {
		failError(c, err)
		return
	}
	ok(c, mapEquipment(*updated))
}

// UpdateStatus 报废/修改设备状态。
func (h *EquipmentHandler) UpdateStatus(c *gin.Context) {
	id, valid := parseID(c)
	if !valid {
		return
	}
	var req dto.UpdateEquipmentStatusRequest
	if !bindJSON(c, &req) {
		return
	}
	status := constants.AssetStatus(req.Status)
	if !status.Valid() {
		fail(c, 400, 40000, "资产状态无效")
		return
	}
	if status == constants.AssetStatusRetired {
		if err := h.equipmentService.Retire(c.Request.Context(), id, middleware.GetActor(c)); err != nil {
			failError(c, err)
			return
		}
		ok(c, nil)
		return
	}
	fail(c, 400, 40000, "仅支持 Retired 状态变更")
}

// TransferOwner 转移责任人。
func (h *EquipmentHandler) TransferOwner(c *gin.Context) {
	id, valid := parseID(c)
	if !valid {
		return
	}
	var req dto.TransferOwnerRequest
	if !bindJSON(c, &req) {
		return
	}
	if err := h.equipmentService.TransferOwner(c.Request.Context(), id, req.OwnerID, middleware.GetActor(c)); err != nil {
		failError(c, err)
		return
	}
	ok(c, nil)
}
