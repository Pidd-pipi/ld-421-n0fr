package service

import (
	"io"
	"log/slog"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/labequipment/lab-equipment/database/migrations"
	"github.com/labequipment/lab-equipment/internal/constants"
	"github.com/labequipment/lab-equipment/internal/model"
	"github.com/labequipment/lab-equipment/internal/repository"
	"gorm.io/gorm"
)

type testEnv struct {
	db                 *gorm.DB
	equipmentService   *EquipmentService
	borrowService      *BorrowService
	reservationService *ReservationService
	disposalService    *DisposalService
	categoryService    *CategoryService
	equipmentRepo      repository.EquipmentRepository
	disposalRepo       repository.DisposalRepository
	categoryID         uint
	ownerID            uint
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := migrations.AutoMigrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	categoryRepo := repository.NewCategoryRepository(db)
	userRepo := repository.NewUserRepository(db)
	equipmentRepo := repository.NewEquipmentRepository(db)
	borrowRepo := repository.NewBorrowRepository(db)
	reservationRepo := repository.NewReservationRepository(db)
	disposalRepo := repository.NewDisposalRepository(db)
	auditRepo := repository.NewAuditLogRepository(db)

	auditService := NewAuditService(auditRepo, logger)
	equipmentService := NewEquipmentService(equipmentRepo, categoryRepo, userRepo, disposalRepo, auditService, logger)
	borrowService := NewBorrowService(borrowRepo, equipmentRepo, disposalRepo, auditService, logger)
	reservationService := NewReservationService(reservationRepo, equipmentRepo, disposalRepo, auditService, logger)
	disposalService := NewDisposalService(disposalRepo, equipmentRepo, auditService, logger)
	categoryService := NewCategoryService(categoryRepo, logger)

	role := model.Role{Code: "Admin", Name: "管理员"}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}
	user := model.User{Username: "admin", PasswordHash: "x", Name: "管理员", RoleID: role.ID, Active: true}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	category := model.EquipmentCategory{Name: "分析仪器"}
	if err := db.Create(&category).Error; err != nil {
		t.Fatalf("create category: %v", err)
	}
	equipment := model.Equipment{Name: "设备A", Code: "EQ-A", CategoryID: category.ID, OwnerID: user.ID, Status: constants.AssetStatusAvailable}
	if err := db.Create(&equipment).Error; err != nil {
		t.Fatalf("create equipment: %v", err)
	}

	return &testEnv{
		db:                 db,
		equipmentService:   equipmentService,
		borrowService:      borrowService,
		reservationService: reservationService,
		disposalService:    disposalService,
		categoryService:    categoryService,
		equipmentRepo:      equipmentRepo,
		disposalRepo:       disposalRepo,
		categoryID:         category.ID,
		ownerID:            user.ID,
	}
}
