package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID            int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	TransactionID string         `json:"transaction_id" validate:"required" gorm:"not null;index"`
	TotalPayment  float64        `json:"total_payment" validate:"required,gt=0" gorm:"type:decimal(12,2);not null"`
	Status        string         `json:"status" validate:"required,oneof=pending success failed refunded" gorm:"type:varchar(20);default:'pending'"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Relasi ke Transaction
	Transaction Transaction `json:"-" gorm:"foreignKey:TransactionID;references:TxID"`
}

// Request untuk membuat pembayaran
type CreatePayment struct {
	TransactionID string  `json:"transaction_id" form:"transaction_id" validate:"required"`
	TotalPayment  float64 `json:"total_payment" form:"total_payment" validate:"required,gt=0"`
}

// Request untuk update status pembayaran
type UpdatePayment struct {
	ID     int64  `json:"id" form:"id" validate:"required"`
	Status string `json:"status" form:"status" validate:"required,oneof=pending success failed refunded"`
}

type PaymentListRequest struct {
	Limit  int
	Page   int
	SortBy string
}
