package router

import (
	"github.com/gin-gonic/gin"
	"github.com/labequipment/lab-equipment/internal/middleware"
)

func registerReservationRoutes(group *gin.RouterGroup, deps Dependencies) {
	reservation := group.Group("/reservations")
	{
		reservation.GET("", deps.ReservationHandler.List)
		reservation.GET("/:id", deps.ReservationHandler.Get)
		reservation.POST("", deps.ReservationHandler.Create)
		reservation.POST("/:id/cancel", deps.ReservationHandler.Cancel)
	}
	review := group.Group("/reservations")
	review.Use(middleware.RequireRoles("Admin", "LabManager"))
	{
		review.POST("/:id/approve", deps.ReservationHandler.Approve)
		review.POST("/:id/reject", deps.ReservationHandler.Reject)
	}
}
