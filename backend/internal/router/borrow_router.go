package router

import (
	"github.com/gin-gonic/gin"
	"github.com/labequipment/lab-equipment/internal/middleware"
)

func registerBorrowRoutes(group *gin.RouterGroup, deps Dependencies) {
	borrow := group.Group("/borrows")
	{
		borrow.GET("", deps.BorrowHandler.List)
		borrow.GET("/:id", deps.BorrowHandler.Get)
		borrow.POST("", deps.BorrowHandler.Create)
	}
	review := group.Group("/borrows")
	review.Use(middleware.RequireRoles("Admin", "LabManager"))
	{
		review.POST("/:id/approve", deps.BorrowHandler.Approve)
		review.POST("/:id/reject", deps.BorrowHandler.Reject)
		review.POST("/:id/return", deps.BorrowHandler.Return)
	}
}
