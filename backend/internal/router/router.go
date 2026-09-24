package router

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labequipment/lab-equipment/internal/handler"
	"github.com/labequipment/lab-equipment/internal/middleware"
	"github.com/labequipment/lab-equipment/internal/service"
)

// Dependencies 路由装配所需依赖。
type Dependencies struct {
	AuthHandler        *handler.AuthHandler
	UserHandler        *handler.UserHandler
	CategoryHandler    *handler.CategoryHandler
	EquipmentHandler   *handler.EquipmentHandler
	BorrowHandler      *handler.BorrowHandler
	MaintenanceHandler *handler.MaintenanceHandler
	ReservationHandler *handler.ReservationHandler
	DashboardHandler   *handler.DashboardHandler
	AuditHandler       *handler.AuditHandler
	AuthService        *service.AuthService
	AuditService       *service.AuditService
	Logger             *slog.Logger
}

// NewRouter 装配全部路由。
func NewRouter(deps Dependencies) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middleware.RequestLogger(deps.Logger))
	r.Use(middleware.ErrorHandler(deps.Logger))

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	public := r.Group("/api/v1")
	{
		public.POST("/auth/login", deps.AuthHandler.Login)
	}

	authed := r.Group("/api/v1")
	authed.Use(middleware.Auth(deps.AuthService))
	authed.Use(middleware.AuditLog(deps.AuditService, deps.Logger))
	{
		registerAuthRoutes(authed, deps)
		registerUserRoutes(authed, deps)
		registerCategoryRoutes(authed, deps)
		registerEquipmentRoutes(authed, deps)
		registerBorrowRoutes(authed, deps)
		registerMaintenanceRoutes(authed, deps)
		registerReservationRoutes(authed, deps)
		registerDashboardRoutes(authed, deps)
		registerAuditRoutes(authed, deps)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": 40400, "message": "接口不存在", "data": nil})
	})
	return r
}
