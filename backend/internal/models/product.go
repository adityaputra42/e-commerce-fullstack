package models

import (
	"mime/multipart"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	CategoryID  int64          `json:"category_id" validate:"required" gorm:"not null;index"`
	Name        string         `json:"name" validate:"required,min=3,max=100" gorm:"type:varchar(100);not null;index"`
	Description string         `json:"description" validate:"max=255" gorm:"type:varchar(255)"`
	Images      string         `json:"images" validate:"omitempty,url" gorm:"type:text"`
	Rating      float64        `json:"rating" validate:"gte=0,lte=5" gorm:"type:decimal(2,1);default:0;index"`
	Price       float64        `json:"price" validate:"required,gt=0" gorm:"type:decimal(10,2);not null;index"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime;index"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	ColorVarians []ColorVarian `json:"color_varians" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
type ColorVarian struct {
	ID        int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID int64          `json:"product_id" validate:"required" gorm:"not null;index"`
	Name      string         `json:"name" validate:"required,min=2,max=50" gorm:"type:varchar(50);not null;index"`
	Color     string         `json:"color" validate:"required" gorm:"type:varchar(20);not null"`
	Images    string         `json:"images" validate:"omitempty,url" gorm:"type:text"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	SizeVarians []SizeVarian `json:"size_varians" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type SizeVarian struct {
	ID            int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	ColorVarianID int64          `json:"color_varian_id" validate:"required" gorm:"not null;index"`
	Size          string         `json:"size" validate:"required" gorm:"type:varchar(10);not null;index"`
	Stock         int64          `json:"stock" validate:"gte=0" gorm:"default:0"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}
type CreateProductParam struct {
	CategoryID  int64                      `json:"category_id" form:"category_id" binding:"required"`
	Name        string                     `json:"name" form:"name" binding:"required,min=3,max=100"`
	Description string                     `json:"description" form:"description"`
	Price       float64                    `json:"price" form:"price" binding:"required,gt=0"`
	Image       *multipart.FileHeader      `json:"image" form:"image"`
	ColorVarian []CreateColorVarianRequest `json:"color_varian" form:"color_varian"`
}

type CreateColorVarianRequest struct {
	Name  string                    `json:"name" form:"name" binding:"required,min=2,max=50"`
	Color string                    `json:"color" form:"color" binding:"required"`
	Image *multipart.FileHeader     `json:"image" form:"image" binding:"required"` // hanya 1 gambar per varian
	Sizes []CreateSizeVarianRequest `json:"sizes" form:"sizes"`
}

type CreateSizeVarianRequest struct {
	Size  string `json:"size" form:"size" binding:"required"`
	Stock int64  `json:"stock" form:"stock" binding:"gte=0"`
}

type UpdateProductParam struct {
	ID          int64                      `json:"id" form:"id" binding:"required"`
	CategoryID  *int64                     `json:"category_id,omitempty" form:"category_id"`
	Name        *string                    `json:"name,omitempty" form:"name"`
	Description *string                    `json:"description,omitempty" form:"description"`
	Price       *float64                   `json:"price,omitempty" form:"price"`
	Rating      *float64                   `json:"rating,omitempty" form:"rating"`
	Image       *multipart.FileHeader      `json:"image,omitempty" form:"image"` // jika ingin ubah gambar utama produk
	ColorVarian []UpdateColorVarianRequest `json:"color_varian,omitempty" form:"color_varian"`
}

type UpdateColorVarianRequest struct {
	ID    *int64                    `json:"id,omitempty" form:"id"` // null berarti varian baru
	Name  *string                   `json:"name,omitempty" form:"name"`
	Color *string                   `json:"color,omitempty" form:"color"`
	Image *multipart.FileHeader     `json:"image,omitempty" form:"image"` // ubah 1 gambar
	Sizes []UpdateSizeVarianRequest `json:"sizes,omitempty" form:"sizes"`
}

type UpdateSizeVarianRequest struct {
	ID    *int64  `json:"id,omitempty" form:"id"` // null berarti size baru
	Size  *string `json:"size,omitempty" form:"size"`
	Stock *int64  `json:"stock,omitempty" form:"stock"`
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
	Page       int
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
	Page          int
	SortBy        string
	Search        string
	ColorVarianID int64
}

func (Product) TableName() string {
	return "products"
}

func (SizeVarian) TableName() string {
	return "size_varians"
}

func (ColorVarian) TableName() string {
	return "color_varians"
}

// ToProductResponse converts Product model to ProductResponse (for list)
func (p *Product) ToProductResponse(category Category) ProductResponse {
	return ProductResponse{
		ID: p.ID,
		Category: CategoryResponse{
			ID:        category.ID,
			Name:      category.Name,
			Icon:      category.Icon,
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
		},
		Name:        p.Name,
		Description: p.Description,
		Images:      p.Images,
		Rating:      p.Rating,
		Price:       p.Price,
		UpdatedAt:   p.UpdatedAt,
		CreatedAt:   p.CreatedAt,
	}
}

// ToProductDetailResponse converts Product model to ProductDetailResponse (with variants)
func (p *Product) ToProductDetailResponse(category Category) ProductDetailResponse {
	// Map color variants
	colorVariants := make([]ColorVarianResponse, len(p.ColorVarians))
	for i, cv := range p.ColorVarians {
		colorVariants[i] = cv.ToColorVarianResponse()
	}

	return ProductDetailResponse{
		ID: p.ID,
		Category: CategoryResponse{
			ID:        category.ID,
			Name:      category.Name,
			Icon:      category.Icon,
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
		},
		Name:        p.Name,
		Description: p.Description,
		Images:      p.Images,
		Rating:      p.Rating,
		Price:       p.Price,
		ColorVarian: colorVariants,
		UpdatedAt:   p.UpdatedAt,
		CreatedAt:   p.CreatedAt,
	}
}

func (cv *ColorVarian) ToColorVarianResponse() ColorVarianResponse {

	sizeVariants := make([]SizeVarianResponse, len(cv.SizeVarians))
	for i, sv := range cv.SizeVarians {
		sizeVariants[i] = sv.ToSizeVarianResponse()
	}

	return ColorVarianResponse{
		ID:         cv.ID,
		ProductID:  cv.ProductID,
		Name:       cv.Name,
		Color:      cv.Color,
		Images:     cv.Images,
		SizeVarian: sizeVariants,
		UpdatedAt:  cv.UpdatedAt,
		CreatedAt:  cv.CreatedAt,
	}
}

// ToSizeVarianResponse converts SizeVarian model to SizeVarianResponse
func (sv *SizeVarian) ToSizeVarianResponse() SizeVarianResponse {
	return SizeVarianResponse{
		ID:            sv.ID,
		ColorVarianID: sv.ColorVarianID,
		Size:          sv.Size,
		Stock:         sv.Stock,
		UpdatedAt:     sv.UpdatedAt,
		CreatedAt:     sv.CreatedAt,
	}
}

// ToProductResponseList converts slice of Product to slice of ProductResponse
func ToProductResponseList(products []Product, categoryMap map[int64]Category) []ProductResponse {
	result := make([]ProductResponse, len(products))
	for i, product := range products {
		category := categoryMap[product.CategoryID]
		result[i] = product.ToProductResponse(category)
	}
	return result
}

func ToProductDetailResponseList(products []Product, categoryMap map[int64]Category) []ProductDetailResponse {
	result := make([]ProductDetailResponse, len(products))
	for i, product := range products {
		category := categoryMap[product.CategoryID]
		result[i] = product.ToProductDetailResponse(category)
	}
	return result
}
func BuildCategoryMap(categories []Category) map[int64]Category {
	categoryMap := make(map[int64]Category)
	for _, category := range categories {
		categoryMap[category.ID] = category
	}
	return categoryMap
}
