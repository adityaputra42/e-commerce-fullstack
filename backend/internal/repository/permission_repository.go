package repository

import (
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"

	"gorm.io/gorm"
)

type PermissionRepository interface {
	Create(param models.Permission, tx *gorm.DB) (models.Permission, error)
	Update(param models.Permission, tx *gorm.DB) (models.Permission, error)
	Delete(param models.Permission) error
	FindById(paramId uint) (models.Permission, error)
	FindAll(param models.PermissionListRequest) ([]models.Permission, error)
}

type PermissionRepositoryImpl struct {
}

// Create implements PermissionRepository.
func (a *PermissionRepositoryImpl) Create(param models.Permission, tx *gorm.DB) (models.Permission, error) {
	var result models.Permission

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

// Delete implements PermissionRepository.
func (a *PermissionRepositoryImpl) Delete(param models.Permission) error {
	return database.DB.Delete(&param).Error
}

// FindAll implements PermissionRepository.
func (a *PermissionRepositoryImpl) FindAll(param models.PermissionListRequest) ([]models.Permission, error) {

	offset := (param.Page - 1) * param.Limit

	var Permissions []models.Permission
	db := database.DB

	if param.Search != "" {
		db = db.Where("name ILIKE ?", "%"+param.Search+"%")
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

	if err := db.Preload("permissions").Find(&Permissions).Error; err != nil {
		return nil, err
	}

	return Permissions, nil
}

// FindById implements PermissionRepository.
func (a *PermissionRepositoryImpl) FindById(paramId uint) (models.Permission, error) {
	Permission := models.Permission{}
	err := database.DB.Model(&models.User{}).Take(&Permission, "id =?", paramId).Error

	return Permission, err
}

// Update implements PermissionRepository.
func (a *PermissionRepositoryImpl) Update(param models.Permission, tx *gorm.DB) (models.Permission, error) {
	var result models.Permission

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

func NewPermissionRepository() PermissionRepository {
	return &PermissionRepositoryImpl{}
}
