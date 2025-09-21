package repository

import (
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(param models.Product, tx *gorm.DB) (models.Product, error)
	UpdateProduct(param models.Product, tx *gorm.DB) (models.Product, error)
	DeleteProduct(param models.Product) error
	FindProductById(paramId uint) (models.Product, error)
	FindAllProduct(param models.ProductListRequest) ([]models.Product, error)
	CreateColorVarian(param models.ColorVarian, tx *gorm.DB) (models.ColorVarian, error)
	UpdateColorVarian(param models.ColorVarian, tx *gorm.DB) (models.ColorVarian, error)
	DeleteColorVarian(param models.ColorVarian) error
	FindColorVarianById(paramId uint) (models.ColorVarian, error)
	FindAllColorVarian(param models.ColorVarianListRequest) ([]models.ColorVarian, error)
	CreateSizeVarian(param models.SizeVarian, tx *gorm.DB) (models.SizeVarian, error)
	UpdateSizeVarian(param models.SizeVarian, tx *gorm.DB) (models.SizeVarian, error)
	DeleteSizeVarian(param models.SizeVarian) error
	FindSizeVarianById(paramId uint) (models.SizeVarian, error)
	FindAllSizeVarian(param models.SizeVarianListRequest) ([]models.SizeVarian, error)
}

type ProductRepositoryImpl struct {
}

// FindAllProduct implements ProductRepository.
func (a *ProductRepositoryImpl) FindAllProduct(param models.ProductListRequest) ([]models.Product, error) {
	panic("unimplemented")
}

// CreateColorVarian implements ProductRepository.
func (a *ProductRepositoryImpl) CreateColorVarian(param models.ColorVarian, tx *gorm.DB) (models.ColorVarian, error) {
	var result models.ColorVarian

	db := database.DB
	if tx != nil {
		db = tx
	}

	err := db.Create(&param).Error
	if err != nil {
		return result, err
	}

	err = db.First(&result, param.ID).Error
	return result, err
}

// CreateSizeVarian implements ProductRepository.
func (a *ProductRepositoryImpl) CreateSizeVarian(param models.SizeVarian, tx *gorm.DB) (models.SizeVarian, error) {
	var result models.SizeVarian

	db := database.DB
	if tx != nil {
		db = tx
	}

	err := db.Create(&param).Error
	if err != nil {
		return result, err
	}

	err = db.First(&result, param.ID).Error
	return result, err
}

// DeleteColorVarian implements ProductRepository.
func (a *ProductRepositoryImpl) DeleteColorVarian(param models.ColorVarian) error {
	return database.DB.Delete(&param).Error
}

// DeleteSizeVarian implements ProductRepository.
func (a *ProductRepositoryImpl) DeleteSizeVarian(param models.SizeVarian) error {
	return database.DB.Delete(&param).Error
}

// FindAllColorVarian implements ProductRepository.
func (a *ProductRepositoryImpl) FindAllColorVarian(param models.ColorVarianListRequest) ([]models.ColorVarian, error) {
	var colorVarian []models.ColorVarian
	db := database.DB

	if param.SortBy != "" {
		db = db.Order(param.SortBy)
	}

	if param.Limit > 0 {
		db = db.Limit(param.Limit)
	}

	if param.Offset > 0 {
		db = db.Offset(param.Offset)
	}

	if err := db.Preload("Products").Find(&colorVarian).Error; err != nil {
		return nil, err
	}

	return colorVarian, nil
}

// FindAllSizeVarian implements ProductRepository.
func (a *ProductRepositoryImpl) FindAllSizeVarian(param models.SizeVarianListRequest) ([]models.SizeVarian, error) {
	var sizeVarian []models.SizeVarian
	db := database.DB

	if param.SortBy != "" {
		db = db.Order(param.SortBy)
	}

	if param.Limit > 0 {
		db = db.Limit(param.Limit)
	}

	if param.Offset > 0 {
		db = db.Offset(param.Offset)
	}

	if err := db.Preload("Products").Find(&sizeVarian).Error; err != nil {
		return nil, err
	}

	return sizeVarian, nil
}

// FindColorVarianById implements ProductRepository.
func (a *ProductRepositoryImpl) FindColorVarianById(paramId uint) (models.ColorVarian, error) {
	ColorVarian := models.ColorVarian{}
	err := database.DB.Model(&models.ColorVarian{}).Take(&ColorVarian, "id =?", paramId).Error

	return ColorVarian, err
}

// FindSizeVarianById implements ProductRepository.
func (a *ProductRepositoryImpl) FindSizeVarianById(paramId uint) (models.SizeVarian, error) {
	sizeMovie := models.SizeVarian{}
	err := database.DB.Model(&models.SizeVarian{}).Take(&sizeMovie, "id =?", paramId).Error

	return sizeMovie, err
}

// UpdateColorVarian implements ProductRepository.
func (a *ProductRepositoryImpl) UpdateColorVarian(param models.ColorVarian, tx *gorm.DB) (models.ColorVarian, error) {
	var result models.ColorVarian

	db := database.DB
	if tx != nil {
		db = tx
	}

	err := db.Model(&param).Updates(param).Error
	if err != nil {
		return result, err
	}

	err = db.First(&result, param.ID).Error
	return result, err

}

// UpdateSizeVarian implements ProductRepository.
func (a *ProductRepositoryImpl) UpdateSizeVarian(param models.SizeVarian, tx *gorm.DB) (models.SizeVarian, error) {
	var result models.SizeVarian

	db := database.DB
	if tx != nil {
		db = tx
	}

	err := db.Model(&param).Updates(param).Error
	if err != nil {
		return result, err
	}

	err = db.First(&result, param.ID).Error
	return result, err

}

// Create implements ProductRepository.
func (a *ProductRepositoryImpl) CreateProduct(param models.Product, tx *gorm.DB) (models.Product, error) {
	var result models.Product

	db := database.DB
	if tx != nil {
		db = tx
	}

	err := db.Create(&param).Error
	if err != nil {
		return result, err
	}

	err = db.First(&result, param.ID).Error
	return result, err
}

// Delete implements ProductRepository.
func (a *ProductRepositoryImpl) DeleteProduct(param models.Product) error {
	return database.DB.Delete(&param).Error
}

// FindAll implements ProductRepository.
func (a *ProductRepositoryImpl) FindProductAll(param models.ProductListRequest) ([]models.Product, error) {
	var Products []models.Product
	db := database.DB

	if param.SortBy != "" {
		db = db.Order(param.SortBy)
	}

	if param.Limit > 0 {
		db = db.Limit(param.Limit)
	}

	if param.Offset > 0 {
		db = db.Offset(param.Offset)
	}

	if err := db.Preload("Products").Find(&Products).Error; err != nil {
		return nil, err
	}

	return Products, nil
}

// FindById implements ProductRepository.
func (a *ProductRepositoryImpl) FindProductById(paramId uint) (models.Product, error) {
	Product := models.Product{}
	err := database.DB.Model(&models.User{}).Take(&Product, "id =?", paramId).Error

	return Product, err
}

// Update implements ProductRepository.
func (a *ProductRepositoryImpl) UpdateProduct(param models.Product, tx *gorm.DB) (models.Product, error) {
	var result models.Product

	db := database.DB
	if tx != nil {
		db = tx
	}

	err := db.Model(&param).Updates(param).Error
	if err != nil {
		return result, err
	}

	err = db.First(&result, param.ID).Error
	return result, err

}

func NewProductRepository() ProductRepository {
	return &ProductRepositoryImpl{}
}
