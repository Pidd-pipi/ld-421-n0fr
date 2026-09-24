package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labequipment/lab-equipment/internal/constants"
)

// RequireRoles 校验当前用户角色，任一匹配即可通过。
func RequireRoles(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}
	return func(c *gin.Context) {
		claims := GetClaims(c)
		if _, ok := allowed[claims.Role]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    constants.CodeForbidden,
				"message": "无权限执行此操作",
				"data":    nil,
			})
			return
		}
		c.Next()
	}
}
