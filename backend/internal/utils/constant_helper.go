package utils

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func GetHeaderKey() string {
	return authorizationHeaderKey
}
func GetTypeBearer() string {
	return authorizationTypeBearer
}
func GetPayloadKey() string {
	return authorizationPayloadKey
}

const (
	WaitingPayment       = "waiting_payment"
	WaitingConfirPayment = "waiting_confirm_payment"
	Confirmed            = "confirmed"
)

const (
	Pending    = "pending"
	Packaged   = "packaged"
	OnDelivery = "on_delivery"
	Delivered  = "delivered"
	Received   = "received"
	Cancelled  = "cancelled"
)

func ValidateStatusOrder(item string) bool {
	listStatus := []string{OnDelivery, Delivered, Received, Cancelled}
	for _, v := range listStatus {
		if v == item {
			return true
		}
	}
	return false
}

// SEBELUMNYA map ini pakai key/value UPPERCASE ("WAITING_PAYMENT", "PAID",
// dst), padahal status transaksi yang benar-benar disimpan di database dan
// divalidasi oleh handler (lihat transaction_handler.go validStatuses) semua
// lowercase snake_case ("waiting_payment", "paid", dst). Karena map key tidak
// pernah match, IsValidStatusTransition SELALU mengembalikan false — artinya
// TIDAK ADA transisi status transaksi yang pernah berhasil lewat jalur ini.
// Sekarang disamakan dengan vocabulary yang benar-benar dipakai di seluruh
// aplikasi, dan dibangun dari StatusTransition (lihat status_transition.go)
// bukan diimplementasikan manual lagi.
var transactionStatusTransition = NewStatusTransition(map[string][]string{
	"waiting_payment": {"paid", "cancelled"},
	"paid":            {"processing", "refunded"},
	"processing":      {"shipped", "cancelled"},
	"shipped":         {"completed"},
	"completed":       {"refunded"},
	"cancelled":       {},
	"refunded":        {},
})

func IsValidStatusTransition(from, to string) bool {
	return transactionStatusTransition.IsValid(from, to)
}
