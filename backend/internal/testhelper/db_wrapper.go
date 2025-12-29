package testhelper

import (
	"e-commerce/backend/internal/database"

	"gorm.io/gorm"
)

// DBWrapper allows us to temporarily replace the global database.DB for testing
type DBWrapper struct {
	OriginalDB *gorm.DB
}

// SetTestDB temporarily replaces database.DB with a test database
func SetTestDB(testDB *gorm.DB) *DBWrapper {
	wrapper := &DBWrapper{
		OriginalDB: database.DB,
	}
	database.DB = testDB
	return wrapper
}

// Restore restores the original database.DB
func (w *DBWrapper) Restore() {
	database.DB = w.OriginalDB
}
