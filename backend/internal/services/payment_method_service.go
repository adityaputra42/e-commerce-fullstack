package services

import (
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/repository"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
)

type PaymentMethodService interface {
	CreatePaymentMethod(param models.CreatePaymentMethod, file *multipart.FileHeader) (*models.PaymentMethodResponse, error)
	UpdatePaymentMethod(param models.UpdatePaymentMethod, file *multipart.FileHeader) (*models.PaymentMethodResponse, error)
	FindAllPaymentMethod(param models.PaymentMethodListRequest) (*[]models.PaymentMethodResponse, error)
	FindById(id int64) (*models.PaymentMethodResponse, error)
	DeletePaymentMethod(id int64) error
}

type PaymentMethodServiceImpl struct {
	paymentMethodRepo repository.PaymentMethodRepository
}

// CreatePaymentMethod implements PaymentMethodService.
func (pm *PaymentMethodServiceImpl) CreatePaymentMethod(param models.CreatePaymentMethod, file *multipart.FileHeader) (*models.PaymentMethodResponse, error) {
	// Validasi input
	if param.AccountName == "" {
		return nil, errors.New("account name is required")
	}

	if param.AccountNumber == "" {
		return nil, errors.New("account number is required")
	}

	if param.BankName == "" {
		return nil, errors.New("bank name is required")
	}

	if file == nil {
		return nil, errors.New("bank images is required")
	}

	// Save file
	filePath, err := pm.saveFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to save bank image: %w", err)
	}

	pmParam := models.PaymentMethod{
		AccountName:   param.AccountName,
		AccountNumber: param.AccountNumber,
		BankName:      param.BankName,
		BankImages:    filePath,
	}

	paymentMethod, err := pm.paymentMethodRepo.Create(pmParam, nil)
	if err != nil {

		os.Remove(filePath)
		return nil, fmt.Errorf("failed to create payment method: %w", err)
	}
	result := paymentMethod.ToResponsePaymentMethod()
	return result, nil
}

// UpdatePaymentMethod implements PaymentMethodService.
func (pm *PaymentMethodServiceImpl) UpdatePaymentMethod(param models.UpdatePaymentMethod, file *multipart.FileHeader) (*models.PaymentMethodResponse, error) {
	// Validasi payment method exists
	existingPM, err := pm.paymentMethodRepo.FindById(uint(param.ID))
	if err != nil {
		return nil, errors.New("payment method not found")
	}

	if param.AccountName == "" {
		return nil, errors.New("account name is required")
	}

	if param.AccountNumber == "" {
		return nil, errors.New("account number is required")
	}

	if param.BankName == "" {
		return nil, errors.New("bank name is required")
	}

	if file == nil {
		return nil, errors.New("bank images is required")
	}

	// Save new file
	filePath, err := pm.saveFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to save bank image: %w", err)
	}

	existingPM.AccountName = param.AccountName
	existingPM.AccountNumber = param.AccountNumber
	existingPM.BankImages = filePath
	existingPM.BankName = param.BankName
	updatedPM, err := pm.paymentMethodRepo.Update(*existingPM, nil)
	if err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to update payment method: %w", err)
	}

	if existingPM.BankImages != "" && existingPM.BankImages != filePath {
		os.Remove(existingPM.BankImages)
	}
	result := updatedPM.ToResponsePaymentMethod()
	return result, nil
}

// FindAllPaymentMethod implements PaymentMethodService.
func (pm *PaymentMethodServiceImpl) FindAllPaymentMethod(param models.PaymentMethodListRequest) (*[]models.PaymentMethodResponse, error) {

	paymentMethods, err := pm.paymentMethodRepo.FindAll(param)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment method list: %w", err)
	}
	var response []models.PaymentMethodResponse

	for _, v := range paymentMethods {
		response = append(response, *v.ToResponsePaymentMethod())
	}
	return &response, nil
}

// FindById implements PaymentMethodService.
func (pm *PaymentMethodServiceImpl) FindById(id int64) (*models.PaymentMethodResponse, error) {
	if id <= 0 {
		return nil, errors.New("invalid payment method id")
	}

	paymentMethod, err := pm.paymentMethodRepo.FindById(uint(id))
	if err != nil {
		return nil, errors.New("payment method not found")
	}
	result := paymentMethod.ToResponsePaymentMethod()
	return result, nil
}

// DeletePaymentMethod implements PaymentMethodService.
func (pm *PaymentMethodServiceImpl) DeletePaymentMethod(id int64) error {
	if id <= 0 {
		return errors.New("invalid payment method id")
	}

	paymentMethod, err := pm.paymentMethodRepo.FindById(uint(id))
	if err != nil {
		return errors.New("payment method not found")
	}

	err = pm.paymentMethodRepo.Delete(*paymentMethod)
	if err != nil {
		return fmt.Errorf("failed to delete payment method: %w", err)
	}

	if paymentMethod.BankImages != "" {
		os.Remove(paymentMethod.BankImages)
	}

	return nil
}

// Helper function untuk save file
func (pm *PaymentMethodServiceImpl) saveFile(file *multipart.FileHeader) (string, error) {
	folder := "./assets/bank/"

	// Create directory if not exists
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	filePath := filepath.Join(folder, file.Filename)

	// Save file menggunakan multipart
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err = dst.ReadFrom(src); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return filePath, nil
}

func NewPaymentMethodService(paymentMethodRepo repository.PaymentMethodRepository) PaymentMethodService {
	return &PaymentMethodServiceImpl{
		paymentMethodRepo: paymentMethodRepo,
	}
}
