package repository

import (
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"
	"math"
)

type UserRepository interface {
	Create(param models.User) (models.User, error)
	Update(param *models.User) (models.User, error)
	Delete(param models.User) error
	FindById(paramId int64) (models.User, error)
	FindByEmail(email string) (models.User, error)
	FindAll(param models.UserListRequest) (*models.UserListResponse, error)
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
func (u *UserRepositoryImpl) Create(param models.User) (models.User, error) {
	var result models.User

	db := database.DB

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
func (u *UserRepositoryImpl) FindAll(param models.UserListRequest) (*models.UserListResponse, error) {
	offset := (param.Page - 1) * param.Limit

	query := database.DB.Model(&models.User{}).Preload("Role")

	if param.SortBy == "" {
		param.SortBy = "created_at"
	}

	orderClause := param.SortBy
	query = query.Order(orderClause)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var users []models.User
	if err := query.Offset(offset).Limit(param.Limit).Find(&users).Error; err != nil {
		return nil, err
	}

	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *user.ToResponse()
	}

	totalPages := int(math.Ceil(float64(total) / float64(param.Limit)))

	return &models.UserListResponse{
		Users:      userResponses,
		Total:      total,
		Page:       param.Page,
		Limit:      param.Limit,
		TotalPages: totalPages,
	}, nil
}

// FindById implements UserRepository.
func (u *UserRepositoryImpl) FindById(paramId int64) (models.User, error) {
	user := models.User{}
	err := database.DB.Preload("Role.Permissions.Address").First(&user, user.ID).Error

	return user, err
}

// Update implements UserRepository.
func (u *UserRepositoryImpl) Update(param *models.User) (models.User, error) {
	var result models.User

	db := database.DB

	err := db.Model(&param).Updates(param).Error
	if err != nil {
		return result, err
	}

	err = db.Preload("Role.Permissions").First(&result, param.ID).Error
	return result, err

}

func NewUserReposiory() UserRepository {
	return &UserRepositoryImpl{}
}
