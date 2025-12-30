.PHONY: test test-unit test-integration test-coverage test-db-up test-db-down test-clean mock-gen

# Run all tests
test: test-unit

# Run unit tests (service and handler tests with mocks)
test-unit:
	@echo "Running unit tests..."
	@go test -v -race -coverprofile=coverage.out ./internal/services/... ./internal/handler/...

# Run integration tests (repository tests with real database)
test-integration: test-db-up
	@echo "Waiting for database to be ready..."
	@sleep 5
	@echo "Running integration tests..."
	@TEST_DB_HOST=localhost TEST_DB_PORT=3307 TEST_DB_USER=testuser TEST_DB_PASSWORD=testpassword TEST_DB_NAME=ecommerce_test \
		go test -v -race ./internal/repository/...

# Run all tests with coverage
test-coverage: test-db-up
	@echo "Running all tests with coverage..."
	@sleep 5
	@TEST_DB_HOST=localhost TEST_DB_PORT=3307 TEST_DB_USER=testuser TEST_DB_PASSWORD=testpassword TEST_DB_NAME=ecommerce_test \
		go test -v -race -coverprofile=coverage.out -covermode=atomic ./internal/...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Start test database
test-db-up:
	@echo "Starting test database..."
	@docker compose -f docker-compose.test.yml up -d
	@echo "Test database started on port 3307"

# Stop test database
test-db-down:
	@echo "Stopping test database..."
	@docker compose -f docker-compose.test.yml down
	@echo "Test database stopped"

# Clean test database and volumes
test-clean: test-db-down
	@echo "Cleaning test database volumes..."
	@docker compose -f docker-compose.test.yml down -v
	@rm -f coverage.out coverage.html
	@echo "Test environment cleaned"

# Generate mocks
mock-gen:
	@echo "Generating mocks..."
	@go run go.uber.org/mock/mockgen@latest -source=internal/repository/user_repository.go -destination=internal/mocks/mock_user_repository.go -package=mocks
	@go run go.uber.org/mock/mockgen@latest -source=internal/repository/role_reposiotry.go -destination=internal/mocks/mock_role_repository.go -package=mocks
	@go run go.uber.org/mock/mockgen@latest -source=internal/repository/activity_log_repository.go -destination=internal/mocks/mock_activity_log_repository.go -package=mocks
	@go run go.uber.org/mock/mockgen@latest -source=internal/repository/product_repository.go -destination=internal/mocks/mock_product_repository.go -package=mocks
	@go run go.uber.org/mock/mockgen@latest -source=internal/repository/category_repository.go -destination=internal/mocks/mock_category_repository.go -package=mocks
	@go run go.uber.org/mock/mockgen@latest -source=internal/repository/order_repository.go -destination=internal/mocks/mock_order_repository.go -package=mocks
	@go run go.uber.org/mock/mockgen@latest -source=internal/services/user_service.go -destination=internal/mocks/mock_user_service.go -package=mocks
	@go run go.uber.org/mock/mockgen@latest -source=internal/services/product_service.go -destination=internal/mocks/mock_product_service.go -package=mocks
	@go run go.uber.org/mock/mockgen@latest -source=internal/services/auth_service.go -destination=internal/mocks/mock_auth_service.go -package=mocks
	@go run go.uber.org/mock/mockgen@latest -source=internal/services/order_service.go -destination=internal/mocks/mock_order_service.go -package=mocks
	@go run go.uber.org/mock/mockgen@latest -source=internal/services/transaction_service.go -destination=internal/mocks/mock_transaction_service.go -package=mocks
	@echo "Mocks generated successfully"

# Run specific test
test-specific:
	@echo "Running specific test: $(TEST)"
	@go test -v -run $(TEST) ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	@go test -v ./...

# Check test coverage percentage
test-coverage-check: test-coverage
	@echo "Checking coverage..."
	@go tool cover -func=coverage.out | grep total | awk '{print "Total coverage: " $$3}'