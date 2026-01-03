package repository

import (
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"
	"time"
)

type PasswordResetTokenRepository interface {
	Create(param *models.PasswordResetToken) (models.PasswordResetToken, error)
	Update(param *models.PasswordResetToken) (models.PasswordResetToken, error)
	Delete(param *models.PasswordResetToken) error
	FindByToken(token string) (models.PasswordResetToken, error)
	FindAll(param *models.PasswordResetTokenListRequest) ([]models.PasswordResetToken, error)
}

type PasswordResetTokenRepositoryImpl struct {
}

// Create implements PasswordResetTokenRepository.
func (p *PasswordResetTokenRepositoryImpl) Create(param *models.PasswordResetToken) (models.PasswordResetToken, error) {
	var result models.PasswordResetToken

	db := database.DB

	err := db.Create(&param).Error
	if err != nil {
		return result, err
	}

	// Reload the record to get the complete data including auto-generated fields
	err = db.
		Select("id", "user_id", "token", "expires_at", "used_at", "created_at").
		First(&result, param.ID).Error
	return result, err
}

// Delete implements PasswordResetTokenRepository.
func (p *PasswordResetTokenRepositoryImpl) Delete(param *models.PasswordResetToken) error {
	return database.DB.Delete(&param).Error
}

// FindAll implements PasswordResetTokenRepository.
func (p *PasswordResetTokenRepositoryImpl) FindAll(param *models.PasswordResetTokenListRequest) ([]models.PasswordResetToken, error) {
	offset := (param.Page - 1) * param.Limit

	var tokens []models.PasswordResetToken
	db := database.DB.
		Select("id", "user_id", "token", "expires_at", "used_at", "created_at")

	if param.UserId != nil {
		db = db.Where("user_id = ?", param.UserId)
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

	if err := db.Find(&tokens).Error; err != nil {
		return nil, err
	}

	return tokens, nil
}

// FindByToken implements PasswordResetTokenRepository.
func (p *PasswordResetTokenRepositoryImpl) FindByToken(token string) (models.PasswordResetToken, error) {
	resetPassword := models.PasswordResetToken{}
	err := database.DB.
		Select("id", "user_id", "token", "expires_at", "used_at", "created_at").
		Where("token = ? AND expires_at > ?", token, time.Now()).
		First(&resetPassword).Error

	return resetPassword, err
}

// Update implements PasswordResetTokenRepository.
func (p *PasswordResetTokenRepositoryImpl) Update(param *models.PasswordResetToken) (models.PasswordResetToken, error) {
	var result models.PasswordResetToken

	db := database.DB

	err := db.Model(&param).Updates(param).Error
	if err != nil {
		return result, err
	}

	err = db.
		Select("id", "user_id", "token", "expires_at", "used_at", "created_at").
		First(&result, param.ID).Error
	return result, err
}

func NewPasswordResetTokenRepository() PasswordResetTokenRepository {
	return &PasswordResetTokenRepositoryImpl{}
}
