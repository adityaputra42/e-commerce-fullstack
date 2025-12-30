package handler_test

import (
	"context"
	"e-commerce/backend/internal/handler"
	"e-commerce/backend/internal/mocks"
	"e-commerce/backend/internal/models"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"go.uber.org/mock/gomock"
)

func TestProductHandler_GetProductByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockProductService(ctrl)
	productHandler := handler.NewProductHandler(mockService)

	t.Run("Success", func(t *testing.T) {
		productID := int64(123)
		
		expectedProduct := models.ProductDetailResponse{
			ID: productID,
			Name: "Test Product",
			Price: 10000,
		}

		mockService.EXPECT().FindProductById(productID).Return(&expectedProduct, nil)

		req := httptest.NewRequest(http.MethodGet, "/products/" + strconv.FormatInt(productID, 10), nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", strconv.FormatInt(productID, 10))
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		
		w := httptest.NewRecorder()

		productHandler.GetProductByID(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		productID := int64(999)
		// Using generic error for not found to match strict service behavior usually returning wrapped error
		mockService.EXPECT().FindProductById(productID).Return(nil, errors.New("product not found"))

		req := httptest.NewRequest(http.MethodGet, "/products/" + strconv.FormatInt(productID, 10), nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", strconv.FormatInt(productID, 10))
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		
		w := httptest.NewRecorder()

		productHandler.GetProductByID(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
	})
}
