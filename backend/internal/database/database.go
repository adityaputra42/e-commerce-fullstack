package database

import (
	"e-commerce/backend/internal/config"
	"e-commerce/backend/internal/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return err
	}

	log.Println("Database connected successfully!")
	return nil
}

func Migrate() error {
	err := DB.AutoMigrate(
		&models.SeedTracker{},
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.ActivityLog{},
		&models.PasswordResetToken{},
		&models.Address{},
		&models.Category{},
		&models.Product{},
		&models.ColorVarian{},
		&models.SizeVarian{},
		&models.Shipping{},
		&models.PaymentMethod{},
		&models.Transaction{},
		&models.Order{},
		&models.Payment{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
