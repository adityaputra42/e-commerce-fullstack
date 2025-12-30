package handler_test

import (
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

func TestOrderHandler_CancelOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockOrderService(ctrl)
	orderHandler := handler.NewOrderHandler(mockService)

	t.Run("Success", func(t *testing.T) {
		orderID := "ORD-123"
		userID := int64(1)

		expectedOrder := models.OrderResponse{
			ID:     orderID,
			Status: "cancelled",
		}

		mockService.EXPECT().CancelOrder(orderID, userID).Return(&expectedOrder, nil)

		req := httptest.NewRequest(http.MethodPatch, "/orders/" + orderID + "/cancel", nil)
		
		// Setup URL Param
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", orderID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		
		// Setup Auth Context (middleware simulation)
		req = req.WithContext(context.WithValue(req.Context(), "user_id", userID))
		
		w := httptest.NewRecorder()

		orderHandler.CancelOrder(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
		
		var resp struct {
			Data models.OrderResponse `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &resp)
		if resp.Data.Status != "cancelled" {
			t.Errorf("Expected status cancelled, got %s", resp.Data.Status)
		}
	})

	t.Run("Unauthorized", func(t *testing.T) {
		orderID := "ORD-123"
		userID := int64(1)

		mockService.EXPECT().CancelOrder(orderID, userID).Return(nil, errors.New("unauthorized: order does not belong to user"))

		req := httptest.NewRequest(http.MethodPatch, "/orders/" + orderID + "/cancel", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", orderID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		req = req.WithContext(context.WithValue(req.Context(), "user_id", userID))

		w := httptest.NewRecorder()

		orderHandler.CancelOrder(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", w.Code)
		}
	})
	
	t.Run("Unauthenticated", func(t *testing.T) {
		orderID := "ORD-123"
		
		req := httptest.NewRequest(http.MethodPatch, "/orders/" + orderID + "/cancel", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", orderID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		// No user_id in context

		w := httptest.NewRecorder()

		orderHandler.CancelOrder(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})
}

func TestOrderHandler_GetAllOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockOrderService(ctrl)
	orderHandler := handler.NewOrderHandler(mockService)

	t.Run("Success", func(t *testing.T) {
		mockService.EXPECT().FindAllOrder(gomock.Any()).Return(&[]models.OrderResponse{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/orders", nil)
		w := httptest.NewRecorder()

		orderHandler.GetAllOrders(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})
}
