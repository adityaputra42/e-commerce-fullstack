package services

import (
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/repository"
	"e-commerce/backend/internal/utils"
	"errors"
	"time"

	"gorm.io/gorm"
)

type AuthService interface {
	SignIn(req models.LoginRequest, ipAddress, userAgent string) (*models.TokenResponse, error)
	SignUp(req models.RegisterRequest) (*models.TokenResponse, error)
	ForgotPassword(req models.ForgotPasswordRequest) error
	ResetPassword(req models.ResetPasswordRequest) error
	RefreshToken(req models.RefreshTokenRequest) (*models.TokenResponse, error)
	generateTokenResponse(user *models.User) (*models.TokenResponse, error)
}
type AuthServiceImpl struct {
	jwtService      *utils.JWTService
	userRepo        repository.UserRepository
	acitvityLogRepo repository.ActivityLogRepository
	roleRepo        repository.RoleRepository
}

// generateTokenResponse implements AuthService.
func (a *AuthServiceImpl) generateTokenResponse(user *models.User) (*models.TokenResponse, error) {
	accessToken, expiresAt, err := a.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, _, err := a.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User:         *user.ToResponse(),
	}, nil
}

// ForgotPassword implements AuthService.
func (a *AuthServiceImpl) ForgotPassword(req models.ForgotPasswordRequest) error {
	panic("unimplemented")
}

// RefreshToken implements AuthService.
func (a *AuthServiceImpl) RefreshToken(req models.RefreshTokenRequest) (*models.TokenResponse, error) {
	claims, err := a.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	user, err := a.userRepo.FindById(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	return a.generateTokenResponse(&user)
}

// ResetPassword implements AuthService.
func (a *AuthServiceImpl) ResetPassword(req models.ResetPasswordRequest) error {
	panic("unimplemented")
}

// SignIn implements AuthService.
func (a *AuthServiceImpl) SignIn(req models.LoginRequest, ipAddress, userAgent string) (*models.TokenResponse, error) {
	var userResult models.User
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		userResult, err := a.userRepo.FindByEmail(req.Email)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("invalid credentials")
			}
			return err
		}

		if !userResult.IsActive {
			return errors.New("account is deactivated")
		}

		if err := utils.CheckPassword(req.Password, userResult.PasswordHash); err != nil {
			return errors.New("invalid credentials")
		}

		now := time.Now()
		userResult.LastLoginAt = &now
		database.DB.Save(&userResult)

		activityLog := models.ActivityLog{
			UserID:    userResult.ID,
			Action:    "login",
			Resource:  "auth",
			Details:   "User logged in successfully",
			IPAddress: ipAddress,
			UserAgent: userAgent,
		}
		_, err = a.acitvityLogRepo.Create(activityLog, tx)
		if err != nil {
			return errors.New("failed create activity log")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return a.generateTokenResponse(&userResult)
}

// SignUp implements AuthService.
func (a *AuthServiceImpl) SignUp(req models.RegisterRequest) (*models.TokenResponse, error) {
	var userResult models.User
	err := database.DB.Transaction(func(tx *gorm.DB) error {

		_, err := a.userRepo.FindByEmail(req.Email)

		if err == nil {
			return errors.New("user with this email already exists")
		}

		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return err
		}

		defaultRole, err := a.roleRepo.FindByName("User")
		if err != nil {
			return errors.New("default role not found")
		}

		user := models.User{
			Email:        req.Email,
			Username:     req.Username,
			PasswordHash: hashedPassword,
			FirstName:    req.FirstName,
			LastName:     req.LastName,
			RoleID:       defaultRole.ID,
			IsActive:     true,
		}

		userResult, err = a.userRepo.Create(user, tx)
		return err

	})
	if err != nil {
		return nil, err
	}

	return a.generateTokenResponse(&userResult)
}

func NewAuthService(jwtService *utils.JWTService, userRepo repository.UserRepository, acitvityLogRepo repository.ActivityLogRepository, roleRepo repository.RoleRepository) AuthService {
	return &AuthServiceImpl{
		jwtService:      jwtService,
		userRepo:        userRepo,
		acitvityLogRepo: acitvityLogRepo,
		roleRepo:        roleRepo,
	}
}
