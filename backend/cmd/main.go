package main

import (
	"e-commerce/backend/internal/config"
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/di"
	"e-commerce/backend/internal/routes"
	"e-commerce/backend/internal/utils"
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

// @title E-Commerce API
// @version 1.0
// @description This is a sample server for an e-commerce application.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @query.collection.format multi

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and then your token.

func initLogger() *logrus.Logger {
	logger := logrus.New()

	logger.SetOutput(os.Stdout)

	logger.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&utils.LogFormatter{})

	return logger
}

func main() {
	// Initialize logger
	logger := initLogger()

	cfg := config.Load()

	logger.Info("========================================")
	logger.Info("🔧 Initializing Application...")
	logger.Info("========================================")

	logger.Info("📦 Connecting to database...")
	config.InitSupabase(*cfg)
	if err := database.Connect(cfg); err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	logger.Info("✅ Database connected successfully")

	logger.Info("🔄 Running database migrations...")
	if err := database.Migrate(); err != nil {
		logger.WithError(err).Fatal("Failed to migrate database")
	}
	logger.Info("✅ Migrations completed")

	logger.Info("🌱 Running database seeders...")
	if err := database.SeedDatabase(cfg); err != nil {
		logger.WithError(err).Fatal("Failed to seed database")
	}
	logger.Info("✅ Database seeded successfully")

	handler := di.InitializeAllHandler(cfg)

	router := routes.SetupRoutes(handler, logger, cfg.CORS)

	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	logger.WithFields(logrus.Fields{
		"port": port,
	}).Info("🚀 Server starting...")

	fmt.Printf("\n")
	fmt.Printf("╔════════════════════════════════════════╗\n")
	fmt.Printf("║  Server is running on port %-4s        ║\n", port)
	fmt.Printf("╚════════════════════════════════════════╝\n")
	fmt.Printf("\n")

	if err := http.ListenAndServe(":"+port, router); err != nil {
		logger.WithError(err).Fatal("Failed to start server")
	}
}
