package handler_test

import (
	"bytes"
	"e-commerce/backend/internal/handler"
	"e-commerce/backend/internal/mocks"
	"e-commerce/backend/internal/models"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestAuthHandler_SignIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockAuthService(ctrl)
	authHandler := handler.NewAuthHandler(mockService)

	t.Run("Success", func(t *testing.T) {
		reqBody := models.LoginRequest{
			Email:    "test@example.com",
			Password: "password",
		}
		
		expectedResp := &models.TokenResponse{
			AccessToken: "access_token",
			RefreshToken: "refresh_token",
		}

		// Expect SignIn with any ip/agent
		mockService.EXPECT().SignIn(reqBody, gomock.Any(), gomock.Any()).Return(expectedResp, nil)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/auth/signin", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		authHandler.SignIn(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		reqBody := models.LoginRequest{
			Email:    "test@example.com",
			Password: "wrong",
		}
		
		mockService.EXPECT().SignIn(reqBody, gomock.Any(), gomock.Any()).Return(nil, errors.New("invalid credentials"))

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/auth/signin", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		authHandler.SignIn(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})
}

func TestAuthHandler_SignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockAuthService(ctrl)
	authHandler := handler.NewAuthHandler(mockService)

	t.Run("Success", func(t *testing.T) {
		reqBody := models.RegisterRequest{
			Email:    "new@example.com",
			Password: "password",
			Username: "newuser",
			FirstName: "New",
			LastName: "User",
		}
		
		expectedResp := &models.TokenResponse{
			AccessToken: "access_token",
			User: models.UserResponse{Email: "new@example.com"},
		}

		mockService.EXPECT().SignUp(reqBody).Return(expectedResp, nil)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		authHandler.SignUp(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", w.Code)
		}
	})
}
