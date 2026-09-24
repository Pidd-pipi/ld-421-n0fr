package repository

import (
	"context"
	"fmt"

	"github.com/labequipment/lab-equipment/internal/model"
	"gorm.io/gorm"
)

// AuditLogRepository 审计日志仓储接口。
type AuditLogRepository interface {
	Create(ctx context.Context, log *model.AuditLog) error
	List(ctx context.Context, filter AuditLogFilter) ([]model.AuditLog, int64, error)
}

// AuditLogFilter 审计日志筛选条件。
type AuditLogFilter struct {
	Keyword string
	Pagination
}

type auditLogRepository struct {
	db *gorm.DB
}

// NewAuditLogRepository 构造审计日志仓储。
func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(ctx context.Context, log *model.AuditLog) error {
	if err := r.db.WithContext(ctx).Create(log).Error; err != nil {
		return fmt.Errorf("create audit log: %w", err)
	}
	return nil
}

func (r *auditLogRepository) List(ctx context.Context, filter AuditLogFilter) ([]model.AuditLog, int64, error) {
	filter.Normalize()
	query := r.db.WithContext(ctx).Model(&model.AuditLog{})
	if filter.Keyword != "" {
		like := "%" + filter.Keyword + "%"
		query = query.Where("action LIKE ? OR user_name LIKE ? OR resource_type LIKE ?", like, like, like)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count audit logs: %w", err)
	}
	var logs []model.AuditLog
	if err := query.Order("id DESC").Offset(filter.Offset()).Limit(filter.PageSize).Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("list audit logs: %w", err)
	}
	return logs, total, nil
}
