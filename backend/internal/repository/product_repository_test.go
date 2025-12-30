package repository_test

import (
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/repository"
	"e-commerce/backend/internal/testhelper"
	"testing"
)

func TestProductRepository_ProductOperations(t *testing.T) {
	// Setup
	tx := testhelper.BeginTestTransaction(t, testDB)
	defer testhelper.RollbackTestTransaction(tx)

	repo := repository.NewProductRepository()
	
	// Create dependency (Category)
	category := testhelper.CreateTestCategory(tx, "Electronics")

	// 1. Test CreateProduct
	product := models.Product{
		Name:        "Smartphone X",
		Description: "Latest smartphone",
		CategoryID:  category.ID,
		Price:       999.99,
		Images:      "image.jpg",
		Rating:      5.0,
	}

	createdProduct, err := repo.CreateProduct(product, tx)
	if err != nil {
		t.Fatalf("CreateProduct failed: %v", err)
	}
	if createdProduct.ID == 0 {
		t.Errorf("Expected product ID to be set")
	}
	if createdProduct.Name != product.Name {
		t.Errorf("Expected name %s, got %s", product.Name, createdProduct.Name)
	}

	// 2. Test FindProductById
	foundProduct, err := repo.FindProductById(createdProduct.ID, tx)
	if err != nil {
		t.Fatalf("FindProductById failed: %v", err)
	}
	if foundProduct.ID != createdProduct.ID {
		t.Errorf("Expected ID %d, got %d", createdProduct.ID, foundProduct.ID)
	}

	// 3. Test UpdateProduct
	createdProduct.Name = "Smartphone X Pro"
	updatedProduct, err := repo.UpdateProduct(createdProduct, tx)
	if err != nil {
		t.Fatalf("UpdateProduct failed: %v", err)
	}
	if updatedProduct.Name != "Smartphone X Pro" {
		t.Errorf("Expected updated name 'Smartphone X Pro', got %s", updatedProduct.Name)
	}

	// 4. Test FindAllProduct
	listRequest := models.ProductListRequest{
		Limit: 10,
		Page:  1,
	}
	products, total, err := repo.FindAllProduct(listRequest, tx)
	if err != nil {
		t.Fatalf("FindAllProduct failed: %v", err)
	}
	if total != 1 {
		t.Errorf("Expected total 1, got %d", total)
	}
	if len(products) != 1 {
		t.Errorf("Expected 1 product, got %d", len(products))
	}
	if products[0].Name != "Smartphone X Pro" {
		t.Errorf("Expected product name 'Smartphone X Pro', got %s", products[0].Name)
	}

	// 5. Test DeleteProduct
	err = repo.DeleteProduct(createdProduct.ID, tx)
	if err != nil {
		t.Fatalf("DeleteProduct failed: %v", err)
	}

	// Verify deletion (Soft Delete check)
	// FindProductById implements "Where deleted_at IS NULL" so it should return error or nil
	_, err = repo.FindProductById(createdProduct.ID, tx)
	if err == nil {
		t.Errorf("Expected error when finding deleted product, got nil")
	}
}

func TestProductRepository_VariantOperations(t *testing.T) {
	// Setup
	tx := testhelper.BeginTestTransaction(t, testDB)
	defer testhelper.RollbackTestTransaction(tx)

	repo := repository.NewProductRepository()
	
	category := testhelper.CreateTestCategory(tx, "Fashion")
	product := testhelper.CreateTestProduct(tx, "T-Shirt", category.ID, 19.99)

	// 1. Test CreateColorVarian
	colorVarian := models.ColorVarian{
		ProductID: product.ID,
		Name:      "Red Edition",
		Color:     "Red",
	}
	
	createdColor, err := repo.CreateColorVarian(colorVarian, tx)
	if err != nil {
		t.Fatalf("CreateColorVarian failed: %v", err)
	}
	if createdColor.ID == 0 {
		t.Errorf("Expected color varian ID to be set")
	}

	// 2. Test CreateSizeVarian
	sizeVarian := models.SizeVarian{
		ColorVarianID: createdColor.ID,
		Size:          "XL",
		Stock:         100,
	}
	
	createdSize, err := repo.CreateSizeVarian(sizeVarian, tx)
	if err != nil {
		t.Fatalf("CreateSizeVarian failed: %v", err)
	}
	if createdSize.ID == 0 {
		t.Errorf("Expected size varian ID to be set")
	}

	// 3. Verify Product Preload (FindProductById should load variants)
	foundProduct, err := repo.FindProductById(product.ID, tx)
	if err != nil {
		t.Fatalf("FindProductById failed: %v", err)
	}
	
	if len(foundProduct.ColorVarians) != 1 {
		t.Errorf("Expected 1 color variant, got %d", len(foundProduct.ColorVarians))
	}
	if len(foundProduct.ColorVarians[0].SizeVarians) != 1 {
		t.Errorf("Expected 1 size variant, got %d", len(foundProduct.ColorVarians[0].SizeVarians))
	}
	
	if foundProduct.ColorVarians[0].Name != "Red Edition" {
		t.Errorf("Expected variant name 'Red Edition', got %s", foundProduct.ColorVarians[0].Name)
	}
	if foundProduct.ColorVarians[0].SizeVarians[0].Size != "XL" {
		t.Errorf("Expected size 'XL', got %s", foundProduct.ColorVarians[0].SizeVarians[0].Size)
	}
}
