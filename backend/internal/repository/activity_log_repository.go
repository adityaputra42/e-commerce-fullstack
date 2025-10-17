package repository

import (
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"

	"gorm.io/gorm"
)

type ActivityLogRepository interface {
	FindAll(param models.ActivityLogListRequest) ([]models.ActivityLog, error)
	Create(param models.ActivityLog, tx *gorm.DB) (models.ActivityLog, error)
}
type ActivityLogRepositoryImpl struct {
}

// Create implements ActivityLogRepository.
func (a *ActivityLogRepositoryImpl) Create(param models.ActivityLog, tx *gorm.DB) (models.ActivityLog, error) {
	var result models.ActivityLog

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

// FindAll implements ActivityLogRepository.
func (a *ActivityLogRepositoryImpl) FindAll(param models.ActivityLogListRequest) ([]models.ActivityLog, error) {
	var result []models.ActivityLog
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

	if err := db.Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func NewActivityLogRepository() ActivityLogRepository {
	return &ActivityLogRepositoryImpl{}
}
