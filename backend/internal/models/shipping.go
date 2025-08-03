package models

import (
	"time"

	"gorm.io/gorm"
)

type Shipping struct {
	ID        int64          `json:"id"`
	Name      string         `json:"name"`
	Price     float64        `json:"price"`
	State     string         `json:"state"`
	UpdatedAt time.Time      `json:"updated_at"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateShipping struct {
	Name  string  `json:"name"`
	Price float32 `json:"price"`
	State string  `json:"state"`
}

type UpadateShipping struct {
	Id    int64   `json:"id"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
	State string  `json:"state"`
}

type ShippingResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	State     string    `json:"state"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type TransactionResponse struct {
	TxID          string                `json:"tx_id"`
	Address       AddressResponse       `json:"address"`
	Shipping      ShippingResponse      `json:"shipping"`
	PaymentMethod PaymentMethodResponse `json:"payment_method"`
	ShippingPrice float64               `json:"shipping_price"`
	TotalPrice    float64               `json:"total_price"`
	Status        string                `json:"status"`
	Orders        []OrderResponse       `json:"orders"`
	UpdatedAt     time.Time             `json:"updated_at"`
	CreatedAt     time.Time             `json:"created_at"`
}
