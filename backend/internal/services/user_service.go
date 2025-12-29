package services

import (
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/repository"
	"e-commerce/backend/internal/utils"
	"errors"

	"gorm.io/gorm"
)

type UserService interface {
	GetUsers(req models.UserListRequest) (*models.UserListResponse, error)
	GetUserById(id uint) (*models.User, error)
	CreateUser(req *models.UserInput) (*models.User, error)
	UpdateUser(id uint, req *models.UserUpdateInput) (*models.User, error)
	DeleteUser(id uint) error
	ActivateUser(id uint) (*models.User, error)
	DeactivateUser(id uint) (*models.User, error)
	GetUserActivityLogs(req *models.ActivityLogListRequest) (*models.ActivityLogListResponse, error)
	BulkUserActions(req *BulkActionRequest) error
	UpdatePassword(userID uint, req *models.PasswordUpdateInput) error
}

type BulkActionRequest struct {
	UserIDs []uint `json:"user_ids" validate:"required"`
	Action  string `json:"action" validate:"required,oneof=activate deactivate delete"`
}

type UserServiceImpl struct {
	userRepo        repository.UserRepository
	acitvityLogRepo repository.ActivityLogRepository
	roleRepo        repository.RoleRepository
}

func (u *UserServiceImpl) ActivateUser(id uint) (*models.User, error) {
	user, err := u.userRepo.FindById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user.IsActive = true
	updatedUser, err := u.userRepo.Update(&user)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

// BulkUserActions implements UserService.
func (u *UserServiceImpl) BulkUserActions(req *BulkActionRequest) error {
	if len(req.UserIDs) == 0 {
		return errors.New("no user IDs provided")
	}

	switch req.Action {
	case "activate":
		return database.DB.Model(&models.User{}).Where("id IN ?", req.UserIDs).Update("is_active", true).Error
	case "deactivate":
		return database.DB.Model(&models.User{}).Where("id IN ?", req.UserIDs).Update("is_active", false).Error
	case "delete":
		return database.DB.Where("id IN ?", req.UserIDs).Delete(&models.User{}).Error
	default:
		return errors.New("invalid action")
	}
}

func (u *UserServiceImpl) CreateUser(req *models.UserInput) (*models.User, error) {

	_, err := u.userRepo.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("user with this email or username already exists")
	}

	role, err := u.roleRepo.FindById(req.RoleID)
	if err != nil {
		return nil, errors.New("role not found")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		RoleID:       role.ID,
		IsActive:     true,
	}

	userResult, err := u.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return &userResult, nil
}

// DeactivateUser implements UserService.
func (u *UserServiceImpl) DeactivateUser(id uint) (*models.User, error) {
	user, err := u.userRepo.FindById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user.IsActive = false
	updatedUser, err := u.userRepo.Update(&user)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

// DeleteUser implements UserService.
func (u *UserServiceImpl) DeleteUser(id uint) error {
	user, err := u.userRepo.FindById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	return u.userRepo.Delete(user)
}

// GetUserActivityLogs implements UserService.
func (u *UserServiceImpl) GetUserActivityLogs(req *models.ActivityLogListRequest) (*models.ActivityLogListResponse, error) {
	activityLogs, err := u.acitvityLogRepo.FindAll(req)
	if err != nil {
		return nil, err
	}
	return activityLogs, nil
}

// GetUserById implements UserService.
func (u *UserServiceImpl) GetUserById(id uint) (*models.User, error) {
	user, err := u.userRepo.FindById(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetUsers implements UserService.
func (u *UserServiceImpl) GetUsers(req models.UserListRequest) (*models.UserListResponse, error) {

	users, err := u.userRepo.FindAll(req)
	if err != nil {
		return nil, err
	}
	return users, nil

}

// UpdatePassword implements UserService.
func (u *UserServiceImpl) UpdatePassword(userID uint, req *models.PasswordUpdateInput) error {

	user, err := u.userRepo.FindById(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}
	if !utils.CheckPasswordHash(req.CurrentPassword, user.PasswordHash) {
		return errors.New("current password is incorrect")
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("failed to hash new password")
	}
	user.PasswordHash = hashedPassword
	_, err = u.userRepo.Update(&user)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUser implements UserService.
func (u *UserServiceImpl) UpdateUser(id uint, req *models.UserUpdateInput) (*models.User, error) {

	user, err := u.userRepo.FindById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if req.Email != "" && req.Email != user.Email {
		var existingUser models.User
		if err := database.DB.Where("email = ? AND id != ?", req.Email, id).First(&existingUser).Error; err == nil {
			return nil, errors.New("user with this email already exists")
		}
		user.Email = req.Email
	}

	if req.Username != "" && req.Username != user.Username {
		var existingUser models.User
		if err := database.DB.Where("username = ? AND id != ?", req.Username, id).First(&existingUser).Error; err == nil {
			return nil, errors.New("user with this username already exists")
		}
		user.Username = req.Username
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.RoleID != 0 {
		_, err := u.roleRepo.FindById(req.RoleID)
		if err != nil {
			return nil, errors.New("role not found")
		}
		user.RoleID = req.RoleID
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	updatedUser, err := u.userRepo.Update(&user)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func NewUserService(userRepo repository.UserRepository,
	acitvityLogRepo repository.ActivityLogRepository, roleRepo repository.RoleRepository) UserService {
	return &UserServiceImpl{
		userRepo: userRepo, acitvityLogRepo: acitvityLogRepo, roleRepo: roleRepo,
	}
}
