package router

import (
	"github.com/gin-gonic/gin"
)

func registerAuthRoutes(group *gin.RouterGroup, deps Dependencies) {
	auth := group.Group("/auth")
	{
		auth.GET("/me", deps.AuthHandler.Me)
	}
}
