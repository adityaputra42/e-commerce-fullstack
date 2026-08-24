package services

import (
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/repository"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type OrderService interface {
	UpdateOrder(param models.UpdateOrder) (*models.OrderResponse, error)
	FindAllOrder(param models.OrderListRequest) ([]models.OrderResponse, error)
	FindById(id string, userId int64) (*models.OrderResponse, error)
	DeleteOrder(id string) error
	CancelOrder(id string, userId int64) (*models.OrderResponse, error)
}

type OrderServiceImpl struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
}

func (o *OrderServiceImpl) CancelOrder(id string, userId int64) (*models.OrderResponse, error) {
	var result models.OrderResponse

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		order, err := o.orderRepo.FindByIdLocking(tx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("order not found")
			}
			return fmt.Errorf("error mencari order: %w", err)
		}

		if order.UserID != userId {
			return errors.New("unauthorized: order does not belong to user")
		}

		if !isValidStatusForCancel(order.Status) {
			return errors.New("cannot cancel order with current status")
		}

		sizeVariant, err := o.productRepo.FindSizeVarianLocked(tx, uint(order.SizeVarianID))
		if err != nil {
			return fmt.Errorf("gagal mengambil size variant untuk restock: %w", err)
		}

		sizeVariant.Stock += order.Quantity
		if _, err := o.productRepo.UpdateSizeVarian(*sizeVariant, tx); err != nil {
			return fmt.Errorf("gagal mengembalikan stock: %w", err)
		}

		order.Status = "cancelled"
		updatedOrder, err := o.orderRepo.Update(order, tx)
		if err != nil {
			return fmt.Errorf("failed to cancel order: %w", err)
		}

		result = updatedOrder.ToOrderResponse()
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &result, nil
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
func (o *OrderServiceImpl) FindAllOrder(param models.OrderListRequest) ([]models.OrderResponse, error) {
	orders, err := o.orderRepo.FindAll(param)
	if err != nil {
		return nil, fmt.Errorf("failed to get order list: %w", err)
	}

	var orderResponse []models.OrderResponse

	for _, order := range orders {
		orderResponse = append(orderResponse, order.ToOrderResponse())
	}

	return orderResponse, nil
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

func (o *OrderServiceImpl) UpdateOrder(param models.UpdateOrder) (*models.OrderResponse, error) {
	var result models.OrderResponse

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		order, err := o.orderRepo.FindByIdLocking(tx, param.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("order not found")
			}
			return fmt.Errorf("error mencari order: %w", err)
		}

		becomingCancelled := param.Status == "cancelled" && order.Status != "cancelled"

		if becomingCancelled {
			sizeVariant, err := o.productRepo.FindSizeVarianLocked(tx, uint(order.SizeVarianID))
			if err != nil {
				return fmt.Errorf("gagal mengambil size variant untuk restock: %w", err)
			}
			sizeVariant.Stock += order.Quantity
			if _, err := o.productRepo.UpdateSizeVarian(*sizeVariant, tx); err != nil {
				return fmt.Errorf("gagal mengembalikan stock: %w", err)
			}
		}

		order.Status = param.Status
		updatedOrder, err := o.orderRepo.Update(order, tx)
		if err != nil {
			return fmt.Errorf("failed to update order: %w", err)
		}

		result = updatedOrder.ToOrderResponse()
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func isValidStatusForCancel(status string) bool {
	validStatuses := []string{"pending", "confirmed", "processing"}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

func NewOrderService(OrderRepo repository.OrderRepository, ProductRepo repository.ProductRepository) OrderService {
	return &OrderServiceImpl{orderRepo: OrderRepo, productRepo: ProductRepo}
}
