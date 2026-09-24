package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/labequipment/lab-equipment/internal/constants"
	"github.com/labequipment/lab-equipment/internal/dto"
	"github.com/labequipment/lab-equipment/internal/middleware"
	"github.com/labequipment/lab-equipment/internal/model"
	"github.com/labequipment/lab-equipment/internal/repository"
	"github.com/labequipment/lab-equipment/internal/service"
)

// ReservationHandler 预约记录处理器。
type ReservationHandler struct {
	reservationService *service.ReservationService
}

// NewReservationHandler 构造预约处理器。
func NewReservationHandler(reservationService *service.ReservationService) *ReservationHandler {
	return &ReservationHandler{reservationService: reservationService}
}

// List 分页查询预约。
func (h *ReservationHandler) List(c *gin.Context) {
	page, pageSize := parsePage(c)
	equipmentID, _ := strconv.ParseUint(c.Query("equipment_id"), 10, 64)
	userID, _ := strconv.ParseUint(c.Query("user_id"), 10, 64)
	filter := repository.ReservationFilter{
		EquipmentID: uint(equipmentID),
		UserID:      uint(userID),
		Status:      constants.ReservationStatus(c.Query("status")),
		Pagination:  repository.Pagination{Page: page, PageSize: pageSize},
	}
	list, total, err := h.reservationService.List(c.Request.Context(), filter)
	if err != nil {
		failError(c, err)
		return
	}
	items := make([]dto.ReservationResponse, 0, len(list))
	for _, reservation := range list {
		items = append(items, mapReservation(reservation))
	}
	ok(c, dto.PageResult{List: items, Total: total, Page: page, PageSize: pageSize})
}

// Get 获取预约详情。
func (h *ReservationHandler) Get(c *gin.Context) {
	id, valid := parseID(c)
	if !valid {
		return
	}
	reservation, err := h.reservationService.Get(c.Request.Context(), id)
	if err != nil {
		failError(c, err)
		return
	}
	ok(c, mapReservation(*reservation))
}

// Create 创建预约。
func (h *ReservationHandler) Create(c *gin.Context) {
	var req dto.CreateReservationRequest
	if !bindJSON(c, &req) {
		return
	}
	startTime, err := dto.ParseDateTime(req.StartTime)
	if err != nil {
		fail(c, 400, 40000, "预约开始时间格式错误")
		return
	}
	endTime, err := dto.ParseDateTime(req.EndTime)
	if err != nil {
		fail(c, 400, 40000, "预约结束时间格式错误")
		return
	}
	reservation := &model.Reservation{
		EquipmentID: req.EquipmentID,
		StartTime:   startTime,
		EndTime:     endTime,
		Purpose:     req.Purpose,
	}
	created, err := h.reservationService.Create(c.Request.Context(), reservation, middleware.GetActor(c))
	if err != nil {
		failError(c, err)
		return
	}
	ok(c, mapReservation(*created))
}

// Approve 审批通过。
func (h *ReservationHandler) Approve(c *gin.Context) {
	id, valid := parseID(c)
	if !valid {
		return
	}
	if err := h.reservationService.Approve(c.Request.Context(), id, middleware.GetActor(c)); err != nil {
		failError(c, err)
		return
	}
	ok(c, nil)
}

// Reject 驳回预约。
func (h *ReservationHandler) Reject(c *gin.Context) {
	id, valid := parseID(c)
	if !valid {
		return
	}
	if err := h.reservationService.Reject(c.Request.Context(), id, middleware.GetActor(c)); err != nil {
		failError(c, err)
		return
	}
	ok(c, nil)
}

// Cancel 取消预约。
func (h *ReservationHandler) Cancel(c *gin.Context) {
	id, valid := parseID(c)
	if !valid {
		return
	}
	if err := h.reservationService.Cancel(c.Request.Context(), id, middleware.GetActor(c)); err != nil {
		failError(c, err)
		return
	}
	ok(c, nil)
}
