package router

import (
	"github.com/gin-gonic/gin"
	"github.com/labequipment/lab-equipment/internal/middleware"
)

func registerDisposalRoutes(group *gin.RouterGroup, deps Dependencies) {
	disposal := group.Group("/disposals")
	{
		disposal.GET("", deps.DisposalHandler.List)
		disposal.GET("/:id", deps.DisposalHandler.Get)
	}
	// 查询设备的处置审批进度：/equipment/:id/disposal
	progress := group.Group("/equipment")
	{
		progress.GET("/:id/disposal", deps.DisposalHandler.GetByEquipment)
	}
	manage := group.Group("/equipment")
	manage.Use(middleware.RequireRoles("Admin", "LabManager"))
	{
		manage.POST("/:id/disposal", deps.DisposalHandler.Submit)
	}
	review := group.Group("/disposals")
	review.Use(middleware.RequireRoles("Admin", "LabManager"))
	{
		review.POST("/:id/review", deps.DisposalHandler.Review)
	}
}
