package models

import (
	"time"

	"gorm.io/gorm"
)

// Shipping (DB Model)
type Shipping struct {
	ID        int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string         `json:"name" validate:"required,min=3,max=100" gorm:"type:varchar(100);not null"`
	Price     float64        `json:"price" validate:"required,gt=0" gorm:"type:decimal(12,2);not null"`
	State     string         `json:"state" validate:"required,oneof=active inactive" gorm:"type:varchar(20);default:'active'"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Create Shipping (Request Payload)
type CreateShipping struct {
	Name  string  `json:"name" form:"name" validate:"required,min=3,max=100"`
	Price float64 `json:"price" form:"price" validate:"required,gt=0"`
	State string  `json:"state" form:"state" validate:"omitempty,oneof=active inactive"`
}

// Update Shipping (Request Payload)
type UpdateShipping struct {
	ID    int64   `json:"id" form:"id" validate:"required,gt=0"`
	Name  string  `json:"name" form:"name" validate:"omitempty,min=3,max=100"`
	Price float64 `json:"price" form:"price" validate:"omitempty,gt=0"`
	State string  `json:"state" form:"state" validate:"omitempty,oneof=active inactive"`
}

// Response (API Output)
type ShippingResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	State     string    `json:"state"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// Transaction Response (gabungan beberapa relasi)
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

type ShippingListRequest struct {
	Limit  int
	Page   int
	SortBy string
	Search string
}
