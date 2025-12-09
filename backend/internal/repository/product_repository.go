package repository

import (
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepository interface {
	// Product
	CreateProduct(param models.Product, tx *gorm.DB) (models.Product, error)
	UpdateProduct(param models.Product, tx *gorm.DB) (models.Product, error)
	DeleteProduct(paramid int64, tx *gorm.DB) error
	FindProductById(id int64, tx *gorm.DB) (*models.Product, error)
	FindAllProduct(param models.ProductListRequest, tx *gorm.DB) ([]models.Product, int64, error)

	// ColorVarian
	CreateColorVarian(param models.ColorVarian, tx *gorm.DB) (models.ColorVarian, error)
	UpdateColorVarian(param models.ColorVarian, tx *gorm.DB) (models.ColorVarian, error)
	DeleteColorVarian(param int64, tx *gorm.DB) error
	FindColorVarianById(id int64, tx *gorm.DB) (models.ColorVarian, error)
	FindAllColorVarian(param models.ColorVarianListRequest, tx *gorm.DB) ([]models.ColorVarian, error)
	FindColorVarianByProductId(productId int64, tx *gorm.DB) ([]models.ColorVarian, error)

	// SizeVarian
	CreateSizeVarian(param models.SizeVarian, tx *gorm.DB) (models.SizeVarian, error)
	UpdateSizeVarian(param models.SizeVarian, tx *gorm.DB) (models.SizeVarian, error)
	DeleteSizeVarian(param int64, tx *gorm.DB) error
	FindSizeVarianById(id int64, tx *gorm.DB) (models.SizeVarian, error)
	FindAllSizeVarian(param models.SizeVarianListRequest, tx *gorm.DB) ([]models.SizeVarian, error)
	FindSizeVarianByColorVarianId(colorVarianId int64, tx *gorm.DB) ([]models.SizeVarian, error)
	FindByNameAndCategory(name string, categoryID int64, tx *gorm.DB) (*models.Product, error)
	FindSizeVarianLocked(tx *gorm.DB, id uint) (*models.SizeVarian, error)
}

type ProductRepositoryImpl struct{}

// FindSizeVarianLocked implements ProductRepository.
func (r *ProductRepositoryImpl) FindSizeVarianLocked(tx *gorm.DB, id uint) (*models.SizeVarian, error) {
	db := getDB(tx)

	var sizeVarian models.SizeVarian

	err := db.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&sizeVarian, id).Error

	if err != nil {
		return nil, err
	}

	return &sizeVarian, nil
}

func getDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return database.DB
}

func (r *ProductRepositoryImpl) FindByNameAndCategory(name string, categoryID int64, tx *gorm.DB) (*models.Product, error) {
	db := getDB(tx)

	var product models.Product

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
	db := getDB(tx)
	return db.Delete(&models.Product{}, param).Error
}

func (r *ProductRepositoryImpl) FindProductById(id int64, tx *gorm.DB) (*models.Product, error) {
	db := getDB(tx)

	var product models.Product

	// Use preloads with the chosen db (tx or global)
	err := db.
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

func (r *ProductRepositoryImpl) FindAllProduct(param models.ProductListRequest, tx *gorm.DB) ([]models.Product, int64, error) {
	db := getDB(tx)

	var products []models.Product
	var total int64

	// Base query with soft delete filter
	query := db.Model(&models.Product{}).Where("deleted_at IS NULL")

	// Filter by category
	if param.CategoryID != 0 {
		query = query.Where("category_id = ?", param.CategoryID)
	}

	// Search by name or description (case-insensitive)
	if param.Search != "" {
		searchQuery := "%" + param.Search + "%"
		// Use ILIKE for Postgres; if using MySQL replace with LIKE and lower comparison
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchQuery, searchQuery)
	}

	// Count total before limit/offset
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
	db := getDB(tx)
	return db.Delete(&models.ColorVarian{}, param).Error
}

func (r *ProductRepositoryImpl) FindColorVarianById(id int64, tx *gorm.DB) (models.ColorVarian, error) {
	db := getDB(tx)
	var cv models.ColorVarian
	err := db.Preload("SizeVarians", "deleted_at IS NULL").First(&cv, "id = ? AND deleted_at IS NULL", id).Error
	return cv, err
}

func (r *ProductRepositoryImpl) FindAllColorVarian(param models.ColorVarianListRequest, tx *gorm.DB) ([]models.ColorVarian, error) {
	db := getDB(tx)

	var list []models.ColorVarian
	q := db.Model(&models.ColorVarian{})

	if param.ProductID != 0 {
		q = q.Where("product_id = ?", param.ProductID)
	}

	if param.Search != "" {
		searchQ := "%" + param.Search + "%"
		q = q.Where("name LIKE ? OR color LIKE ?", searchQ, searchQ)
	}

	if param.SortBy != "" {
		q = q.Order(param.SortBy)
	}
	if param.Limit > 0 {
		q = q.Limit(param.Limit)
	}
	if param.Offset > 0 {
		q = q.Offset(param.Offset)
	}

	if err := q.Preload("SizeVarians", "deleted_at IS NULL").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ProductRepositoryImpl) FindColorVarianByProductId(productId int64, tx *gorm.DB) ([]models.ColorVarian, error) {
	req := models.ColorVarianListRequest{ProductID: productId}
	return r.FindAllColorVarian(req, tx)
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
	db := getDB(tx)
	return db.Delete(&models.SizeVarian{}, param).Error
}

func (r *ProductRepositoryImpl) FindSizeVarianById(id int64, tx *gorm.DB) (models.SizeVarian, error) {
	db := getDB(tx)
	var sv models.SizeVarian
	err := db.First(&sv, "id = ? AND deleted_at IS NULL", id).Error
	return sv, err
}

func (r *ProductRepositoryImpl) FindAllSizeVarian(param models.SizeVarianListRequest, tx *gorm.DB) ([]models.SizeVarian, error) {
	db := getDB(tx)

	offset := (param.Page - 1) * param.Limit

	var list []models.SizeVarian
	q := db.Model(&models.SizeVarian{})

	if param.ColorVarianID != 0 {
		q = q.Where("color_varian_id = ?", param.ColorVarianID)
	}

	if param.Search != "" {
		search := "%" + param.Search + "%"
		q = q.Where("size LIKE ?", search)
	}

	if param.SortBy != "" {
		q = q.Order(param.SortBy)
	}
	if param.Limit > 0 {
		q = q.Limit(param.Limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}

	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ProductRepositoryImpl) FindSizeVarianByColorVarianId(colorVarianId int64, tx *gorm.DB) ([]models.SizeVarian, error) {
	req := models.SizeVarianListRequest{ColorVarianID: colorVarianId}
	return r.FindAllSizeVarian(req, tx)
}

func NewProductRepository() ProductRepository {
	return &ProductRepositoryImpl{}
}
