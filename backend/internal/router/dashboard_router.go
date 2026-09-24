package router

import (
	"github.com/gin-gonic/gin"
)

func registerDashboardRoutes(group *gin.RouterGroup, deps Dependencies) {
	dashboard := group.Group("/dashboard")
	{
		dashboard.GET("/stats", deps.DashboardHandler.Stats)
	}
}
