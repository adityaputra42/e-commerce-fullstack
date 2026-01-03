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
	FindById(paramId uint) (*models.Shipping, error)
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

	err = db.
		Select("id", "name", "price", "state", "created_at", "updated_at").
		First(&result, param.ID).Error
	return result, err
}

// Delete implements ShippingRepository.
func (a *ShippingRepositoryImpl) Delete(param models.Shipping) error {
	return database.DB.Delete(&param).Error
}

// FindAll implements ShippingRepository.
func (a *ShippingRepositoryImpl) FindAll(param models.ShippingListRequest) ([]models.Shipping, error) {
	offset := (param.Page - 1) * param.Limit
	var Shippings []models.Shipping
	db := database.DB.
		Select("id", "name", "price", "state", "created_at", "updated_at")

	if param.SortBy != "" {
		db = db.Order(param.SortBy)
	}

	if param.Limit > 0 {
		db = db.Limit(param.Limit)
	}

	if offset > 0 {
		db = db.Offset(offset)
	}

	if err := db.Find(&Shippings).Error; err != nil {
		return nil, err
	}

	return Shippings, nil
}

// FindById implements ShippingRepository.
func (a *ShippingRepositoryImpl) FindById(paramId uint) (*models.Shipping, error) {
	Shipping := models.Shipping{}
	err := database.DB.
		Select("id", "name", "price", "state", "created_at", "updated_at").
		First(&Shipping, "id = ?", paramId).Error

	if err != nil {
		return nil, err
	}

	return &Shipping, nil
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

	err = db.
		Select("id", "name", "price", "state", "created_at", "updated_at").
		First(&result, param.ID).Error
	return result, err
}

func NewShippingRepository() ShippingRepository {
	return &ShippingRepositoryImpl{}
}