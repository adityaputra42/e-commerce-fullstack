package main

import (
	"e-commerce/backend/internal/config"
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/di"
	"e-commerce/backend/internal/routes"
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

func initLogger() *logrus.Logger {
	logger := logrus.New()

	// Set output ke stdout
	logger.SetOutput(os.Stdout)

	// Set formatter (pilih salah satu)
	// Format JSON - bagus untuk production
	// logger.SetFormatter(&logrus.JSONFormatter{
	// 	TimestampFormat: "2006-01-02 15:04:05",
	// 	PrettyPrint:     false,
	// })

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})

	// Set level logging
	logger.SetLevel(logrus.DebugLevel)

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
