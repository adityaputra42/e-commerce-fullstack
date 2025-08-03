package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID            string         `json:"id"`
	TransactionID string         `json:"transaction_id"`
	ProductID     int64          `json:"product_id"`
	ColorVarianID int64          `json:"color_varian_id"`
	SizeVarianID  int64          `json:"size_varian_id"`
	UnitPrice     float64        `json:"unit_price"`
	Subtotal      float64        `json:"subtotal"`
	Quantity      int64          `json:"quantity"`
	Status        string         `json:"status"`
	UpdatedAt     time.Time      `json:"updated_at"`
	CreatedAt     time.Time      `json:"created_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type UpdateOrder struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type CreateOrder struct {
	ProductID     int64   `json:"product_id"`
	ColorVarianID int64   `json:"color_varian_id"`
	SizeVarianID  int64   `json:"size_varian_id"`
	UnitPrice     float64 `json:"unit_price"`
	Subtotal      float64 `json:"subtotal"`
	Quantity      int64   `json:"quantity"`
}

type OrderResponse struct {
	ID            string               `json:"id"`
	TransactionID string               `json:"transaction_id"`
	Product       ProductOrderResponse `json:"product"`
	Size          string               `json:"size"`
	UnitPrice     float64              `json:"unit_price"`
	Subtotal      float64              `json:"subtotal"`
	Quantity      int64                `json:"quantity"`
	Status        string               `json:"status"`
	UpdatedAt     time.Time            `json:"updated_at"`
	CreatedAt     time.Time            `json:"created_at"`
}

type ProductOrderResponse struct {
	ID          int64                    `json:"id"`
	Name        string                   `json:"name"`
	Category    Category                 `json:"category"`
	Description string                   `json:"description"`
	Images      string                   `json:"images"`
	Rating      float64                  `json:"rating"`
	Price       float64                  `json:"price"`
	ColorVarian ColorVarianOrderResponse `json:"color_varian"`
	UpdatedAt   time.Time                `json:"updated_at"`
	CreatedAt   time.Time                `json:"created_at"`
}

type ColorVarianOrderResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	Images    string    `json:"images"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
