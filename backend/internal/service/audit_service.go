package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/labequipment/lab-equipment/internal/model"
	"github.com/labequipment/lab-equipment/internal/repository"
)

// AuditService 负责写入操作审计日志。
type AuditService struct {
	repo   repository.AuditLogRepository
	logger *slog.Logger
}

// NewAuditService 构造审计服务。
func NewAuditService(repo repository.AuditLogRepository, logger *slog.Logger) *AuditService {
	return &AuditService{repo: repo, logger: logger}
}

// Log 记录一条审计日志，失败时透传错误。
func (s *AuditService) Log(ctx context.Context, actor Actor, action, resourceType string, resourceID uint, detail string) error {
	entry := &model.AuditLog{
		UserID:       actor.UserID,
		UserName:     actor.Username,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Detail:       detail,
		IP:           actor.IP,
	}
	if err := s.repo.Create(ctx, entry); err != nil {
		return fmt.Errorf("audit %s: %w", action, err)
	}
	return nil
}

// List 分页查询审计日志。
func (s *AuditService) List(ctx context.Context, filter repository.AuditLogFilter) ([]model.AuditLog, int64, error) {
	items, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("list audit logs: %w", err)
	}
	return items, total, nil
}
