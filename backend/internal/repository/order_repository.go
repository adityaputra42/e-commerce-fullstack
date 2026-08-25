package repository

import (
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderRepository interface {
	Create(param models.Order, tx *gorm.DB) (models.Order, error)
	Update(param models.Order, tx *gorm.DB) (models.Order, error)
	Delete(param models.Order) error
	FindById(paramId string) (models.Order, error)
	// FindByIdLocking mengambil order dengan row lock (SELECT ... FOR UPDATE)
	// di dalam transaction tx. Dipakai saat mengubah status order (cancel,
	// update) bersamaan dengan restock stock, supaya tidak race dengan
	// checkout lain yang menahan lock size variant yang sama.
	FindByIdLocking(tx *gorm.DB, paramId string) (models.Order, error)
	FindAll(param models.OrderListRequest) ([]models.Order, error)
	FindAllByTxId(txId string) ([]models.Order, error)
	// FindAllByTxIdLocking sama seperti FindAllByTxId, tapi membaca lewat tx
	// yang diberikan (bukan database.DB global) dengan row lock. WAJIB
	// dipakai di alur orkestrasi lintas-entity (lihat
	// OrderFulfillmentService) supaya pembacaan order tetap berada di
	// koneksi/transaction yang sama dengan penulisannya — membaca lewat
	// database.DB terpisah di tengah transaction yang belum commit berisiko
	// baca data basi (koneksi berbeda) dan tidak melindungi dari race saat
	// restock.
	FindAllByTxIdLocking(tx *gorm.DB, txId string) ([]models.Order, error)
}

type OrderRepositoryImpl struct {
}

// FindByIdLocking implements OrderRepository.
func (a *OrderRepositoryImpl) FindByIdLocking(tx *gorm.DB, paramId string) (models.Order, error) {
	var order models.Order
	err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", paramId).
		First(&order).Error

	return order, err
}

// FindAllByTxId implements OrderRepository.
func (a *OrderRepositoryImpl) FindAllByTxId(txId string) ([]models.Order, error) {
	var Orders []models.Order
	db := database.DB

	if err := db.
		Select("id", "transaction_id", "product_id", "color_varian_id", "size_varian_id", "quantity", "unit_price", "subtotal", "created_at", "updated_at").
		Preload("Product", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "category_id", "name", "description", "created_at", "updated_at")
		}).
		Preload("ColorVarian", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "product_id", "name", "color", "images", "created_at", "updated_at")
		}).
		Preload("SizeVarian", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "color_varian_id", "size", "stock", "created_at", "updated_at")
		}).
		Where("transaction_id = ?", txId).
		Find(&Orders).Error; err != nil {
		return nil, err
	}

	return Orders, nil
}

// FindAllByTxIdLocking implements OrderRepository. Sama seperti
// FindAllByTxId, tapi membaca lewat tx yang diberikan (bukan database.DB
// global) dan mengunci baris (SELECT ... FOR UPDATE) — lihat komentar di
// interface untuk alasannya. Tidak perlu preload relasi (Product/
// ColorVarian/SizeVarian) di sini karena dipakai untuk orkestrasi status,
// bukan untuk ditampilkan ke user.
func (a *OrderRepositoryImpl) FindAllByTxIdLocking(tx *gorm.DB, txId string) ([]models.Order, error) {
	var orders []models.Order
	err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("transaction_id = ?", txId).
		Find(&orders).Error

	return orders, err
}

// Create implements OrderRepository.
func (a *OrderRepositoryImpl) Create(param models.Order, tx *gorm.DB) (models.Order, error) {
	var result models.Order

	db := database.DB
	if tx != nil {
		db = tx
	}

	err := db.Create(&param).Error
	if err != nil {
		return result, err
	}

	err = db.
		Select("id", "transaction_id", "product_id", "color_varian_id", "size_varian_id", "quantity", "unit_price", "subtotal", "created_at", "updated_at").
		First(&result, "id = ?", param.ID).Error

	return result, err
}

// Delete implements OrderRepository.
func (a *OrderRepositoryImpl) Delete(param models.Order) error {
	return database.DB.Delete(&param).Error
}

// FindAll implements OrderRepository.
func (a *OrderRepositoryImpl) FindAll(param models.OrderListRequest) ([]models.Order, error) {
	offset := (param.Page - 1) * param.Limit

	var Orders []models.Order
	db := database.DB.
		Select("id", "transaction_id", "product_id", "color_varian_id", "size_varian_id", "quantity", "unit_price", "subtotal", "created_at", "updated_at")

	if param.SortBy != "" {
		db = db.Order(param.SortBy)
	}

	if param.Limit > 0 {
		db = db.Limit(int(param.Limit))
	}

	if offset > 0 {
		db = db.Offset(int(offset))
	}

	if err := db.
		Preload("Product", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "category_id", "name", "description", "created_at", "updated_at")
		}).
		Preload("ColorVarian", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "product_id", "name", "color", "images", "created_at", "updated_at")
		}).
		Preload("SizeVarian", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "color_varian_id", "size", "stock", "created_at", "updated_at")
		}).
		Find(&Orders).Error; err != nil {
		return nil, err
	}

	return Orders, nil
}

// FindById implements OrderRepository.
func (a *OrderRepositoryImpl) FindById(paramId string) (models.Order, error) {
	Order := models.Order{}
	err := database.DB.
		Select("id", "transaction_id", "product_id", "color_varian_id", "size_varian_id", "quantity", "unit_price", "subtotal", "created_at", "updated_at").
		Preload("Product", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "category_id", "name", "description", "created_at", "updated_at")
		}).
		Preload("ColorVarian", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "product_id", "name", "color", "images", "created_at", "updated_at")
		}).
		Preload("SizeVarian", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "color_varian_id", "size", "stock", "created_at", "updated_at")
		}).
		First(&Order, "id = ?", paramId).Error

	return Order, err
}

// Update implements OrderRepository.
func (a *OrderRepositoryImpl) Update(param models.Order, tx *gorm.DB) (models.Order, error) {
	var result models.Order

	db := database.DB
	if tx != nil {
		db = tx
	}

	err := db.Model(&param).Updates(param).Error
	if err != nil {
		return result, err
	}

	err = db.
		Select("id", "transaction_id", "product_id", "color_varian_id", "size_varian_id", "quantity", "unit_price", "subtotal", "created_at", "updated_at").
		First(&result, "id = ?", param.ID).Error

	return result, err
}

func NewOrderRepository() OrderRepository {
	return &OrderRepositoryImpl{}
}
