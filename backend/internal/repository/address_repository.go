package repository

import (
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"

	"gorm.io/gorm"
)

type AddressRepository interface {
	Create(param models.Address, tx *gorm.DB) (*models.Address, error)
	Update(param models.Address, tx *gorm.DB) (*models.Address, error)
	Delete(param models.Address) error
	FindById(paramId uint) (*models.Address, error)
	FindAll(param models.AddressListRequest) ([]*models.Address, error)
	FindAllByUserId(param int64) ([]*models.Address, error)
	CountByUser(param int64) (int64, error)
}

type AddressRepositoryImpl struct {
}

// CountByUser implements AddressRepository.
func (a *AddressRepositoryImpl) CountByUser(param int64) (int64, error) {
	var count int64
	err := database.DB.
		Model(&models.Address{}).
		Where("user_id = ?", param).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// FindAllByUserId implements AddressRepository.
func (a *AddressRepositoryImpl) FindAllByUserId(param int64) ([]*models.Address, error) {
	var addresses []*models.Address
	db := database.DB

	if err := db.
		Select("id", "user_id", "recipient_name", "recipient_phone_number", "province", "city", "district", "village", "postal_code", "full_address", "created_at", "updated_at").
		Where("user_id = ?", param).
		Find(&addresses).Error; err != nil {
		return nil, err
	}
	return addresses, nil
}

// Create implements AddressRepository.
func (a *AddressRepositoryImpl) Create(param models.Address, tx *gorm.DB) (*models.Address, error) {
	var result models.Address

	db := database.DB
	if tx != nil {
		db = tx
	}

	err := db.Create(&param).Error
	if err != nil {
		return nil, err
	}

	err = db.
		Select("id", "user_id", "recipient_name", "recipient_phone_number", "province", "city", "district", "village", "postal_code", "full_address", "created_at", "updated_at").
		First(&result, param.ID).Error
	return &result, err
}

// Delete implements AddressRepository.
func (a *AddressRepositoryImpl) Delete(param models.Address) error {
	return database.DB.Delete(&param).Error
}

// FindAll implements AddressRepository.
func (a *AddressRepositoryImpl) FindAll(param models.AddressListRequest) ([]*models.Address, error) {
	offset := (param.Page - 1) * param.Limit

	var addresses []*models.Address
	db := database.DB.
		Select("id", "user_id", "recipient_name", "recipient_phone_number", "province", "city", "district", "village", "postal_code", "full_address", "created_at", "updated_at")

	if param.UserId != nil {
		db = db.Where("user_id = ?", param.UserId)
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
func (a *AddressRepositoryImpl) FindById(paramId uint) (*models.Address, error) {
	var address models.Address
	err := database.DB.
		Select("id", "user_id", "recipient_name", "recipient_phone_number", "province", "city", "district", "village", "postal_code", "full_address", "created_at", "updated_at").
		First(&address, "id = ?", paramId).Error

	if err != nil {
		return nil, err
	}

	return &address, nil
}

// Update implements AddressRepository.
func (a *AddressRepositoryImpl) Update(param models.Address, tx *gorm.DB) (*models.Address, error) {
	var result models.Address

	db := database.DB
	if tx != nil {
		db = tx
	}

	err := db.Model(&param).Updates(param).Error
	if err != nil {
		return nil, err
	}

	err = db.
		Select("id", "user_id", "recipient_name", "recipient_phone_number", "province", "city", "district", "village", "postal_code", "full_address", "created_at", "updated_at").
		First(&result, param.ID).Error
	return &result, err
}

func NewAddressRepository() AddressRepository {
	return &AddressRepositoryImpl{}
}
