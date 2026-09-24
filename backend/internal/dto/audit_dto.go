package dto

import "time"

// AuditLogResponse 审计日志返回。
type AuditLogResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"userId"`
	UserName     string    `json:"userName"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resourceType"`
	ResourceID   uint      `json:"resourceId"`
	Detail       string    `json:"detail"`
	IP           string    `json:"ip"`
	CreatedAt    time.Time `json:"createdAt"`
}
