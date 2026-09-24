package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labequipment/lab-equipment/database/migrations"
	"github.com/labequipment/lab-equipment/database/seeds"
	"github.com/labequipment/lab-equipment/internal/config"
	"github.com/labequipment/lab-equipment/internal/handler"
	"github.com/labequipment/lab-equipment/internal/logger"
	"github.com/labequipment/lab-equipment/internal/repository"
	"github.com/labequipment/lab-equipment/internal/router"
	"github.com/labequipment/lab-equipment/internal/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	log := logger.New(cfg.AppEnv)
	if err := run(cfg, log); err != nil {
		log.Error("server exited", "error", err)
		os.Exit(1)
	}
}

func connectDB(cfg *config.Config, log *slog.Logger) (*gorm.DB, error) {
	const attempts = 30
	const retryInterval = 2 * time.Second
	dsn := cfg.Database().DSN()
	var lastErr error
	for attempt := 1; attempt <= attempts; attempt++ {
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: gormlogger.Default.LogMode(gormlogger.Warn),
		})
		if err == nil {
			sqlDB, dbErr := db.DB()
			if dbErr == nil {
				sqlDB.SetMaxOpenConns(25)
				sqlDB.SetMaxIdleConns(5)
				sqlDB.SetConnMaxLifetime(5 * time.Minute)
				if pingErr := sqlDB.Ping(); pingErr == nil {
					return db, nil
				} else {
					lastErr = pingErr
					if closeErr := sqlDB.Close(); closeErr != nil {
						log.Warn("close database", "error", closeErr)
					}
				}
			} else {
				lastErr = dbErr
			}
		} else {
			lastErr = err
		}
		log.Warn("waiting for database", "attempt", attempt, "error", lastErr)
		time.Sleep(retryInterval)
	}
	return nil, fmt.Errorf("connect database: %w", lastErr)
}

func run(cfg *config.Config, log *slog.Logger) error {
	db, err := connectDB(cfg, log)
	if err != nil {
		return err
	}

	if err := migrations.AutoMigrate(db); err != nil {
		return err
	}
	if err := seeds.Seed(db); err != nil {
		return err
	}

	// 仓储层
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	equipmentRepo := repository.NewEquipmentRepository(db)
	borrowRepo := repository.NewBorrowRepository(db)
	maintenanceRepo := repository.NewMaintenanceRepository(db)
	reservationRepo := repository.NewReservationRepository(db)
	auditRepo := repository.NewAuditLogRepository(db)

	// 服务层
	auditService := service.NewAuditService(auditRepo, log)
	authService := service.NewAuthService(userRepo, cfg.JWT(), log)
	userService := service.NewUserService(userRepo, log)
	categoryService := service.NewCategoryService(categoryRepo, log)
	equipmentService := service.NewEquipmentService(equipmentRepo, categoryRepo, userRepo, auditService, log)
	borrowService := service.NewBorrowService(borrowRepo, equipmentRepo, auditService, log)
	maintenanceService := service.NewMaintenanceService(maintenanceRepo, equipmentRepo, auditService, log)
	reservationService := service.NewReservationService(reservationRepo, equipmentRepo, auditService, log)
	dashboardService := service.NewDashboardService(equipmentRepo, borrowRepo, reservationRepo, log)

	// 处理器层
	deps := router.Dependencies{
		AuthHandler:        handler.NewAuthHandler(authService),
		UserHandler:        handler.NewUserHandler(userService),
		CategoryHandler:    handler.NewCategoryHandler(categoryService),
		EquipmentHandler:   handler.NewEquipmentHandler(equipmentService),
		BorrowHandler:      handler.NewBorrowHandler(borrowService),
		MaintenanceHandler: handler.NewMaintenanceHandler(maintenanceService),
		ReservationHandler: handler.NewReservationHandler(reservationService),
		DashboardHandler:   handler.NewDashboardHandler(dashboardService),
		AuditHandler:       handler.NewAuditHandler(auditService),
		AuthService:        authService,
		AuditService:       auditService,
		Logger:             log,
	}

	engine := router.NewRouter(deps)
	server := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           engine,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Info("server started", "port", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("listen and serve", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return server.Shutdown(ctx)
}
