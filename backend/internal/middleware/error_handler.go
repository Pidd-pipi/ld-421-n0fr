package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labequipment/lab-equipment/internal/constants"
)

// ErrorHandler 捕获 panic，统一转换为 JSON 响应。
func ErrorHandler(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Error("panic recovered", "panic", rec, "path", c.Request.URL.Path)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    constants.CodeInternal,
					"message": "服务器内部错误",
					"data":    nil,
				})
			}
		}()
		c.Next()
	}
}
