package repository

import (
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(param models.Payment, tx *gorm.DB) (models.Payment, error)
	Update(param models.Payment, tx *gorm.DB) (models.Payment, error)
	Delete(param models.Payment) error
	FindById(paramId uint) (models.Payment, error)
	FindAll(param models.PaymentListRequest) ([]models.Payment, error)
}

type PaymentRepositoryImpl struct {
}

// Create implements PaymentRepository.
func (a *PaymentRepositoryImpl) Create(param models.Payment, tx *gorm.DB) (models.Payment, error) {
	var result models.Payment

	db := database.DB
	if tx != nil {
		db = tx
	}

	err := db.Create(&param).Error
	if err != nil {
		return result, err
	}

	err = db.First(&result, param.ID).Error
	return result, err
}

// Delete implements PaymentRepository.
func (a *PaymentRepositoryImpl) Delete(param models.Payment) error {
	return database.DB.Delete(&param).Error
}

// FindAll implements PaymentRepository.
func (a *PaymentRepositoryImpl) FindAll(param models.PaymentListRequest) ([]models.Payment, error) {
	var Payments []models.Payment
	db := database.DB

	if param.SortBy != "" {
		db = db.Order(param.SortBy)
	}

	if param.Limit > 0 {
		db = db.Limit(param.Limit)
	}

	if param.Offset > 0 {
		db = db.Offset(param.Offset)
	}

	if err := db.Preload("Payments").Find(&Payments).Error; err != nil {
		return nil, err
	}

	return Payments, nil
}

// FindById implements PaymentRepository.
func (a *PaymentRepositoryImpl) FindById(paramId uint) (models.Payment, error) {
	Payment := models.Payment{}
	err := database.DB.Model(&models.User{}).Take(&Payment, "id =?", paramId).Error

	return Payment, err
}

// Update implements PaymentRepository.
func (a *PaymentRepositoryImpl) Update(param models.Payment, tx *gorm.DB) (models.Payment, error) {
	var result models.Payment

	db := database.DB
	if tx != nil {
		db = tx
	}

	err := db.Model(&param).Updates(param).Error
	if err != nil {
		return result, err
	}

	err = db.First(&result, param.ID).Error
	return result, err

}

func NewPaymentRepository() PaymentRepository {
	return &PaymentRepositoryImpl{}
}
