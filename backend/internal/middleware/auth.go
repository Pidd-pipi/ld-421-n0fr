package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/labequipment/lab-equipment/internal/constants"
	apperrors "github.com/labequipment/lab-equipment/internal/errors"
	"github.com/labequipment/lab-equipment/internal/service"
)

const claimsKey = "authClaims"

// Auth 解析并校验 JWT，将 Claims 注入上下文。
func Auth(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    constants.CodeUnauthorized,
				"message": "缺少认证信息",
				"data":    nil,
			})
			return
		}
		token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		claims, err := authService.ParseToken(token)
		if err != nil {
			code := constants.CodeUnauthorized
			msg := "token 无效或已过期"
			var be *apperrors.BusinessError
			if ok := asBusinessError(err, &be); ok {
				code = be.Code
				msg = be.Message
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    code,
				"message": msg,
				"data":    nil,
			})
			return
		}
		c.Set(claimsKey, claims)
		c.Next()
	}
}

// GetClaims 从上下文读取 JWT Claims。
func GetClaims(c *gin.Context) *service.Claims {
	if v, ok := c.Get(claimsKey); ok {
		if claims, ok := v.(*service.Claims); ok {
			return claims
		}
	}
	return &service.Claims{}
}

// GetActor 从上下文构造操作者。
func GetActor(c *gin.Context) service.Actor {
	claims := GetClaims(c)
	return service.Actor{
		UserID:   claims.UserID,
		Username: claims.Username,
		Role:     claims.Role,
		IP:       c.ClientIP(),
	}
}

func asBusinessError(err error, target **apperrors.BusinessError) bool {
	for err != nil {
		if be, ok := err.(*apperrors.BusinessError); ok {
			*target = be
			return true
		}
		unwrapper, ok := err.(interface{ Unwrap() error })
		if !ok {
			return false
		}
		err = unwrapper.Unwrap()
	}
	return false
}
