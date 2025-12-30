package services

import (
	"crypto/rand"
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/repository"
	"e-commerce/backend/internal/utils"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	SignIn(req models.LoginRequest, ipAddress, userAgent string) (*models.TokenResponse, error)
	SignUp(req models.RegisterRequest) (*models.TokenResponse, error)
	ForgotPassword(req models.ForgotPasswordRequest) (string, error)
	ResetPassword(req models.ResetPasswordRequest) error
	RefreshToken(req models.RefreshTokenRequest) (*models.TokenResponse, error)
}
type AuthServiceImpl struct {
	jwtService         *utils.JWTService
	userRepo           repository.UserRepository
	acitvityLogRepo    repository.ActivityLogRepository
	roleRepo           repository.RoleRepository
	passwordRepository repository.PasswordResetTokenRepository
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
func (a *AuthServiceImpl) ForgotPassword(req models.ForgotPasswordRequest) (string, error) {

	user, err := a.userRepo.FindByEmail(req.Email)
	if err != nil {
		return "", errors.New("user not found")
	}

	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	tokenStr := hex.EncodeToString(token)

	resetToken := models.PasswordResetToken{
		UserID:    user.ID,
		Token:     tokenStr,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	_, err = a.passwordRepository.Create(&resetToken)
	if err != nil {
		return "", err
	}

	// Here you would typically send an email with the reset link.
	// For this example, we'll just return the token.
	// log.Printf("Password reset token for %s: %s", email, tokenStr)

	return tokenStr, nil
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

	resetToken, err := a.passwordRepository.FindByToken(req.Token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	var user models.User
	if err := database.DB.First(&user, resetToken.UserID).Error; err != nil {
		return errors.New("user not found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	_, err = a.userRepo.Update(&user)
	if err != nil {
		return err
	}

	err = a.passwordRepository.Delete(&resetToken)

	if err != nil {
		return err
	}

	return nil
}

// SignIn implements AuthService.
func (a *AuthServiceImpl) SignIn(req models.LoginRequest, ipAddress, userAgent string) (*models.TokenResponse, error) {

	var userResult models.User

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var err error

		userResult, err = a.userRepo.FindByEmail(req.Email)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("invalid credentials")
			}
			return err
		}

		if userResult.ID == 0 {
			log.Printf("❌ LOGIN ERROR: user.ID is 0 for email: %s", req.Email)
			return errors.New("invalid credentials")
		}

		if !userResult.IsActive {
			return errors.New("account is deactivated")
		}

		if err := utils.CheckPassword(req.Password, userResult.PasswordHash); err != nil {
			return errors.New("invalid credentials")
		}

		now := time.Now()
		userResult.LastLoginAt = &now

		if err := tx.Save(&userResult).Error; err != nil {
			return err
		}

		activityLog := models.ActivityLog{
			UserID:    userResult.ID,
			Action:    "login",
			Resource:  "auth",
			Details:   "User logged in successfully",
			IPAddress: ipAddress,
			UserAgent: userAgent,
		}

		if _, err := a.acitvityLogRepo.Create(activityLog, tx); err != nil {
			return err
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

	_, err := a.userRepo.FindByEmail(req.Email)

	if err == nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	defaultRole, err := a.roleRepo.FindByName("User")
	if err != nil {
		return nil, errors.New("default role not found")
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

	userResult, err := a.userRepo.Create(user)

	if err != nil {
		return nil, err
	}

	return a.generateTokenResponse(&userResult)
}

func NewAuthService(jwtService *utils.JWTService, userRepo repository.UserRepository, acitvityLogRepo repository.ActivityLogRepository, roleRepo repository.RoleRepository, passwordRepository repository.PasswordResetTokenRepository) AuthService {
	return &AuthServiceImpl{
		jwtService:         jwtService,
		userRepo:           userRepo,
		acitvityLogRepo:    acitvityLogRepo,
		roleRepo:           roleRepo,
		passwordRepository: passwordRepository,
	}
}
