package repository

import (
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(param models.Order, tx *gorm.DB) (models.Order, error)
	Update(param models.Order, tx *gorm.DB) (models.Order, error)
	Delete(param models.Order) error
	FindById(paramId uint) (models.Order, error)
	FindAll(param models.OrderListRequest) ([]models.Order, error)
}

type OrderRepositoryImpl struct {
}

// Create implements OrderRepository.
func (a *OrderRepositoryImpl) Create(param models.Order, tx *gorm.DB) (models.Order, error) {
	var result models.Order

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

// Delete implements OrderRepository.
func (a *OrderRepositoryImpl) Delete(param models.Order) error {
	return database.DB.Delete(&param).Error
}

// FindAll implements OrderRepository.
func (a *OrderRepositoryImpl) FindAll(param models.OrderListRequest) ([]models.Order, error) {

	offset := (param.Page - 1) * param.Limit

	var Orders []models.Order
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

	if err := db.Preload("Orders").Find(&Orders).Error; err != nil {
		return nil, err
	}

	return Orders, nil
}

// FindById implements OrderRepository.
func (a *OrderRepositoryImpl) FindById(paramId uint) (models.Order, error) {
	Order := models.Order{}
	err := database.DB.Model(&models.User{}).Take(&Order, "id =?", paramId).Error

	return Order, err
}

// Update implements OrderRepository.
func (a *OrderRepositoryImpl) Update(param models.Order, tx *gorm.DB) (models.Order, error) {
	var result models.Order

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

func NewOrderRepository() OrderRepository {
	return &OrderRepositoryImpl{}
}
