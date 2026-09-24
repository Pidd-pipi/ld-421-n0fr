package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/labequipment/lab-equipment/internal/model"
	"github.com/labequipment/lab-equipment/internal/repository"
)

// UserService 用户服务。
type UserService struct {
	userRepo repository.UserRepository
	logger   *slog.Logger
}

// NewUserService 构造用户服务。
func NewUserService(userRepo repository.UserRepository, logger *slog.Logger) *UserService {
	return &UserService{userRepo: userRepo, logger: logger}
}

// List 分页查询用户。
func (s *UserService) List(ctx context.Context, filter repository.UserFilter) ([]model.User, int64, error) {
	users, total, err := s.userRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	return users, total, nil
}
