package repository

import (
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	// Product
	CreateProduct(param models.Product, tx *gorm.DB) (models.Product, error)
	UpdateProduct(param models.Product, tx *gorm.DB) (models.Product, error)
	DeleteProduct(paramid int64, tx *gorm.DB) error
	FindProductById(id int64) (*models.Product, error)
	FindAllProduct(param models.ProductListRequest) ([]models.Product, int64, error)

	// ColorVarian
	CreateColorVarian(param models.ColorVarian, tx *gorm.DB) (models.ColorVarian, error)
	UpdateColorVarian(param models.ColorVarian, tx *gorm.DB) (models.ColorVarian, error)
	DeleteColorVarian(param int64, tx *gorm.DB) error
	FindColorVarianById(id int64) (models.ColorVarian, error)
	FindAllColorVarian(param models.ColorVarianListRequest) ([]models.ColorVarian, error)
	FindColorVarianByProductId(productId int64) ([]models.ColorVarian, error)

	// SizeVarian
	CreateSizeVarian(param models.SizeVarian, tx *gorm.DB) (models.SizeVarian, error)
	UpdateSizeVarian(param models.SizeVarian, tx *gorm.DB) (models.SizeVarian, error)
	DeleteSizeVarian(param int64, tx *gorm.DB) error
	FindSizeVarianById(id int64) (models.SizeVarian, error)
	FindAllSizeVarian(param models.SizeVarianListRequest) ([]models.SizeVarian, error)
	FindSizeVarianByColorVarianId(colorVarianId int64) ([]models.SizeVarian, error)
	FindByNameAndCategory(name string, categoryID int64, tx *gorm.DB) (*models.Product, error)
}

type ProductRepositoryImpl struct{}

func getDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return database.DB
}

func (r *ProductRepositoryImpl) FindByNameAndCategory(name string, categoryID int64, tx *gorm.DB) (*models.Product, error) {
	var product models.Product

	db := database.DB
	if tx != nil {
		db = tx
	}

	err := db.Where("name = ? AND category_id = ? AND deleted_at IS NULL", name, categoryID).
		First(&product).Error

	if err != nil {
		return nil, err
	}

	return &product, nil
}
func (r *ProductRepositoryImpl) CreateProduct(param models.Product, tx *gorm.DB) (models.Product, error) {
	db := getDB(tx)
	if err := db.Create(&param).Error; err != nil {
		return models.Product{}, err
	}
	return param, nil
}

func (r *ProductRepositoryImpl) UpdateProduct(param models.Product, tx *gorm.DB) (models.Product, error) {
	db := getDB(tx)
	if err := db.Save(&param).Error; err != nil {
		return models.Product{}, err
	}
	return param, nil
}

func (r *ProductRepositoryImpl) DeleteProduct(param int64, tx *gorm.DB) error {
	db := database.DB
	if tx != nil {
		db = tx
	}

	return db.Delete(&models.Product{}, param).Error
}

func (r *ProductRepositoryImpl) FindProductById(id int64) (*models.Product, error) {
	var product models.Product

	err := database.DB.
		Preload("ColorVarians", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL").Order("created_at ASC")
		}).
		Preload("ColorVarians.SizeVarians", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL").Order("size ASC")
		}).
		Where("deleted_at IS NULL").
		First(&product, id).Error

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepositoryImpl) FindAllProduct(param models.ProductListRequest) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	// Base query dengan filter deleted_at
	query := database.DB.Model(&models.Product{}).Where("deleted_at IS NULL")

	// Filter by category
	if param.CategoryID != 0 {
		query = query.Where("category_id = ?", param.CategoryID)
	}

	// Search by name atau description (case-insensitive)
	if param.Search != "" {
		searchQuery := "%" + param.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchQuery, searchQuery)
	}

	// Count total sebelum limit/offset
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	sortBy := "created_at DESC"
	if param.SortBy != "" {
		sortBy = param.SortBy
	}
	query = query.Order(sortBy)

	// Pagination
	if param.Limit > 0 {
		query = query.Limit(param.Limit)
	}
	if param.Offset > 0 {
		query = query.Offset(param.Offset)
	}

	// Execute query (tanpa preload Category)
	if err := query.Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *ProductRepositoryImpl) CreateColorVarian(param models.ColorVarian, tx *gorm.DB) (models.ColorVarian, error) {
	db := getDB(tx)
	if err := db.Create(&param).Error; err != nil {
		return models.ColorVarian{}, err
	}
	return param, nil
}

func (r *ProductRepositoryImpl) UpdateColorVarian(param models.ColorVarian, tx *gorm.DB) (models.ColorVarian, error) {
	db := getDB(tx)
	if err := db.Save(&param).Error; err != nil {
		return models.ColorVarian{}, err
	}
	return param, nil
}

func (r *ProductRepositoryImpl) DeleteColorVarian(param int64, tx *gorm.DB) error {

	db := database.DB
	if tx != nil {
		db = tx
	}
	return db.Delete(&models.ColorVarian{}, param).Error
}

func (r *ProductRepositoryImpl) FindColorVarianById(id int64) (models.ColorVarian, error) {
	var cv models.ColorVarian
	err := database.DB.Preload("SizeVarians", "deleted_at IS NULL").First(&cv, "id = ? AND deleted_at IS NULL", id).Error
	return cv, err
}

func (r *ProductRepositoryImpl) FindAllColorVarian(param models.ColorVarianListRequest) ([]models.ColorVarian, error) {
	var list []models.ColorVarian
	db := database.DB.Model(&models.ColorVarian{})

	if param.ProductID != 0 {
		db = db.Where("product_id = ?", param.ProductID)
	}

	if param.Search != "" {
		q := "%" + param.Search + "%"
		db = db.Where("name LIKE ? OR color LIKE ?", q, q)
	}

	if param.SortBy != "" {
		db = db.Order(param.SortBy)
	}
	if param.Limit > 0 {
		db = db.Limit(param.Limit)
	}
	if param.Offset > 0 {
		db = db.Offset(param.Offset)
	}

	if err := db.Preload("SizeVarians", "deleted_at IS NULL").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ProductRepositoryImpl) FindColorVarianByProductId(productId int64) ([]models.ColorVarian, error) {
	req := models.ColorVarianListRequest{ProductID: productId}
	return r.FindAllColorVarian(req)
}

// -------------------- SizeVarian --------------------

func (r *ProductRepositoryImpl) CreateSizeVarian(param models.SizeVarian, tx *gorm.DB) (models.SizeVarian, error) {
	db := getDB(tx)
	if err := db.Create(&param).Error; err != nil {
		return models.SizeVarian{}, err
	}
	return param, nil
}

func (r *ProductRepositoryImpl) UpdateSizeVarian(param models.SizeVarian, tx *gorm.DB) (models.SizeVarian, error) {
	db := getDB(tx)
	if err := db.Save(&param).Error; err != nil {
		return models.SizeVarian{}, err
	}
	return param, nil
}

func (r *ProductRepositoryImpl) DeleteSizeVarian(param int64, tx *gorm.DB) error {
	db := database.DB
	if tx != nil {
		db = tx
	}
	return db.Delete(&models.SizeVarian{}, param).Error
}

func (r *ProductRepositoryImpl) FindSizeVarianById(id int64) (models.SizeVarian, error) {
	var sv models.SizeVarian
	err := database.DB.First(&sv, "id = ? AND deleted_at IS NULL", id).Error
	return sv, err
}

func (r *ProductRepositoryImpl) FindAllSizeVarian(param models.SizeVarianListRequest) ([]models.SizeVarian, error) {

	offset := (param.Page - 1) * param.Limit

	var list []models.SizeVarian
	db := database.DB.Model(&models.SizeVarian{})

	if param.ColorVarianID != 0 {
		db = db.Where("color_varian_id = ?", param.ColorVarianID)
	}

	if param.Search != "" {
		q := "%" + param.Search + "%"
		db = db.Where("size LIKE ?", q)
	}

	if param.SortBy != "" {
		db = db.Order(param.SortBy)
	}
	if param.Limit > 0 {
		db = db.Limit(param.Limit)
	}
	if offset > 0 {
		db = db.Offset(offset)
	}

	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ProductRepositoryImpl) FindSizeVarianByColorVarianId(colorVarianId int64) ([]models.SizeVarian, error) {
	req := models.SizeVarianListRequest{ColorVarianID: colorVarianId}
	return r.FindAllSizeVarian(req)
}

func NewProductRepository() ProductRepository {
	return &ProductRepositoryImpl{}
}
