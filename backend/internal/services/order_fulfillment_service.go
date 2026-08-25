package services

import (
	"e-commerce/backend/internal/repository"
	"fmt"

	"gorm.io/gorm"
)

// OrderFulfillmentService adalah SATU-SATUNYA tempat yang bertanggung jawab
// atas "apa yang terjadi ke Order-order sebuah Transaction ketika status
// Transaction itu berubah". Sebelumnya (lihat laporan arsitektur awal, item
// #5) tidak ada tempat seperti ini sama sekali:
//   - TransactionServiceImpl.UpdateTransaction cuma mengubah Transaction.Status,
//     tidak pernah menyentuh Order sama sekali.
//   - PaymentServiceImpl.UpdatePayment cuma mengubah Payment.Status, tidak
//     pernah menyentuh Transaction atau Order sama sekali.
//
// Akibatnya: payment yang sukses tidak pernah mendorong Order maju secara
// otomatis, dan transaction yang dibatalkan tidak pernah mengembalikan stock
// Order-order di dalamnya. Setiap orkestrasi yang seharusnya otomatis jadi
// tanggung jawab manual admin lewat panel terpisah-pisah — rawan human error
// dan race condition.
//
// Service ini TIDAK menggantikan TransactionService/PaymentService/
// OrderService untuk operasi CRUD sederhana — service itu tetap ada dan
// tetap dipakai untuk itu. Service ini HANYA menangani efek samping lintas-
// entity yang harus terjadi bersamaan, dipanggil DI DALAM DB transaction
// yang sama dengan update Transaction/Payment (lihat pemanggilnya di
// transaction_service.go dan payment_service.go).
type OrderFulfillmentService interface {
	// SyncOrdersToTransactionStatus mengubah status semua Order milik sebuah
	// Transaction supaya konsisten dengan status Transaction yang baru, dan
	// merestock Order yang dibatalkan/refund. WAJIB dipanggil dengan tx yang
	// sama dengan yang dipakai untuk update Transaction/Payment, supaya
	// semuanya atomic (kalau salah satu gagal, semuanya di-rollback).
	SyncOrdersToTransactionStatus(tx *gorm.DB, transactionID string, newTransactionStatus string) error
}

type OrderFulfillmentServiceImpl struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
}

func NewOrderFulfillmentService(orderRepo repository.OrderRepository, productRepo repository.ProductRepository) OrderFulfillmentService {
	return &OrderFulfillmentServiceImpl{orderRepo: orderRepo, productRepo: productRepo}
}

// orderStatusTerminal adalah status Order yang dianggap final — Order yang
// sudah di status ini TIDAK disentuh lagi oleh sinkronisasi otomatis (supaya
// tidak, misalnya, "membatalkan" order yang sudah completed).
var orderStatusTerminal = map[string]bool{
	"completed": true,
	"cancelled": true,
}

// mapTransactionStatusToOrderStatus adalah PEMETAAN EKSPLISIT dan SENGAJA
// dari status Transaction ke status Order yang seharusnya mengikuti.
//
// CATATAN PENTING (asumsi bisnis yang saya buat, mohon dikonfirmasi ulang
// oleh product owner): vocabulary status Transaction dan Order di codebase
// ini TIDAK 1:1 (Transaction punya "waiting_payment"/"processing"/"refunded"
// yang tidak ada padanannya di Order, yang cuma punya
// pending/paid/shipped/completed/cancelled). Saya buat pemetaan yang paling
// masuk akal secara bisnis, tapi ini keputusan produk yang idealnya
// dikonfirmasi, bukan cuma ditebak dari nama field:
//   - Transaction "paid"      -> Order "paid" (pembayaran diterima, mulai diproses)
//   - Transaction "processing" -> TIDAK memaksa perubahan Order (tidak ada
//     padanan langsung; order tetap di status sebelumnya)
//   - Transaction "shipped"   -> Order "shipped"
//   - Transaction "completed" -> Order "completed"
//   - Transaction "cancelled" -> Order "cancelled" + RESTOCK
//   - Transaction "refunded"  -> Order "cancelled" + RESTOCK (Order model
//     tidak punya status "refunded" tersendiri; dipetakan ke "cancelled"
//     karena efek stock-nya sama — barang kembali ke gudang)
//
// Status yang tidak ada di map ini (mis. "waiting_payment") sengaja tidak
// memicu perubahan apa pun ke Order.
var mapTransactionStatusToOrderStatus = map[string]string{
	"paid":      "paid",
	"shipped":   "shipped",
	"completed": "completed",
	"cancelled": "cancelled",
	"refunded":  "cancelled",
}

// restockStatuses adalah status Transaction yang, saat dipetakan ke Order,
// juga harus mengembalikan stock size-variant terkait (uang tidak jadi/
// dikembalikan -> barang tidak jadi terjual -> stock harus balik).
var restockStatuses = map[string]bool{
	"cancelled": true,
	"refunded":  true,
}

func (s *OrderFulfillmentServiceImpl) SyncOrdersToTransactionStatus(tx *gorm.DB, transactionID string, newTransactionStatus string) error {
	newOrderStatus, shouldSync := mapTransactionStatusToOrderStatus[newTransactionStatus]
	if !shouldSync {
		// Status ini sengaja tidak punya padanan Order (lihat komentar di
		// atas) — tidak ada yang perlu dilakukan.
		return nil
	}

	orders, err := s.orderRepo.FindAllByTxIdLocking(tx, transactionID)
	if err != nil {
		return fmt.Errorf("gagal mengambil orders untuk transaction %s: %w", transactionID, err)
	}

	needsRestock := restockStatuses[newTransactionStatus]

	for _, order := range orders {
		if orderStatusTerminal[order.Status] {
			// Order sudah final (completed/cancelled) — jangan disentuh lagi.
			continue
		}
		if order.Status == newOrderStatus {
			continue
		}

		if needsRestock {
			sizeVariant, err := s.productRepo.FindSizeVarianLocked(tx, uint(order.SizeVarianID))
			if err != nil {
				return fmt.Errorf("gagal mengambil size variant untuk restock order %s: %w", order.ID, err)
			}
			sizeVariant.Stock += order.Quantity
			if _, err := s.productRepo.UpdateSizeVarian(*sizeVariant, tx); err != nil {
				return fmt.Errorf("gagal restock order %s: %w", order.ID, err)
			}
		}

		order.Status = newOrderStatus
		if _, err := s.orderRepo.Update(order, tx); err != nil {
			return fmt.Errorf("gagal update status order %s: %w", order.ID, err)
		}
	}

	return nil
}
