package repository

import (
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(param models.Category, tx *gorm.DB) (models.Category, error)
	Update(param models.Category, tx *gorm.DB) (models.Category, error)
	Delete(param models.Category) error
	FindById(paramId int64) (models.Category, error)
	FindByIds(paramId []int64) ([]models.Category, error)
	FindAll(param models.CategoryListRequest) ([]models.Category, error)
}

type CategoryRepositoryImpl struct {
}

// FindByIds implements CategoryRepository.
func (a *CategoryRepositoryImpl) FindByIds(paramId []int64) ([]models.Category, error) {
	var categories []models.Category

	if len(paramId) == 0 {
		return categories, nil
	}

	err := database.DB.
		Where("id IN ? AND deleted_at IS NULL", paramId).
		Find(&categories).Error

	if err != nil {
		return nil, err
	}

	return categories, nil
}

// Create implements CategoryRepository.
func (a *CategoryRepositoryImpl) Create(param models.Category, tx *gorm.DB) (models.Category, error) {
	var result models.Category

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

// Delete implements CategoryRepository.
func (a *CategoryRepositoryImpl) Delete(param models.Category) error {
	return database.DB.Delete(&param).Error
}

// FindAll implements CategoryRepository.
func (a *CategoryRepositoryImpl) FindAll(param models.CategoryListRequest) ([]models.Category, error) {

	offset := (param.Page - 1) * param.Limit

	var Categorys []models.Category
	db := database.DB

	if param.Search != "" {
		db = db.Where("name ILIKE ?", "%"+param.Search+"%")
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

	if err := db.Preload("Categorys").Find(&Categorys).Error; err != nil {
		return nil, err
	}

	return Categorys, nil
}

// FindById implements CategoryRepository.
func (a *CategoryRepositoryImpl) FindById(paramId int64) (models.Category, error) {
	Category := models.Category{}
	err := database.DB.Model(&models.Category{}).Take(&Category, "id =?", paramId).Error

	return Category, err
}

// Update implements CategoryRepository.
func (a *CategoryRepositoryImpl) Update(param models.Category, tx *gorm.DB) (models.Category, error) {
	var result models.Category

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

func NewCategoryRepository() CategoryRepository {
	return &CategoryRepositoryImpl{}
}
