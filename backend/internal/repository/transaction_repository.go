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

// FindByIdLocking implements [TransactionRepository].
func (r *TransactionRepositoryImpl) FindByIdLocking(tx *gorm.DB, id string) (*models.Transaction, error) {
	var trx models.Transaction

	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
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

	err = db.First(&result, param.TxID).Error
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
	db := database.DB

	if param.SortBy != "" {
		db = db.Order(param.SortBy)
	}

	if param.Limit > 0 {
		db = db.Limit(param.Limit)
	}

	if offset > 0 {
		db = db.Offset(offset)
	}

	if err := db.
		Find(&Transactions).Error; err != nil {
		return nil, err
	}

	return Transactions, nil
}

// FindById implements TransactionRepository.
func (a *TransactionRepositoryImpl) FindById(paramId string) (models.Transaction, error) {
	Transaction := models.Transaction{}
	err := database.DB.
		Preload("Address").
		Preload("Shipping").
		Preload("PaymentMethod").
		Preload("Orders", func(db *gorm.DB) *gorm.DB {
			return db.
				Preload("Product").
				Preload("ColorVarian").
				Preload("SizeVarian")
		}).
		First(&Transaction, "tx_id =?", paramId).Error

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

	err = db.Preload("Address").
		Preload("Shipping").
		Preload("PaymentMethod").
		Preload("Orders", func(db *gorm.DB) *gorm.DB {
			return db.
				Preload("Product").
				Preload("ColorVarian").
				Preload("SizeVarian")
		}).First(&result, param.TxID).Error
	return result, err

}

func NewTransactionRepository() TransactionRepository {
	return &TransactionRepositoryImpl{}
}
