// Command migrate menjalankan database migration + seeding sebagai
// deployment step yang eksplisit dan terpisah dari boot aplikasi utama.
//
// Dibuat karena SEBELUMNYA migrate+seed selalu jalan otomatis di dalam
// cmd/main.go setiap kali aplikasi start — rawan race condition kalau lebih
// dari satu instance start bersamaan, dan menghilangkan kontrol eksplisit
// atas kapan migrasi sungguhan dijalankan.
//
// Pemakaian:
//
//	go run cmd/migrate/main.go
//
// Lalu di cmd/main.go / .env production, set RUN_MIGRATIONS_ON_BOOT=false
// supaya aplikasi utama tidak lagi migrate+seed otomatis saat start.
package main

import (
	"e-commerce/backend/internal/config"
	"e-commerce/backend/internal/database"
	"log"
)

func main() {
	cfg := config.Load()

	log.Println("📦 Connecting to database...")
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("✅ Database connected")

	log.Println("🔄 Running database migrations...")
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("✅ Migrations completed")

	log.Println("🌱 Running database seeders...")
	if err := database.SeedDatabase(cfg); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}
	log.Println("✅ Database seeded successfully")
}
