package services

import (
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/repository"
	"e-commerce/backend/internal/utils"
	"fmt"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

type ProductService interface {
	CreateProduct(param models.CreateProductParam) (*models.ProductDetailResponse, error)
	FindProductById(id int64) (*models.ProductDetailResponse, error)
	FindAllProduct(param models.ProductListRequest) (*[]models.ProductResponse, error)
	UpdateProduct(param models.UpdateProductParam) (*models.ProductDetailResponse, error)
	DeleteProduct(id int64) error
	AddColorVarianProduct(productId int64, param models.CreateColorVarianRequest) (*models.ProductDetailResponse, error)
	updateSizeVariants(colorVarianID int64, sizesParam []models.UpdateSizeVarianRequest, tx *gorm.DB) error
}

type ProductServiceImpl struct {
	categoryRepo repository.CategoryRepository
	productRepo  repository.ProductRepository
}

// updateSizeVariants implements ProductService.
func (p *ProductServiceImpl) updateSizeVariants(colorVarianID int64, sizesParam []models.UpdateSizeVarianRequest, tx *gorm.DB) error {

	var existingSizes []models.SizeVarian
	err := tx.Where("color_varian_id = ? AND deleted_at IS NULL", colorVarianID).
		Find(&existingSizes).Error
	if err != nil {
		return fmt.Errorf("error mengambil size variants: %w", err)
	}

	existingSizeMap := make(map[int64]*models.SizeVarian)
	for i := range existingSizes {
		existingSizeMap[existingSizes[i].ID] = &existingSizes[i]
	}

	updatedSizeIDs := make(map[int64]bool)

	for j, sizeParam := range sizesParam {
		if sizeParam.ID != nil {
			existingSize, exists := existingSizeMap[*sizeParam.ID]
			if !exists {
				return fmt.Errorf("size variant dengan ID %d tidak ditemukan", *sizeParam.ID)
			}

			if sizeParam.Size != nil {
				existingSize.Size = *sizeParam.Size
			}
			if sizeParam.Stock != nil {
				existingSize.Stock = *sizeParam.Stock
			}

			_, err := p.productRepo.UpdateSizeVarian(*existingSize, tx)
			if err != nil {
				return fmt.Errorf("gagal mengupdate size variant ke-%d: %w", j+1, err)
			}

			updatedSizeIDs[*sizeParam.ID] = true

		} else {
			if sizeParam.Size == nil || sizeParam.Stock == nil {
				return fmt.Errorf("size dan stock wajib diisi untuk size variant baru")
			}

			newSizeVarian := models.SizeVarian{
				ColorVarianID: colorVarianID,
				Size:          *sizeParam.Size,
				Stock:         *sizeParam.Stock,
			}

			createdSize, err := p.productRepo.CreateSizeVarian(newSizeVarian, tx)
			if err != nil {
				return fmt.Errorf("gagal membuat size variant ke-%d: %w", j+1, err)
			}

			updatedSizeIDs[createdSize.ID] = true
		}
	}

	for id := range existingSizeMap {
		if !updatedSizeIDs[id] {
			if err := p.productRepo.DeleteSizeVarian(id, tx); err != nil {
				return fmt.Errorf("gagal menghapus size variant: %w", err)
			}
		}
	}

	return nil
}

func sanitizeFileName(name string) string {
	name = strings.ReplaceAll(name, " ", "_")
	reg := regexp.MustCompile("[^a-zA-Z0-9_-]+")
	return reg.ReplaceAllString(name, "")
}
func (p *ProductServiceImpl) AddColorVarianProduct(productId int64, param models.CreateColorVarianRequest) (*models.ProductDetailResponse, error) {
	var product models.Product
	var category models.Category

	err := database.DB.Transaction(func(tx *gorm.DB) error {

		productResult, err := p.productRepo.FindProductById(productId)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("product dengan ID %d tidak ditemukan", productId)
			}
			return fmt.Errorf("error mencari product: %w", err)
		}
		product = *productResult

		categoryResult, err := p.categoryRepo.FindById((product.CategoryID))
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("category dengan ID %d tidak ditemukan", product.CategoryID)
			}
			return fmt.Errorf("error mengambil category: %w", err)
		}
		category = categoryResult

		for _, cv := range product.ColorVarians {
			if strings.EqualFold(cv.Name, param.Name) {
				return fmt.Errorf("color varian '%s' sudah ada di product ini", param.Name)
			}
		}
		if param.Image == nil {
			return fmt.Errorf("gambar wajib diisi untuk color varian")
		}

		if len(param.Sizes) == 0 {
			return fmt.Errorf("minimal harus ada 1 ukuran untuk color varian")
		}

		sizeMap := make(map[string]bool)
		for _, sv := range param.Sizes {
			if sizeMap[strings.ToUpper(sv.Size)] {
				return fmt.Errorf("ukuran '%s' duplikat dalam color varian '%s'", sv.Size, param.Name)
			}
			sizeMap[strings.ToUpper(sv.Size)] = true
		}

		folderName := fmt.Sprintf("product/%s/colors", sanitizeFileName(product.Name))
		colorImageURL, err := utils.UploadToSupabase(param.Image, folderName)
		if err != nil {
			return fmt.Errorf("gagal upload gambar color varian: %w", err)
		}

		colorVariant := models.ColorVarian{
			ProductID: product.ID,
			Name:      param.Name,
			Color:     param.Color,
			Images:    colorImageURL,
		}

		createdColorVariant, err := p.productRepo.CreateColorVarian(colorVariant, tx)
		if err != nil {
			return fmt.Errorf("gagal membuat color varian: %w", err)
		}

		for i, sizeParam := range param.Sizes {
			sizeVariant := models.SizeVarian{
				ColorVarianID: createdColorVariant.ID,
				Size:          sizeParam.Size,
				Stock:         sizeParam.Stock,
			}

			_, err := p.productRepo.CreateSizeVarian(sizeVariant, tx)
			if err != nil {
				return fmt.Errorf("gagal membuat size variant ke-%d: %w", i+1, err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	productWithRelations, err := p.productRepo.FindProductById(productId)
	if err != nil {
		return nil, fmt.Errorf("gagal memuat data product: %w", err)
	}

	result := productWithRelations.ToProductDetailResponse(category)
	return &result, nil
}

// CreateProduct implements ProductService.
func (p *ProductServiceImpl) CreateProduct(param models.CreateProductParam) (*models.ProductDetailResponse, error) {
	var product models.Product
	var category models.Category
	var colorVariants []models.ColorVarian

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		categoryResult, err := p.categoryRepo.FindById(param.CategoryID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("category dengan ID %d tidak ditemukan", param.CategoryID)
			}
			return fmt.Errorf("error validasi category: %w", err)
		}
		category = categoryResult

		var mainImageURL string
		if param.Image != nil {
			folderName := fmt.Sprintf("product/%s", sanitizeFileName(param.Name))
			uploadedURL, err := utils.UploadToSupabase(param.Image, folderName)
			if err != nil {
				return fmt.Errorf("gagal upload gambar produk: %w", err)
			}
			mainImageURL = uploadedURL
		}

		productParam := models.Product{
			CategoryID:  param.CategoryID,
			Name:        param.Name,
			Description: param.Description,
			Price:       param.Price,
			Images:      mainImageURL,
		}

		product, err = p.productRepo.CreateProduct(productParam, tx)
		if err != nil {
			return fmt.Errorf("gagal membuat produk: %w", err)
		}

		for i, colorParam := range param.ColorVarian {

			folderName := fmt.Sprintf("product/%s/colors", sanitizeFileName(param.Name))
			colorImageURL, err := utils.UploadToSupabase(colorParam.Image, folderName)
			if err != nil {
				return fmt.Errorf("gagal upload gambar varian warna ke-%d: %w", i+1, err)
			}

			colorVariantParam := models.ColorVarian{
				ProductID: product.ID,
				Name:      colorParam.Name,
				Color:     colorParam.Color,
				Images:    colorImageURL,
			}

			colorVariant, err := p.productRepo.CreateColorVarian(colorVariantParam, tx)
			if err != nil {
				return fmt.Errorf("gagal membuat varian warna ke-%d: %w", i+1, err)
			}

			var sizeVariants []models.SizeVarian
			for j, sizeParam := range colorParam.Sizes {
				sizeVariantParam := models.SizeVarian{
					ColorVarianID: colorVariant.ID,
					Size:          sizeParam.Size,
					Stock:         sizeParam.Stock,
				}

				sizeVariant, err := p.productRepo.CreateSizeVarian(sizeVariantParam, tx)
				if err != nil {
					return fmt.Errorf("gagal membuat varian ukuran ke-%d untuk warna '%s': %w", j+1, colorParam.Name, err)
				}

				sizeVariants = append(sizeVariants, sizeVariant)
			}

			colorVariant.SizeVarians = sizeVariants
			colorVariants = append(colorVariants, colorVariant)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	product.ColorVarians = colorVariants
	result := product.ToProductDetailResponse(category)

	return &result, nil
}

// DeleteProduct implements ProductService.
func (p *ProductServiceImpl) DeleteProduct(id int64) error {
	var imageURLs []string

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		product, err := p.productRepo.FindProductById(id)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("produk dengan ID %d tidak ditemukan", id)
			}
			return fmt.Errorf("error mencari produk: %w", err)
		}

		if product.Images != "" {
			imageURLs = append(imageURLs, product.Images)
		}
		for _, cv := range product.ColorVarians {
			if cv.Images != "" {
				imageURLs = append(imageURLs, cv.Images)
			}
		}

		if err := p.productRepo.DeleteProduct(id, tx); err != nil {
			return fmt.Errorf("gagal menghapus produk: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	if len(imageURLs) > 0 {
		go func() {
			if err := utils.DeleteMultipleFromSupabase(imageURLs); err != nil {
				fmt.Printf("Warning: gagal menghapus images: %v\n", err)
			}
		}()
	}

	return nil
}

// FindAllProduct implements ProductService.
func (p *ProductServiceImpl) FindAllProduct(param models.ProductListRequest) (*[]models.ProductResponse, error) {

	products, _, err := p.productRepo.FindAllProduct(param)
	if err != nil {
		return nil, fmt.Errorf("error mengambil data produk: %w", err)
	}

	categoryIDs := make([]int64, 0)
	categoryIDMap := make(map[int64]bool)
	for _, product := range products {
		if !categoryIDMap[product.CategoryID] {
			categoryIDs = append(categoryIDs, product.CategoryID)
			categoryIDMap[product.CategoryID] = true
		}
	}

	categories, err := p.categoryRepo.FindByIds(categoryIDs)
	if err != nil {
		return nil, fmt.Errorf("error mengambil data kategori: %w", err)
	}

	categoryMap := models.BuildCategoryMap(categories)

	response := models.ToProductResponseList(products, categoryMap)

	return &response, nil
}

// FindProductById implements ProductService.
func (p *ProductServiceImpl) FindProductById(id int64) (*models.ProductDetailResponse, error) {
	product, err := p.productRepo.FindProductById(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("produk dengan ID %d tidak ditemukan", id)
		}
		return nil, fmt.Errorf("error mengambil data produk: %w", err)
	}

	category, err := p.categoryRepo.FindById(product.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("error mengambil data kategori: %w", err)
	}

	response := product.ToProductDetailResponse(category)
	return &response, nil
}

// UpdateProduct implements ProductService.
func (p *ProductServiceImpl) UpdateProduct(param models.UpdateProductParam) (*models.ProductDetailResponse, error) {
	var product models.Product
	var category models.Category

	err := database.DB.Transaction(func(tx *gorm.DB) error {

		existingProduct, err := p.productRepo.FindProductById(param.ID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("produk dengan ID %d tidak ditemukan", param.ID)
			}
			return fmt.Errorf("error mencari produk: %w", err)
		}
		product = *existingProduct

		if param.CategoryID != nil {
			categoryResult, err := p.categoryRepo.FindById(*param.CategoryID)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return fmt.Errorf("category dengan ID %d tidak ditemukan", *param.CategoryID)
				}
				return fmt.Errorf("error validasi category: %w", err)
			}
			category = categoryResult
			product.CategoryID = *param.CategoryID
		} else {
			categoryResult, err := p.categoryRepo.FindById(product.CategoryID)
			if err != nil {
				return fmt.Errorf("error mengambil data kategori: %w", err)
			}
			category = categoryResult
		}

		// 3. Validasi duplicate name jika nama diubah
		if param.Name != nil && *param.Name != product.Name {
			duplicate, err := p.productRepo.FindByNameAndCategory(*param.Name, product.CategoryID, tx)
			if err != nil && err != gorm.ErrRecordNotFound {
				return fmt.Errorf("error validasi produk: %w", err)
			}
			if duplicate != nil {
				return fmt.Errorf("produk dengan nama '%s' sudah ada di kategori '%s'", *param.Name, category.Name)
			}
			product.Name = *param.Name
		}

		// 4. Update basic fields
		if param.Description != nil {
			product.Description = *param.Description
		}
		if param.Price != nil {
			product.Price = *param.Price
		}
		if param.Rating != nil {
			product.Rating = *param.Rating
		}

		if param.Image != nil {
			newImageURL, err := utils.ReplaceFile(product.Images, param.Image, fmt.Sprintf("product/%s", sanitizeFileName(product.Name)))
			if err != nil {
				return fmt.Errorf("gagal mengganti gambar produk: %w", err)
			}
			product.Images = newImageURL
		}

		product, err = p.productRepo.UpdateProduct(product, tx)
		if err != nil {
			return fmt.Errorf("gagal mengupdate produk: %w", err)
		}

		if len(param.ColorVarian) > 0 {
			existingColorVariants, err := p.productRepo.FindColorVarianByProductId(product.ID)
			if err != nil {
				return fmt.Errorf("error mengambil color variants: %w", err)
			}

			existingColorMap := make(map[int64]*models.ColorVarian)
			for i := range existingColorVariants {
				existingColorMap[existingColorVariants[i].ID] = &existingColorVariants[i]
			}

			updatedColorIDs := make(map[int64]bool)

			for i, colorParam := range param.ColorVarian {
				if colorParam.ID != nil {
					existingColor, exists := existingColorMap[*colorParam.ID]
					if !exists {
						return fmt.Errorf("color variant dengan ID %d tidak ditemukan", *colorParam.ID)
					}

					if colorParam.Name != nil {
						existingColor.Name = *colorParam.Name
					}
					if colorParam.Color != nil {
						existingColor.Color = *colorParam.Color
					}

					if colorParam.Image != nil {
						folderName := fmt.Sprintf("product/%s/colors", sanitizeFileName(product.Name))
						newImageURL, err := utils.ReplaceFile(existingColor.Images, colorParam.Image, folderName)
						if err != nil {
							return fmt.Errorf("gagal mengganti gambar varian warna: %w", err)
						}
						existingColor.Images = newImageURL
					}

					_, err := p.productRepo.UpdateColorVarian(*existingColor, tx)
					if err != nil {
						return fmt.Errorf("gagal mengupdate varian warna ke-%d: %w", i+1, err)
					}

					updatedColorIDs[*colorParam.ID] = true

					if len(colorParam.Sizes) > 0 {
						err := p.updateSizeVariants(*colorParam.ID, colorParam.Sizes, tx)
						if err != nil {
							return err
						}
					}

				} else {
					if colorParam.Image == nil {
						return fmt.Errorf("gambar wajib diisi untuk varian warna baru")
					}

					folderName := fmt.Sprintf("product/%s/colors", sanitizeFileName(product.Name))
					colorImageURL, err := utils.UploadToSupabase(colorParam.Image, folderName)
					if err != nil {
						return fmt.Errorf("gagal upload gambar varian warna ke-%d: %w", i+1, err)
					}

					newColorVarian := models.ColorVarian{
						ProductID: product.ID,
						Name:      *colorParam.Name,
						Color:     *colorParam.Color,
						Images:    colorImageURL,
					}

					createdColor, err := p.productRepo.CreateColorVarian(newColorVarian, tx)
					if err != nil {
						return fmt.Errorf("gagal membuat varian warna ke-%d: %w", i+1, err)
					}

					for j, sizeParam := range colorParam.Sizes {
						newSizeVarian := models.SizeVarian{
							ColorVarianID: createdColor.ID,
							Size:          *sizeParam.Size,
							Stock:         *sizeParam.Stock,
						}

						_, err := p.productRepo.CreateSizeVarian(newSizeVarian, tx)
						if err != nil {
							return fmt.Errorf("gagal membuat size variant ke-%d untuk warna baru: %w", j+1, err)
						}
					}

					updatedColorIDs[createdColor.ID] = true
				}
			}

			for id := range existingColorMap {
				if !updatedColorIDs[id] {
					if err := p.productRepo.DeleteColorVarian(id, tx); err != nil {
						return fmt.Errorf("gagal menghapus color variant: %w", err)
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	productWithRelations, err := p.productRepo.FindProductById(param.ID)
	if err != nil {
		return nil, fmt.Errorf("gagal memuat data produk: %w", err)
	}

	response := productWithRelations.ToProductDetailResponse(category)
	return &response, nil
}

func NewProductService(categoryRepo repository.CategoryRepository,
	productRepo repository.ProductRepository) ProductService {
	return &ProductServiceImpl{categoryRepo: categoryRepo, productRepo: productRepo}
}
