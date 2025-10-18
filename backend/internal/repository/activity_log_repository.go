package repository

import (
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"
	"math"

	"gorm.io/gorm"
)

type ActivityLogRepository interface {
	FindAll(param *models.ActivityLogListRequest) (*models.ActivityLogListResponse, error)
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
func (a *ActivityLogRepositoryImpl) FindAll(param *models.ActivityLogListRequest) (*models.ActivityLogListResponse, error) {
	offset := (param.Page - 1) * param.Limit

	var total int64
	if err := database.DB.Model(&models.ActivityLog{}).Where("user_id = ?", param.UserId).Count(&total).Error; err != nil {
		return nil, err
	}

	var activityLogs []models.ActivityLog
	if err := database.DB.Preload("User").Where("user_id = ?", param.UserId).
		Order("created_at desc").Offset(offset).Limit(param.Limit).Find(&activityLogs).Error; err != nil {
		return nil, err
	}

	activities := make([]models.ActivityLogResponse, len(activityLogs))
	for i, log := range activityLogs {
		activities[i] = *log.ToResponse()
	}

	totalPages := int(math.Ceil(float64(total) / float64(param.Limit)))

	return &models.ActivityLogListResponse{
		Activities: activities,
		Total:      total,
		Page:       param.Page,
		Limit:      param.Limit,
		TotalPages: totalPages,
	}, nil
}

func NewActivityLogRepository() ActivityLogRepository {
	return &ActivityLogRepositoryImpl{}
}
