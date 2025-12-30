package main

import (
	"e-commerce/backend/internal/config"
	"e-commerce/backend/internal/di"
	"e-commerce/backend/internal/routes"
	"testing"
)

func TestAppStartup(t *testing.T) {
	// Setup Dummy Config
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: "8080",
		},
		Database: config.DatabaseConfig{
			Host: "localhost",
		},
		JWT: config.JWTConfig{
			Secret: "testsecret",
		},
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"*"},
		},
	}

	// 1. Initialize Handlers (Dependency Injection)
	// This tests that our DI wiring is correct and constructors don't panic
	handler := di.InitializeAllHandler(cfg)
	if handler == nil {
		t.Fatal("InitializeAllHandler returned nil")
	}

	// Verify critical handlers are initialized
	// Note: We check pointers to ensure DI worked
	if handler.AuthHandler == nil {
		t.Error("AuthHandler is nil")
	}
	if handler.UserHandler == nil {
		t.Error("UserHandler is nil")
	}
	if handler.ProductHandler == nil {
		t.Error("ProductHandler is nil")
	}

	// 2. Setup Routes
	// This tests that our routing setup logic is correct
	logger := initLogger() // Accessing unexported function from main.go
	router := routes.SetupRoutes(handler, logger, cfg.CORS)

	if router == nil {
		t.Error("SetupRoutes returned nil")
	}
}
