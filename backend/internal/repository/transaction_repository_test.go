package repository_test

import (
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/repository"
	"e-commerce/backend/internal/testhelper"
	"testing"
	"time"
)

func TestTransactionRepository_CRUD(t *testing.T) {
	// Setup
	tx := testhelper.BeginTestTransaction(t, testDB)
	defer testhelper.RollbackTestTransaction(tx)

	// Replace global DB
	dbWrapper := testhelper.SetTestDB(tx)
	defer dbWrapper.Restore()

	repo := repository.NewTransactionRepository()

	// Setup Dependencies
	role := testhelper.CreateTestRole(tx, "User")
	user := testhelper.CreateTestUser(tx, "trx_user@example.com", role.ID)

	// Create Address
	address := models.Address{
		UserID:               int64(user.ID),
		RecipientName:        "Home",
		RecipientPhoneNumber: "08123456789",
		Province:             "Province",
		City:                 "City",
		District:             "District",
		Village:              "Village",
		PostalCode:           "12345",
		FullAddress:          "Street 1",
	}
	tx.Create(&address)

	// Create Shipping
	shipping := models.Shipping{
		Name: "JNE",
	}
	tx.Create(&shipping)

	// Create Payment Method
	pm := models.PaymentMethod{
		BankName:      "Bank Transfer",
		AccountName:   "User",
		AccountNumber: "1234567890",
	}
	tx.Create(&pm)

	trx := models.Transaction{
		TxID:            "TRX-TEST-002",
		AddressID:       address.ID,
		ShippingID:      shipping.ID,
		PaymentMethodID: pm.ID,
		TotalPrice:      500000,
		ShippingPrice:   10000,
		Status:          "pending",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}


	// 1. Test Create
	createdTrx, err := repo.Create(trx, tx)
	if err != nil {
		t.Fatalf("Create Transaction failed: %v", err)
	}
	if createdTrx.TxID != trx.TxID {
		t.Errorf("Expected TxID %s, got %s", trx.TxID, createdTrx.TxID)
	}

	// 2. Test FindById
	foundTrx, err := repo.FindById(trx.TxID)
	if err != nil {
		t.Fatalf("FindById failed: %v", err)
	}
	if foundTrx.TxID != trx.TxID {
		t.Errorf("Expected TxID %s, got %s", trx.TxID, foundTrx.TxID)
	}

	// 3. Test Update
	trx.Status = "success"
	updatedTrx, err := repo.Update(trx, tx)
	if err != nil {
		t.Fatalf("Update Transaction failed: %v", err)
	}
	if updatedTrx.Status != "success" {
		t.Errorf("Expected status success, got %s", updatedTrx.Status)
	}

	// 4. Test FindAll
	req := models.TransactionListRequest{
		Limit: 10,
		Page: 1,
	}
	trxs, err := repo.FindAll(req)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	if len(trxs) == 0 {
		t.Errorf("Expected at least 1 transaction")
	}

	// 5. Test FindByIdLocking
	lockedTrx, err := repo.FindByIdLocking(tx, trx.TxID)
	if err != nil {
		t.Fatalf("FindByIdLocking failed: %v", err)
	}
	if lockedTrx.TxID != trx.TxID {
		t.Errorf("Expected TxID %s, got %s", trx.TxID, lockedTrx.TxID)
	}

	// 6. Test Delete
	err = repo.Delete(trx)
	if err != nil {
		t.Fatalf("Delete Transaction failed: %v", err)
	}
}
