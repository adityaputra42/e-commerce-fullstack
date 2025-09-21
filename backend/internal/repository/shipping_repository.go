package repository

import (
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"

	"gorm.io/gorm"
)

type ShippingRepository interface {
	Create(param models.Shipping, tx *gorm.DB) (models.Shipping, error)
	Update(param models.Shipping, tx *gorm.DB) (models.Shipping, error)
	Delete(param models.Shipping) error
	FindById(paramId uint) (models.Shipping, error)
	FindAll(param models.ShippingListRequest) ([]models.Shipping, error)
}

type ShippingRepositoryImpl struct {
}

// Create implements ShippingRepository.
func (a *ShippingRepositoryImpl) Create(param models.Shipping, tx *gorm.DB) (models.Shipping, error) {
	var result models.Shipping

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

// Delete implements ShippingRepository.
func (a *ShippingRepositoryImpl) Delete(param models.Shipping) error {
	return database.DB.Delete(&param).Error
}

// FindAll implements ShippingRepository.
func (a *ShippingRepositoryImpl) FindAll(param models.ShippingListRequest) ([]models.Shipping, error) {
	var Shippings []models.Shipping
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

	if err := db.Preload("Shippings").Find(&Shippings).Error; err != nil {
		return nil, err
	}

	return Shippings, nil
}

// FindById implements ShippingRepository.
func (a *ShippingRepositoryImpl) FindById(paramId uint) (models.Shipping, error) {
	Shipping := models.Shipping{}
	err := database.DB.Model(&models.User{}).Take(&Shipping, "id =?", paramId).Error

	return Shipping, err
}

// Update implements ShippingRepository.
func (a *ShippingRepositoryImpl) Update(param models.Shipping, tx *gorm.DB) (models.Shipping, error) {
	var result models.Shipping

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

func NewShippingRepository() ShippingRepository {
	return &ShippingRepositoryImpl{}
}
