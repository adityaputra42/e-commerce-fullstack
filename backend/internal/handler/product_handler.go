package handler

import (
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/services"
	"e-commerce/backend/internal/utils"
	"encoding/json"
	"net/http"
)

type ProductHandler struct {
	productService services.ProductService
}

func NewProductHandler(productService services.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// CreateAddress handles POST /Product
func (h *ProductHandler) CreateAddress(w http.ResponseWriter, r *http.Request) {

	var param models.CreateProductParam
	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.productService.CreateProduct(param)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Product created successfully",
		"data":    result,
	})
}
