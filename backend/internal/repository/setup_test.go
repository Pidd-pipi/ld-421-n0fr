package repository

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/labequipment/lab-equipment/database/migrations"
	"gorm.io/gorm"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := migrations.AutoMigrate(db); err != nil {
		t.Fatalf("migrate sqlite: %v", err)
	}
	return db
}
