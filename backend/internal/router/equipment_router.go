package router

import (
	"github.com/gin-gonic/gin"
	"github.com/labequipment/lab-equipment/internal/middleware"
)

func registerEquipmentRoutes(group *gin.RouterGroup, deps Dependencies) {
	equipment := group.Group("/equipment")
	{
		equipment.GET("", deps.EquipmentHandler.List)
		equipment.GET("/:id", deps.EquipmentHandler.Get)
	}
	manage := group.Group("/equipment")
	manage.Use(middleware.RequireRoles("Admin", "LabManager"))
	{
		manage.POST("", deps.EquipmentHandler.Create)
		manage.PUT("/:id", deps.EquipmentHandler.Update)
		manage.PATCH("/:id/status", deps.EquipmentHandler.UpdateStatus)
		manage.PATCH("/:id/owner", deps.EquipmentHandler.TransferOwner)
	}
}
