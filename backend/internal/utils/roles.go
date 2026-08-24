package utils

// Role names — SATU-SATUNYA sumber kebenaran untuk nama role di seluruh
// aplikasi. Sebelumnya nama role ditulis sebagai string literal tersebar di
// banyak file (seeder, rbac_service, middleware, auth_service) dengan
// casing yang tidak konsisten ("admin" vs "Admin" vs "Super Admin"), yang
// menyebabkan pengecekan role gagal secara diam-diam. Semua kode yang perlu
// merujuk nama role WAJIB memakai konstanta di bawah ini, bukan string
// literal baru.
const (
	RoleSuperAdmin = "super_admin"
	RoleAdmin      = "admin"
	RoleVendor     = "vendor"
	RoleCustomer   = "customer"
)

// RoleHierarchyLevel mengembalikan level hierarki role (semakin besar semakin
// tinggi kewenangannya). Dipakai oleh RBACService.HasRole untuk perbandingan
// "role user >= role yang dibutuhkan".
func RoleHierarchyLevel(roleName string) int {
	switch roleName {
	case RoleSuperAdmin:
		return 4
	case RoleAdmin:
		return 3
	case RoleVendor:
		return 2
	case RoleCustomer:
		return 1
	default:
		return 0
	}
}
