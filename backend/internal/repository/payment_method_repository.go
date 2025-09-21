package repository

import (
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"

	"gorm.io/gorm"
)

type PaymentMethodRepository interface {
	Create(param models.PaymentMethod, tx *gorm.DB) (models.PaymentMethod, error)
	Update(param models.PaymentMethod, tx *gorm.DB) (models.PaymentMethod, error)
	Delete(param models.PaymentMethod) error
	FindById(paramId uint) (models.PaymentMethod, error)
	FindAll(param models.PaymentMethodListRequest) ([]models.PaymentMethod, error)
}

type PaymentMethodRepositoryImpl struct {
}

// Create implements PaymentMethodRepository.
func (a *PaymentMethodRepositoryImpl) Create(param models.PaymentMethod, tx *gorm.DB) (models.PaymentMethod, error) {
	var result models.PaymentMethod

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

// Delete implements PaymentMethodRepository.
func (a *PaymentMethodRepositoryImpl) Delete(param models.PaymentMethod) error {
	return database.DB.Delete(&param).Error
}

// FindAll implements PaymentMethodRepository.
func (a *PaymentMethodRepositoryImpl) FindAll(param models.PaymentMethodListRequest) ([]models.PaymentMethod, error) {
	var PaymentMethods []models.PaymentMethod
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

	if err := db.Preload("PaymentMethods").Find(&PaymentMethods).Error; err != nil {
		return nil, err
	}

	return PaymentMethods, nil
}

// FindById implements PaymentMethodRepository.
func (a *PaymentMethodRepositoryImpl) FindById(paramId uint) (models.PaymentMethod, error) {
	PaymentMethod := models.PaymentMethod{}
	err := database.DB.Model(&models.User{}).Take(&PaymentMethod, "id =?", paramId).Error

	return PaymentMethod, err
}

// Update implements PaymentMethodRepository.
func (a *PaymentMethodRepositoryImpl) Update(param models.PaymentMethod, tx *gorm.DB) (models.PaymentMethod, error) {
	var result models.PaymentMethod

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

func NewPaymentMethodRepository() PaymentMethodRepository {
	return &PaymentMethodRepositoryImpl{}
}
