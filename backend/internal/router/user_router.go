package router

import (
	"github.com/gin-gonic/gin"
	"github.com/labequipment/lab-equipment/internal/middleware"
)

func registerUserRoutes(group *gin.RouterGroup, deps Dependencies) {
	users := group.Group("/users")
	users.Use(middleware.RequireRoles("Admin", "LabManager", "Researcher"))
	{
		users.GET("", deps.UserHandler.List)
	}
}
