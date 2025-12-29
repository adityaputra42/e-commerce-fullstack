package testhelper

import (
	"e-commerce/backend/internal/models"
	"fmt"
	"log"
	"os"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var TestDB *gorm.DB

// SetupTestDB initializes the test database connection
func SetupTestDB() *gorm.DB {
	// Get test database configuration from environment or use defaults
	dbHost := getEnv("TEST_DB_HOST", "localhost")
	dbPort := getEnv("TEST_DB_PORT", "3307")
	dbUser := getEnv("TEST_DB_USER", "testuser")
	dbPassword := getEnv("TEST_DB_PASSWORD", "testpassword")
	dbName := getEnv("TEST_DB_NAME", "ecommerce_test")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	TestDB = db
	return db
}

// CleanupTestDB cleans up all tables in the test database
func CleanupTestDB(db *gorm.DB) {
	// Drop all tables in reverse order of dependencies
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")

	tables := []interface{}{
		&models.ActivityLog{},
		&models.Order{},
		&models.Transaction{},
		&models.Payment{},
		&models.PaymentMethod{},
		&models.Address{},
		&models.ColorVarian{},
		&models.SizeVarian{},
		&models.Product{},
		&models.Category{},
		&models.PasswordResetToken{},
		&models.User{},
		&models.Permission{},
		&models.Role{},
		&models.SeedTracker{},
	}

	for _, table := range tables {
		db.Migrator().DropTable(table)
	}

	db.Exec("SET FOREIGN_KEY_CHECKS = 1")
}

// MigrateTestDB runs migrations on the test database
func MigrateTestDB(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Role{},
		&models.Permission{},
		&models.User{},
		&models.PasswordResetToken{},
		&models.Category{},
		&models.Product{},
		&models.ColorVarian{},
		&models.SizeVarian{},
		&models.Address{},
		&models.PaymentMethod{},
		&models.Payment{},
		&models.Transaction{},
		&models.Order{},
		&models.ActivityLog{},
		&models.SeedTracker{},
	)
}

// SetupTestSuite sets up the test database and runs migrations
func SetupTestSuite(t *testing.T) *gorm.DB {
	db := SetupTestDB()
	CleanupTestDB(db)

	if err := MigrateTestDB(db); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// TeardownTestSuite cleans up the test database
func TeardownTestSuite(db *gorm.DB) {
	CleanupTestDB(db)
	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

// BeginTestTransaction starts a new transaction for a test
func BeginTestTransaction(t *testing.T, db *gorm.DB) *gorm.DB {
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("Failed to begin transaction: %v", tx.Error)
	}
	return tx
}

// RollbackTestTransaction rolls back a test transaction
func RollbackTestTransaction(tx *gorm.DB) {
	tx.Rollback()
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// CreateTestRole creates a test role
func CreateTestRole(db *gorm.DB, name string) *models.Role {
	role := &models.Role{
		Name:        name,
		Description: "Test role: " + name,
	}
	db.Create(role)
	return role
}

// CreateTestUser creates a test user
func CreateTestUser(db *gorm.DB, email string, roleID uint) *models.User {
	user := &models.User{
		Email:        email,
		Username:     email,
		PasswordHash: "$2a$10$test.hash.password",
		FirstName:    "Test",
		LastName:     "User",
		RoleID:       roleID,
		IsActive:     true,
	}
	db.Create(user)
	return user
}

// CreateTestCategory creates a test category
func CreateTestCategory(db *gorm.DB, name string) *models.Category {
	category := &models.Category{
		Name: name,
		Icon: "https://example.com/icon.png",
	}
	db.Create(category)
	return category
}

// CreateTestProduct creates a test product
func CreateTestProduct(db *gorm.DB, name string, categoryID int64, price float64) *models.Product {
	product := &models.Product{
		Name:        name,
		Description: "Test product: " + name,
		Price:       price,
		CategoryID:  categoryID,
		Images:      "https://example.com/product.png",
		Rating:      4.5,
	}
	db.Create(product)
	return product
}
