package models

import (
	"time"

	"gorm.io/gorm"
)

// Address Model (DB)
type Address struct {
	ID                   int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID               int64          `json:"user_id" validate:"required" gorm:"not null;index"`
	RecipientName        string         `json:"recipient_name" validate:"required,min=3,max=100" gorm:"type:varchar(100);not null"`
	RecipientPhoneNumber string         `json:"recipient_phone_number" validate:"required,e164" gorm:"type:varchar(20);not null"`
	Province             string         `json:"province" validate:"required" gorm:"type:varchar(100);not null"`
	City                 string         `json:"city" validate:"required" gorm:"type:varchar(100);not null"`
	District             string         `json:"district" validate:"required" gorm:"type:varchar(100);not null"`
	Village              string         `json:"village" validate:"required" gorm:"type:varchar(100);not null"`
	PostalCode           string         `json:"postal_code" validate:"required,len=5,numeric" gorm:"type:varchar(10);not null"`
	FullAddress          string         `json:"full_address" validate:"required,min=10,max=255" gorm:"type:varchar(255);not null"`
	UpdatedAt            time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedAt            time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt            gorm.DeletedAt `json:"-" gorm:"index"`

	User User `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Address) TableName() string {
	return "addresses"
}

// Create Address (Payload)
type CreateAddress struct {
	RecipientName        string `json:"recipient_name" validate:"required,min=3,max=100"`
	RecipientPhoneNumber string `json:"recipient_phone_number" validate:"required,e164"`
	Province             string `json:"province" validate:"required"`
	City                 string `json:"city" validate:"required"`
	District             string `json:"district" validate:"required"`
	Village              string `json:"village" validate:"required"`
	PostalCode           string `json:"postal_code" validate:"required,len=5,numeric"`
	FullAddress          string `json:"full_address" validate:"required,min=10,max=255"`
}

// Update Address (Payload)
type UpdateAddress struct {
	RecipientName        string `json:"recipient_name" validate:"omitempty,min=3,max=100"`
	RecipientPhoneNumber string `json:"recipient_phone_number" validate:"omitempty,e164"`
	Province             string `json:"province" validate:"omitempty"`
	City                 string `json:"city" validate:"omitempty"`
	District             string `json:"district" validate:"omitempty"`
	Village              string `json:"village" validate:"omitempty"`
	PostalCode           string `json:"postal_code" validate:"omitempty,len=5,numeric"`
	FullAddress          string `json:"full_address" validate:"omitempty,min=10,max=255"`
}

// Response Struct
type AddressResponse struct {
	ID                   int64     `json:"id"`
	RecipientName        string    `json:"recipient_name"`
	RecipientPhoneNumber string    `json:"recipient_phone_number"`
	Province             string    `json:"province"`
	City                 string    `json:"city"`
	District             string    `json:"district"`
	Village              string    `json:"village"`
	PostalCode           string    `json:"postal_code"`
	FullAddress          string    `json:"full_address"`
	UpdatedAt            time.Time `json:"updated_at"`
	CreatedAt            time.Time `json:"created_at"`
}

type AddressListRequest struct {
	UserId *uint
	Limit  int
	Page   int
	SortBy string
}

type AddressListResponse struct {
	Activities []AddressResponse `json:"activities"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

// Convert Address model ke Response
func (a *Address) ToResponse() *AddressResponse {
	return &AddressResponse{
		ID:                   a.ID,
		RecipientName:        a.RecipientName,
		RecipientPhoneNumber: a.RecipientPhoneNumber,
		Province:             a.Province,
		City:                 a.City,
		District:             a.District,
		Village:              a.Village,
		PostalCode:           a.PostalCode,
		FullAddress:          a.FullAddress,
		UpdatedAt:            a.UpdatedAt,
		CreatedAt:            a.CreatedAt,
	}
}
