package main

import (
	"e-commerce/backend/internal/config"
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/di"
	"e-commerce/backend/internal/routes"
	"fmt"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()

	log.Println("========================================")
	log.Println("🔧 Initializing Application...")
	log.Println("========================================")

	log.Println("📦 Connecting to database...")
	config.InitSupabase(*cfg)
	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("🔄 Running database migrations...")
	if err := database.Migrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("✅ Migrations completed")

	log.Println("🌱 Running database seeders...")
	if err := database.SeedDatabase(cfg); err != nil {
		log.Fatal("Failed to seed database:", err)
	}
	handler := di.InitializeAllHandler(cfg)
	router := routes.SetupRoutes(handler)

	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is running on port %s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
