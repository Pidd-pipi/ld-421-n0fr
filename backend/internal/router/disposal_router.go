package router

import (
	"github.com/gin-gonic/gin"
	"github.com/labequipment/lab-equipment/internal/middleware"
)

func registerDisposalRoutes(group *gin.RouterGroup, deps Dependencies) {
	group.GET("/equipment/:id/disposals", deps.DisposalHandler.List)
	manage := group.Group("/equipment")
	manage.Use(middleware.RequireRoles("Admin", "LabManager"))
	{
		manage.POST("/:id/disposals", deps.DisposalHandler.Submit)
		manage.POST("/:id/disposals/:disposalId/approve", deps.DisposalHandler.Approve)
	}
}
