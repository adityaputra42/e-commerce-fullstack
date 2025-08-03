package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID            int64          `json:"id"`
	TransactionID string         `json:"transaction_id"`
	TotalPayment  float64        `json:"total_payment"`
	Status        string         `json:"status"`
	UpdatedAt     time.Time      `json:"updated_at"`
	CreatedAt     time.Time      `json:"created_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreatePayment struct {
	TransactionID string  `json:"transaction_id"`
	TotalPayment  float64 `json:"total_payment"`
}
