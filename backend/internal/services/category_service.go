package services

import (
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/repository"
	"fmt"

	"gorm.io/gorm"
)

type CategoryService interface {
	Create(param models.Category) (models.Category, error)
	Update(param models.Category) (models.Category, error)
	Delete(id int64) error
	GetById(id int64) (models.Category, error)
	GetByIds(ids []int64) ([]models.Category, error)
	GetAll(param models.CategoryListRequest) ([]models.Category, error)
}

type CategoryServiceImpl struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &CategoryServiceImpl{
		repo: repo,
	}
}

// Create implements CategoryService.
func (s *CategoryServiceImpl) Create(param models.Category) (models.Category, error) {
	if param.Name == "" {
		return models.Category{}, fmt.Errorf("category name cannot be empty")
	}

	result, err := s.repo.Create(param, nil)
	if err != nil {
		return models.Category{}, fmt.Errorf("failed to create category: %w", err)
	}

	return result, nil
}

// Update implements CategoryService.
func (s *CategoryServiceImpl) Update(param models.Category) (models.Category, error) {
	if param.ID == 0 {
		return models.Category{}, fmt.Errorf("category id cannot be empty")
	}

	if param.Name == "" {
		return models.Category{}, fmt.Errorf("category name cannot be empty")
	}

	// Check if category exists
	existing, err := s.repo.FindById(param.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.Category{}, fmt.Errorf("category not found")
		}
		return models.Category{}, fmt.Errorf("failed to check category: %w", err)
	}

	// Preserve fields that shouldn't be updated
	param.CreatedAt = existing.CreatedAt

	result, err := s.repo.Update(param, nil)
	if err != nil {
		return models.Category{}, fmt.Errorf("failed to update category: %w", err)
	}

	return result, nil
}

// Delete implements CategoryService.
func (s *CategoryServiceImpl) Delete(id int64) error {
	if id == 0 {
		return fmt.Errorf("category id cannot be empty")
	}

	// Check if category exists
	category, err := s.repo.FindById(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("category not found")
		}
		return fmt.Errorf("failed to check category: %w", err)
	}

	err = s.repo.Delete(category)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

// GetById implements CategoryService.
func (s *CategoryServiceImpl) GetById(id int64) (models.Category, error) {
	if id == 0 {
		return models.Category{}, fmt.Errorf("category id cannot be empty")
	}

	category, err := s.repo.FindById(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.Category{}, fmt.Errorf("category not found")
		}
		return models.Category{}, fmt.Errorf("failed to get category: %w", err)
	}

	return category, nil
}

// GetByIds implements CategoryService.
func (s *CategoryServiceImpl) GetByIds(ids []int64) ([]models.Category, error) {
	if len(ids) == 0 {
		return []models.Category{}, fmt.Errorf("category ids cannot be empty")
	}

	categories, err := s.repo.FindByIds(ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	return categories, nil
}

// GetAll implements CategoryService.
func (s *CategoryServiceImpl) GetAll(param models.CategoryListRequest) ([]models.Category, error) {
	if param.Page < 1 {
		param.Page = 1
	}

	if param.Limit < 1 {
		param.Limit = 10
	}

	if param.Limit > 100 {
		param.Limit = 100
	}

	categories, err := s.repo.FindAll(param)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	return categories, nil
}
