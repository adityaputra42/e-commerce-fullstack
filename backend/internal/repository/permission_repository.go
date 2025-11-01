package repository

import (
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"
)

type PermissionRepository interface {
	Create(param models.Permission) (models.Permission, error)
	Update(param models.Permission) (models.Permission, error)
	Delete(param models.Permission) error
	FindById(paramId uint) (models.Permission, error)
	FindAll() ([]models.Permission, error)
	FindAllById(listId []uint) (*[]models.Permission, error)
}

type PermissionRepositoryImpl struct {
}

// FindAllById implements PermissionRepository.
func (a *PermissionRepositoryImpl) FindAllById(listId []uint) (*[]models.Permission, error) {
	var permissions []models.Permission
	if err := database.DB.Where("id IN ?", listId).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return &permissions, nil
}

// Create implements PermissionRepository.
func (a *PermissionRepositoryImpl) Create(param models.Permission) (models.Permission, error) {
	var result models.Permission

	db := database.DB

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
func (a *PermissionRepositoryImpl) FindAll() ([]models.Permission, error) {

	var Permissions []models.Permission
	db := database.DB

	if err := db.Find(&Permissions).Error; err != nil {
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
func (a *PermissionRepositoryImpl) Update(param models.Permission) (models.Permission, error) {
	var result models.Permission

	db := database.DB

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
