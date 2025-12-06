package services

import (
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/repository"
)

type OrderService interface {
	UpdateOrder(param models.UpdateOrder) (*models.OrderResponse, error)
	FindAllOrder(param models.OrderListRequest) (*[]models.OrderResponse, error)
	FindById(id, userId int64) (*models.OrderResponse, error)
	DeleteAddress(id int64) error
	CancelOrder(id, userId int64) (*models.OrderResponse, error)
}

type OrderServiceImpl struct {
	orderRepo repository.OrderRepository
}

// CancelOrder implements [OrderService].
func (o *OrderServiceImpl) CancelOrder(id int64, userId int64) (*models.OrderResponse, error) {
	panic("unimplemented")
}

// DeleteAddress implements [OrderService].
func (o *OrderServiceImpl) DeleteAddress(id int64) error {
	panic("unimplemented")
}

// FindAllOrder implements [OrderService].
func (o *OrderServiceImpl) FindAllOrder(param models.OrderListRequest) (*[]models.OrderResponse, error) {
	panic("unimplemented")
}

// FindById implements [OrderService].
func (o *OrderServiceImpl) FindById(id int64, userId int64) (*models.OrderResponse, error) {
	panic("unimplemented")
}

// UpdateOrder implements [OrderService].
func (o *OrderServiceImpl) UpdateOrder(param models.UpdateOrder) (*models.OrderResponse, error) {
	panic("unimplemented")
}

func NewOrderService(OrderRepo repository.OrderRepository) OrderService {
	return &OrderServiceImpl{}
}
