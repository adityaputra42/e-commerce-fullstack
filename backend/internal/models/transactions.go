package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	TxID            string         `json:"tx_id"`
	AddressID       int64          `json:"address_id"`
	ShippingID      int64          `json:"shipping_id"`
	PaymentMethodID int64          `json:"payment_method_id"`
	ShippingPrice   float64        `json:"shipping_price"`
	TotalPrice      float64        `json:"total_price"`
	Status          string         `json:"status"`
	UpdatedAt       time.Time      `json:"updated_at"`
	CreatedAt       time.Time      `json:"created_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateTransaction struct {
	AddressID       int64         `json:"address_id"`
	ShippingID      int64         `json:"shipping_id"`
	PaymentMethodID int64         `json:"payment_method_id"`
	ShippingPrice   float64       `json:"shipping_price"`
	TotalPrice      float64       `json:"total_price"`
	ProductOrders   []CreateOrder `json:"product_orders"`
}


type UpdateTransaction struct {
	TxID   string `json:"tx_id"`
	Status string `json:"status"`
}

