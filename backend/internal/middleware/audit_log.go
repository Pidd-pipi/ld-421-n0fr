package middleware

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/labequipment/lab-equipment/internal/service"
)

// AuditLog 记录所有写操作请求到审计日志（业务服务另有语义化埋点）。
func AuditLog(audit *service.AuditService, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			return
		}
		actor := GetActor(c)
		action := fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path)
		if err := audit.Log(c.Request.Context(), actor, action, "http", 0, "请求入口"); err != nil {
			logger.Error("audit middleware log failed", "error", err)
		}
	}
}
