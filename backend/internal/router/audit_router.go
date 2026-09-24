package router

import (
	"github.com/gin-gonic/gin"
	"github.com/labequipment/lab-equipment/internal/middleware"
)

func registerAuditRoutes(group *gin.RouterGroup, deps Dependencies) {
	audit := group.Group("/audit-logs")
	audit.Use(middleware.RequireRoles("Admin"))
	{
		audit.GET("", deps.AuditHandler.List)
	}
}
