package utils

import (
	"net/http"
	"time"
)

// RefreshTokenCookieName adalah nama cookie tempat refresh token disimpan.
const RefreshTokenCookieName = "refresh_token"

// SetRefreshTokenCookie menyimpan refresh token sebagai httpOnly cookie —
// TIDAK BISA dibaca lewat JavaScript (document.cookie), jadi kalau ada XSS
// yang berhasil inject script ke halaman, script itu tidak bisa mencuri
// refresh token untuk dipakai ambil-alih sesi di luar browser korban.
//
// SEBELUMNYA refresh token dikembalikan di JSON response dan frontend
// menyimpannya di localStorage — bisa dibaca skrip apa pun yang jalan di
// halaman itu (termasuk dependency yang di-compromise atau XSS lewat konten
// user). secure menentukan apakah cookie cuma dikirim lewat HTTPS; harus
// true di production (lihat cfg.Server.Env), boleh false untuk dev lokal
// HTTP biasa.
func SetRefreshTokenCookie(w http.ResponseWriter, token string, expiresAt time.Time, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

// ClearRefreshTokenCookie menghapus cookie refresh token (dipakai saat logout).
func ClearRefreshTokenCookie(w http.ResponseWriter, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

// GetRefreshTokenFromCookie membaca refresh token dari cookie httpOnly.
func GetRefreshTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(RefreshTokenCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
