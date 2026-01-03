package repository

import (
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionRepository interface {
	Create(param models.Transaction, tx *gorm.DB) (models.Transaction, error)
	Update(param models.Transaction, tx *gorm.DB) (models.Transaction, error)
	Delete(param models.Transaction) error
	FindById(paramId string) (models.Transaction, error)
	FindByIdLocking(tx *gorm.DB, id string) (*models.Transaction, error)
	FindAll(param models.TransactionListRequest) ([]models.Transaction, error)
}

type TransactionRepositoryImpl struct {
}

// FindByIdLocking implements TransactionRepository.
func (r *TransactionRepositoryImpl) FindByIdLocking(tx *gorm.DB, id string) (*models.Transaction, error) {
	var trx models.Transaction

	err := tx.
		Select("tx_id", "address_id", "shipping_id", "payment_method_id", "total_price", "shipping_price", "status", "created_at", "updated_at").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("tx_id = ?", id).
		First(&trx).Error

	if err != nil {
		return nil, err
	}

	return &trx, nil
}

// Create implements TransactionRepository.
func (a *TransactionRepositoryImpl) Create(param models.Transaction, tx *gorm.DB) (models.Transaction, error) {
	var result models.Transaction

	db := database.DB
	if tx != nil {
		db = tx
	}

	err := db.Create(&param).Error
	if err != nil {
		return result, err
	}

	err = db.
		Select("tx_id", "address_id", "shipping_id", "payment_method_id", "total_price", "shipping_price", "status", "created_at", "updated_at").
		First(&result, "tx_id = ?", param.TxID).Error

	return result, err
}

// Delete implements TransactionRepository.
func (a *TransactionRepositoryImpl) Delete(param models.Transaction) error {
	return database.DB.Delete(&param).Error
}

// FindAll implements TransactionRepository.
func (a *TransactionRepositoryImpl) FindAll(param models.TransactionListRequest) ([]models.Transaction, error) {
	offset := (param.Page - 1) * param.Limit

	var Transactions []models.Transaction
	db := database.DB.
		Select("tx_id", "address_id", "shipping_id", "payment_method_id", "total_price", "shipping_price", "status", "created_at", "updated_at")

	if param.SortBy != "" {
		db = db.Order(param.SortBy)
	}

	if param.Limit > 0 {
		db = db.Limit(param.Limit)
	}

	if offset > 0 {
		db = db.Offset(offset)
	}

	if err := db.Find(&Transactions).Error; err != nil {
		return nil, err
	}

	return Transactions, nil
}

// FindById implements TransactionRepository.
func (a *TransactionRepositoryImpl) FindById(paramId string) (models.Transaction, error) {
	Transaction := models.Transaction{}
	err := database.DB.
		Select("tx_id", "address_id", "shipping_id", "payment_method_id", "total_price", "shipping_price", "status", "created_at", "updated_at").
		Preload("Address", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "user_id", "recipient_name", "recipient_phone_number", "province", "city", "district", "village", "postal_code", "full_address", "created_at", "updated_at")
		}).
		Preload("Shipping", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "price", "state", "created_at", "updated_at")
		}).
		Preload("PaymentMethod", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "account_name", "account_number", "bank_name", "bank_images", "is_active", "created_at", "updated_at")
		}).
		Preload("Orders", func(db *gorm.DB) *gorm.DB {
			return db.
				Select("id", "transaction_id", "product_id", "color_varian_id", "size_varian_id", "quantity", "unit_price", "subtotal", "created_at", "updated_at").
				Preload("Product", func(db *gorm.DB) *gorm.DB {
					return db.Select("id", "category_id", "name", "description", "created_at", "updated_at")
				}).
				Preload("ColorVarian", func(db *gorm.DB) *gorm.DB {
					return db.Select("id", "product_id", "name", "color", "images", "created_at", "updated_at")
				}).
				Preload("SizeVarian", func(db *gorm.DB) *gorm.DB {
					return db.Select("id", "color_varian_id", "size", "stock", "created_at", "updated_at")
				})
		}).
		First(&Transaction, "tx_id = ?", paramId).Error

	return Transaction, err
}

// Update implements TransactionRepository.
func (a *TransactionRepositoryImpl) Update(param models.Transaction, tx *gorm.DB) (models.Transaction, error) {
	var result models.Transaction

	db := database.DB
	if tx != nil {
		db = tx
	}

	err := db.Model(&param).Updates(param).Error
	if err != nil {
		return result, err
	}

	err = db.
		Select("tx_id", "address_id", "shipping_id", "payment_method_id", "total_price", "shipping_price", "status", "created_at", "updated_at").
		Preload("Address", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "user_id", "recipient_name", "recipient_phone_number", "province", "city", "district", "village", "postal_code", "full_address", "created_at", "updated_at")
		}).
		Preload("Shipping", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "price", "state", "created_at", "updated_at")
		}).
		Preload("PaymentMethod", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "account_name", "account_number", "bank_name", "bank_images", "is_active", "created_at", "updated_at")
		}).
		Preload("Orders", func(dbf *gorm.DB) *gorm.DB {
			return dbf.
				Select("id", "transaction_id", "product_id", "color_varian_id", "size_varian_id", "quantity", "unit_price", "subtotal", "created_at", "updated_at").
				Preload("Product", func(db *gorm.DB) *gorm.DB {
					return db.Select("id", "category_id", "name", "description", "created_at", "updated_at")
				}).
				Preload("ColorVarian", func(db *gorm.DB) *gorm.DB {
					return db.Select("id", "product_id", "name", "color", "images", "created_at", "updated_at")
				}).
				Preload("SizeVarian", func(db *gorm.DB) *gorm.DB {
					return db.Select("id", "color_varian_id", "size", "stock", "created_at", "updated_at")
				})
		}).
		First(&result, "tx_id = ?", param.TxID).Error

	return result, err
}

func NewTransactionRepository() TransactionRepository {
	return &TransactionRepositoryImpl{}
}
