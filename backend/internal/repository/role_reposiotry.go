package repository

import (
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"

	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(param models.Role, tx *gorm.DB) (models.Role, error)
	Update(param models.Role, tx *gorm.DB) (models.Role, error)
	Delete(param models.Role) error
	FindById(paramId uint) (models.Role, error)
	FindByName(name string) (models.Role, error)
	FindAll(param models.RoleListRequest) ([]models.Role, error)
}

type RoleRepositoryImpl struct {
}

// FindByName implements RoleRepository.
func (a *RoleRepositoryImpl) FindByName(name string) (models.Role, error) {
	role := models.Role{}
	err := database.DB.Model(&models.User{}).Take(&role, "name =?", name).Error

	return role, err
}

// Create implements RoleRepository.
func (a *RoleRepositoryImpl) Create(param models.Role, tx *gorm.DB) (models.Role, error) {
	var result models.Role

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

// Delete implements RoleRepository.
func (a *RoleRepositoryImpl) Delete(param models.Role) error {
	return database.DB.Delete(&param).Error
}

// FindAll implements RoleRepository.
func (a *RoleRepositoryImpl) FindAll(param models.RoleListRequest) ([]models.Role, error) {

	offset := (param.Page - 1) * param.Limit

	var Roles []models.Role
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

	if err := db.Preload("permissions").Find(&Roles).Error; err != nil {
		return nil, err
	}

	return Roles, nil
}

// FindById implements RoleRepository.
func (a *RoleRepositoryImpl) FindById(paramId uint) (models.Role, error) {
	role := models.Role{}
	err := database.DB.Model(&models.User{}).Take(&role, "id =?", paramId).Error

	return role, err
}

// Update implements RoleRepository.
func (a *RoleRepositoryImpl) Update(param models.Role, tx *gorm.DB) (models.Role, error) {
	var result models.Role

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

func NewRoleRepository() RoleRepository {
	return &RoleRepositoryImpl{}
}
