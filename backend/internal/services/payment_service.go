package services

import (
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/repository"
	"e-commerce/backend/internal/utils"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type PaymentService interface {
	CreatePayment(param models.CreatePayment) (*models.PaymentResponse, error)
	UpdatePayment(param models.UpdatePayment) (*models.PaymentResponse, error)
	FindAllPayment(param models.PaymentListRequest) (*[]models.PaymentResponse, error)
	FindById(id int64) (*models.PaymentResponse, error)
	DeletePayment(id int64) error
}

type PaymentServiceImpl struct {
	paymentRepo        repository.PaymentRepository
	transactionRepo    repository.TransactionRepository
	fulfillmentService OrderFulfillmentService
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

// paymentStatusToTransactionStatus memetakan status Payment ke status
// Transaction yang seharusnya ikut berubah. Status yang tidak ada di map ini
// (mis. "rejected" — rejected mengizinkan resubmit ke "pending", jadi
// Transaction tidak perlu diapa-apakan) sengaja tidak memicu cascade.
//
// CATATAN (asumsi bisnis, sama seperti di order_fulfillment_service.go):
// Payment "completed" TIDAK dipetakan ke Transaction status apa pun di sini
// — "completed" pada Payment diartikan sebagai "dana sudah settle/
// terkonfirmasi penuh", bukan otomatis berarti "barang sudah sampai".
// Transaction tetap maju lewat progres shippingnya sendiri (paid -> shipped
// -> completed), independen dari kapan uangnya settle. Kalau asumsi ini
// keliru secara bisnis, ini satu-satunya tempat yang perlu diubah.
var paymentStatusToTransactionStatus = map[string]string{
	"confirmed": "paid",
	"cancelled": "cancelled",
	"refunded":  "refunded",
}

// UpdatePayment implements PaymentService.
//
// SEBELUMNYA fungsi ini HANYA mengubah Payment.Status — Transaction dan
// Order terkait tidak pernah ikut disentuh sama sekali, bahkan tidak
// dibungkus dalam DB transaction (`p.paymentRepo.Update(existingPayment,
// nil)`, tx-nya nil). Payment yang dikonfirmasi tidak pernah mendorong
// Transaction jadi "paid", dan Payment yang dibatalkan/refund tidak pernah
// mengembalikan stock Order-order di dalamnya. Sekarang semuanya jadi SATU
// operasi atomic lewat OrderFulfillmentService.
func (p *PaymentServiceImpl) UpdatePayment(param models.UpdatePayment) (*models.PaymentResponse, error) {
	if param.ID <= 0 {
		return nil, errors.New("invalid payment id")
	}

	var result models.Payment

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		existingPayment, err := p.paymentRepo.FindByIdLocking(tx, uint(param.ID))
		if err != nil {
			return errors.New("payment not found")
		}

		if existingPayment.Status == param.Status {
			return errors.New("payment status is already the same, no changes needed")
		}

		if err := p.validateStatusTransition(existingPayment.Status, param.Status); err != nil {
			return err
		}
		existingPayment.Status = param.Status

		updatedPayment, err := p.paymentRepo.Update(existingPayment, tx)
		if err != nil {
			return fmt.Errorf("failed to update payment: %w", err)
		}
		result = updatedPayment

		// Cascade ke Transaction + Order kalau status Payment ini memang
		// seharusnya mendorong Transaction maju/mundur (lihat
		// paymentStatusToTransactionStatus di atas untuk penjelasan mapping).
		newTransactionStatus, shouldCascade := paymentStatusToTransactionStatus[param.Status]
		if !shouldCascade {
			return nil
		}

		transaction, err := p.transactionRepo.FindByIdLocking(tx, existingPayment.TransactionID)
		if err != nil {
			return fmt.Errorf("gagal mengambil transaction terkait payment: %w", err)
		}

		if !utils.IsValidStatusTransition(transaction.Status, newTransactionStatus) {
			// Transaction sudah di status yang tidak kompatibel dengan
			// cascade ini (mis. sudah "completed" atau "cancelled" duluan
			// lewat jalur lain). Payment tetap berhasil diupdate — ini
			// cuma soal sinkronisasi sekunder yang sudah tidak relevan lagi
			// — tapi dicatat di log supaya kelihatan kalau ini sering
			// terjadi (bisa jadi tanda ada race/urutan operasi yang aneh).
			log.Printf("PaymentService: skip cascade transaction %s dari %s ke %s (transisi tidak valid)", transaction.TxID, transaction.Status, newTransactionStatus)
			return nil
		}

		transaction.Status = newTransactionStatus
		if _, err := p.transactionRepo.Update(*transaction, tx); err != nil {
			return fmt.Errorf("gagal update transaction terkait payment: %w", err)
		}

		if err := p.fulfillmentService.SyncOrdersToTransactionStatus(tx, transaction.TxID, newTransactionStatus); err != nil {
			return fmt.Errorf("gagal menyinkronkan order terkait payment: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result.ToResponsePayment(), nil
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

// Payment punya vocabulary status sendiri (beda dari Transaction), tapi
// logic validasinya sekarang dibangun dari StatusTransition yang sama
// dipakai transaction (lihat utils/status_transition.go) — bukan
// diimplementasikan manual terpisah lagi.
var paymentStatusTransition = utils.NewStatusTransition(map[string][]string{
	"pending":   {"confirmed", "rejected", "cancelled"},
	"confirmed": {"completed", "refunded"},
	"rejected":  {"pending"}, // Allow resubmit
	"cancelled": {},          // Cannot change from cancelled
	"completed": {"refunded"},
	"refunded":  {}, // Final state
})

// validateStatusTransition adalah helper validasi status transition.
func (p *PaymentServiceImpl) validateStatusTransition(currentStatus, newStatus string) error {
	return paymentStatusTransition.Validate(currentStatus, newStatus)
}

func NewPaymentService(paymentRepo repository.PaymentRepository, transactionRepo repository.TransactionRepository, orderRepo repository.OrderRepository, productRepo repository.ProductRepository) PaymentService {
	return &PaymentServiceImpl{
		paymentRepo:        paymentRepo,
		transactionRepo:    transactionRepo,
		fulfillmentService: NewOrderFulfillmentService(orderRepo, productRepo),
	}
}
