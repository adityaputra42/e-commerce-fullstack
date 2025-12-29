package services_test

import (
	"e-commerce/backend/internal/mocks"
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/services"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestUserService_GetUserById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockActivityLogRepo := mocks.NewMockActivityLogRepository(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	service := services.NewUserService(mockUserRepo, mockActivityLogRepo, mockRoleRepo)

	tests := []struct {
		name    string
		userID  uint
		mockFn  func()
		wantErr bool
		errMsg  string
	}{
		{
			name:   "Get user successfully",
			userID: 1,
			mockFn: func() {
				mockUserRepo.EXPECT().
					FindById(uint(1)).
					Return(models.User{
						ID:        1,
						Email:     "test@example.com",
						Username:  "testuser",
						FirstName: "Test",
						LastName:  "User",
						IsActive:  true,
					}, nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name:   "User not found",
			userID: 999,
			mockFn: func() {
				mockUserRepo.EXPECT().
					FindById(uint(999)).
					Return(models.User{}, gorm.ErrRecordNotFound).
					Times(1)
			},
			wantErr: true,
			errMsg:  "user not found",
		},
		{
			name:   "Database error",
			userID: 1,
			mockFn: func() {
				mockUserRepo.EXPECT().
					FindById(uint(1)).
					Return(models.User{}, errors.New("database error")).
					Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			result, err := service.GetUserById(tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errMsg != "" {
				if err.Error() != tt.errMsg {
					t.Errorf("GetUserById() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}

			if !tt.wantErr {
				if result == nil {
					t.Errorf("GetUserById() result should not be nil")
				}
				if result.ID != tt.userID {
					t.Errorf("GetUserById() ID = %v, want %v", result.ID, tt.userID)
				}
			}
		})
	}
}

func TestUserService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockActivityLogRepo := mocks.NewMockActivityLogRepository(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	service := services.NewUserService(mockUserRepo, mockActivityLogRepo, mockRoleRepo)

	tests := []struct {
		name    string
		input   *models.UserInput
		mockFn  func()
		wantErr bool
		errMsg  string
	}{
		{
			name: "Create user successfully",
			input: &models.UserInput{
				Email:     "newuser@example.com",
				Username:  "newuser",
				Password:  "password123",
				FirstName: "New",
				LastName:  "User",
				RoleID:    1,
			},
			mockFn: func() {
				// Check if email exists
				mockUserRepo.EXPECT().
					FindByEmail("newuser@example.com").
					Return(models.User{}, gorm.ErrRecordNotFound).
					Times(1)

				// Check if role exists
				mockRoleRepo.EXPECT().
					FindById(uint(1)).
					Return(models.Role{ID: 1, Name: "User"}, nil).
					Times(1)

				// Create user
				mockUserRepo.EXPECT().
					Create(gomock.Any()).
					Return(models.User{
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
			wantErr: false,
		},
		{
			name: "Create user with existing email",
			input: &models.UserInput{
				Email:     "existing@example.com",
				Username:  "newuser",
				Password:  "password123",
				FirstName: "New",
				LastName:  "User",
				RoleID:    1,
			},
			mockFn: func() {
				mockUserRepo.EXPECT().
					FindByEmail("existing@example.com").
					Return(models.User{ID: 1, Email: "existing@example.com"}, nil).
					Times(1)
			},
			wantErr: true,
			errMsg:  "user with this email or username already exists",
		},
		{
			name: "Create user with invalid role",
			input: &models.UserInput{
				Email:     "newuser@example.com",
				Username:  "newuser",
				Password:  "password123",
				FirstName: "New",
				LastName:  "User",
				RoleID:    999,
			},
			mockFn: func() {
				mockUserRepo.EXPECT().
					FindByEmail("newuser@example.com").
					Return(models.User{}, gorm.ErrRecordNotFound).
					Times(1)

				mockRoleRepo.EXPECT().
					FindById(uint(999)).
					Return(models.Role{}, gorm.ErrRecordNotFound).
					Times(1)
			},
			wantErr: true,
			errMsg:  "role not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			result, err := service.CreateUser(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errMsg != "" {
				if err.Error() != tt.errMsg {
					t.Errorf("CreateUser() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}

			if !tt.wantErr {
				if result == nil {
					t.Errorf("CreateUser() result should not be nil")
				}
				if result.Email != tt.input.Email {
					t.Errorf("CreateUser() Email = %v, want %v", result.Email, tt.input.Email)
				}
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockActivityLogRepo := mocks.NewMockActivityLogRepository(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	service := services.NewUserService(mockUserRepo, mockActivityLogRepo, mockRoleRepo)

	tests := []struct {
		name    string
		userID  uint
		input   *models.UserUpdateInput
		mockFn  func()
		wantErr bool
		errMsg  string
	}{
		{
			name:   "Update user successfully",
			userID: 1,
			input: &models.UserUpdateInput{
				FirstName: "Updated",
				LastName:  "Name",
			},
			mockFn: func() {
				mockUserRepo.EXPECT().
					FindById(uint(1)).
					Return(models.User{
						ID:        1,
						Email:     "test@example.com",
						Username:  "testuser",
						FirstName: "Test",
						LastName:  "User",
						RoleID:    1,
						IsActive:  true,
					}, nil).
					Times(1)

				mockUserRepo.EXPECT().
					Update(gomock.Any()).
					Return(models.User{
						ID:        1,
						Email:     "test@example.com",
						Username:  "testuser",
						FirstName: "Updated",
						LastName:  "Name",
						RoleID:    1,
						IsActive:  true,
					}, nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name:   "Update non-existing user",
			userID: 999,
			input: &models.UserUpdateInput{
				FirstName: "Updated",
			},
			mockFn: func() {
				mockUserRepo.EXPECT().
					FindById(uint(999)).
					Return(models.User{}, gorm.ErrRecordNotFound).
					Times(1)
			},
			wantErr: true,
			errMsg:  "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			result, err := service.UpdateUser(tt.userID, tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errMsg != "" {
				if err.Error() != tt.errMsg {
					t.Errorf("UpdateUser() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}

			if !tt.wantErr {
				if result == nil {
					t.Errorf("UpdateUser() result should not be nil")
				}
			}
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockActivityLogRepo := mocks.NewMockActivityLogRepository(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	service := services.NewUserService(mockUserRepo, mockActivityLogRepo, mockRoleRepo)

	tests := []struct {
		name    string
		userID  uint
		mockFn  func()
		wantErr bool
		errMsg  string
	}{
		{
			name:   "Delete user successfully",
			userID: 1,
			mockFn: func() {
				mockUserRepo.EXPECT().
					FindById(uint(1)).
					Return(models.User{
						ID:    1,
						Email: "test@example.com",
					}, nil).
					Times(1)

				mockUserRepo.EXPECT().
					Delete(gomock.Any()).
					Return(nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name:   "Delete non-existing user",
			userID: 999,
			mockFn: func() {
				mockUserRepo.EXPECT().
					FindById(uint(999)).
					Return(models.User{}, gorm.ErrRecordNotFound).
					Times(1)
			},
			wantErr: true,
			errMsg:  "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			err := service.DeleteUser(tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errMsg != "" {
				if err.Error() != tt.errMsg {
					t.Errorf("DeleteUser() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestUserService_ActivateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockActivityLogRepo := mocks.NewMockActivityLogRepository(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	service := services.NewUserService(mockUserRepo, mockActivityLogRepo, mockRoleRepo)

	tests := []struct {
		name    string
		userID  uint
		mockFn  func()
		wantErr bool
	}{
		{
			name:   "Activate user successfully",
			userID: 1,
			mockFn: func() {
				mockUserRepo.EXPECT().
					FindById(uint(1)).
					Return(models.User{
						ID:       1,
						Email:    "test@example.com",
						IsActive: false,
					}, nil).
					Times(1)

				mockUserRepo.EXPECT().
					Update(gomock.Any()).
					DoAndReturn(func(user *models.User) (models.User, error) {
						if !user.IsActive {
							t.Errorf("ActivateUser() should set IsActive to true")
						}
						return *user, nil
					}).
					Times(1)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			result, err := service.ActivateUser(tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("ActivateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result == nil {
					t.Errorf("ActivateUser() result should not be nil")
				}
			}
		})
	}
}

func TestUserService_DeactivateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockActivityLogRepo := mocks.NewMockActivityLogRepository(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	service := services.NewUserService(mockUserRepo, mockActivityLogRepo, mockRoleRepo)

	tests := []struct {
		name    string
		userID  uint
		mockFn  func()
		wantErr bool
	}{
		{
			name:   "Deactivate user successfully",
			userID: 1,
			mockFn: func() {
				mockUserRepo.EXPECT().
					FindById(uint(1)).
					Return(models.User{
						ID:       1,
						Email:    "test@example.com",
						IsActive: true,
					}, nil).
					Times(1)

				mockUserRepo.EXPECT().
					Update(gomock.Any()).
					DoAndReturn(func(user *models.User) (models.User, error) {
						if user.IsActive {
							t.Errorf("DeactivateUser() should set IsActive to false")
						}
						return *user, nil
					}).
					Times(1)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			result, err := service.DeactivateUser(tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeactivateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result == nil {
					t.Errorf("DeactivateUser() result should not be nil")
				}
			}
		})
	}
}

func TestUserService_GetUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockActivityLogRepo := mocks.NewMockActivityLogRepository(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	service := services.NewUserService(mockUserRepo, mockActivityLogRepo, mockRoleRepo)

	tests := []struct {
		name    string
		request models.UserListRequest
		mockFn  func()
		wantErr bool
	}{
		{
			name: "Get users successfully",
			request: models.UserListRequest{
				Page:  1,
				Limit: 10,
			},
			mockFn: func() {
				mockUserRepo.EXPECT().
					FindAll(gomock.Any()).
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
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			result, err := service.GetUsers(tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result == nil {
					t.Errorf("GetUsers() result should not be nil")
				}
				if len(result.Users) != 2 {
					t.Errorf("GetUsers() users count = %v, want %v", len(result.Users), 2)
				}
			}
		})
	}
}
