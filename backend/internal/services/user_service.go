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
	GetUserActivityLogs(req models.AddressListRequest) (*models.ActivityLogListResponse, error)
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
	panic("unimplemented")
}

// CreateUser implements UserService.
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
	panic("unimplemented")
}

// GetUserActivityLogs implements UserService.
func (u *UserServiceImpl) GetUserActivityLogs(req models.AddressListRequest) (*models.ActivityLogListResponse, error) {
	panic("unimplemented")
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

	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *user.ToResponse()
	}

	// totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &models.UserListResponse{
		Users: userResponses,
		// Total:      total,
		Page:  req.Offset,
		Limit: req.Limit,
		// TotalPages: totalPages,
	}, nil

}

// UpdatePassword implements UserService.
func (u *UserServiceImpl) UpdatePassword(userID uint, req *models.PasswordUpdateInput) error {
	panic("unimplemented")
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
