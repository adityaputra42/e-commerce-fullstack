package repository_test

import (
	"e-commerce/backend/internal/testhelper"
	"testing"

	"gorm.io/gorm"
)

var testDB *gorm.DB

// TestMain sets up and tears down the test suite for repository package
func TestMain(m *testing.M) {
	// Setup
	testDB = testhelper.SetupTestDB()
	testhelper.CleanupTestDB(testDB)
	testhelper.MigrateTestDB(testDB)

	// Run tests
	code := m.Run()

	// Teardown
	testhelper.TeardownTestSuite(testDB)
	
	// Exit is handled by calling os.Exit if needed, but returning code allows defers to run if we were in main
	// panic if failed
	if code != 0 {
		panic("Repository Tests failed")
	}
}
