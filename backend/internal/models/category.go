package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        int64          `json:"id"`
	Name      string         `json:"name"`
	Icon      string         `json:"icon"`
	UpdatedAt time.Time      `json:"updated_at"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type CategoriesParam struct {
	Name string `form:"name"`
	Icon string `form:"icon"`
}

type UpdateCategory struct {
	ID   int64  `form:"id"`
	Name string `form:"name"`
	Icon string `form:"icon"`
}

type CategoryResponnse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Icon      string    `json:"icon"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
