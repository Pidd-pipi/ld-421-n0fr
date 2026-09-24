package migrations

import (
	"fmt"

	"github.com/labequipment/lab-equipment/internal/model"
	"gorm.io/gorm"
)

// AutoMigrate 自动迁移全部数据表。
func AutoMigrate(db *gorm.DB) error {
	models := []any{
		&model.Role{},
		&model.User{},
		&model.EquipmentCategory{},
		&model.Equipment{},
		&model.BorrowRecord{},
		&model.MaintenanceRecord{},
		&model.Reservation{},
		&model.AuditLog{},
	}
	for _, m := range models {
		if err := db.AutoMigrate(m); err != nil {
			return fmt.Errorf("auto migrate %T: %w", m, err)
		}
	}
	return nil
}
