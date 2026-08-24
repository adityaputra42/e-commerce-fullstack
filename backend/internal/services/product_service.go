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
	UpdateSizeVariants(colorVarianID int64, sizesParam []models.UpdateSizeVarianRequest, tx *gorm.DB) error
}

type ProductServiceImpl struct {
	categoryRepo repository.CategoryRepository
	productRepo  repository.ProductRepository
}

// imageChange melacak satu file gambar yang diupload selama sebuah operasi.
// OldURL diisi kalau operasi ini MENGGANTI gambar yang sudah ada (kosong
// kalau ini gambar baru). Dipakai untuk membersihkan file di Supabase Storage
// SETELAH tahu hasil akhir DB transaction:
//   - Transaction SUKSES  -> hapus OldURL (gambar lama yang sudah tergantikan)
//   - Transaction GAGAL   -> hapus NewURL (gambar baru jadi orphan karena
//     tidak jadi dipakai; OldURL TIDAK disentuh karena DB masih menunjuk ke
//     situ)
//
// SEBELUMNYA utils.ReplaceFile() menghapus file lama secara async begitu
// upload baru sukses — TANPA peduli transaction DB-nya nanti gagal atau
// tidak. Kalau transaction gagal setelah itu, file lama sudah hilang
// permanen padahal DB (ter-rollback) masih menunjuk ke URL yang sudah tidak
// ada. Pola imageChange ini membuat penghapusan file selalu menunggu hasil
// transaction dulu.
type imageChange struct {
	OldURL string
	NewURL string
}

// cleanupImageChanges menghapus file di Supabase berdasarkan hasil transaction.
func cleanupImageChanges(changes []imageChange, txSucceeded bool) {
	if len(changes) == 0 {
		return
	}
	var urls []string
	for _, c := range changes {
		if txSucceeded {
			if c.OldURL != "" {
				urls = append(urls, c.OldURL)
			}
		} else {
			if c.NewURL != "" {
				urls = append(urls, c.NewURL)
			}
		}
	}
	if len(urls) == 0 {
		return
	}
	go func() {
		if err := utils.DeleteMultipleFromSupabase(urls); err != nil {
			fmt.Printf("Warning: gagal membersihkan file gambar: %v\n", err)
		}
	}()
}

// updateSizeVariants implements ProductService.
func (p *ProductServiceImpl) UpdateSizeVariants(colorVarianID int64, sizesParam []models.UpdateSizeVarianRequest, tx *gorm.DB) error {

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
// AddColorVarianProduct implements ProductService.
//
// Sama seperti CreateProduct: upload gambar dipindah keluar dari transaction,
// dan dibersihkan lewat cleanupImageChanges kalau transaction gagal.
func (p *ProductServiceImpl) AddColorVarianProduct(productId int64, param models.CreateColorVarianRequest) (*models.ProductDetailResponse, error) {
	product, err := p.productRepo.FindProductById(productId, nil)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product dengan ID %d tidak ditemukan", productId)
		}
		return nil, fmt.Errorf("error mencari product: %w", err)
	}

	category, err := p.categoryRepo.FindById(product.CategoryID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("category dengan ID %d tidak ditemukan", product.CategoryID)
		}
		return nil, fmt.Errorf("error mengambil category: %w", err)
	}

	for _, cv := range product.ColorVarians {
		if strings.EqualFold(cv.Name, param.Name) {
			return nil, fmt.Errorf("color varian '%s' sudah ada di product ini", param.Name)
		}
	}
	if param.Image == nil {
		return nil, fmt.Errorf("gambar wajib diisi untuk color varian")
	}

	if len(param.Sizes) == 0 {
		return nil, fmt.Errorf("minimal harus ada 1 ukuran untuk color varian")
	}

	sizeMap := make(map[string]bool)
	for _, sv := range param.Sizes {
		if sizeMap[strings.ToUpper(sv.Size)] {
			return nil, fmt.Errorf("ukuran '%s' duplikat dalam color varian '%s'", sv.Size, param.Name)
		}
		sizeMap[strings.ToUpper(sv.Size)] = true
	}

	folderName := fmt.Sprintf("product/%s/colors", sanitizeFileName(product.Name))
	colorImageURL, err := utils.UploadToSupabase(param.Image, folderName)
	if err != nil {
		return nil, fmt.Errorf("gagal upload gambar color varian: %w", err)
	}
	imageChanges := []imageChange{{NewURL: colorImageURL}}

	txErr := database.DB.Transaction(func(tx *gorm.DB) error {
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

	cleanupImageChanges(imageChanges, txErr == nil)

	if txErr != nil {
		return nil, txErr
	}

	productWithRelations, err := p.productRepo.FindProductById(productId, nil)
	if err != nil {
		return nil, fmt.Errorf("gagal memuat data product: %w", err)
	}

	result := productWithRelations.ToProductDetailResponse(category)
	return &result, nil
}

// CreateProduct implements ProductService.
//
// SEBELUMNYA seluruh upload gambar (produk utama + tiap color variant)
// dijalankan DI DALAM database.DB.Transaction(...). Dampaknya: (1) koneksi
// DB dipegang selama proses upload ke Supabase berlangsung, memboroskan
// connection pool untuk I/O jaringan yang tidak ada hubungannya dengan DB;
// (2) kalau upload warna ke-3 gagal setelah warna 1-2 sukses, transaction
// di-rollback tapi file 1-2 sudah terlanjur ada di storage dan tidak pernah
// dibersihkan (orphan file permanen). Sekarang: semua upload dilakukan
// SEBELUM transaction dibuka (transaction jadi murni operasi DB, cepat),
// dan kalau transaction tetap gagal, semua file yang sudah terlanjur
// diupload dibersihkan lewat cleanupImageChanges.
func (p *ProductServiceImpl) CreateProduct(param models.CreateProductParam) (*models.ProductDetailResponse, error) {
	var product models.Product
	var colorVariants []models.ColorVarian

	// Validasi category dulu (read-only, tidak butuh tx) sebelum melakukan
	// upload apa pun — supaya request yang pasti gagal validasi tidak sampai
	// membuang waktu/kuota upload.
	category, err := p.categoryRepo.FindById(param.CategoryID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("category dengan ID %d tidak ditemukan", param.CategoryID)
		}
		return nil, fmt.Errorf("error validasi category: %w", err)
	}

	var imageChanges []imageChange

	var imageURL string
	if param.Image != nil {
		folder := fmt.Sprintf("product/%s", sanitizeFileName(param.Name))
		imageURL, err = utils.UploadToSupabase(param.Image, folder)
		if err != nil {
			return nil, fmt.Errorf("gagal upload gambar produk: %w", err)
		}
		imageChanges = append(imageChanges, imageChange{NewURL: imageURL})
	}

	type colorUpload struct {
		Name, Color, URL string
		Sizes            []models.CreateSizeVarianRequest
	}
	colorUploads := make([]colorUpload, 0, len(param.ColorVarian))
	for i, c := range param.ColorVarian {
		folder := fmt.Sprintf("product/%s/colors", sanitizeFileName(param.Name))
		colorURL, err := utils.UploadToSupabase(c.Image, folder)
		if err != nil {
			// Bersihkan semua file yang sudah terlanjur diupload sebelum ini.
			cleanupImageChanges(imageChanges, false)
			return nil, fmt.Errorf("gagal upload gambar varian warna ke-%d: %w", i+1, err)
		}
		imageChanges = append(imageChanges, imageChange{NewURL: colorURL})
		colorUploads = append(colorUploads, colorUpload{Name: c.Name, Color: c.Color, URL: colorURL, Sizes: c.Sizes})
	}

	txErr := database.DB.Transaction(func(tx *gorm.DB) error {
		pData := models.Product{
			CategoryID:  param.CategoryID,
			Name:        param.Name,
			Description: param.Description,
			Price:       param.Price,
			Images:      imageURL,
		}

		var err error
		product, err = p.productRepo.CreateProduct(pData, tx)
		if err != nil {
			return fmt.Errorf("gagal membuat product: %w", err)
		}

		for i, cu := range colorUploads {
			colorData := models.ColorVarian{
				ProductID: product.ID,
				Name:      cu.Name,
				Color:     cu.Color,
				Images:    cu.URL,
			}

			colorVariant, err := p.productRepo.CreateColorVarian(colorData, tx)
			if err != nil {
				return fmt.Errorf("gagal membuat varian warna ke-%d: %w", i+1, err)
			}

			for j, s := range cu.Sizes {
				sizeData := models.SizeVarian{
					ColorVarianID: colorVariant.ID,
					Size:          s.Size,
					Stock:         s.Stock,
				}

				_, err := p.productRepo.CreateSizeVarian(sizeData, tx)
				if err != nil {
					return fmt.Errorf("gagal membuat size %d for warna %s: %w", j+1, cu.Name, err)
				}
			}

			colorVariants = append(colorVariants, colorVariant)
		}

		return nil
	})

	// Transaction gagal: semua file yang sudah diupload jadi orphan, hapus.
	// Transaction sukses: tidak ada file lama yang perlu dihapus (semua baru).
	cleanupImageChanges(imageChanges, txErr == nil)

	if txErr != nil {
		return nil, txErr
	}

	product.ColorVarians = colorVariants
	result := product.ToProductDetailResponse(category)
	return &result, nil
}

// DeleteProduct implements ProductService.
func (p *ProductServiceImpl) DeleteProduct(id int64) error {
	var imageURLs []string

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		product, err := p.productRepo.FindProductById(id, tx)
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

	products, _, err := p.productRepo.FindAllProduct(param, nil)
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
	product, err := p.productRepo.FindProductById(id, nil)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("produk dengan ID %d tidak ditemukan", id)
		}
		return nil, fmt.Errorf("error mengambil produk: %w", err)
	}

	category, err := p.categoryRepo.FindById(product.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("error mengambil kategori: %w", err)
	}

	resp := product.ToProductDetailResponse(category)
	return &resp, nil
}

// UpdateProduct implements ProductService.
func (p *ProductServiceImpl) UpdateProduct(param models.UpdateProductParam) (*models.ProductDetailResponse, error) {
	var product models.Product
	var category models.Category
	var imageChanges []imageChange

	err := database.DB.Transaction(func(tx *gorm.DB) error {

		existing, err := p.productRepo.FindProductById(param.ID, tx)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("produk ID %d tidak ditemukan", param.ID)
			}
			return fmt.Errorf("error load produk: %w", err)
		}
		product = *existing

		// UPDATE CATEGORY
		if param.CategoryID != nil {
			catRes, err := p.categoryRepo.FindById(*param.CategoryID)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return fmt.Errorf("category ID %d tidak ditemukan", *param.CategoryID)
				}
				return fmt.Errorf("error validate category: %w", err)
			}
			category = catRes
			product.CategoryID = *param.CategoryID
		} else {
			catRes, err := p.categoryRepo.FindById(product.CategoryID)
			if err != nil {
				return err
			}
			category = catRes
		}

		// DUPLICATE NAME CHECK
		if param.Name != nil && *param.Name != product.Name {
			dup, err := p.productRepo.FindByNameAndCategory(*param.Name, product.CategoryID, tx)
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}
			if dup != nil {
				return fmt.Errorf("produk '%s' sudah ada di kategori '%s'", *param.Name, category.Name)
			}
			product.Name = *param.Name
		}

		// BASIC UPDATE
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
			// SEBELUMNYA pakai utils.ReplaceFile(...) yang langsung menghapus
			// file lama secara async begitu upload baru sukses — TIDAK peduli
			// transaction ini nanti berhasil atau di-rollback. Sekarang cuma
			// upload gambar baru; penghapusan file lama ditunda sampai
			// transaction benar-benar sukses (lihat cleanupImageChanges di
			// bawah, dipanggil setelah blok Transaction()).
			folder := fmt.Sprintf("product/%s", sanitizeFileName(product.Name))
			newURL, err := utils.UploadToSupabase(param.Image, folder)
			if err != nil {
				return err
			}
			imageChanges = append(imageChanges, imageChange{OldURL: product.Images, NewURL: newURL})
			product.Images = newURL
		}

		product, err = p.productRepo.UpdateProduct(product, tx)
		if err != nil {
			return err
		}

		// HANDLE COLOR VARIANTS
		if len(param.ColorVarian) > 0 {
			existingColors, err := p.productRepo.FindColorVarianByProductId(product.ID, tx)
			if err != nil {
				return err
			}

			// map untuk tracking delete
			colorMap := make(map[int64]*models.ColorVarian)
			for i := range existingColors {
				colorMap[existingColors[i].ID] = &existingColors[i]
			}

			updatedIDs := make(map[int64]bool)

			for _, cv := range param.ColorVarian {
				if cv.ID != nil {
					// UPDATE OLD VARIANT
					exColor := colorMap[*cv.ID]
					if exColor == nil {
						return fmt.Errorf("color varian ID %d tidak ditemukan", *cv.ID)
					}

					if cv.Name != nil {
						exColor.Name = *cv.Name
					}
					if cv.Color != nil {
						exColor.Color = *cv.Color
					}

					if cv.Image != nil {
						folder := fmt.Sprintf("product/%s/colors", sanitizeFileName(product.Name))
						newURL, err := utils.UploadToSupabase(cv.Image, folder)
						if err != nil {
							return err
						}
						imageChanges = append(imageChanges, imageChange{OldURL: exColor.Images, NewURL: newURL})
						exColor.Images = newURL
					}

					_, err := p.productRepo.UpdateColorVarian(*exColor, tx)
					if err != nil {
						return err
					}

					updatedIDs[*cv.ID] = true

					// UPDATE SIZES
					if len(cv.Sizes) > 0 {
						if err := p.UpdateSizeVariants(*cv.ID, cv.Sizes, tx); err != nil {
							return err
						}
					}

				} else {
					// CREATE NEW COLOR
					if cv.Image == nil {
						return fmt.Errorf("gambar warna baru wajib diisi")
					}

					url, err := utils.UploadToSupabase(cv.Image, fmt.Sprintf("product/%s/colors", sanitizeFileName(product.Name)))
					if err != nil {
						return err
					}
					imageChanges = append(imageChanges, imageChange{NewURL: url})

					newColor := models.ColorVarian{
						ProductID: product.ID,
						Name:      *cv.Name,
						Color:     *cv.Color,
						Images:    url,
					}

					created, err := p.productRepo.CreateColorVarian(newColor, tx)
					if err != nil {
						return err
					}

					updatedIDs[created.ID] = true

					for _, s := range cv.Sizes {
						_, err := p.productRepo.CreateSizeVarian(models.SizeVarian{
							ColorVarianID: created.ID,
							Size:          *s.Size,
							Stock:         *s.Stock,
						}, tx)
						if err != nil {
							return err
						}
					}
				}
			}

			// DELETE unused color variants
			for id := range colorMap {
				if !updatedIDs[id] {
					if err := p.productRepo.DeleteColorVarian(id, tx); err != nil {
						return err
					}
				}
			}
		}

		return nil
	})

	// err == nil (sukses)   -> hapus file LAMA yang sudah tergantikan.
	// err != nil (rollback) -> hapus file BARU yang jadi orphan (DB masih
	// menunjuk ke file lama yang tidak pernah tersentuh).
	cleanupImageChanges(imageChanges, err == nil)

	if err != nil {
		return nil, err
	}

	pd, err := p.productRepo.FindProductById(param.ID, nil)
	if err != nil {
		return nil, err
	}

	resp := pd.ToProductDetailResponse(category)
	return &resp, nil
}

func NewProductService(categoryRepo repository.CategoryRepository,
	productRepo repository.ProductRepository) ProductService {
	return &ProductServiceImpl{categoryRepo: categoryRepo, productRepo: productRepo}
}
