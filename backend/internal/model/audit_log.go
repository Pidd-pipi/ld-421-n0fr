package model

// AuditLog 操作审计日志实体。
type AuditLog struct {
	Base
	UserID       uint   `gorm:"index" json:"userId"`
	UserName     string `gorm:"size:64" json:"userName"`
	Action       string `gorm:"size:128;index" json:"action"`
	ResourceType string `gorm:"size:64" json:"resourceType"`
	ResourceID   uint   `gorm:"index" json:"resourceId"`
	Detail       string `gorm:"size:1024" json:"detail"`
	IP           string `gorm:"size:64" json:"ip"`
}

func (AuditLog) TableName() string { return "audit_logs" }
