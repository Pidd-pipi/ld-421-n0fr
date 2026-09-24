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

// MaintenanceHandler 维护记录处理器。
type MaintenanceHandler struct {
	maintenanceService *service.MaintenanceService
}

// NewMaintenanceHandler 构造维护处理器。
func NewMaintenanceHandler(maintenanceService *service.MaintenanceService) *MaintenanceHandler {
	return &MaintenanceHandler{maintenanceService: maintenanceService}
}

// List 分页查询维护记录。
func (h *MaintenanceHandler) List(c *gin.Context) {
	page, pageSize := parsePage(c)
	equipmentID, _ := strconv.ParseUint(c.Query("equipment_id"), 10, 64)
	filter := repository.MaintenanceFilter{
		EquipmentID: uint(equipmentID),
		Pagination:  repository.Pagination{Page: page, PageSize: pageSize},
	}
	list, total, err := h.maintenanceService.List(c.Request.Context(), filter)
	if err != nil {
		failError(c, err)
		return
	}
	items := make([]dto.MaintenanceResponse, 0, len(list))
	for _, record := range list {
		items = append(items, mapMaintenance(record))
	}
	ok(c, dto.PageResult{List: items, Total: total, Page: page, PageSize: pageSize})
}

// Stats 维护统计。
func (h *MaintenanceHandler) Stats(c *gin.Context) {
	stats, err := h.maintenanceService.Stats(c.Request.Context())
	if err != nil {
		failError(c, err)
		return
	}
	resultCount := make(map[string]int64, len(stats.ResultCount))
	for k, v := range stats.ResultCount {
		resultCount[string(k)] = v
	}
	ok(c, dto.MaintenanceStatsResponse{
		TotalRecords: stats.TotalRecords,
		TotalCost:    stats.TotalCost,
		ResultCount:  resultCount,
	})
}

// Create 创建维护记录。
func (h *MaintenanceHandler) Create(c *gin.Context) {
	var req dto.CreateMaintenanceRequest
	if !bindJSON(c, &req) {
		return
	}
	maintenanceDate, err := dto.ParseDate(req.MaintenanceDate)
	if err != nil {
		fail(c, 400, 40000, "维护日期格式错误")
		return
	}
	var nextDate *time.Time
	if req.NextMaintenanceDate != "" {
		t, err := dto.ParseDate(req.NextMaintenanceDate)
		if err != nil {
			fail(c, 400, 40000, "下次维护日期格式错误")
			return
		}
		nextDate = &t
	}
	record := &model.MaintenanceRecord{
		EquipmentID:         req.EquipmentID,
		Type:                constants.MaintenanceType(req.Type),
		Content:             req.Content,
		MaintenanceDate:     maintenanceDate,
		NextMaintenanceDate: nextDate,
		Cost:                req.Cost,
		MaintainerID:        req.MaintainerID,
	}
	created, err := h.maintenanceService.Create(c.Request.Context(), record, middleware.GetActor(c))
	if err != nil {
		failError(c, err)
		return
	}
	ok(c, mapMaintenance(*created))
}

// Execute 执行维护并记录结果。
func (h *MaintenanceHandler) Execute(c *gin.Context) {
	id, valid := parseID(c)
	if !valid {
		return
	}
	var req dto.ExecuteMaintenanceRequest
	if !bindJSON(c, &req) {
		return
	}
	result := constants.MaintenanceResult(req.Result)
	if err := h.maintenanceService.Execute(c.Request.Context(), id, result, middleware.GetActor(c)); err != nil {
		failError(c, err)
		return
	}
	ok(c, nil)
}
