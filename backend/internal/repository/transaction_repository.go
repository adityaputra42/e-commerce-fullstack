package repository

import (
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(param models.Transaction, tx *gorm.DB) (models.Transaction, error)
	Update(param models.Transaction, tx *gorm.DB) (models.Transaction, error)
	Delete(param models.Transaction) error
	FindById(paramId uint) (models.Transaction, error)
	FindAll(param models.TransactionListRequest) ([]models.Transaction, error)
}

type TransactionRepositoryImpl struct {
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

	if err := db.Preload("Transactions").Find(&Transactions).Error; err != nil {
		return nil, err
	}

	return Transactions, nil
}

// FindById implements TransactionRepository.
func (a *TransactionRepositoryImpl) FindById(paramId uint) (models.Transaction, error) {
	Transaction := models.Transaction{}
	err := database.DB.Model(&models.User{}).Take(&Transaction, "id =?", paramId).Error

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

	err = db.First(&result, param.TxID).Error
	return result, err

}

func NewTransactionRepository() TransactionRepository {
	return &TransactionRepositoryImpl{}
}
