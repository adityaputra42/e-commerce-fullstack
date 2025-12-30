package handler_test

import (
	"bytes"
	"context"
	"e-commerce/backend/internal/handler"
	"e-commerce/backend/internal/mocks"
	"e-commerce/backend/internal/models"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"go.uber.org/mock/gomock"
)

func TestTransactionHandler_CreateTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTransactionService(ctrl)
	transactionHandler := handler.NewTransactionHandler(mockService)

	t.Run("Success", func(t *testing.T) {
		input := models.CreateTransaction{
			AddressID: 1,
			ShippingID: 1,
			PaymentMethodID: 1,
			ShippingPrice: 5000,
			TotalPrice: 15000,
			ProductOrders: []models.CreateOrder{
				{
					ProductID: 1, 
					Quantity: 1,
					ColorVarianID: 1,
					SizeVarianID: 1,
					UnitPrice: 10000,
					Subtotal: 10000,
				},
			},
		}
		
		expectedResp := models.TransactionResponse{
			TxID: "TRX-123",
			Status: "pending",
		}

		mockService.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(&expectedResp, nil)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
		
		w := httptest.NewRecorder()

		transactionHandler.CreateTransaction(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", w.Code)
		}
	})

	t.Run("ValidationError", func(t *testing.T) {
		// Missing fields
		input := models.CreateTransaction{}
		
		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
		
		w := httptest.NewRecorder()

		transactionHandler.CreateTransaction(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for validation error, got %d", w.Code)
		}
	})
}

func TestTransactionHandler_GetTransactionByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTransactionService(ctrl)
	transactionHandler := handler.NewTransactionHandler(mockService)

	t.Run("Success", func(t *testing.T) {
		txID := "TRX-123"
		mockService.EXPECT().FindTransactionById(txID).Return(&models.TransactionResponse{TxID: txID}, nil) // Pointer

		req := httptest.NewRequest(http.MethodGet, "/transactions/" + txID, nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("tx_id", txID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()

		transactionHandler.GetTransactionByID(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		txID := "TRX-999"
		mockService.EXPECT().FindTransactionById(txID).Return(nil, errors.New("transaction not found")) // Nil

		req := httptest.NewRequest(http.MethodGet, "/transactions/" + txID, nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("tx_id", txID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()

		transactionHandler.GetTransactionByID(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
	})
}
