package services_test

import (
	"e-commerce/backend/internal/mocks"
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/services"
	"e-commerce/backend/internal/testhelper"
	"testing"

	"go.uber.org/mock/gomock"
)

// TestOrderService_CancelOrder butuh koneksi database nyata (bukan cuma mock
// repository) karena CancelOrder membungkus operasinya dalam
// database.DB.Transaction(...) untuk row-locking saat restock — lihat
// komentar di order_service.go. Repository tetap di-mock supaya assertion
// query/update tetap presisi, sementara database.DB sendiri diarahkan ke test
// DB lewat testhelper.SetTestDB, sama seperti pola di TestAuthService_SignIn.
func TestOrderService_CancelOrder(t *testing.T) {
	testDB := testhelper.SetupTestSuite(t)
	tx := testhelper.BeginTestTransaction(t, testDB)
	defer testhelper.RollbackTestTransaction(tx)

	dbWrapper := testhelper.SetTestDB(tx)
	defer dbWrapper.Restore()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	service := services.NewOrderService(mockOrderRepo, mockProductRepo)

	t.Run("Success", func(t *testing.T) {
		orderID := "ORD-123"
		userID := int64(1)

		existingOrder := models.Order{
			ID:           orderID,
			UserID:       userID,
			Status:       "pending",
			SizeVarianID: 1,
			Quantity:     2,
		}
		sizeVarian := models.SizeVarian{ID: 1, Stock: 5}

		mockOrderRepo.EXPECT().FindByIdLocking(gomock.Any(), orderID).Return(existingOrder, nil)
		mockProductRepo.EXPECT().FindSizeVarianLocked(gomock.Any(), uint(1)).Return(&sizeVarian, nil)
		mockProductRepo.EXPECT().
			UpdateSizeVarian(gomock.Any(), gomock.Any()).
			DoAndReturn(func(sv models.SizeVarian, tx interface{}) (models.SizeVarian, error) {
				// Restock harus menambahkan quantity order ke stock lama (5 + 2 = 7).
				if sv.Stock != 7 {
					t.Errorf("Expected restocked stock 7, got %d", sv.Stock)
				}
				return sv, nil
			})
		mockOrderRepo.EXPECT().
			Update(gomock.Any(), gomock.Any()).
			DoAndReturn(func(o models.Order, tx interface{}) (models.Order, error) {
				o.Status = "cancelled"
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
		orderID := "ORD-124"
		userID := int64(1)
		otherUser := int64(2)

		existingOrder := models.Order{
			ID:     orderID,
			UserID: otherUser,
			Status: "pending",
		}

		mockOrderRepo.EXPECT().FindByIdLocking(gomock.Any(), orderID).Return(existingOrder, nil)

		_, err := service.CancelOrder(orderID, userID)

		if err == nil {
			t.Errorf("Expected error Unauthorized")
		}
	})

	t.Run("Invalid Status", func(t *testing.T) {
		orderID := "ORD-125"
		userID := int64(1)

		existingOrder := models.Order{
			ID:     orderID,
			UserID: userID,
			Status: "completed",
		}

		mockOrderRepo.EXPECT().FindByIdLocking(gomock.Any(), orderID).Return(existingOrder, nil)

		_, err := service.CancelOrder(orderID, userID)

		if err == nil {
			t.Errorf("Expected error Invalid Status")
		}
	})
}

func TestOrderService_FindAllOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	service := services.NewOrderService(mockOrderRepo, mockProductRepo)

	req := models.OrderListRequest{
		Page:  1,
		Limit: 10,
	}

	orders := []models.Order{
		{ID: "1", Status: "pending"},
		{ID: "2", Status: "completed"},
	}

	mockOrderRepo.EXPECT().FindAll(req).Return(orders, nil)

	result, err := service.FindAllOrder(req)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 orders, got %d", len(result))
	}
}
