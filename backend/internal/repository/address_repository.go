package repository

import (
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"

	"gorm.io/gorm"
)

type AddressRepository interface {
	Create(param models.Address, tx *gorm.DB) (models.Address, error)
	Update(param models.Address, tx *gorm.DB) (models.Address, error)
	Delete(param models.Address) error
	FindById(paramId uint) (models.Address, error)
	FindAll(param models.AddressListRequest) ([]models.Address, error)
}

type AddressRepositoryImpl struct {
}

// Create implements AddressRepository.
func (a *AddressRepositoryImpl) Create(param models.Address, tx *gorm.DB) (models.Address, error) {
	var result models.Address

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

// Delete implements AddressRepository.
func (a *AddressRepositoryImpl) Delete(param models.Address) error {
	return database.DB.Delete(&param).Error
}

// FindAll implements AddressRepository.
func (a *AddressRepositoryImpl) FindAll(param models.AddressListRequest) ([]models.Address, error) {

	offset := (param.Page - 1) * param.Limit

	var addresses []models.Address
	db := database.DB

	if param.UserId != nil {
		db = db.Where("user_id = ?", &param.UserId)
	}

	if param.SortBy != "" {
		db = db.Order(param.SortBy)
	}

	if param.Limit > 0 {
		db = db.Limit(param.Limit)
	}

	if offset > 0 {
		db = db.Offset(offset)
	}

	if err := db.Find(&addresses).Error; err != nil {
		return nil, err
	}

	return addresses, nil
}

// FindById implements AddressRepository.
func (a *AddressRepositoryImpl) FindById(paramId uint) (models.Address, error) {
	administrasi := models.Address{}
	err := database.DB.Model(&models.User{}).Take(&administrasi, "id =?", administrasi.UserID).Error

	return administrasi, err
}

// Update implements AddressRepository.
func (a *AddressRepositoryImpl) Update(param models.Address, tx *gorm.DB) (models.Address, error) {
	var result models.Address

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

func NewAddressRepository() AddressRepository {
	return &AddressRepositoryImpl{}
}
