package models

import (
	"mime/multipart"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          int64          `json:"id"`
	CategoryID  int64          `json:"category_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Images      string         `json:"images"`
	Rating      float64        `json:"rating"`
	Price       float64        `json:"price"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type ColorVarian struct {
	ID        int64          `json:"id"`
	ProductID int64          `json:"product_id"`
	Name      string         `json:"name"`
	Color     string         `json:"color"`
	Images    string         `json:"images"`
	UpdatedAt time.Time      `json:"updated_at"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type SizeVarian struct {
	ID            int64          `json:"id"`
	ColorVarianID int64          `json:"color_varian_id"`
	Size          string         `json:"size"`
	Stock         int64          `json:"stock"`
	UpdatedAt     time.Time      `json:"updated_at"`
	CreatedAt     time.Time      `json:"created_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateProduct struct {
	CategoryID  int64                 `form:"category_id"`
	Name        string                `form:"name"`
	Description string                `form:"description"`
	Images      *multipart.FileHeader `form:"images"`
	Rating      float32               `form:"rating"`
	Price       float64               `form:"price"`
	ColorVarian string                `form:"color_varians"`
}

type CreateColorVarianProduct struct {
	ProductId int64                 `json:"product_id"`
	Name      string                `json:"name"`
	Color     string                `json:"color"`
	Images    *multipart.FileHeader `json:"images"`
	Sizes     string                `json:"sizes"`
}

type CreateSizeVarianProduct struct {
	ColorVarianId int64  `json:"color_varian_id"`
	Size          string `json:"size"`
	Stock         int64  `json:"stock"`
}

type UpdateSizeVarianProduct struct {
	ID    int64  `json:"id"`
	Size  string `json:"size"`
	Stock int64  `json:"stock"`
}

type UpdateColorVarianProduct struct {
	Id     int64                 `json:"id"`
	Name   string                `json:"name"`
	Color  string                `json:"color"`
	Images *multipart.FileHeader `json:"-"`
	Sizes  string                `json:"sizes"`
}

type UpdateProduct struct {
	ID          int64                 `form:"id"`
	CategoryID  int64                 `form:"category_id"`
	Name        string                `form:"name"`
	Description string                `form:"description"`
	Images      *multipart.FileHeader `form:"images"`
	Rating      float32               `form:"rating"`
	Price       float64               `form:"price"`
	ColorVarian string                `form:"color_varians"`
}

type ProductResponse struct {
	ID          int64             `json:"id"`
	Category    CategoryResponnse `json:"category"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Images      string            `json:"images"`
	Rating      float64           `json:"rating"`
	Price       float64           `json:"price"`
	UpdatedAt   time.Time         `json:"updated_at"`
	CreatedAt   time.Time         `json:"created_at"`
}

type ProductDetailResponse struct {
	ID          int64                 `json:"id"`
	Category    CategoryResponnse     `json:"category"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Images      string                `json:"images"`
	Rating      float64               `json:"rating"`
	Price       float64               `json:"price"`
	ColorVarian []ColorVarianResponse `json:"color_varian"`
	UpdatedAt   time.Time             `json:"updated_at"`
	CreatedAt   time.Time             `json:"created_at"`
}

type ColorVarianResponse struct {
	ID         int64                `json:"id"`
	ProductID  int64                `json:"product_id"`
	Name       string               `json:"name"`
	Color      string               `json:"color"`
	Images     string               `json:"images"`
	SizeVarian []SizeVarianResponse `json:"size_varian"`
	UpdatedAt  time.Time            `json:"updated_at"`
	CreatedAt  time.Time            `json:"created_at"`
}

type SizeVarianResponse struct {
	ID            int64     `json:"id"`
	ColorVarianID int64     `json:"color_varian_id"`
	Size          string    `json:"size"`
	Stock         int64     `json:"stock"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedAt     time.Time `json:"created_at"`
}
