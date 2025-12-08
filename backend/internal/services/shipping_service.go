package services

import (
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/repository"
	"errors"
	"fmt"
)

type ShippingService interface {
	CreateShipping(param models.CreateShipping) (*models.ShippingResponse, error)
	UpdateShipping(param models.UpdateShipping) (*models.ShippingResponse, error)
	FindAllShipping(param models.ShippingListRequest) (*[]models.ShippingResponse, error)
	FindById(id int64) (*models.ShippingResponse, error)
	DeleteShipping(id int64) error
}

type ShippingServiceImpl struct {
	shippingRepo repository.ShippingRepository
}

// CreateShipping implements ShippingService.
func (s *ShippingServiceImpl) CreateShipping(param models.CreateShipping) (*models.ShippingResponse, error) {
	// Validasi input
	if param.Name == "" {
		return nil, errors.New("shipping name is required")
	}

	if param.Price < 0 {
		return nil, errors.New("shipping price must be greater than or equal to 0")
	}

	if param.State == "" {
		return nil, errors.New("shipping state is required")
	}
	shippingParam := models.Shipping{
		Name:  param.Name,
		Price: param.Price,
		State: param.State,
	}

	// Create shipping
	shipping, err := s.shippingRepo.Create(shippingParam, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create shipping: %w", err)
	}

	result := shipping.ToResponseShipping()
	return result, nil
}

// UpdateShipping implements ShippingService.
func (s *ShippingServiceImpl) UpdateShipping(param models.UpdateShipping) (*models.ShippingResponse, error) {
	// Validasi shipping exists
	shipping, err := s.shippingRepo.FindById(uint(param.ID))
	if err != nil {
		return nil, errors.New("shipping not found")
	}

	// Validasi input
	if param.Name == "" {
		return nil, errors.New("shipping name is required")
	}

	if param.Price < 0 {
		return nil, errors.New("shipping price must be greater than or equal to 0")
	}

	if param.State == "" {
		return nil, errors.New("shipping state is required")
	}

	shipping.Name = param.Name
	shipping.Price = param.Price
	shipping.State = param.State

	updatedShipping, err := s.shippingRepo.Update(*shipping, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to update shipping: %w", err)
	}

	result := updatedShipping.ToResponseShipping()

	return result, nil
}

// FindAllShipping implements ShippingService.
func (s *ShippingServiceImpl) FindAllShipping(param models.ShippingListRequest) (*[]models.ShippingResponse, error) {

	shippings, err := s.shippingRepo.FindAll(param)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipping list: %w", err)
	}
	var results []models.ShippingResponse
	for _, v := range shippings {
		results = append(results, *v.ToResponseShipping())
	}
	return &results, nil
}

// FindById implements ShippingService.
func (s *ShippingServiceImpl) FindById(id int64) (*models.ShippingResponse, error) {
	if id <= 0 {
		return nil, errors.New("invalid shipping id")
	}

	shipping, err := s.shippingRepo.FindById(uint(id))
	if err != nil {
		return nil, errors.New("shipping not found")
	}
	result := shipping.ToResponseShipping()
	return result, nil
}

// DeleteShipping implements ShippingService.
func (s *ShippingServiceImpl) DeleteShipping(id int64) error {
	if id <= 0 {
		return errors.New("invalid shipping id")
	}

	shipping, err := s.shippingRepo.FindById(uint(id))
	if err != nil {
		return errors.New("shipping not found")
	}

	err = s.shippingRepo.Delete(*shipping)
	if err != nil {
		return fmt.Errorf("failed to delete shipping: %w", err)
	}

	return nil
}

func NewShippingService(shippingRepo repository.ShippingRepository) ShippingService {
	return &ShippingServiceImpl{
		shippingRepo: shippingRepo,
	}
}
