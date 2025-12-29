# Testing Guide

Panduan lengkap untuk menjalankan unit test dan integration test pada backend e-commerce.

## Struktur Testing

Testing dibagi menjadi 3 layer:

### 1. **Repository Layer** (Integration Tests)
- Menggunakan **MySQL database di Docker**
- Test CRUD operations dengan database nyata
- Lokasi: `internal/repository/*_test.go`

### 2. **Service Layer** (Unit Tests)
- Menggunakan **gomock** untuk mock repository
- Test business logic tanpa database
- Lokasi: `internal/services/*_test.go`

### 3. **Handler Layer** (Unit Tests)
- Menggunakan **gomock** untuk mock service
- Menggunakan **httptest** untuk HTTP testing
- Test HTTP handlers dan responses
- Lokasi: `internal/handler/*_test.go`

## Prerequisites

1. **Docker** dan **Docker Compose** terinstall
2. **Go 1.24+** terinstall
3. **Make** terinstall

## Setup

### 1. Install Dependencies

```bash
go mod download
```

### 2. Install gomock

```bash
go get -u go.uber.org/mock/mockgen@latest
```

### 3. Generate Mocks

```bash
make mock-gen
```

## Menjalankan Tests

### Unit Tests (Service + Handler)

Test ini **TIDAK memerlukan database** karena menggunakan mocks:

```bash
make test-unit
```

atau langsung dengan go:

```bash
go test -v ./internal/services/... ./internal/handler/...
```

### Integration Tests (Repository)

Test ini **memerlukan MySQL database di Docker**:

```bash
# Start database dan run tests
make test-integration

# Atau manual:
# 1. Start database
make test-db-up

# 2. Wait for database to be ready (5-10 seconds)
sleep 5

# 3. Run tests
TEST_DB_HOST=localhost TEST_DB_PORT=3307 TEST_DB_USER=testuser \
TEST_DB_PASSWORD=testpassword TEST_DB_NAME=ecommerce_test \
go test -v ./internal/repository/...
```

### All Tests dengan Coverage

```bash
make test-coverage
```

Ini akan:
1. Start MySQL database di Docker
2. Run semua tests (unit + integration)
3. Generate coverage report di `coverage.html`

### Test Specific Function

```bash
make test-specific TEST=TestUserRepository_Create
```

atau:

```bash
go test -v -run TestUserRepository_Create ./internal/repository/...
```

## Database Management

### Start Test Database

```bash
make test-db-up
```

Database akan berjalan di:
- **Host**: localhost
- **Port**: 3307
- **Database**: ecommerce_test
- **User**: testuser
- **Password**: testpassword

### Stop Test Database

```bash
make test-db-down
```

### Clean Test Database (termasuk volumes)

```bash
make test-clean
```

## Environment Variables untuk Testing

Untuk repository tests, Anda bisa override konfigurasi database:

```bash
export TEST_DB_HOST=localhost
export TEST_DB_PORT=3307
export TEST_DB_USER=testuser
export TEST_DB_PASSWORD=testpassword
export TEST_DB_NAME=ecommerce_test
```

## Coverage Report

Setelah menjalankan `make test-coverage`, buka `coverage.html` di browser:

```bash
# Linux
xdg-open coverage.html

# macOS
open coverage.html

# Windows
start coverage.html
```

Atau lihat coverage di terminal:

```bash
go tool cover -func=coverage.out
```

## Contoh Test Cases

### Repository Test (Integration)

```go
func TestUserRepository_Create(t *testing.T) {
    // Setup database transaction
    tx := testhelper.BeginTestTransaction(t, testDB)
    defer testhelper.RollbackTestTransaction(tx)
    
    dbWrapper := testhelper.SetTestDB(tx)
    defer dbWrapper.Restore()

    repo := repository.NewUserReposiory()
    
    // Test dengan database nyata
    user := models.User{
        Email: "test@example.com",
        // ...
    }
    
    result, err := repo.Create(user)
    // assertions...
}
```

### Service Test (Unit dengan Mock)

```go
func TestUserService_GetUserById(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockUserRepo := mocks.NewMockUserRepository(ctrl)
    service := services.NewUserService(mockUserRepo, ...)

    // Setup mock expectations
    mockUserRepo.EXPECT().
        FindById(uint(1)).
        Return(models.User{ID: 1}, nil).
        Times(1)

    // Test
    result, err := service.GetUserById(1)
    // assertions...
}
```

### Handler Test (Unit dengan Mock)

```go
func TestUserHandler_GetUserById(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockUserService := mocks.NewMockUserService(ctrl)
    handler := handler.NewUserHandler(mockUserService)

    // Setup mock
    mockUserService.EXPECT().
        GetUserById(uint(1)).
        Return(&models.User{ID: 1}, nil).
        Times(1)

    // Create HTTP request
    req := httptest.NewRequest(http.MethodGet, "/api/users/1", nil)
    w := httptest.NewRecorder()

    // Test
    handler.GetUserById(w, req)
    
    // Assert HTTP response
    assert.Equal(t, http.StatusOK, w.Code)
}
```

## Best Practices

### 1. **Isolasi Tests**
- Setiap test harus independen
- Gunakan transactions untuk repository tests
- Rollback setelah setiap test

### 2. **Mock Expectations**
- Selalu set `.Times(1)` untuk memastikan method dipanggil
- Gunakan `gomock.Any()` untuk parameter yang tidak penting
- Gunakan `DoAndReturn()` untuk logic kompleks

### 3. **Test Naming**
- Format: `Test<Component>_<Method>_<Scenario>`
- Contoh: `TestUserService_CreateUser_WithExistingEmail`

### 4. **Table-Driven Tests**
- Gunakan slice of test cases
- Lebih mudah menambah test cases baru
- Lebih readable

### 5. **Cleanup**
- Selalu cleanup resources (defer)
- Rollback transactions
- Restore mocked objects

## Troubleshooting

### Database Connection Error

```bash
# Pastikan database sudah running
docker ps | grep mysql-test

# Restart database
make test-db-down
make test-db-up
sleep 5
```

### Mock Generation Error

```bash
# Install mockgen
go install go.uber.org/mock/mockgen@latest

# Regenerate mocks
make mock-gen
```

### Test Timeout

Tambahkan timeout flag:

```bash
go test -v -timeout 30s ./...
```

## CI/CD Integration

Contoh GitHub Actions workflow:

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      mysql:
        image: mysql:8.0
        env:
          MYSQL_ROOT_PASSWORD: testpassword
          MYSQL_DATABASE: ecommerce_test
        ports:
          - 3307:3306
        options: >-
          --health-cmd="mysqladmin ping"
          --health-interval=10s
          --health-timeout=5s
          --health-retries=5

    steps:
      - uses: actions/checkout@v2
      
      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.24
      
      - name: Install dependencies
        run: go mod download
      
      - name: Generate mocks
        run: make mock-gen
      
      - name: Run tests
        run: make test-coverage
        env:
          TEST_DB_HOST: 127.0.0.1
          TEST_DB_PORT: 3307
          TEST_DB_USER: root
          TEST_DB_PASSWORD: testpassword
          TEST_DB_NAME: ecommerce_test
      
      - name: Upload coverage
        uses: codecov/codecov-action@v2
        with:
          file: ./coverage.out
```

## Makefile Commands Reference

| Command | Description |
|---------|-------------|
| `make test` | Run all unit tests |
| `make test-unit` | Run unit tests (service + handler) |
| `make test-integration` | Run integration tests (repository) |
| `make test-coverage` | Run all tests with coverage report |
| `make test-db-up` | Start MySQL test database |
| `make test-db-down` | Stop MySQL test database |
| `make test-clean` | Clean database and coverage files |
| `make mock-gen` | Generate all mocks |
| `make test-specific TEST=<name>` | Run specific test |
| `make test-verbose` | Run tests with verbose output |

## Resources

- [Go Testing](https://golang.org/pkg/testing/)
- [gomock Documentation](https://github.com/uber-go/mock)
- [GORM Testing](https://gorm.io/docs/testing.html)
- [httptest Package](https://golang.org/pkg/net/http/httptest/)
