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
	"gorm.io/gorm"
)

func TestUserHandler_GetUserById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	userHandler := handler.NewUserHandler(mockUserService)

	tests := []struct {
		name           string
		userID         string
		mockFn         func()
		expectedStatus int
	}{
		{
			name:   "Get user successfully",
			userID: "1",
			mockFn: func() {
				mockUserService.EXPECT().
					GetUserById(uint(1)).
					Return(&models.User{
						ID:        1,
						Email:     "test@example.com",
						Username:  "testuser",
						FirstName: "Test",
						LastName:  "User",
						IsActive:  true,
					}, nil).
					Times(1)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "User not found",
			userID: "999",
			mockFn: func() {
				mockUserService.EXPECT().
					GetUserById(uint(999)).
					Return(nil, errors.New("user not found")).
					Times(1)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Invalid user ID",
			userID:         "invalid",
			mockFn:         func() {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			req := httptest.NewRequest(http.MethodGet, "/api/users/"+tt.userID, nil)
			w := httptest.NewRecorder()

			userHandler.GetUserById(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("GetUserById() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.NewDecoder(w.Body).Decode(&response)
				if err != nil {
					t.Errorf("Failed to decode response: %v", err)
				}

				if response["success"] != true {
					t.Errorf("GetUserById() success = %v, want true", response["success"])
				}
			}
		})
	}
}

func TestUserHandler_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	userHandler := handler.NewUserHandler(mockUserService)

	tests := []struct {
		name           string
		input          models.UserInput
		mockFn         func()
		expectedStatus int
	}{
		{
			name: "Create user successfully",
			input: models.UserInput{
				Email:     "newuser@example.com",
				Username:  "newuser",
				Password:  "password123",
				FirstName: "New",
				LastName:  "User",
				RoleID:    1,
			},
			mockFn: func() {
				mockUserService.EXPECT().
					CreateUser(gomock.Any()).
					Return(&models.User{
						ID:        1,
						Email:     "newuser@example.com",
						Username:  "newuser",
						FirstName: "New",
						LastName:  "User",
						RoleID:    1,
						IsActive:  true,
					}, nil).
					Times(1)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Create user with existing email",
			input: models.UserInput{
				Email:     "existing@example.com",
				Username:  "newuser",
				Password:  "password123",
				FirstName: "New",
				LastName:  "User",
				RoleID:    1,
			},
			mockFn: func() {
				mockUserService.EXPECT().
					CreateUser(gomock.Any()).
					Return(nil, errors.New("user with this email or username already exists")).
					Times(1)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			userHandler.CreateUser(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("CreateUser() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusCreated {
				var response map[string]interface{}
				err := json.NewDecoder(w.Body).Decode(&response)
				if err != nil {
					t.Errorf("Failed to decode response: %v", err)
				}

				if response["success"] != true {
					t.Errorf("CreateUser() success = %v, want true", response["success"])
				}
			}
		})
	}
}

func TestUserHandler_UpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	userHandler := handler.NewUserHandler(mockUserService)

	tests := []struct {
		name           string
		userID         string
		input          models.UserUpdateInput
		mockFn         func()
		expectedStatus int
	}{
		{
			name:   "Update user successfully",
			userID: "1",
			input: models.UserUpdateInput{
				FirstName: "Updated",
				LastName:  "Name",
			},
			mockFn: func() {
				mockUserService.EXPECT().
					UpdateUser(uint(1), gomock.Any()).
					Return(&models.User{
						ID:        1,
						Email:     "test@example.com",
						Username:  "testuser",
						FirstName: "Updated",
						LastName:  "Name",
						IsActive:  true,
					}, nil).
					Times(1)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Update non-existing user",
			userID: "999",
			input: models.UserUpdateInput{
				FirstName: "Updated",
			},
			mockFn: func() {
				mockUserService.EXPECT().
					UpdateUser(uint(999), gomock.Any()).
					Return(nil, errors.New("user not found")).
					Times(1)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPut, "/api/users/"+tt.userID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			userHandler.UpdateUser(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("UpdateUser() status = %v, want %v", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestUserHandler_DeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	userHandler := handler.NewUserHandler(mockUserService)

	tests := []struct {
		name           string
		userID         string
		mockFn         func()
		expectedStatus int
	}{
		{
			name:   "Delete user successfully",
			userID: "1",
			mockFn: func() {
				mockUserService.EXPECT().
					DeleteUser(uint(1)).
					Return(nil).
					Times(1)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Delete non-existing user",
			userID: "999",
			mockFn: func() {
				mockUserService.EXPECT().
					DeleteUser(uint(999)).
					Return(errors.New("user not found")).
					Times(1)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid user ID",
			userID:         "invalid",
			mockFn:         func() {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			req := httptest.NewRequest(http.MethodDelete, "/api/users/"+tt.userID, nil)
			w := httptest.NewRecorder()

			userHandler.DeleteUser(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("DeleteUser() status = %v, want %v", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestUserHandler_ActivateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	userHandler := handler.NewUserHandler(mockUserService)

	tests := []struct {
		name           string
		userID         string
		mockFn         func()
		expectedStatus int
	}{
		{
			name:   "Activate user successfully",
			userID: "1",
			mockFn: func() {
				mockUserService.EXPECT().
					ActivateUser(uint(1)).
					Return(&models.User{
						ID:       1,
						Email:    "test@example.com",
						IsActive: true,
					}, nil).
					Times(1)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Activate non-existing user",
			userID: "999",
			mockFn: func() {
				mockUserService.EXPECT().
					ActivateUser(uint(999)).
					Return(nil, gorm.ErrRecordNotFound).
					Times(1)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			req := httptest.NewRequest(http.MethodPut, "/api/users/"+tt.userID+"/activate", nil)
			w := httptest.NewRecorder()

			userHandler.ActivateUser(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("ActivateUser() status = %v, want %v", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestUserHandler_DeactivateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	userHandler := handler.NewUserHandler(mockUserService)

	tests := []struct {
		name           string
		userID         string
		mockFn         func()
		expectedStatus int
	}{
		{
			name:   "Deactivate user successfully",
			userID: "1",
			mockFn: func() {
				mockUserService.EXPECT().
					DeactivateUser(uint(1)).
					Return(&models.User{
						ID:       1,
						Email:    "test@example.com",
						IsActive: false,
					}, nil).
					Times(1)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			req := httptest.NewRequest(http.MethodPut, "/api/users/"+tt.userID+"/deactivate", nil)
			w := httptest.NewRecorder()

			userHandler.DeactivateUser(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("DeactivateUser() status = %v, want %v", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestUserHandler_GetUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	userHandler := handler.NewUserHandler(mockUserService)

	tests := []struct {
		name           string
		queryParams    string
		mockFn         func()
		expectedStatus int
	}{
		{
			name:        "Get users successfully",
			queryParams: "?page=1&limit=10",
			mockFn: func() {
				mockUserService.EXPECT().
					GetUsers(gomock.Any()).
					Return(&models.UserListResponse{
						Users: []models.UserResponse{
							{ID: 1, Email: "user1@example.com"},
							{ID: 2, Email: "user2@example.com"},
						},
						Total:      2,
						Page:       1,
						Limit:      10,
						TotalPages: 1,
					}, nil).
					Times(1)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "Get users with default pagination",
			queryParams: "",
			mockFn: func() {
				mockUserService.EXPECT().
					GetUsers(gomock.Any()).
					Return(&models.UserListResponse{
						Users:      []models.UserResponse{},
						Total:      0,
						Page:       1,
						Limit:      10,
						TotalPages: 0,
					}, nil).
					Times(1)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			req := httptest.NewRequest(http.MethodGet, "/api/users"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			userHandler.GetUsers(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("GetUsers() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.NewDecoder(w.Body).Decode(&response)
				if err != nil {
					t.Errorf("Failed to decode response: %v", err)
				}

				if response["success"] != true {
					t.Errorf("GetUsers() success = %v, want true", response["success"])
				}
			}
		})
	}
}

func TestUserHandler_BulkUserActions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	userHandler := handler.NewUserHandler(mockUserService)

	tests := []struct {
		name           string
		input          interface{}
		mockFn         func()
		expectedStatus int
	}{
		{
			name: "Bulk activate users successfully",
			input: map[string]interface{}{
				"user_ids": []uint{1, 2, 3},
				"action":   "activate",
			},
			mockFn: func() {
				mockUserService.EXPECT().
					BulkUserActions(gomock.Any()).
					Return(nil).
					Times(1)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Bulk action with invalid input",
			input: map[string]interface{}{
				"invalid": "data",
			},
			mockFn: func() {
				mockUserService.EXPECT().
					BulkUserActions(gomock.Any()).
					Return(errors.New("invalid action")).
					Times(1)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/api/users/bulk", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			userHandler.BulkUserActions(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("BulkUserActions() status = %v, want %v", w.Code, tt.expectedStatus)
			}
		})
	}
}
