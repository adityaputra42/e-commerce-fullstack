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
	FindAllShipping(param models.ShippingListRequest) ([]models.ShippingResponse, error)
	FindByID(id int64) (*models.ShippingResponse, error)
	DeleteShipping(id int64) error
}

type ShippingServiceImpl struct {
	shippingRepo repository.ShippingRepository
}

func NewShippingService(shippingRepo repository.ShippingRepository) ShippingService {
	return &ShippingServiceImpl{
		shippingRepo: shippingRepo,
	}
}

// --- helpers ---
func validateShipping(name, state string, price float64) error {
	if name == "" {
		return errors.New("shipping name is required")
	}
	if state == "" {
		return errors.New("shipping state is required")
	}
	if price < 0 {
		return errors.New("shipping price must be >= 0")
	}
	return nil
}

// --- services ---

func (s *ShippingServiceImpl) CreateShipping(param models.CreateShipping) (*models.ShippingResponse, error) {
	if err := validateShipping(param.Name, param.State, param.Price); err != nil {
		return nil, err
	}

	shipping := models.Shipping{
		Name:  param.Name,
		Price: param.Price,
		State: param.State,
	}

	created, err := s.shippingRepo.Create(shipping, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create shipping: %w", err)
	}

	return created.ToResponseShipping(), nil
}

func (s *ShippingServiceImpl) UpdateShipping(param models.UpdateShipping) (*models.ShippingResponse, error) {
	if param.ID <= 0 {
		return nil, errors.New("invalid shipping id")
	}

	existing, err := s.shippingRepo.FindById(uint(param.ID))
	if err != nil {
		return nil, err
	}

	if err := validateShipping(param.Name, param.State, param.Price); err != nil {
		return nil, err
	}

	existing.Name = param.Name
	existing.Price = param.Price
	existing.State = param.State

	updated, err := s.shippingRepo.Update(*existing, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to update shipping: %w", err)
	}

	return updated.ToResponseShipping(), nil
}

func (s *ShippingServiceImpl) FindAllShipping(param models.ShippingListRequest) ([]models.ShippingResponse, error) {
	shippings, err := s.shippingRepo.FindAll(param)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipping list: %w", err)
	}

	results := make([]models.ShippingResponse, 0, len(shippings))
	for _, v := range shippings {
		results = append(results, *v.ToResponseShipping())
	}

	return results, nil
}

func (s *ShippingServiceImpl) FindByID(id int64) (*models.ShippingResponse, error) {
	if id <= 0 {
		return nil, errors.New("invalid shipping id")
	}

	shipping, err := s.shippingRepo.FindById(uint(id))
	if err != nil {
		return nil, err
	}

	return shipping.ToResponseShipping(), nil
}

func (s *ShippingServiceImpl) DeleteShipping(id int64) error {
	if id <= 0 {
		return errors.New("invalid shipping id")
	}

	shipping, err := s.shippingRepo.FindById(uint(id))
	if err != nil {
		return err
	}

	if err := s.shippingRepo.Delete(*shipping); err != nil {
		return fmt.Errorf("failed to delete shipping: %w", err)
	}

	return nil
}
