package services_test

import (
	"e-commerce/backend/internal/mocks"
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/services"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

func TestOrderService_CancelOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepository(ctrl)
	service := services.NewOrderService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		orderID := "ORD-123"
		userID := int64(1)
		
		existingOrder := models.Order{
			ID:     orderID,
			UserID: userID,
			Status: "pending",
		}

		mockRepo.EXPECT().FindById(orderID).Return(existingOrder, nil)
		
		mockRepo.EXPECT().
			Update(gomock.Any(), nil).
			DoAndReturn(func(o models.Order, tx interface{}) (models.Order, error) {
				o.Status = "cancelled" 
				o.UpdatedAt = time.Now()
				return o, nil
			})

		result, err := service.CancelOrder(orderID, userID)
		
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if result.Status != "cancelled" {
			t.Errorf("Expected status cancelled, got %s", result.Status)
		}
	})

	t.Run("Unauthorized", func(t *testing.T) {
		orderID := "ORD-123"
		userID := int64(1)
		otherUser := int64(2)
		
		existingOrder := models.Order{
			ID:     orderID,
			UserID: otherUser,
			Status: "pending",
		}

		mockRepo.EXPECT().FindById(orderID).Return(existingOrder, nil)

		_, err := service.CancelOrder(orderID, userID)
		
		if err == nil {
			t.Errorf("Expected error Unauthorized")
		}
	})

	t.Run("Invalid Status", func(t *testing.T) {
		orderID := "ORD-123"
		userID := int64(1)
		
		existingOrder := models.Order{
			ID:     orderID,
			UserID: userID,
			Status: "completed", 
		}

		mockRepo.EXPECT().FindById(orderID).Return(existingOrder, nil)

		_, err := service.CancelOrder(orderID, userID)
		
		if err == nil {
			t.Errorf("Expected error Invalid Status")
		}
	})
}

func TestOrderService_FindAllOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepository(ctrl)
	service := services.NewOrderService(mockRepo)

	req := models.OrderListRequest{
		Page:  1,
		Limit: 10,
	}

	orders := []models.Order{
		{ID: "1", Status: "pending"},
		{ID: "2", Status: "completed"},
	}

	mockRepo.EXPECT().FindAll(req).Return(orders, nil)

	result, err := service.FindAllOrder(req)
	
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 orders, got %d", len(result))
	}
}
