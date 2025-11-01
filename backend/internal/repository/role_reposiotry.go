package repository

import (
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"
	"errors"
)

type RoleRepository interface {
	Create(param models.Role) (models.Role, error)
	Update(param models.Role) (models.Role, error)
	Delete(param models.Role) error
	FindById(paramId uint) (models.Role, error)
	FindByName(name string) (models.Role, error)
	FindByNameAndId(name string, id uint) (*models.Role, error)
	AddPermission(id uint, permissions *[]models.Permission) (*models.Role, error)
	UpdatePermission(id uint, permissions *[]models.Permission) (*models.Role, error)
	FindAll() ([]models.Role, error)
}

type RoleRepositoryImpl struct {
}

// UpdatePermission implements RoleRepository.
func (a *RoleRepositoryImpl) UpdatePermission(id uint, permissions *[]models.Permission) (*models.Role, error) {
	role := models.Role{}
	err := database.DB.Model(&models.Role{}).Take(&role, "id =?", id).Error
	if err != nil {
		return nil, err
	}

	if err := database.DB.Model(&role).Association("Permissions").Replace(&permissions); err != nil {
		return nil, err
	}
	return &role, nil
}

// FindByNameAndId implements RoleRepository.
func (a *RoleRepositoryImpl) FindByNameAndId(name string, id uint) (*models.Role, error) {
	var existingRole models.Role
	if err := database.DB.Where("name = ? AND id != ?", name, id).First(&existingRole).Error; err == nil {
		return nil, errors.New("role with this name already exists")
	}
	return &existingRole, nil
}

// AddPermission implements RoleRepository.
func (a *RoleRepositoryImpl) AddPermission(id uint, permissions *[]models.Permission) (*models.Role, error) {
	role := models.Role{}
	err := database.DB.Model(&models.Role{}).Take(&role, "id =?", id).Error
	if err != nil {
		return nil, err
	}

	if err := database.DB.Model(&role).Association("Permissions").Append(&permissions); err != nil {
		return nil, err
	}
	return &role, nil
}

// FindByName implements RoleRepository.
func (a *RoleRepositoryImpl) FindByName(name string) (models.Role, error) {
	role := models.Role{}
	err := database.DB.Model(&models.User{}).Take(&role, "name =?", name).Error

	return role, err
}

// Create implements RoleRepository.
func (a *RoleRepositoryImpl) Create(param models.Role) (models.Role, error) {
	var result models.Role

	db := database.DB

	err := db.Create(&param).Error
	if err != nil {
		return result, err
	}

	err = db.First(&result, param.ID).Error
	return result, err
}

// Delete implements RoleRepository.
func (a *RoleRepositoryImpl) Delete(param models.Role) error {
	return database.DB.Delete(&param).Error
}

// FindAll implements RoleRepository.
func (a *RoleRepositoryImpl) FindAll() ([]models.Role, error) {
	var roles []models.Role
	if err := database.DB.Preload("Permissions").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// FindById implements RoleRepository.
func (a *RoleRepositoryImpl) FindById(paramId uint) (models.Role, error) {
	role := models.Role{}
	err := database.DB.Preload("permissions").Model(&models.Role{}).Take(&role, "id =?", paramId).Error

	return role, err
}

// Update implements RoleRepository.
func (a *RoleRepositoryImpl) Update(param models.Role) (models.Role, error) {
	var result models.Role

	db := database.DB

	err := db.Model(&param).Updates(param).Error
	if err != nil {
		return result, err
	}

	err = db.First(&result, param.ID).Error
	return result, err

}

func NewRoleRepository() RoleRepository {
	return &RoleRepositoryImpl{}
}
