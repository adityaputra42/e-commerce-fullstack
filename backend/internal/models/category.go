package models

import (
	"time"

	"gorm.io/gorm"
)

// Category Model
type Category struct {
	ID        int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string         `json:"name" validate:"required,min=3,max=100" gorm:"type:varchar(100);unique;not null"`
	Icon      string         `json:"icon" validate:"omitempty,url" gorm:"type:text"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type CategoriesParam struct {
	Name string `form:"name" json:"name" validate:"required,min=3,max=100"`
	Icon string `form:"icon" json:"icon" validate:"omitempty,url"`
}

type UpdateCategory struct {
	ID   int64  `form:"id" json:"id" validate:"required,gt=0"`
	Name string `form:"name" json:"name" validate:"omitempty,min=3,max=100"`
	Icon string `form:"icon" json:"icon" validate:"omitempty,url"`
}

type CategoryResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Icon      string    `json:"icon"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type CategoryListRequest struct {
	Limit  int
	Page   int
	SortBy string
	Search string
}
