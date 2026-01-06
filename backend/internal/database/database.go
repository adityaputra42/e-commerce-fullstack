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
		return fmt.Errorf("failed to connect database: %w", err)
	}

	log.Println("Database connected successfully!")
	return nil
}

func Migrate() error {
	err := DB.AutoMigrate(
		&models.SeedTracker{},
		&models.Role{},
		&models.Permission{},
		&models.User{},
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

	if err := createCompositeIndexes(); err != nil {
		return fmt.Errorf("failed to create composite indexes: %w", err)
	}

	return nil
}

func createCompositeIndexes() error {
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_order_user_status ON orders(user_id, status)",
		"CREATE INDEX IF NOT EXISTS idx_order_transaction_status ON orders(transaction_id, status)",
		"CREATE INDEX IF NOT EXISTS idx_order_product_variant ON orders(product_id, color_varian_id, size_varian_id)",
		"CREATE INDEX IF NOT EXISTS idx_order_created_status ON orders(created_at, status)",

		"CREATE INDEX IF NOT EXISTS idx_payment_tx_status ON payments(transaction_id, status)",
		"CREATE INDEX IF NOT EXISTS idx_payment_status_created ON payments(status, created_at)",

		"CREATE INDEX IF NOT EXISTS idx_transaction_status_created ON transactions(status, created_at)",
		"CREATE INDEX IF NOT EXISTS idx_transaction_address_shipping ON transactions(address_id, shipping_id)",

		"CREATE INDEX IF NOT EXISTS idx_user_role_active ON users(role_id, is_active)",
		"CREATE INDEX IF NOT EXISTS idx_user_active_created ON users(is_active, created_at)",

		"CREATE INDEX IF NOT EXISTS idx_address_user_created ON addresses(user_id, created_at)",

		"CREATE INDEX IF NOT EXISTS idx_payment_method_active_created ON payment_methods(is_active, created_at)",

		"CREATE INDEX IF NOT EXISTS idx_shipping_state ON shippings(state)",
	}

	for _, indexSQL := range indexes {
		if err := DB.Exec(indexSQL).Error; err != nil {
			log.Printf("Warning: could not create index: %v\n", err)
		}
	}

	log.Println("Composite indexes created/verified successfully")
	return nil
}

// Utility function untuk informasi database
func GetDBStats() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("Error getting database stats: %v\n", err)
		return
	}

	stats := sqlDB.Stats()
	log.Printf("DB Stats - Open Connections: %d, In Use: %d, Idle: %d, Wait Count: %d, Wait Duration: %v, Max Idle Closed: %d, Max Lifetime Closed: %d\n",
		stats.OpenConnections,
		stats.InUse,
		stats.Idle,
		stats.WaitCount,
		stats.WaitDuration,
		stats.MaxIdleClosed,
		stats.MaxLifetimeClosed,
	)
}

func GetDB() *gorm.DB {
	return DB
}
