package handler

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/labequipment/lab-equipment/internal/dto"
	"github.com/labequipment/lab-equipment/internal/middleware"
	"github.com/labequipment/lab-equipment/internal/service"
)

// DisposalHandler 报废处置审批处理器。
type DisposalHandler struct {
	disposalService *service.DisposalService
}

// NewDisposalHandler 构造报废处置审批处理器。
func NewDisposalHandler(disposalService *service.DisposalService) *DisposalHandler {
	return &DisposalHandler{disposalService: disposalService}
}

// Submit 提交报废申请。
func (h *DisposalHandler) Submit(c *gin.Context) {
	id, valid := parseID(c)
	if !valid {
		return
	}
	var req dto.SubmitDisposalRequest
	if !bindJSON(c, &req) {
		return
	}
	reason := strings.TrimSpace(req.Reason)
	if reason == "" {
		fail(c, 400, 40000, "报废原因不能为空")
		return
	}
	item, err := h.disposalService.Submit(c.Request.Context(), id, reason, middleware.GetActor(c))
	if err != nil {
		failError(c, err)
		return
	}
	ok(c, mapDisposalItem(*item))
}

// Approve 审批报废申请：无阻塞项则通过并报废设备，否则列出阻塞项并退回。
func (h *DisposalHandler) Approve(c *gin.Context) {
	id, valid := parseID(c)
	if !valid {
		return
	}
	disposalID, valid := parseParamID(c, "disposalId")
	if !valid {
		return
	}
	item, err := h.disposalService.Approve(c.Request.Context(), id, disposalID, middleware.GetActor(c))
	if err != nil {
		failError(c, err)
		return
	}
	ok(c, mapDisposalItem(*item))
}

// List 查询设备的报废申请列表（含审批进度、原因与阻塞项）。
func (h *DisposalHandler) List(c *gin.Context) {
	id, valid := parseID(c)
	if !valid {
		return
	}
	items, err := h.disposalService.ListByEquipment(c.Request.Context(), id)
	if err != nil {
		failError(c, err)
		return
	}
	result := make([]dto.DisposalResponse, 0, len(items))
	for _, item := range items {
		result = append(result, mapDisposalItem(item))
	}
	ok(c, result)
}
