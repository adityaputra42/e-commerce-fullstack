package services_test

import (
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/repository"
	"e-commerce/backend/internal/services"
	"e-commerce/backend/internal/testhelper"
	"testing"

	"github.com/google/uuid"
)
func TestOrderFulfillmentService_SyncOrdersToTransactionStatus(t *testing.T) {
	testDB := testhelper.SetupTestSuite(t)
	tx := testhelper.BeginTestTransaction(t, testDB)
	defer testhelper.RollbackTestTransaction(tx)

	role := testhelper.CreateTestRole(tx, "customer-fulfillment-test")
	user := testhelper.CreateTestUser(tx, "fulfillment-test@example.com", role.ID)
	category := testhelper.CreateTestCategory(tx, "Fulfillment Test Category")
	product := testhelper.CreateTestProduct(tx, "Fulfillment Test Product", category.ID, 100000)

	colorVarian := &models.ColorVarian{ProductID: product.ID, Name: "Black", Color: "#000000", Images: "https://example.com/black.png"}
	if err := tx.Create(colorVarian).Error; err != nil {
		t.Fatalf("failed to create color varian: %v", err)
	}

	sizeVarian := &models.SizeVarian{ColorVarianID: colorVarian.ID, Size: "M", Stock: 10}
	if err := tx.Create(sizeVarian).Error; err != nil {
		t.Fatalf("failed to create size varian: %v", err)
	}

	address := &models.Address{
		UserID:               int64(user.ID),
		RecipientName:        "Test Recipient",
		RecipientPhoneNumber: "+6281234567890",
		Province:             "Jawa Barat",
		City:                 "Bandung",
		District:             "Coblong",
		Village:              "Dago",
		PostalCode:           "40135",
		FullAddress:          "Jl. Test No. 123, Bandung",
	}
	if err := tx.Create(address).Error; err != nil {
		t.Fatalf("failed to create address: %v", err)
	}

	shipping := &models.Shipping{Name: "Test Shipping", Price: 10000, State: "active"}
	if err := tx.Create(shipping).Error; err != nil {
		t.Fatalf("failed to create shipping: %v", err)
	}

	paymentMethod := &models.PaymentMethod{
		AccountName:   "Test Account",
		AccountNumber: "1234567890",
		BankName:      "Test Bank",
		IsActive:      true,
	}
	if err := tx.Create(paymentMethod).Error; err != nil {
		t.Fatalf("failed to create payment method: %v", err)
	}

	transactionID := uuid.New().String()
	if err := tx.Create(&models.Transaction{
		TxID:            transactionID,
		AddressID:       address.ID,
		ShippingID:      shipping.ID,
		PaymentMethodID: paymentMethod.ID,
		ShippingPrice:   10000,
		TotalPrice:      210000,
		Status:          "paid",
	}).Error; err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}

	order := &models.Order{
		ID:            uuid.New().String(),
		UserID:        int64(user.ID),
		TransactionID: transactionID,
		ProductID:     product.ID,
		ColorVarianID: colorVarian.ID,
		SizeVarianID:  sizeVarian.ID,
		UnitPrice:     100000,
		Subtotal:      200000,
		Quantity:      3,
		Status:        "paid",
	}
	if err := tx.Create(order).Error; err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	orderRepo := repository.NewOrderRepository()
	productRepo := repository.NewProductRepository()
	fulfillmentService := services.NewOrderFulfillmentService(orderRepo, productRepo)

	t.Run("Transaction shipped -> Order ikut shipped, TANPA restock", func(t *testing.T) {
		if err := fulfillmentService.SyncOrdersToTransactionStatus(tx, transactionID, "shipped"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var updatedOrder models.Order
		tx.First(&updatedOrder, "id = ?", order.ID)
		if updatedOrder.Status != "shipped" {
			t.Errorf("expected order status 'shipped', got %q", updatedOrder.Status)
		}

		var sv models.SizeVarian
		tx.First(&sv, sizeVarian.ID)
		if sv.Stock != 10 {
			t.Errorf("expected stock unchanged at 10 (not a restock scenario), got %d", sv.Stock)
		}
	})

	t.Run("Transaction cancelled -> Order cancelled DAN stock direstock", func(t *testing.T) {
		if err := fulfillmentService.SyncOrdersToTransactionStatus(tx, transactionID, "cancelled"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var updatedOrder models.Order
		tx.First(&updatedOrder, "id = ?", order.ID)
		if updatedOrder.Status != "cancelled" {
			t.Errorf("expected order status 'cancelled', got %q", updatedOrder.Status)
		}

		var sv models.SizeVarian
		tx.First(&sv, sizeVarian.ID)
		// Stock awal 10, order quantity 3 -> setelah restock harus 13.
		if sv.Stock != 13 {
			t.Errorf("expected stock restocked to 13 (10 + quantity 3), got %d", sv.Stock)
		}
	})

	t.Run("Order sudah cancelled (terminal) tidak disentuh lagi", func(t *testing.T) {
		// Panggil lagi dengan status lain — order yang sudah cancelled
		// (terminal) tidak boleh berubah atau kena restock kedua kalinya.
		if err := fulfillmentService.SyncOrdersToTransactionStatus(tx, transactionID, "shipped"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var sv models.SizeVarian
		tx.First(&sv, sizeVarian.ID)
		if sv.Stock != 13 {
			t.Errorf("expected stock to remain 13 (no double-restock), got %d", sv.Stock)
		}
	})
}
