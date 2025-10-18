package models

import (
	"time"

	"gorm.io/gorm"
)

// PaymentMethod Model (DB)
type PaymentMethod struct {
	ID            int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	AccountName   string         `json:"account_name" validate:"required,min=3,max=100" gorm:"type:varchar(100);not null"`
	AccountNumber string         `json:"account_number" validate:"required,numeric,min=5,max=30" gorm:"type:varchar(30);not null"`
	BankName      string         `json:"bank_name" validate:"required,min=3,max=100" gorm:"type:varchar(100);not null"`
	BankImages    string         `json:"bank_images" validate:"omitempty,url" gorm:"type:text"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// Create Payment Method (Request Payload)
type CreatePaymentMethod struct {
	AccountName   string `json:"account_name" form:"account_name" validate:"required,min=3,max=100"`
	AccountNumber string `json:"account_number" form:"account_number" validate:"required,numeric,min=5,max=30"`
	BankName      string `json:"bank_name" form:"bank_name" validate:"required,min=3,max=100"`
	BankImages    string `json:"bank_images" form:"bank_images" validate:"omitempty,url"`
}

// Update Payment Method (Request Payload)
type UpdatePaymentMethod struct {
	ID            int64  `json:"id" form:"id" validate:"required,gt=0"`
	AccountName   string `json:"account_name" form:"account_name" validate:"omitempty,min=3,max=100"`
	AccountNumber string `json:"account_number" form:"account_number" validate:"omitempty,numeric,min=5,max=30"`
	BankName      string `json:"bank_name" form:"bank_name" validate:"omitempty,min=3,max=100"`
	BankImages    string `json:"bank_images" form:"bank_images" validate:"omitempty,url"`
}

// Response Struct (API Output)
type PaymentMethodResponse struct {
	ID            int64     `json:"id"`
	AccountName   string    `json:"account_name"`
	AccountNumber string    `json:"account_number"`
	BankName      string    `json:"bank_name"`
	BankImages    string    `json:"bank_images"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedAt     time.Time `json:"created_at"`
}

type PaymentMethodListRequest struct {
	Limit  int
	Page   int
	SortBy string
}
