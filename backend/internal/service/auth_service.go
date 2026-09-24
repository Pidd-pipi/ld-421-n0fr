package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labequipment/lab-equipment/internal/config"
	apperrors "github.com/labequipment/lab-equipment/internal/errors"
	"github.com/labequipment/lab-equipment/internal/model"
	"github.com/labequipment/lab-equipment/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// Claims JWT 载荷。
type Claims struct {
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// AuthService 认证服务。
type AuthService struct {
	userRepo repository.UserRepository
	jwtCfg   config.JWTConfig
	logger   *slog.Logger
}

// NewAuthService 构造认证服务。
func NewAuthService(userRepo repository.UserRepository, jwtCfg config.JWTConfig, logger *slog.Logger) *AuthService {
	return &AuthService{userRepo: userRepo, jwtCfg: jwtCfg, logger: logger}
}

// Login 校验用户名密码并签发 token。
func (s *AuthService) Login(ctx context.Context, username, password string) (*model.User, string, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, "", apperrors.NewBusinessError(40100, 401, "用户名或密码错误")
	}
	if err != nil {
		return nil, "", fmt.Errorf("login find user: %w", err)
	}
	if !user.Active {
		return nil, "", apperrors.NewBusinessError(40300, 403, "用户已被禁用")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", apperrors.NewBusinessError(40100, 401, "用户名或密码错误")
	}
	roleCode := ""
	if user.Role != nil {
		roleCode = user.Role.Code
	}
	now := time.Now()
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     roleCode,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.jwtCfg.ExpireDuration())),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "lab-equipment",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.jwtCfg.Secret))
	if err != nil {
		return nil, "", fmt.Errorf("sign token: %w", err)
	}
	return user, signed, nil
}

// ParseToken 解析并校验 JWT。
func (s *AuthService) ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.jwtCfg.Secret), nil
	})
	if err != nil || !token.Valid {
		return nil, apperrors.NewBusinessError(40100, 401, "token 无效或已过期")
	}
	return claims, nil
}
