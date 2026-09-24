package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/labequipment/lab-equipment/internal/dto"
	"github.com/labequipment/lab-equipment/internal/middleware"
	"github.com/labequipment/lab-equipment/internal/service"
)

// AuthHandler 认证处理器。
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler 构造认证处理器。
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login 登录并签发 token。
// @Summary 登录
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.LoginRequest true "登录请求"
// @Success 200 {object} response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if !bindJSON(c, &req) {
		return
	}
	user, token, err := h.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		failError(c, err)
		return
	}
	ok(c, dto.LoginResponse{Token: token, User: mapUser(*user)})
}

// Me 返回当前登录用户信息。
func (h *AuthHandler) Me(c *gin.Context) {
	claims := middleware.GetClaims(c)
	ok(c, gin.H{"userId": claims.UserID, "username": claims.Username, "role": claims.Role})
}
