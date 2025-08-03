package main

import (
	"e-commerce/backend/internal/config"
	"e-commerce/backend/internal/database"
	"log"
)

func main() {
	cfg := config.Load()

	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

}
