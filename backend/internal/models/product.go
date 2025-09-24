package models

import (
	"mime/multipart"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	CategoryID  int64          `json:"category_id" validate:"required" gorm:"not null"`
	Name        string         `json:"name" validate:"required,min=3,max=100" gorm:"type:varchar(100);not null"`
	Description string         `json:"description" validate:"max=255" gorm:"type:varchar(255)"`
	Images      string         `json:"images" validate:"omitempty,url" gorm:"type:text"`
	Rating      float64        `json:"rating" validate:"gte=0,lte=5" gorm:"type:decimal(2,1);default:0"`
	Price       float64        `json:"price" validate:"required,gt=0" gorm:"type:decimal(10,2);not null"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	ColorVarians []ColorVarian `json:"color_varians" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
type ColorVarian struct {
	ID        int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID int64          `json:"product_id" validate:"required" gorm:"not null"`
	Name      string         `json:"name" validate:"required,min=2,max=50" gorm:"type:varchar(50);not null"`
	Color     string         `json:"color" validate:"required" gorm:"type:varchar(20);not null"`
	Images    string         `json:"images" validate:"omitempty,url" gorm:"type:text"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	SizeVarians []SizeVarian `json:"size_varians" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type SizeVarian struct {
	ID            int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	ColorVarianID int64          `json:"color_varian_id" validate:"required" gorm:"not null"`
	Size          string         `json:"size" validate:"required" gorm:"type:varchar(10);not null"`
	Stock         int64          `json:"stock" validate:"gte=0" gorm:"default:0"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
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
	ID          int64            `json:"id"`
	Category    CategoryResponse `json:"category"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Images      string           `json:"images"`
	Rating      float64          `json:"rating"`
	Price       float64          `json:"price"`
	UpdatedAt   time.Time        `json:"updated_at"`
	CreatedAt   time.Time        `json:"created_at"`
}

type ProductDetailResponse struct {
	ID          int64                 `json:"id"`
	Category    CategoryResponse      `json:"category"`
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
type ProductListRequest struct {
	Limit      int
	Offset     int
	SortBy     string
	Search     string
	CategoryID int64 
}

type ColorVarianListRequest struct {
	Limit     int
	Offset    int
	SortBy    string
	Search    string
	ProductID int64 
}

type SizeVarianListRequest struct {
	Limit         int
	Offset        int
	SortBy        string
	Search        string
	ColorVarianID int64 
}

func (Product) TableName() string {
	return "products"
}

func (SizeVarianResponse) TableName() string {
	return "size_varian"
}

func (ColorVarian) TableName() string {
	return "color_varian"
}
