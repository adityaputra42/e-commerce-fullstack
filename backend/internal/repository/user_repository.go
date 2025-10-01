package repository

import (
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(param models.User, tx *gorm.DB) (models.User, error)
	Update(param *models.User, tx *gorm.DB) (models.User, error)
	Delete(param models.User) error
	FindById(paramId uint) (models.User, error)
	FindByEmail(email string) (models.User, error)
	FindAll(param models.UserListRequest) ([]models.User, error)
}

type UserRepositoryImpl struct {
}

// FindByEmail implements UserRepository.
func (u *UserRepositoryImpl) FindByEmail(email string) (models.User, error) {
	user := models.User{}
	err := database.DB.Preload("Role.Permissions").Where("email = ?", email).First(&user).Error
	return user, err
}

// Create implements UserRepository.
func (u *UserRepositoryImpl) Create(param models.User, tx *gorm.DB) (models.User, error) {
	var result models.User

	db := database.DB
	if tx != nil {
		db = tx
	}

	err := db.Create(&param).Error
	if err != nil {
		return result, err
	}

	err = db.Preload("Role.Permissions").First(&result, param.ID).Error
	return result, err
}

// Delete implements UserRepository.
func (u *UserRepositoryImpl) Delete(param models.User) error {
	return database.DB.Delete(&param).Error
}

// FindAll implements UserRepository.
func (u *UserRepositoryImpl) FindAll(param models.UserListRequest) ([]models.User, error) {
	var users []models.User
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

	if param.Offset > 0 {
		db = db.Offset(param.Offset)
	}

	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// FindById implements UserRepository.
func (u *UserRepositoryImpl) FindById(paramId uint) (models.User, error) {
	user := models.User{}
	err := database.DB.Preload("Role.Permissions").First(&user, user.ID).Error

	return user, err
}

// Update implements UserRepository.
func (u *UserRepositoryImpl) Update(param *models.User, tx *gorm.DB) (models.User, error) {
	var result models.User

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

func NewUserReposiory() UserRepository {
	return &UserRepositoryImpl{}
}
