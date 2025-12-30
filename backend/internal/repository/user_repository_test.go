package repository_test

import (
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/repository"
	"e-commerce/backend/internal/testhelper"
	"testing"
)



func TestUserRepository_Create(t *testing.T) {
	// Setup
	tx := testhelper.BeginTestTransaction(t, testDB)
	defer testhelper.RollbackTestTransaction(tx)
	
	dbWrapper := testhelper.SetTestDB(tx)
	defer dbWrapper.Restore()

	repo := repository.NewUserReposiory()
	
	// Create test role first
	role := testhelper.CreateTestRole(tx, "Test Role")

	tests := []struct {
		name    string
		user    models.User
		wantErr bool
	}{
		{
			name: "Create user successfully",
			user: models.User{
				Email:        "test@example.com",
				Username:     "testuser",
				PasswordHash: "$2a$10$test.hash.password",
				FirstName:    "Test",
				LastName:     "User",
				RoleID:       role.ID,
				IsActive:     true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.Create(tt.user)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result.Email != tt.user.Email {
					t.Errorf("Create() email = %v, want %v", result.Email, tt.user.Email)
				}
				if result.Username != tt.user.Username {
					t.Errorf("Create() username = %v, want %v", result.Username, tt.user.Username)
				}
				if result.ID == 0 {
					t.Errorf("Create() ID should not be 0")
				}
			}
		})
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	// Setup
	tx := testhelper.BeginTestTransaction(t, testDB)
	defer testhelper.RollbackTestTransaction(tx)
	
	dbWrapper := testhelper.SetTestDB(tx)
	defer dbWrapper.Restore()

	repo := repository.NewUserReposiory()
	
	// Create test data
	role := testhelper.CreateTestRole(tx, "Test Role")
	testUser := testhelper.CreateTestUser(tx, "findme@example.com", role.ID)

	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "Find existing user by email",
			email:   testUser.Email,
			wantErr: false,
		},
		{
			name:    "Find non-existing user by email",
			email:   "notfound@example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.FindByEmail(tt.email)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result.Email != tt.email {
					t.Errorf("FindByEmail() email = %v, want %v", result.Email, tt.email)
				}
				if result.Role.ID == 0 {
					t.Errorf("FindByEmail() should preload Role")
				}
			}
		})
	}
}

func TestUserRepository_FindById(t *testing.T) {
	// Setup
	tx := testhelper.BeginTestTransaction(t, testDB)
	defer testhelper.RollbackTestTransaction(tx)
	
	dbWrapper := testhelper.SetTestDB(tx)
	defer dbWrapper.Restore()

	repo := repository.NewUserReposiory()
	
	// Create test data
	role := testhelper.CreateTestRole(tx, "Test Role")
	testUser := testhelper.CreateTestUser(tx, "findbyid@example.com", role.ID)

	tests := []struct {
		name    string
		id      uint
		wantErr bool
	}{
		{
			name:    "Find existing user by ID",
			id:      testUser.ID,
			wantErr: false,
		},
		{
			name:    "Find non-existing user by ID",
			id:      99999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.FindById(tt.id)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("FindById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result.ID != tt.id {
					t.Errorf("FindById() ID = %v, want %v", result.ID, tt.id)
				}
				if result.Role.ID == 0 {
					t.Errorf("FindById() should preload Role")
				}
			}
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	// Setup
	tx := testhelper.BeginTestTransaction(t, testDB)
	defer testhelper.RollbackTestTransaction(tx)
	
	dbWrapper := testhelper.SetTestDB(tx)
	defer dbWrapper.Restore()

	repo := repository.NewUserReposiory()
	
	// Create test data
	role := testhelper.CreateTestRole(tx, "Test Role")
	testUser := testhelper.CreateTestUser(tx, "update@example.com", role.ID)

	tests := []struct {
		name      string
		updateFn  func(*models.User)
		wantErr   bool
		checkFn   func(*testing.T, models.User)
	}{
		{
			name: "Update user first name",
			updateFn: func(u *models.User) {
				u.FirstName = "Updated"
			},
			wantErr: false,
			checkFn: func(t *testing.T, result models.User) {
				if result.FirstName != "Updated" {
					t.Errorf("Update() FirstName = %v, want %v", result.FirstName, "Updated")
				}
			},
		},
		{
			name: "Update user is_active status",
			updateFn: func(u *models.User) {
				u.IsActive = false
			},
			wantErr: false,
			checkFn: func(t *testing.T, result models.User) {
				if result.IsActive != false {
					t.Errorf("Update() IsActive = %v, want %v", result.IsActive, false)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get fresh copy of user
			user, _ := repo.FindById(testUser.ID)
			tt.updateFn(&user)

			result, err := repo.Update(&user)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.checkFn != nil {
				tt.checkFn(t, result)
			}
		})
	}
}

func TestUserRepository_Delete(t *testing.T) {
	// Setup
	tx := testhelper.BeginTestTransaction(t, testDB)
	defer testhelper.RollbackTestTransaction(tx)
	
	dbWrapper := testhelper.SetTestDB(tx)
	defer dbWrapper.Restore()

	repo := repository.NewUserReposiory()
	
	// Create test data
	role := testhelper.CreateTestRole(tx, "Test Role")
	testUser := testhelper.CreateTestUser(tx, "delete@example.com", role.ID)

	tests := []struct {
		name    string
		user    models.User
		wantErr bool
	}{
		{
			name:    "Delete existing user",
			user:    *testUser,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(tt.user)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify user is deleted
			if !tt.wantErr {
				_, err := repo.FindById(tt.user.ID)
				if err == nil {
					t.Errorf("Delete() user still exists after deletion")
				}
			}
		})
	}
}

func TestUserRepository_FindAll(t *testing.T) {
	// Setup
	tx := testhelper.BeginTestTransaction(t, testDB)
	defer testhelper.RollbackTestTransaction(tx)
	
	dbWrapper := testhelper.SetTestDB(tx)
	defer dbWrapper.Restore()

	repo := repository.NewUserReposiory()
	
	// Create test data
	role := testhelper.CreateTestRole(tx, "Test Role")
	testhelper.CreateTestUser(tx, "user1@example.com", role.ID)
	testhelper.CreateTestUser(tx, "user2@example.com", role.ID)
	testhelper.CreateTestUser(tx, "user3@example.com", role.ID)

	tests := []struct {
		name      string
		request   models.UserListRequest
		wantCount int
		wantErr   bool
	}{
		{
			name: "Get all users with pagination",
			request: models.UserListRequest{
				Page:  1,
				Limit: 10,
			},
			wantCount: 3,
			wantErr:   false,
		},
		{
			name: "Get users with limit",
			request: models.UserListRequest{
				Page:  1,
				Limit: 2,
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "Get second page",
			request: models.UserListRequest{
				Page:  2,
				Limit: 2,
			},
			wantCount: 1,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.FindAll(tt.request)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(result.Users) != tt.wantCount {
					t.Errorf("FindAll() count = %v, want %v", len(result.Users), tt.wantCount)
				}
				if result.Total != 3 {
					t.Errorf("FindAll() total = %v, want %v", result.Total, 3)
				}
				if result.Page != tt.request.Page {
					t.Errorf("FindAll() page = %v, want %v", result.Page, tt.request.Page)
				}
			}
		})
	}
}
