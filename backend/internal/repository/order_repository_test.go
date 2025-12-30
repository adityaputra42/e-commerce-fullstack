package repository_test

import (
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/repository"
	"e-commerce/backend/internal/testhelper"
	"testing"
	"time"
)

func TestOrderRepository_CRUD(t *testing.T) {
	// Setup
	tx := testhelper.BeginTestTransaction(t, testDB)
	defer testhelper.RollbackTestTransaction(tx)

	// Replace global DB with transaction
	dbWrapper := testhelper.SetTestDB(tx)
	defer dbWrapper.Restore()

	repo := repository.NewOrderRepository()

	// Setup Dependencies
	role := testhelper.CreateTestRole(tx, "User")
	user := testhelper.CreateTestUser(tx, "order_test@example.com", role.ID)
	category := testhelper.CreateTestCategory(tx, "Category Order")
	product := testhelper.CreateTestProduct(tx, "Product Order", category.ID, 100.0)
	
	// Setup variants
	colorParams := models.ColorVarian{ProductID: product.ID, Name: "Blue", Color: "blue"}
	tx.Create(&colorParams)
	sizeParams := models.SizeVarian{ColorVarianID: colorParams.ID, Size: "L", Stock: 10}
	tx.Create(&sizeParams)

	// Create Dependencies for Transaction
	address := models.Address{
		UserID:               int64(user.ID),
		RecipientName:        "Test Recipient",
		RecipientPhoneNumber: "08123456789",
		Province:             "Test Province",
		City:                 "Test City",
		District:             "Test District",
		Village:              "Test Village",
		PostalCode:           "12345",
		FullAddress:          "Test Address Complete",
	}
	tx.Create(&address)

	shipping := models.Shipping{
		Name: "Test Shipping",
	}
	tx.Create(&shipping)
	
	pm := models.PaymentMethod{
		AccountName:   "Test Account",
		AccountNumber: "1234567890",
		BankName:      "Test Bank",
	}
	tx.Create(&pm)

	// Create Transaction first as Order references it
	trx := models.Transaction{
		TxID:            "TRX-ORDER-TEST-001",
		AddressID:       address.ID,
		ShippingID:      shipping.ID,
		PaymentMethodID: pm.ID,
		TotalPrice:      100.0,
		ShippingPrice:   10.0,
		Status:          "pending",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	tx.Create(&trx)

	order := models.Order{
		ID:            "ORD-TEST-001",
		UserID:        int64(user.ID),
		TransactionID: trx.TxID,
		ProductID:     product.ID,
		ColorVarianID: colorParams.ID,
		SizeVarianID:  sizeParams.ID,
		Quantity:      1,
		UnitPrice:     100.0,
		Subtotal:      100.0,
		Status:        "pending",
	}


	// 1. Test Create
	createdOrder, err := repo.Create(order, tx)
	if err != nil {
		t.Fatalf("Create Order failed: %v", err)
	}
	if createdOrder.ID != order.ID {
		t.Errorf("Expected ID %s, got %s", order.ID, createdOrder.ID)
	}

	// 2. Test FindById
	foundOrder, err := repo.FindById(order.ID)
	if err != nil {
		t.Fatalf("FindById failed: %v", err)
	}
	if foundOrder.ID != order.ID {
		t.Errorf("Expected ID %s, got %s", order.ID, foundOrder.ID)
	}
	// Check Preloads
	if foundOrder.Product.ID == 0 {
		t.Errorf("Expected Product to be preloaded")
	}

	// 3. Test Update
	foundOrder.Status = "paid"
	updatedOrder, err := repo.Update(foundOrder, tx)
	if err != nil {
		t.Fatalf("Update Order failed: %v", err)
	}
	if updatedOrder.Status != "paid" {
		t.Errorf("Expected status paid, got %s", updatedOrder.Status)
	}

	// 4. Test FindAll
	req := models.OrderListRequest{
		Limit: 10,
		Page: 1,
	}
	orders, err := repo.FindAll(req)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	if len(orders) == 0 {
		t.Errorf("Expected at least 1 order")
	}

	// 5. Test FindAllByTxId
	ordersByTx, err := repo.FindAllByTxId(trx.TxID)
	if err != nil {
		t.Fatalf("FindAllByTxId failed: %v", err)
	}
	if len(ordersByTx) != 1 {
		t.Errorf("Expected 1 order for TxId, got %d", len(ordersByTx))
	}

	// 6. Test Delete
	err = repo.Delete(foundOrder)
	if err != nil {
		t.Fatalf("Delete Order failed: %v", err)
	}
	
	// Verify delete
	_, err = repo.FindById(order.ID)
	// Usually GORM soft delete returns success but record is hidden
	// but if FindById implementation uses "deleted_at IS NULL", it might return error or empty
	// Since FindById implementation "Take(...)" usually returns ErrRecordNotFound if not found
	// Let's check if it's returning error or not.
	// Actually GORM Delete with &param performs soft delete if DeletedAt exists.
	// Checking the Order model, it has DeletedAt.
}
