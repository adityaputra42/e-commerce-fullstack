package services

import (
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/repository"
	"errors"
	"fmt"
)

type OrderService interface {
	UpdateOrder(param models.UpdateOrder) (*models.OrderResponse, error)
	FindAllOrder(param models.OrderListRequest) (*[]models.OrderResponse, error)
	FindById(id string, userId int64) (*models.OrderResponse, error)
	DeleteOrder(id string) error
	CancelOrder(id string, userId int64) (*models.OrderResponse, error)
}

type OrderServiceImpl struct {
	orderRepo repository.OrderRepository
}

func (o *OrderServiceImpl) CancelOrder(id string, userId int64) (*models.OrderResponse, error) {
	order, err := o.orderRepo.FindById(id)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if order.UserID != userId {
		return nil, errors.New("unauthorized: order does not belong to user")
	}

	if !isValidStatusForCancel(order.Status) {
		return nil, errors.New("cannot cancel order with current status")
	}

	order.Status = "cancelled"

	updatedOrder, err := o.orderRepo.Update(order, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel order: %w", err)
	}

	orderResponse := updatedOrder.ToOrderResponse()

	return &orderResponse, nil
}

// DeleteOrder implements [OrderService].
func (o *OrderServiceImpl) DeleteOrder(id string) error {
	order, err := o.orderRepo.FindById(id)
	if err != nil {
		return errors.New("order not found")
	}

	err = o.orderRepo.Delete(order)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	return nil
}

// FindAllOrder implements [OrderService].
func (o *OrderServiceImpl) FindAllOrder(param models.OrderListRequest) (*[]models.OrderResponse, error) {
	orders, err := o.orderRepo.FindAll(param)
	if err != nil {
		return nil, fmt.Errorf("failed to get order list: %w", err)
	}

	var orderResponse []models.OrderResponse

	for _, order := range orders {
		orderResponse = append(orderResponse, order.ToOrderResponse())
	}

	return &orderResponse, nil
}

// FindById implements [OrderService].
func (o *OrderServiceImpl) FindById(id string, userId int64) (*models.OrderResponse, error) {
	order, err := o.orderRepo.FindById(id)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if order.UserID != userId {
		return nil, errors.New("unauthorized: order does not belong to user")
	}
	orderResponse := order.ToOrderResponse()
	return &orderResponse, nil
}

// UpdateOrder implements [OrderService].
func (o *OrderServiceImpl) UpdateOrder(param models.UpdateOrder) (*models.OrderResponse, error) {
	order, err := o.orderRepo.FindById(param.ID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	order.Status = param.Status

	updatedOrder, err := o.orderRepo.Update(order, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	result := updatedOrder.ToOrderResponse()

	return &result, nil
}

// Helper function untuk validasi status yang bisa di-cancel
func isValidStatusForCancel(status string) bool {
	validStatuses := []string{"pending", "confirmed", "processing"}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

func NewOrderService(OrderRepo repository.OrderRepository) OrderService {
	return &OrderServiceImpl{orderRepo: OrderRepo}
}
