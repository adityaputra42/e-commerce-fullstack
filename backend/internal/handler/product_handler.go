package handler

import (
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	productService services.ProductService
}

func NewProductHandler(productService services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct - POST /api/v1/products
// @Summary Create a new product
// @Description Create a new product with images and variants
// @Tags Product
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Product name"
// @Param description formData string true "Product description"
// @Param price formData number true "Product price"
// @Param category_id formData int true "Category ID"
// @Param image formData file true "Main product image"
// @Param color_varian formData string true "JSON string of color variants"
// @Success 201 {object} utils.Response{data=models.Product} "Product created successfully"
// @Failure 400 {object} utils.Response "Invalid input"
// @Router /products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		respondError(w, http.StatusBadRequest, "Gagal parse form data", err)
		return
	}

	categoryID, err := strconv.ParseInt(r.FormValue("category_id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Category ID tidak valid", err)
		return
	}

	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Price tidak valid", err)
		return
	}

	_, mainImageHeader, err := r.FormFile("image")
	if err != nil {
		respondError(w, http.StatusBadRequest, "Gambar produk wajib diisi", err)
		return
	}

	var colorVariants []models.CreateColorVarianRequest
	colorVariantsJSON := r.FormValue("color_varian")
	if colorVariantsJSON == "" {
		respondError(w, http.StatusBadRequest, "Color variants wajib diisi", nil)
		return
	}

	if err := json.Unmarshal([]byte(colorVariantsJSON), &colorVariants); err != nil {
		respondError(w, http.StatusBadRequest, "Format color variants tidak valid", err)
		return
	}

	for i := range colorVariants {
		fieldName := fmt.Sprintf("color_image_%d", i)
		_, fileHeader, err := r.FormFile(fieldName)
		if err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("Gambar untuk color variant ke-%d wajib diisi", i+1), err)
			return
		}
		colorVariants[i].Image = fileHeader
	}

	param := models.CreateProductParam{
		CategoryID:  categoryID,
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Price:       price,
		Image:       mainImageHeader,
		ColorVarian: colorVariants,
	}

	result, err := h.productService.CreateProduct(param)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Gagal membuat produk", err)
		return
	}

	respondSuccess(w, http.StatusCreated, "Produk berhasil dibuat", result)
}

// GetProductByID - GET /api/v1/products/{id}
// @Summary Get product by ID
// @Description Get detailed information about a product
// @Tags Product
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} utils.Response{data=models.Product} "Product found"
// @Failure 404 {object} utils.Response "Product not found"
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "ID produk tidak valid", err)
		return
	}

	result, err := h.productService.FindProductById(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Produk tidak ditemukan", err)
		return
	}

	respondSuccess(w, http.StatusOK, "Produk ditemukan", result)
}

// GetAllProducts - GET /api/v1/products
// @Summary List products
// @Description Get a paginated list of products with filters
// @Tags Product
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search by name"
// @Param category_id query int false "Filter by category"
// @Param sort_by query string false "Sort by field" default(created_at)
// @Success 200 {object} utils.Response{data=[]models.Product} "Data produk berhasil diambil"
// @Router /products [get]
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 10
	}

	var categoryID int64
	if catStr := r.URL.Query().Get("category_id"); catStr != "" {
		if cat, err := strconv.ParseInt(catStr, 10, 64); err == nil {
			categoryID = cat
		}
	}

	sortBy := r.URL.Query().Get("sort_by")
	if sortBy == "" {
		sortBy = "created_at"
	}

	param := models.ProductListRequest{
		Page:       page,
		Limit:      limit,
		Search:     r.URL.Query().Get("search"),
		CategoryID: categoryID,

		SortBy: sortBy,
	}

	result, err := h.productService.FindAllProduct(param)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Gagal mengambil data produk", err)
		return
	}

	respondSuccess(w, http.StatusOK, "Data produk berhasil diambil", result)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "ID produk tidak valid", err)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		respondError(w, http.StatusBadRequest, "Gagal parse form data", err)
		return
	}

	param := models.UpdateProductParam{
		ID: id,
	}

	// Parse optional fields
	if catStr := r.FormValue("category_id"); catStr != "" {
		catID, err := strconv.ParseInt(catStr, 10, 64)
		if err == nil {
			param.CategoryID = &catID
		}
	}

	if name := r.FormValue("name"); name != "" {
		param.Name = &name
	}

	if desc := r.FormValue("description"); desc != "" {
		param.Description = &desc
	}

	if priceStr := r.FormValue("price"); priceStr != "" {
		price, err := strconv.ParseFloat(priceStr, 64)
		if err == nil {
			param.Price = &price
		}
	}

	if ratingStr := r.FormValue("rating"); ratingStr != "" {
		rating, err := strconv.ParseFloat(ratingStr, 64)
		if err == nil {
			param.Rating = &rating
		}
	}

	// Get main image file header if exists
	if _, fileHeader, err := r.FormFile("image"); err == nil {
		param.Image = fileHeader
	}

	// Parse color variants if exists
	if cvStr := r.FormValue("color_varian"); cvStr != "" {
		var colorVariants []models.UpdateColorVarianRequest
		if err := json.Unmarshal([]byte(cvStr), &colorVariants); err != nil {
			respondError(w, http.StatusBadRequest, "Format color variants tidak valid", err)
			return
		}

		// Get color variant images from form files
		for i := range colorVariants {
			// For existing color variants with ID
			if colorVariants[i].ID != nil {
				fieldName := fmt.Sprintf("color_image_%d", *colorVariants[i].ID)
				if _, fileHeader, err := r.FormFile(fieldName); err == nil {
					colorVariants[i].Image = fileHeader
				}
			} else {
				// For new color variants without ID
				fieldName := fmt.Sprintf("color_image_new_%d", i)
				if _, fileHeader, err := r.FormFile(fieldName); err == nil {
					colorVariants[i].Image = fileHeader
				}
			}
		}
		param.ColorVarian = colorVariants
	}

	result, err := h.productService.UpdateProduct(param)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Gagal update produk", err)
		return
	}

	respondSuccess(w, http.StatusOK, "Produk berhasil diupdate", result)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "ID produk tidak valid", err)
		return
	}

	if err := h.productService.DeleteProduct(id); err != nil {
		respondError(w, http.StatusInternalServerError, "Gagal menghapus produk", err)
		return
	}

	respondSuccess(w, http.StatusOK, "Produk berhasil dihapus", nil)
}

func (h *ProductHandler) AddColorVariant(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "ID produk tidak valid", err)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		respondError(w, http.StatusBadRequest, "Gagal parse form data", err)
		return
	}

	_, imageHeader, err := r.FormFile("image")
	if err != nil {
		respondError(w, http.StatusBadRequest, "Gambar color variant wajib diisi", err)
		return
	}

	var sizes []models.CreateSizeVarianRequest
	sizesJSON := r.FormValue("sizes")
	if sizesJSON == "" {
		respondError(w, http.StatusBadRequest, "Sizes wajib diisi", nil)
		return
	}

	if err := json.Unmarshal([]byte(sizesJSON), &sizes); err != nil {
		respondError(w, http.StatusBadRequest, "Format sizes tidak valid", err)
		return
	}

	param := models.CreateColorVarianRequest{
		Name:  r.FormValue("name"),
		Color: r.FormValue("color"),
		Image: imageHeader,
		Sizes: sizes,
	}

	result, err := h.productService.AddColorVarianProduct(id, param)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Gagal menambahkan color variant", err)
		return
	}

	respondSuccess(w, http.StatusOK, "Color variant berhasil ditambahkan", result)
}

// Helper functions
func respondSuccess(w http.ResponseWriter, code int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func respondError(w http.ResponseWriter, code int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := map[string]interface{}{
		"success": false,
		"message": message,
	}

	if err != nil {
		response["error"] = err.Error()
	}

	json.NewEncoder(w).Encode(response)
}
