package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/labequipment/lab-equipment/internal/dto"
	"github.com/labequipment/lab-equipment/internal/repository"
	"github.com/labequipment/lab-equipment/internal/service"
)

// AuditHandler 审计日志处理器。
type AuditHandler struct {
	auditService *service.AuditService
}

// NewAuditHandler 构造审计处理器。
func NewAuditHandler(auditService *service.AuditService) *AuditHandler {
	return &AuditHandler{auditService: auditService}
}

// List 分页查询审计日志。
func (h *AuditHandler) List(c *gin.Context) {
	page, pageSize := parsePage(c)
	filter := repository.AuditLogFilter{
		Keyword:    c.Query("keyword"),
		Pagination: repository.Pagination{Page: page, PageSize: pageSize},
	}
	logs, total, err := h.auditService.List(c.Request.Context(), filter)
	if err != nil {
		failError(c, err)
		return
	}
	items := make([]dto.AuditLogResponse, 0, len(logs))
	for _, log := range logs {
		items = append(items, mapAudit(log))
	}
	ok(c, dto.PageResult{List: items, Total: total, Page: page, PageSize: pageSize})
}
