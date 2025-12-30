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
	FindById(id uint) (models.User, error)
	FindByEmail(email string) (models.User, error)
	FindAll(param models.UserListRequest) (*models.UserListResponse, error)
}

type UserRepositoryImpl struct{}

/*
	============================
	  FIND BY EMAIL

============================
*/
func (u *UserRepositoryImpl) FindByEmail(email string) (models.User, error) {
	user := models.User{}

	err := database.DB.
		Preload("Role.Permissions").
		Where("email = ?", email).
		First(&user).Error

	return user, err
}

/*
	============================
	  CREATE USER

============================
*/
func (u *UserRepositoryImpl) Create(param models.User) (models.User, error) {
	var result models.User

	db := database.DB

	if err := db.Create(&param).Error; err != nil {
		return result, err
	}

	err := db.
		Preload("Role.Permissions").
		First(&result, param.ID).Error

	return result, err
}

/*
	============================
	  DELETE USER

============================
*/
func (u *UserRepositoryImpl) Delete(param models.User) error {
	return database.DB.Delete(&param).Error
}

/*
	============================
	  FIND ALL (PAGINATION)

============================
*/
func (u *UserRepositoryImpl) FindAll(param models.UserListRequest) (*models.UserListResponse, error) {
	offset := (param.Page - 1) * param.Limit

	query := database.DB.
		Model(&models.User{}).
		Preload("Role")

	if param.SortBy == "" {
		param.SortBy = "created_at desc"
	}

	query = query.Order(param.SortBy)

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

/*
	============================
	  FIND BY ID  ✅ FIXED

============================
*/
func (u *UserRepositoryImpl) FindById(id uint) (models.User, error) {
	var user models.User
	err := database.DB.
		Preload("Role").
		Preload("Role.Permissions").
		First(&user, id).Error
	return user, err
}

/*
	============================
	  UPDATE USER

============================
*/
func (u *UserRepositoryImpl) Update(param *models.User) (models.User, error) {
	var result models.User

	db := database.DB

	if err := db.Save(param).Error; err != nil {
		return result, err
	}

	err := db.
		Preload("Role.Permissions").
		First(&result, param.ID).Error

	return result, err
}

/*
	============================
	  CONSTRUCTOR

============================
*/
func NewUserReposiory() UserRepository {
	return &UserRepositoryImpl{}
}
