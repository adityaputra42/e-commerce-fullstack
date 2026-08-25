package utils

import (
	"e-commerce/backend/internal/config"
	"fmt"
	"log"
	"net/smtp"
)

// Mailer mengirim email keluar. Diabstraksi lewat interface supaya
// AuthService bisa di-unit-test tanpa benar-benar mengirim email
// (lihat NoopMailer di bawah untuk dev/test).
type Mailer interface {
	SendPasswordResetEmail(toEmail, toName, resetToken string) error
}

// SMTPMailer implementasi Mailer yang benar-benar kirim email lewat SMTP.
type SMTPMailer struct {
	cfg       config.SMTPConfig
	fromEmail string
	// ResetURLBase contoh: "https://toko-kamu.com/reset-password"
	// token akan ditempel sebagai query param: ?token=xxx
	ResetURLBase string
}

func NewSMTPMailer(cfg config.SMTPConfig, fromEmail, resetURLBase string) *SMTPMailer {
	return &SMTPMailer{cfg: cfg, fromEmail: fromEmail, ResetURLBase: resetURLBase}
}

func (m *SMTPMailer) SendPasswordResetEmail(toEmail, toName, resetToken string) error {
	if m.cfg.Host == "" || m.cfg.Port == "" {
		return fmt.Errorf("SMTP belum dikonfigurasi (SMTP_HOST/SMTP_PORT kosong)")
	}

	resetLink := fmt.Sprintf("%s?token=%s", m.ResetURLBase, resetToken)

	subject := "Reset Password Akun Anda"
	body := fmt.Sprintf(
		"Halo %s,\r\n\r\n"+
			"Kami menerima permintaan reset password untuk akun Anda.\r\n"+
			"Klik link berikut untuk mengatur password baru (berlaku 1 jam):\r\n\r\n"+
			"%s\r\n\r\n"+
			"Jika Anda tidak meminta ini, abaikan email ini dan password Anda tidak akan berubah.\r\n",
		toName, resetLink,
	)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		m.fromEmail, toEmail, subject, body)

	addr := fmt.Sprintf("%s:%s", m.cfg.Host, m.cfg.Port)
	auth := smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)

	if err := smtp.SendMail(addr, auth, m.fromEmail, []string{toEmail}, []byte(msg)); err != nil {
		return fmt.Errorf("gagal mengirim email reset password: %w", err)
	}

	return nil
}

// NoopMailer dipakai saat SMTP belum dikonfigurasi (mis. dev lokal) supaya
// aplikasi tetap jalan tanpa mengirim email sungguhan. Token TIDAK PERNAH
// dikembalikan ke caller HTTP — cuma dicatat di log server untuk keperluan
// development. JANGAN pakai ini di production.
type NoopMailer struct{}

func NewNoopMailer() *NoopMailer { return &NoopMailer{} }

func (m *NoopMailer) SendPasswordResetEmail(toEmail, toName, resetToken string) error {
	log.Printf("[DEV ONLY] Password reset token untuk %s: %s (SMTP belum dikonfigurasi, email tidak benar-benar terkirim)", toEmail, resetToken)
	return nil
}
