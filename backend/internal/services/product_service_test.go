package services_test

import (
	"e-commerce/backend/internal/mocks"
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/services"
	"e-commerce/backend/internal/testhelper"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

func TestProductService_FindProductById(t *testing.T) {
	// Setup Controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mocks
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockCategoryRepo := mocks.NewMockCategoryRepository(ctrl)

	// Service
	service := services.NewProductService(mockCategoryRepo, mockProductRepo)

	// Test Case
	t.Run("Found", func(t *testing.T) {
		productID := int64(1)
		categoryID := int64(10)
		
		product := &models.ProductDetailResponse{
			ID:         productID,
			Name:       "Test Product",
			Category:   models.CategoryResponse{ID: categoryID, Name: "Test Category"},
			Price:      100.0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		
		// Expectations
		mockProductRepo.EXPECT().
			FindProductById(productID, nil).
			Return(&models.Product{
				ID: productID,
				Name: "Test Product",
				CategoryID: categoryID,
				Price: 100.0,
			}, nil)

		// Note: FindProductById service actually calls ToProductDetailResponse which might not call CategoryRepo if category is not needed explicitly or passed?
		// Checking implementation of FindProductById in service:
		// It calls productRepo.FindProductById
		// Then it calls categoryRepo.FindById(product.CategoryID)
		
		mockCategoryRepo.EXPECT().
			FindById(categoryID). // Correct method name
			Return(models.Category{ID: categoryID, Name: "Test Category"}, nil)

		// Execute
		result, err := service.FindProductById(productID)
		
		// Assert
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if result.ID != productID {
			t.Errorf("Expected ID %d, got %d", productID, result.ID)
		}
		if result.Category.Name != product.Category.Name {
			t.Errorf("Expected category %s, got %s", product.Category.Name, result.Category.Name)
		}
	})
}

func TestProductService_CreateProduct(t *testing.T) {
	t.Skip("Skipping because of hard dependency on utils.UploadToSupabase")

	// Initialize TestDB locally for this package
	testDB := testhelper.SetupTestDB()
	// No cleanup here to avoid interfering with other tests, or use separate DB? 
	// Ideally service tests should mock everything, but transaction requires this hack.
	
	// Setup Transaction
	tx := testhelper.BeginTestTransaction(t, testDB)
	defer testhelper.RollbackTestTransaction(tx)
	
	// Replace global DB to allow Transaction block to run
	dbWrapper := testhelper.SetTestDB(tx)
	defer dbWrapper.Restore()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockCategoryRepo := mocks.NewMockCategoryRepository(ctrl)

	service := services.NewProductService(mockCategoryRepo, mockProductRepo)

	t.Run("Success", func(t *testing.T) {
		input := models.CreateProductParam{
			Name:       "New Product",
			CategoryID: 1,
			Price:      200.0,
			ColorVarian: []models.CreateColorVarianRequest{
				{
					Name: "Red",
					Color: "red",
					Sizes: []models.CreateSizeVarianRequest{
						{Size: "M", Stock: 10},
					},
				},
			},
		}

		category := models.Category{ID: 1, Name: "Cat"}

		// Expectations
		mockCategoryRepo.EXPECT().FindById(int64(1)).Return(category, nil)
		
		// Expect CreateProduct with Any tx
		mockProductRepo.EXPECT().
			CreateProduct(gomock.Any(), gomock.Any()).
			DoAndReturn(func(p models.Product, tx interface{}) (models.Product, error) {
				p.ID = 100
				return p, nil
			})

		// Expect ColorVarian creation
		mockProductRepo.EXPECT().
			CreateColorVarian(gomock.Any(), gomock.Any()).
			DoAndReturn(func(cv models.ColorVarian, tx interface{}) (models.ColorVarian, error) {
				cv.ID = 200
				return cv, nil
			})

		// Expect SizeVarian creation
		mockProductRepo.EXPECT().
			CreateSizeVarian(gomock.Any(), gomock.Any()).
			DoAndReturn(func(sv models.SizeVarian, tx interface{}) (models.SizeVarian, error) {
				sv.ID = 300
				return sv, nil
			})
			
		mockProductRepo.EXPECT().
			FindProductById(int64(100), gomock.Any()).
			Return(&models.Product{
				ID: 100, 
				Name: "New Product",
				CategoryID: 1,
				Price: 200.0,
				ColorVarians: []models.ColorVarian{
					{ID: 200, Name: "Red", SizeVarians: []models.SizeVarian{{ID: 300, Size: "M"}}},
				},
			}, nil)

		mockCategoryRepo.EXPECT().FindById(int64(1)).Return(category, nil)

		// Execute
		result, err := service.CreateProduct(input)
		
		// Assert
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if result.ID != 100 {
			t.Errorf("Expected product ID 100, got %d", result.ID)
		}
	})
}

