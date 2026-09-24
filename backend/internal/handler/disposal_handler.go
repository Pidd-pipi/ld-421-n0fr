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

// DisposalHandler 设备处置（报废）审批处理器。
type DisposalHandler struct {
	disposalService *service.DisposalService
}

// NewDisposalHandler 构造处置审批处理器。
func NewDisposalHandler(disposalService *service.DisposalService) *DisposalHandler {
	return &DisposalHandler{disposalService: disposalService}
}

// Submit 提交设备报废处置申请。
func (h *DisposalHandler) Submit(c *gin.Context) {
	id, valid := parseID(c)
	if !valid {
		return
	}
	var req dto.SubmitDisposalRequest
	if !bindJSON(c, &req) {
		return
	}
	approval := &model.DisposalApproval{
		EquipmentID: id,
		Reason:      req.Reason,
	}
	created, err := h.disposalService.Submit(c.Request.Context(), approval, middleware.GetActor(c))
	if err != nil {
		failError(c, err)
		return
	}
	ok(c, mapDisposal(*created))
}

// Review 处置审批（通过报废或因阻塞退回）。
func (h *DisposalHandler) Review(c *gin.Context) {
	id, valid := parseID(c)
	if !valid {
		return
	}
	approval, result, err := h.disposalService.Review(c.Request.Context(), id, middleware.GetActor(c))
	if err != nil {
		failError(c, err)
		return
	}
	ok(c, dto.DisposalReviewResponse{
		Approval: mapDisposal(*approval),
		Approved: result.Approved,
		Blockers: mapDisposalBlockers(result.Blockers),
	})
}

// Get 获取处置审批详情。
func (h *DisposalHandler) Get(c *gin.Context) {
	id, valid := parseID(c)
	if !valid {
		return
	}
	approval, err := h.disposalService.Get(c.Request.Context(), id)
	if err != nil {
		failError(c, err)
		return
	}
	ok(c, mapDisposal(*approval))
}

// GetByEquipment 获取设备的处置审批进度（无记录时 data 为 null）。
func (h *DisposalHandler) GetByEquipment(c *gin.Context) {
	id, valid := parseID(c)
	if !valid {
		return
	}
	approval, err := h.disposalService.GetByEquipment(c.Request.Context(), id)
	if err != nil {
		failError(c, err)
		return
	}
	if approval == nil {
		ok(c, nil)
		return
	}
	response := mapDisposal(*approval)
	if approval.Active && approval.Status == constants.DisposalStatusPending {
		blockers, err := h.disposalService.LiveBlockers(c.Request.Context(), id)
		if err != nil {
			failError(c, err)
			return
		}
		response.Blockers = mapDisposalBlockers(blockers)
	}
	ok(c, response)
}

// List 分页查询处置审批。
func (h *DisposalHandler) List(c *gin.Context) {
	page, pageSize := parsePage(c)
	equipmentID, _ := strconv.ParseUint(c.Query("equipment_id"), 10, 64)
	filter := repository.DisposalFilter{
		Status:      constants.DisposalStatus(c.Query("status")),
		EquipmentID: uint(equipmentID),
		Pagination:  repository.Pagination{Page: page, PageSize: pageSize},
	}
	list, total, err := h.disposalService.List(c.Request.Context(), filter)
	if err != nil {
		failError(c, err)
		return
	}
	items := make([]dto.DisposalResponse, 0, len(list))
	for _, approval := range list {
		items = append(items, mapDisposal(approval))
	}
	ok(c, dto.PageResult{List: items, Total: total, Page: page, PageSize: pageSize})
}
