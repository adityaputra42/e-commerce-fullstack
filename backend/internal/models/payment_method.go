package models

import (
	"time"

	"gorm.io/gorm"
)

type PaymentMethod struct {
	ID            int64          `json:"id"`
	AccountName   string         `json:"account_name"`
	AccountNumber string         `json:"account_number"`
	BankName      string         `json:"bank_name"`
	BankImages    string         `json:"bank_images"`
	UpdatedAt     time.Time      `json:"updated_at"`
	CreatedAt     time.Time      `json:"created_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreatePaymentMethod struct {
	AccountName   string `form:"account_name"`
	AccountNumber string `form:"account_number"`
	BankName      string `form:"bank_name"`
	BankImages    string `form:"bank_images"`
}

type UpdatePaymentMethod struct {
	ID            int64  `form:"id"`
	AccountName   string `form:"account_name"`
	AccountNumber string `form:"account_number"`
	BankName      string `form:"bank_name"`
	BankImages    string `form:"bank_images"`
}

type PaymentMethodResponse struct {
	ID            int64     `json:"id"`
	AccountName   string    `json:"account_name"`
	AccountNumber string    `json:"account_number"`
	BankName      string    `json:"bank_name"`
	BankImages    string    `json:"bank_images"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedAt     time.Time `json:"created_at"`
}
