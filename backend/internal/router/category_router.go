package router

import (
	"github.com/gin-gonic/gin"
	"github.com/labequipment/lab-equipment/internal/middleware"
)

func registerCategoryRoutes(group *gin.RouterGroup, deps Dependencies) {
	categories := group.Group("/categories")
	{
		categories.GET("", deps.CategoryHandler.List)
	}
	admin := group.Group("/categories")
	admin.Use(middleware.RequireRoles("Admin", "LabManager"))
	{
		admin.POST("", deps.CategoryHandler.Create)
		admin.PUT("/:id", deps.CategoryHandler.Update)
		admin.DELETE("/:id", deps.CategoryHandler.Delete)
	}
}
