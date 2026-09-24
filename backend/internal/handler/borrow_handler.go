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

// BorrowHandler 借用记录处理器。
type BorrowHandler struct {
	borrowService *service.BorrowService
}

// NewBorrowHandler 构造借用处理器。
func NewBorrowHandler(borrowService *service.BorrowService) *BorrowHandler {
	return &BorrowHandler{borrowService: borrowService}
}

// List 分页查询借用记录。
func (h *BorrowHandler) List(c *gin.Context) {
	page, pageSize := parsePage(c)
	equipmentID, _ := strconv.ParseUint(c.Query("equipment_id"), 10, 64)
	borrowerID, _ := strconv.ParseUint(c.Query("borrower_id"), 10, 64)
	filter := repository.BorrowFilter{
		Status:      constants.BorrowStatus(c.Query("status")),
		EquipmentID: uint(equipmentID),
		BorrowerID:  uint(borrowerID),
		Pagination:  repository.Pagination{Page: page, PageSize: pageSize},
	}
	list, total, err := h.borrowService.List(c.Request.Context(), filter)
	if err != nil {
		failError(c, err)
		return
	}
	items := make([]dto.BorrowResponse, 0, len(list))
	for _, record := range list {
		items = append(items, mapBorrow(record))
	}
	ok(c, dto.PageResult{List: items, Total: total, Page: page, PageSize: pageSize})
}

// Get 获取借用详情。
func (h *BorrowHandler) Get(c *gin.Context) {
	id, valid := parseID(c)
	if !valid {
		return
	}
	record, err := h.borrowService.Get(c.Request.Context(), id)
	if err != nil {
		failError(c, err)
		return
	}
	ok(c, mapBorrow(*record))
}

// Create 提交借用申请。
func (h *BorrowHandler) Create(c *gin.Context) {
	var req dto.CreateBorrowRequest
	if !bindJSON(c, &req) {
		return
	}
	borrowDate, err := dto.ParseDate(req.BorrowDate)
	if err != nil {
		fail(c, 400, 40000, "借用日期格式错误")
		return
	}
	expectedReturnDate, err := dto.ParseDate(req.ExpectedReturnDate)
	if err != nil {
		fail(c, 400, 40000, "预计归还日期格式错误")
		return
	}
	record := &model.BorrowRecord{
		EquipmentID:        req.EquipmentID,
		BorrowDate:         borrowDate,
		ExpectedReturnDate: expectedReturnDate,
		Reason:             req.Reason,
	}
	created, err := h.borrowService.Create(c.Request.Context(), record, middleware.GetActor(c))
	if err != nil {
		failError(c, err)
		return
	}
	ok(c, mapBorrow(*created))
}

// Approve 审批通过。
func (h *BorrowHandler) Approve(c *gin.Context) {
	id, valid := parseID(c)
	if !valid {
		return
	}
	if err := h.borrowService.Approve(c.Request.Context(), id, middleware.GetActor(c)); err != nil {
		failError(c, err)
		return
	}
	ok(c, nil)
}

// Reject 驳回申请。
func (h *BorrowHandler) Reject(c *gin.Context) {
	id, valid := parseID(c)
	if !valid {
		return
	}
	if err := h.borrowService.Reject(c.Request.Context(), id, middleware.GetActor(c)); err != nil {
		failError(c, err)
		return
	}
	ok(c, nil)
}

// Return 确认归还。
func (h *BorrowHandler) Return(c *gin.Context) {
	id, valid := parseID(c)
	if !valid {
		return
	}
	var req dto.ReturnBorrowRequest
	if !bindJSON(c, &req) {
		return
	}
	actualReturnDate, err := dto.ParseDate(req.ActualReturnDate)
	if err != nil {
		fail(c, 400, 40000, "实际归还日期格式错误")
		return
	}
	condition := constants.ReturnCondition(req.ReturnCondition)
	if !condition.Valid() {
		fail(c, 400, 40000, "归还状况无效")
		return
	}
	if err := h.borrowService.Return(c.Request.Context(), id, actualReturnDate, condition, middleware.GetActor(c)); err != nil {
		failError(c, err)
		return
	}
	ok(c, nil)
}
