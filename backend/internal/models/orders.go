package models

import (
	"time"

	"gorm.io/gorm"
)

// Order Model
type Order struct {
	ID            string         `json:"id" gorm:"primaryKey;type:char(36)" validate:"required,uuid4"`
	UserID        int64          `json:"user_id" validate:"required" gorm:"not null;index"`
	TransactionID string         `json:"transaction_id" gorm:"type:char(36);index" validate:"required,uuid4"`
	ProductID     int64          `json:"product_id" gorm:"not null" validate:"required"`
	ColorVarianID int64          `json:"color_varian_id" gorm:"not null" validate:"required"`
	SizeVarianID  int64          `json:"size_varian_id" gorm:"not null" validate:"required"`
	UnitPrice     float64        `json:"unit_price" gorm:"type:decimal(10,2);not null" validate:"required,gt=0"`
	Subtotal      float64        `json:"subtotal" gorm:"type:decimal(10,2);not null" validate:"required,gt=0"`
	Quantity      int64          `json:"quantity" gorm:"not null" validate:"required,gt=0"`
	Status        string         `json:"status" gorm:"type:varchar(20);default:'pending'" validate:"required,oneof=pending paid shipped completed cancelled"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	User        User        `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Product     Product     `json:"product" gorm:"foreignKey:ProductID"`
	ColorVarian ColorVarian `json:"color_varian" gorm:"foreignKey:ColorVarianID"`
	SizeVarian  SizeVarian  `json:"size_varian" gorm:"foreignKey:SizeVarianID"`
}

// Payload untuk Update
type UpdateOrder struct {
	ID     string `json:"id" validate:"required,uuid4"`
	Status string `json:"status" validate:"required,oneof=pending paid shipped completed cancelled"`
}

// Payload untuk Create
type CreateOrder struct {
	ProductID     int64   `json:"product_id" validate:"required"`
	ColorVarianID int64   `json:"color_varian_id" validate:"required"`
	SizeVarianID  int64   `json:"size_varian_id" validate:"required"`
	UnitPrice     float64 `json:"unit_price" validate:"required,gt=0"`
	Subtotal      float64 `json:"subtotal" validate:"required,gt=0"`
	Quantity      int64   `json:"quantity" validate:"required,gt=0"`
}

// Response Struct
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
	CategoryID  int64                    `json:"category_id"`
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

type OrderListRequest struct {
	Limit  int64
	Page   int64
	UserId int64
	SortBy string
}

func (cl *ColorVarian) ToResponseColorVarianOrder() ColorVarianOrderResponse {
	return ColorVarianOrderResponse{
		ID:        cl.ID,
		Name:      cl.Name,
		Color:     cl.Color,
		Images:    cl.Images,
		UpdatedAt: cl.UpdatedAt,
		CreatedAt: cl.CreatedAt,
	}
}

// ToProductResponse converts Product model to ProductResponse (for list)
func (p *Product) ToResponseProductOrder(colorVarian *ColorVarian) ProductOrderResponse {
	return ProductOrderResponse{
		ID:          p.ID,
		CategoryID:  p.CategoryID,
		Name:        p.Name,
		Description: p.Description,
		Images:      p.Images,
		Rating:      p.Rating,
		Price:       p.Price,
		ColorVarian: colorVarian.ToResponseColorVarianOrder(),
		UpdatedAt:   p.UpdatedAt,
		CreatedAt:   p.CreatedAt,
	}
}

// ToProductResponse converts Product model to ProductResponse (for list)
func (p *Order) ToOrderResponse() OrderResponse {
	return OrderResponse{
		ID:            p.ID,
		TransactionID: p.TransactionID,
		Product:       p.Product.ToResponseProductOrder(&p.ColorVarian),
		UpdatedAt:     p.UpdatedAt,
		CreatedAt:     p.CreatedAt,
	}
}
