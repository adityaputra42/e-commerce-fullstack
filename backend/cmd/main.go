package main

import (
	"e-commerce/backend/internal/config"
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/di"
	"e-commerce/backend/internal/routes"
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// =======================
// Logger Initialization
// =======================

func initLogger() *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "msg",
		CallerKey:      "caller",
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

	return logger
}

func main() {
	logger := initLogger()
	defer logger.Sync()

	cfg := config.Load()

	logger.Info("========================================")
	logger.Info("🔧 Initializing Application...")
	logger.Info("========================================")

	logger.Info("📦 Connecting to database...")
	config.InitSupabase(*cfg)

	if err := database.Connect(cfg); err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	logger.Info("✅ Database connected successfully")

	logger.Info("🔄 Running database migrations...")
	if err := database.Migrate(); err != nil {
		logger.Fatal("Failed to migrate database", zap.Error(err))
	}
	logger.Info("✅ Migrations completed")

	logger.Info("🌱 Running database seeders...")
	if err := database.SeedDatabase(cfg); err != nil {
		logger.Fatal("Failed to seed database", zap.Error(err))
	}
	logger.Info("✅ Database seeded successfully")

	handler := di.InitializeAllHandler(cfg)

	router := routes.SetupRoutes(handler, logger, cfg.CORS)

	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}


	fmt.Printf("\n")
	fmt.Printf("╔════════════════════════════════════════╗\n")
	fmt.Printf("║  Server is running on port %-4s        ║\n", port)
	fmt.Printf("║ ══════════════════════════════════════ ║\n")
	fmt.Printf("║  Server is running on port %-4s        ║\n", port)
	fmt.Printf("╚════════════════════════════════════════╝\n")
	fmt.Printf("\n")

	if err := http.ListenAndServe(":"+port, router); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
