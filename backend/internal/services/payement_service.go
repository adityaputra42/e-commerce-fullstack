package services

import (
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/repository"
	"errors"
	"fmt"
)

type PaymentService interface {
	CreatePayment(param models.CreatePayment) (*models.PaymentResponse, error)
	UpdatePayment(param models.UpdatePayment) (*models.PaymentResponse, error)
	FindAllPayment(param models.PaymentListRequest) (*[]models.PaymentResponse, error)
	FindById(id int64) (*models.PaymentResponse, error)
	DeletePayment(id int64) error
}

type PaymentServiceImpl struct {
	paymentRepo     repository.PaymentRepository
	transactionRepo repository.TransactionRepository
}

// CreatePayment implements PaymentService.
func (p *PaymentServiceImpl) CreatePayment(param models.CreatePayment) (*models.PaymentResponse, error) {

	if param.TransactionID == "" {
		return nil, errors.New("transaction id is required")
	}

	if param.TotalPayment <= 0 {
		return nil, errors.New("total payment must be greater than 0")
	}

	transaction, err := p.transactionRepo.FindById(param.TransactionID)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	if transaction.TotalPrice != param.TotalPayment {
		return nil, errors.New("total payment didn't match with transaction total price")
	}

	payParam := models.Payment{TransactionID: transaction.TxID, Status: "pending", TotalPayment: param.TotalPayment}
	payment, err := p.paymentRepo.Create(payParam, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}
	result := payment.ToResponsePayment()
	return result, nil
}

// UpdatePayment implements PaymentService.
func (p *PaymentServiceImpl) UpdatePayment(param models.UpdatePayment) (*models.PaymentResponse, error) {
	if param.ID <= 0 {
		return nil, errors.New("invalid payment id")
	}

	existingPayment, err := p.paymentRepo.FindById(uint(param.ID))
	if err != nil {
		return nil, errors.New("payment not found")
	}

	if existingPayment.Status == param.Status {
		return nil, errors.New("payment status is already the same, no changes needed")
	}

	if err := p.validateStatusTransition(existingPayment.Status, param.Status); err != nil {
		return nil, err
	}
	existingPayment.Status = param.Status

	updatedPayment, err := p.paymentRepo.Update(existingPayment, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}
	result := updatedPayment.ToResponsePayment()
	return result, nil
}

// FindAllPayment implements PaymentService.
func (p *PaymentServiceImpl) FindAllPayment(param models.PaymentListRequest) (*[]models.PaymentResponse, error) {

	payments, err := p.paymentRepo.FindAll(param)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment list: %w", err)
	}
	var paymentResponse []models.PaymentResponse

	for _, payment := range payments {

		paymentResponse = append(paymentResponse, *payment.ToResponsePayment())

	}
	return &paymentResponse, nil
}

// FindById implements PaymentService.
func (p *PaymentServiceImpl) FindById(id int64) (*models.PaymentResponse, error) {
	if id <= 0 {
		return nil, errors.New("invalid payment id")
	}

	// Get payment
	payment, err := p.paymentRepo.FindById(uint(id))
	if err != nil {
		return nil, errors.New("payment not found")
	}

	// Get transaction detail untuk response yang lebih lengkap
	transaction, err := p.transactionRepo.FindById(payment.TransactionID)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	payment.Transaction = transaction

	result := payment.ToResponsePayment()

	return result, nil
}

// DeletePayment implements PaymentService.
func (p *PaymentServiceImpl) DeletePayment(id int64) error {
	if id <= 0 {
		return errors.New("invalid payment id")
	}

	payment, err := p.paymentRepo.FindById(uint(id))
	if err != nil {
		return errors.New("payment not found")
	}

	if payment.Status == "completed" || payment.Status == "confirmed" {
		return errors.New("cannot delete payment with completed or confirmed status")
	}

	// Delete payment
	err = p.paymentRepo.Delete(payment)
	if err != nil {
		return fmt.Errorf("failed to delete payment: %w", err)
	}

	return nil
}

// Helper function untuk validasi status transition
func (p *PaymentServiceImpl) validateStatusTransition(currentStatus, newStatus string) error {
	// Define valid status transitions
	validTransitions := map[string][]string{
		"pending":   {"confirmed", "rejected", "cancelled"},
		"confirmed": {"completed", "refunded"},
		"rejected":  {"pending"}, // Allow resubmit
		"cancelled": {},          // Cannot change from cancelled
		"completed": {"refunded"},
		"refunded":  {}, // Final state
	}

	allowedStatuses, exists := validTransitions[currentStatus]
	if !exists {
		return fmt.Errorf("invalid current status: %s", currentStatus)
	}

	// Check if new status is allowed
	for _, allowed := range allowedStatuses {
		if newStatus == allowed {
			return nil
		}
	}

	return fmt.Errorf("invalid status transition from %s to %s", currentStatus, newStatus)
}

func NewPaymentService(paymentRepo repository.PaymentRepository, transactionRepo repository.TransactionRepository) PaymentService {
	return &PaymentServiceImpl{
		paymentRepo:     paymentRepo,
		transactionRepo: transactionRepo,
	}
}
