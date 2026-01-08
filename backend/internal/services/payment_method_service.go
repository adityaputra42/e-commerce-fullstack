package services

import (
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/repository"
	"e-commerce/backend/internal/utils"
	"errors"
	"fmt"
	"mime/multipart"
)

type PaymentMethodService interface {
	CreatePaymentMethod(param models.CreatePaymentMethod, file *multipart.FileHeader) (*models.PaymentMethodResponse, error)
	UpdatePaymentMethod(param models.UpdatePaymentMethod, file *multipart.FileHeader) (*models.PaymentMethodResponse, error)
	FindAllPaymentMethod(param models.PaymentMethodListRequest) ([]models.PaymentMethodResponse, error)
	FindById(id int64) (*models.PaymentMethodResponse, error)
	DeletePaymentMethod(id int64) error
}

type PaymentMethodServiceImpl struct {
	paymentMethodRepo repository.PaymentMethodRepository
}

const bucketFolder = "payment_method" // folder di Supabase

// -----------------------------------------------------------
// Validation Helper
// -----------------------------------------------------------
func validatePMFields(accountName, accountNumber, bankName string) error {
	if accountName == "" {
		return errors.New("account name is required")
	}
	if accountNumber == "" {
		return errors.New("account number is required")
	}
	if bankName == "" {
		return errors.New("bank name is required")
	}
	return nil
}

// -----------------------------------------------------------
// Create
// -----------------------------------------------------------
func (pm *PaymentMethodServiceImpl) CreatePaymentMethod(param models.CreatePaymentMethod, file *multipart.FileHeader) (*models.PaymentMethodResponse, error) {

	if err := validatePMFields(param.AccountName, param.AccountNumber, param.BankName); err != nil {
		return nil, err
	}

	if file == nil {
		return nil, errors.New("bank image is required")
	}

	// Upload ke Supabase
	imageURL, err := utils.UploadToSupabase(file, bucketFolder)
	if err != nil {
		return nil, fmt.Errorf("failed upload bank image: %w", err)
	}

	newPM := models.PaymentMethod{
		AccountName:   param.AccountName,
		AccountNumber: param.AccountNumber,
		BankName:      param.BankName,
		BankImages:    imageURL,
	}

	created, err := pm.paymentMethodRepo.Create(newPM, nil)
	if err != nil {
		utils.DeleteFromSupabase(imageURL) // rollback
		return nil, fmt.Errorf("failed to create payment method: %w", err)
	}

	return created.ToResponsePaymentMethod(), nil
}

// -----------------------------------------------------------
// Update
// -----------------------------------------------------------
func (pm *PaymentMethodServiceImpl) UpdatePaymentMethod(param models.UpdatePaymentMethod, file *multipart.FileHeader) (*models.PaymentMethodResponse, error) {

	existing, err := pm.paymentMethodRepo.FindById(uint(param.ID))
	if err != nil {
		return nil, errors.New("payment method not found")
	}

	var newImageURL string
	oldImageURL := existing.BankImages

	// Jika ada file baru → replace
	if file != nil {
		newImageURL, err = utils.ReplaceFile(oldImageURL, file, bucketFolder)
		if err != nil {
			return nil, fmt.Errorf("failed replacing bank image: %w", err)
		}
		existing.BankImages = newImageURL
	}

	// Update field lain
	if param.AccountName != "" {
		existing.AccountName = param.AccountName
	}
	if param.AccountNumber != "" {
		existing.AccountNumber = param.AccountNumber
	}
	if param.BankName != "" {
		existing.BankName = param.BankName
	}
	if param.IsActive != nil {
		existing.IsActive = *param.IsActive
	}

	updated, err := pm.paymentMethodRepo.Update(*existing, nil)
	if err != nil {

		if newImageURL != "" {
			utils.DeleteFromSupabase(newImageURL)
		}

		return nil, fmt.Errorf("failed to update payment method: %w", err)
	}

	return updated.ToResponsePaymentMethod(), nil
}

// -----------------------------------------------------------
// Find All
// -----------------------------------------------------------
func (pm *PaymentMethodServiceImpl) FindAllPaymentMethod(param models.PaymentMethodListRequest) ([]models.PaymentMethodResponse, error) {

	all, err := pm.paymentMethodRepo.FindAll(param)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment methods: %w", err)
	}

	resp := make([]models.PaymentMethodResponse, 0, len(all))
	for _, v := range all {
		resp = append(resp, *v.ToResponsePaymentMethod())
	}

	return resp, nil
}

// -----------------------------------------------------------
// Find By ID
// -----------------------------------------------------------
func (pm *PaymentMethodServiceImpl) FindById(id int64) (*models.PaymentMethodResponse, error) {
	if id <= 0 {
		return nil, errors.New("invalid payment method id")
	}

	pmData, err := pm.paymentMethodRepo.FindById(uint(id))
	if err != nil {
		return nil, errors.New("payment method not found")
	}

	return pmData.ToResponsePaymentMethod(), nil
}

// -----------------------------------------------------------
// Delete
// -----------------------------------------------------------
func (pm *PaymentMethodServiceImpl) DeletePaymentMethod(id int64) error {
	if id <= 0 {
		return errors.New("invalid payment method id")
	}

	existing, err := pm.paymentMethodRepo.FindById(uint(id))
	if err != nil {
		return errors.New("payment method not found")
	}

	err = pm.paymentMethodRepo.Delete(*existing)
	if err != nil {
		return fmt.Errorf("failed to delete payment method: %w", err)
	}

	// hapus file image dari Supabase
	utils.DeleteFromSupabase(existing.BankImages)

	return nil
}

// -----------------------------------------------------------
func NewPaymentMethodService(repo repository.PaymentMethodRepository) PaymentMethodService {
	return &PaymentMethodServiceImpl{paymentMethodRepo: repo}
}
