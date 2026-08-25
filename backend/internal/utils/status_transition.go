package utils

import "fmt"

// StatusTransition adalah state-machine generik untuk validasi transisi
// status. SEBELUMNYA logic ini ditulis dua kali dengan cara berbeda:
// IsValidStatusTransition (transaction, di constant_helper.go) dan
// PaymentServiceImpl.validateStatusTransition (payment, private method) —
// sama-sama "map status -> status yang diizinkan", tapi ditulis terpisah,
// tidak dites bareng, dan gampang divergen kalau salah satu diupdate dan
// yang lain lupa. Sekarang keduanya dibangun dari satu tipe yang sama;
// masing-masing entity tetap punya vocabulary status sendiri (Transaction
// dan Payment memang dua state machine yang berbeda), tapi logic
// validasinya satu.
type StatusTransition struct {
	allowed map[string][]string
}

// NewStatusTransition membuat state machine dari map status -> daftar status
// tujuan yang diizinkan.
func NewStatusTransition(allowed map[string][]string) *StatusTransition {
	return &StatusTransition{allowed: allowed}
}

// IsValid mengecek apakah transisi from -> to diizinkan.
func (s *StatusTransition) IsValid(from, to string) bool {
	allowed, ok := s.allowed[from]
	if !ok {
		return false
	}
	for _, a := range allowed {
		if a == to {
			return true
		}
	}
	return false
}

// Validate sama seperti IsValid tapi mengembalikan error yang siap dipakai
// caller. Pesan errornya SENGAJA dipertahankan sama persis dengan
// implementasi lama (payment_service.go) karena beberapa handler (lihat
// payment_handler.go) mengecek substring "invalid status transition" di
// pesan error untuk menentukan HTTP 400 vs 500 — mengubah kata-katanya akan
// diam-diam mengubah response code tanpa ada yang sadar.
func (s *StatusTransition) Validate(from, to string) error {
	if _, ok := s.allowed[from]; !ok {
		return fmt.Errorf("invalid current status: %s", from)
	}
	if !s.IsValid(from, to) {
		return fmt.Errorf("invalid status transition from %s to %s", from, to)
	}
	return nil
}
