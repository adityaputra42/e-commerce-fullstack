package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	TxID            string         `json:"tx_id" gorm:"primaryKey;type:varchar(50)"`
	AddressID       int64          `json:"address_id" validate:"required" gorm:"not null"`
	ShippingID      int64          `json:"shipping_id" validate:"required" gorm:"not null"`
	PaymentMethodID int64          `json:"payment_method_id" validate:"required" gorm:"not null"`
	ShippingPrice   float64        `json:"shipping_price" gorm:"type:decimal(12,2);not null"`
	TotalPrice      float64        `json:"total_price"  gorm:"type:decimal(12,2);not null"`
	Status          string         `json:"status" gorm:"type:varchar(30);not null"`
	CreatedAt       time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
	Address         *Address       `json:"address,omitempty" gorm:"foreignKey:AddressID;references:ID"`
	Shipping        *Shipping      `json:"shipping,omitempty" gorm:"foreignKey:ShippingID;references:ID"`
	PaymentMethod   *PaymentMethod `json:"payment_method,omitempty" gorm:"foreignKey:PaymentMethodID;references:ID"`
	Orders          []Order        `json:"orders,omitempty" gorm:"foreignKey:TransactionID;references:TxID"`
}

type CreateTransaction struct {
	AddressID       int64         `json:"address_id" form:"address_id" validate:"required"`
	ShippingID      int64         `json:"shipping_id" form:"shipping_id" validate:"required"`
	PaymentMethodID int64         `json:"payment_method_id" form:"payment_method_id" validate:"required"`
	ShippingPrice   float64       `json:"shipping_price" form:"shipping_price" validate:"required,gt=0"`
	TotalPrice      float64       `json:"total_price" form:"total_price" validate:"required,gt=0"`
	ProductOrders   []CreateOrder `json:"product_orders" validate:"required,dive"`
}

type UpdateTransaction struct {
	TxID   string `json:"tx_id" form:"tx_id" validate:"required"`
	Status string `json:"status" form:"status" validate:"required,oneof=pending paid shipped cancelled completed"`
}

type TransactionListRequest struct {
	Limit  int
	Page   int
	SortBy string
	Search string
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
