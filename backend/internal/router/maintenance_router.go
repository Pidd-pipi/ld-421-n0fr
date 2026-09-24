package router

import (
	"github.com/gin-gonic/gin"
	"github.com/labequipment/lab-equipment/internal/middleware"
)

func registerMaintenanceRoutes(group *gin.RouterGroup, deps Dependencies) {
	maintenance := group.Group("/maintenance")
	{
		maintenance.GET("", deps.MaintenanceHandler.List)
		maintenance.GET("/stats", deps.MaintenanceHandler.Stats)
	}
	write := group.Group("/maintenance")
	write.Use(middleware.RequireRoles("Admin", "LabManager", "Researcher"))
	{
		write.POST("", deps.MaintenanceHandler.Create)
		write.POST("/:id/execute", deps.MaintenanceHandler.Execute)
	}
}
