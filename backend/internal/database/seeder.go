package database

import (
	"e-commerce/backend/internal/config"
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/utils"
	"fmt"
	"log"
	"time"
)

func SeedDatabase(cfg *config.Config) error {
	// Check if permissions already exist (simple check without seed tracker)
	var permissionCount int64
	DB.Model(&models.Permission{}).Count(&permissionCount)
	if permissionCount > 0 {
		log.Println("Database already seeded, skipping seeding process")
		return nil
	}

	log.Println("Starting database seeding process...")

	if err := seedPermissions(); err != nil {
		return err
	}

	if err := seedRoles(); err != nil {
		return err
	}

	if err := seedDefaultAdmin(cfg); err != nil {
		return err
	}

	if err := seedCategories(); err != nil {
		return err
	}

	if err := seedProducts(); err != nil {
		return err
	}

	if err := seedPaymentMethods(); err != nil {
		return err
	}

	if err := seedShippingMethods(); err != nil {
		return err
	}

	log.Println("Database seeding completed successfully")
	return nil
}

func seedPermissions() error {
	permissions := []models.Permission{
		// User permissions
		{Name: "users.create", Resource: "users", Action: "create", Description: "Create new users"},
		{Name: "users.read", Resource: "users", Action: "read", Description: "View users"},
		{Name: "users.update", Resource: "users", Action: "update", Description: "Update user information"},
		{Name: "users.delete", Resource: "users", Action: "delete", Description: "Delete users"},

		// Role permissions
		{Name: "roles.create", Resource: "roles", Action: "create", Description: "Create new roles"},
		{Name: "roles.read", Resource: "roles", Action: "read", Description: "View roles"},
		{Name: "roles.update", Resource: "roles", Action: "update", Description: "Update role information"},
		{Name: "roles.delete", Resource: "roles", Action: "delete", Description: "Delete roles"},

		// Permission permissions
		{Name: "permissions.read", Resource: "permissions", Action: "read", Description: "View permissions"},

		// Profile permissions
		{Name: "profile.read", Resource: "profile", Action: "read", Description: "View own profile"},
		{Name: "profile.update", Resource: "profile", Action: "update", Description: "Update own profile"},

		// Product permissions
		{Name: "products.create", Resource: "products", Action: "create", Description: "Create new products"},
		{Name: "products.read", Resource: "products", Action: "read", Description: "View products"},
		{Name: "products.update", Resource: "products", Action: "update", Description: "Update product information"},
		{Name: "products.delete", Resource: "products", Action: "delete", Description: "Delete products"},

		// Category permissions
		{Name: "categories.create", Resource: "categories", Action: "create", Description: "Create new categories"},
		{Name: "categories.read", Resource: "categories", Action: "read", Description: "View categories"},
		{Name: "categories.update", Resource: "categories", Action: "update", Description: "Update category information"},
		{Name: "categories.delete", Resource: "categories", Action: "delete", Description: "Delete categories"},

		// Order permissions
		{Name: "orders.create", Resource: "orders", Action: "create", Description: "Create new orders"},
		{Name: "orders.read", Resource: "orders", Action: "read", Description: "View orders"},
		{Name: "orders.update", Resource: "orders", Action: "update", Description: "Update order information"},
		{Name: "orders.delete", Resource: "orders", Action: "delete", Description: "Delete orders"},
		{Name: "orders.read_own", Resource: "orders", Action: "read_own", Description: "View own orders"},

		// Payment Method permissions
		{Name: "payment_methods.create", Resource: "payment_methods", Action: "create", Description: "Create payment methods"},
		{Name: "payment_methods.read", Resource: "payment_methods", Action: "read", Description: "View payment methods"},
		{Name: "payment_methods.update", Resource: "payment_methods", Action: "update", Description: "Update payment methods"},
		{Name: "payment_methods.delete", Resource: "payment_methods", Action: "delete", Description: "Delete payment methods"},

		// Shipping permissions
		{Name: "shipping.create", Resource: "shipping", Action: "create", Description: "Create shipping methods"},
		{Name: "shipping.read", Resource: "shipping", Action: "read", Description: "View shipping methods"},
		{Name: "shipping.update", Resource: "shipping", Action: "update", Description: "Update shipping methods"},
		{Name: "shipping.delete", Resource: "shipping", Action: "delete", Description: "Delete shipping methods"},

		// Transaction permissions
		{Name: "transactions.create", Resource: "transactions", Action: "create", Description: "Create transactions"},
		{Name: "transactions.read", Resource: "transactions", Action: "read", Description: "View transactions"},
		{Name: "transactions.update", Resource: "transactions", Action: "update", Description: "Update transactions"},
		{Name: "transactions.read_own", Resource: "transactions", Action: "read_own", Description: "View own transactions"},

		// Payment permissions
		{Name: "payments.create", Resource: "payments", Action: "create", Description: "Create payments"},
		{Name: "payments.read", Resource: "payments", Action: "read", Description: "View payments"},
		{Name: "payments.update", Resource: "payments", Action: "update", Description: "Update payments"},
		{Name: "payments.delete", Resource: "payments", Action: "delete", Description: "Delete payments"},
		{Name: "payments.read_own", Resource: "payments", Action: "read_own", Description: "View own payments"},

		// Dashboard & Analytics
		{Name: "dashboard.read", Resource: "dashboard", Action: "read", Description: "View dashboard"},
		{Name: "analytics.read", Resource: "analytics", Action: "read", Description: "View analytics"},
		{Name: "activity_logs.read", Resource: "activity_logs", Action: "read", Description: "View activity logs"},
	}

	for _, permission := range permissions {
		var existingPermission models.Permission
		err := DB.Where("name = ?", permission.Name).First(&existingPermission).Error
		if err != nil {
			if err := DB.Create(&permission).Error; err != nil {
				log.Printf("Error creating permission %s: %v", permission.Name, err)
				return err
			}
			log.Printf("Created permission: %s", permission.Name)
		} else {
			log.Printf("Permission already exists: %s", permission.Name)
		}
	}

	return nil
}

func seedRoles() error {
	roles := []models.Role{
		{Name: "Super Admin", Description: "Full system access with all permissions", IsSystemRole: true},
		{Name: "Admin", Description: "Administrative access to manage store", IsSystemRole: true},
		{Name: "Vendor", Description: "Vendor access to manage own products", IsSystemRole: true},
		{Name: "Customer", Description: "Customer access to browse and purchase", IsSystemRole: true},
	}

	rolePermissions := map[string][]string{
		"Super Admin": {
			"users.create", "users.read", "users.update", "users.delete",
			"roles.create", "roles.read", "roles.update", "roles.delete",
			"permissions.read", "profile.read", "profile.update",
			"products.create", "products.read", "products.update", "products.delete",
			"categories.create", "categories.read", "categories.update", "categories.delete",
			"orders.create", "orders.read", "orders.update", "orders.delete",
			"payment_methods.create", "payment_methods.read", "payment_methods.update", "payment_methods.delete",
			"shipping.create", "shipping.read", "shipping.update", "shipping.delete",
			"transactions.create", "transactions.read", "transactions.update",
			"payments.create", "payments.read", "payments.update", "payments.delete",
			"dashboard.read", "analytics.read", "activity_logs.read",
		},
		"Admin": {
			"users.read", "users.update",
			"roles.read", "permissions.read", "profile.read", "profile.update",
			"products.create", "products.read", "products.update", "products.delete",
			"categories.create", "categories.read", "categories.update", "categories.delete",
			"orders.read", "orders.update",
			"payment_methods.read", "payment_methods.update",
			"shipping.read", "shipping.update",
			"transactions.read", "transactions.update",
			"payments.read", "payments.update",
			"dashboard.read", "analytics.read", "activity_logs.read",
		},
		"Vendor": {
			"profile.read", "profile.update",
			"products.create", "products.read", "products.update",
			"categories.read",
			"orders.read",
			"payment_methods.read",
			"shipping.read",
			"transactions.read",
			"payments.read",
			"dashboard.read",
		},
		"Customer": {
			"profile.read", "profile.update",
			"products.read", "categories.read",
			"orders.create", "orders.read_own",
			"payment_methods.read",
			"shipping.read",
			"transactions.create", "transactions.read_own",
			"payments.create", "payments.read_own",
		},
	}

	for _, role := range roles {
		var existingRole models.Role
		if err := DB.Where("name = ?", role.Name).First(&existingRole).Error; err != nil {
			if err := DB.Create(&role).Error; err != nil {
				return err
			}

			var permissions []models.Permission
			permissionNames := rolePermissions[role.Name]
			if err := DB.Where("name IN ?", permissionNames).Find(&permissions).Error; err != nil {
				return err
			}

			if err := DB.Model(&role).Association("Permissions").Append(&permissions); err != nil {
				return err
			}

			log.Printf("Created role: %s with %d permissions", role.Name, len(permissions))
		}
	}

	return nil
}

func seedDefaultAdmin(cfg *config.Config) error {
	if cfg.System.DefaultAdminEmail == "" || cfg.System.DefaultAdminPassword == "" {
		log.Println("Skipping default admin creation - credentials not provided")
		return nil
	}

	var existingAdmin models.User
	if err := DB.Where("email = ?", cfg.System.DefaultAdminEmail).First(&existingAdmin).Error; err == nil {
		log.Println("Default admin already exists")
		return nil
	}

	var superAdminRole models.Role
	if err := DB.Where("name = ?", "Super Admin").First(&superAdminRole).Error; err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(cfg.System.DefaultAdminPassword)
	if err != nil {
		return err
	}

	now := time.Now()
	admin := models.User{
		Email:           cfg.System.DefaultAdminEmail,
		Username:        "admin",
		PasswordHash:    hashedPassword,
		FirstName:       "System",
		LastName:        "Administrator",
		RoleID:          superAdminRole.ID,
		IsActive:        true,
		EmailVerifiedAt: &now,
	}

	if err := DB.Create(&admin).Error; err != nil {
		return err
	}

	log.Printf("Created default admin user: %s", admin.Email)
	return nil
}

func seedCategories() error {
	categories := []models.Category{
		{Name: "Gamis", Icon: "https://example.com/icons/gamis.png"},
		{Name: "Hijab", Icon: "https://example.com/icons/hijab.png"},
		{Name: "Kemeja Wanita", Icon: "https://example.com/icons/kemeja-wanita.png"},
		{Name: "Dress", Icon: "https://example.com/icons/dress.png"},
		{Name: "Tunik", Icon: "https://example.com/icons/tunik.png"},
		{Name: "Celana Wanita", Icon: "https://example.com/icons/celana.png"},
		{Name: "Aksesoris Muslim", Icon: "https://example.com/icons/aksesoris.png"},
		{Name: "Outer & Cardigan", Icon: "https://example.com/icons/outer.png"},
	}

	for _, category := range categories {
		var existingCategory models.Category
		err := DB.Where("name = ?", category.Name).First(&existingCategory).Error
		if err != nil {
			// Jika record not found (belum ada), create baru
			if err := DB.Create(&category).Error; err != nil {
				log.Printf("Error creating category %s: %v", category.Name, err)
				return err
			}
			log.Printf("Created category: %s", category.Name)
		} else {
			// Jika sudah ada, skip
			log.Printf("Category already exists: %s", category.Name)
		}
	}

	return nil
}

// seedColorVariants is now integrated with products
// Colors are created as part of product variants

// seedSizeVariants is now integrated with color variants
// Sizes are created as part of color variant sizes

func seedProducts() error {
	// Get categories
	var gamisCategory, hijabCategory, kemejaCategory, dressCategory models.Category
	DB.Where("name = ?", "Gamis").First(&gamisCategory)
	DB.Where("name = ?", "Hijab").First(&hijabCategory)
	DB.Where("name = ?", "Kemeja Wanita").First(&kemejaCategory)
	DB.Where("name = ?", "Dress").First(&dressCategory)

	// Product 1: Gamis Syari Premium
	var existingProduct1 models.Product
	if err := DB.Where("name = ?", "Gamis Syari Premium").First(&existingProduct1).Error; err != nil {
		product1 := models.Product{
			Name:        "Gamis Syari Premium",
			Description: "Gamis syari berbahan wolfis premium, nyaman dipakai seharian. Cocok untuk acara formal maupun santai",
			Price:       275000,
			CategoryID:  gamisCategory.ID,
			Images:      "https://example.com/products/gamis-syari.jpg",
			Rating:      4.8,
		}
		if err := DB.Create(&product1).Error; err != nil {
			return err
		}

		colors := []struct {
			Name  string
			Color string
		}{
			{"Navy", "#000080"},
			{"Maroon", "#800000"},
			{"Dark Green", "#006400"},
			{"Black", "#000000"},
			{"Dusty Pink", "#DCAE96"},
		}

		sizes := []string{"S", "M", "L", "XL", "XXL"}

		for _, c := range colors {
			colorVariant := models.ColorVarian{
				ProductID: product1.ID,
				Name:      c.Name,
				Color:     c.Color,
				Images:    fmt.Sprintf("https://example.com/products/gamis-syari-%s.jpg", c.Name),
			}
			if err := DB.Create(&colorVariant).Error; err != nil {
				return err
			}

			for _, size := range sizes {
				sizeVariant := models.SizeVarian{
					ColorVarianID: colorVariant.ID,
					Size:          size,
					Stock:         30,
				}
				if err := DB.Create(&sizeVariant).Error; err != nil {
					return err
				}
			}
		}

		log.Printf("Created product: %s with variants", product1.Name)
	}

	// Product 2: Gamis Set Busui Friendly
	var existingProduct2 models.Product
	if err := DB.Where("name = ?", "Gamis Set Busui Friendly").First(&existingProduct2).Error; err != nil {
		product2 := models.Product{
			Name:        "Gamis Set Busui Friendly",
			Description: "Gamis set dengan khimar, busui friendly dengan resleting depan. Bahan katun rayon adem dan tidak mudah kusut",
			Price:       320000,
			CategoryID:  gamisCategory.ID,
			Images:      "https://example.com/products/gamis-busui.jpg",
			Rating:      4.9,
		}
		if err := DB.Create(&product2).Error; err != nil {
			return err
		}

		colors := []struct {
			Name  string
			Color string
		}{
			{"Mint Green", "#98FF98"},
			{"Baby Blue", "#89CFF0"},
			{"Lavender", "#E6E6FA"},
			{"Cream", "#FFFDD0"},
			{"Soft Grey", "#D3D3D3"},
		}

		sizes := []string{"M", "L", "XL", "XXL", "JUMBO"}

		for _, c := range colors {
			colorVariant := models.ColorVarian{
				ProductID: product2.ID,
				Name:      c.Name,
				Color:     c.Color,
				Images:    fmt.Sprintf("https://example.com/products/gamis-busui-%s.jpg", c.Name),
			}
			if err := DB.Create(&colorVariant).Error; err != nil {
				return err
			}

			for _, size := range sizes {
				sizeVariant := models.SizeVarian{
					ColorVarianID: colorVariant.ID,
					Size:          size,
					Stock:         25,
				}
				if err := DB.Create(&sizeVariant).Error; err != nil {
					return err
				}
			}
		}

		log.Printf("Created product: %s with variants", product2.Name)
	}

	// Product 3: Hijab Voal Premium
	var existingProduct3 models.Product
	if err := DB.Where("name = ?", "Hijab Voal Premium").First(&existingProduct3).Error; err != nil {
		product3 := models.Product{
			Name:        "Hijab Voal Premium",
			Description: "Hijab voal import premium, bahan lembut dan adem, tidak licin. Ukuran 115x115cm",
			Price:       45000,
			CategoryID:  hijabCategory.ID,
			Images:      "https://example.com/products/hijab-voal.jpg",
			Rating:      4.7,
		}
		if err := DB.Create(&product3).Error; err != nil {
			return err
		}

		colors := []struct {
			Name  string
			Color string
		}{
			{"Black", "#000000"},
			{"White", "#FFFFFF"},
			{"Navy", "#000080"},
			{"Maroon", "#800000"},
			{"Dusty Pink", "#DCAE96"},
			{"Milo", "#C19A6B"},
			{"Army Green", "#4B5320"},
			{"Chocolate", "#D2691E"},
		}

		for _, c := range colors {
			colorVariant := models.ColorVarian{
				ProductID: product3.ID,
				Name:      c.Name,
				Color:     c.Color,
				Images:    fmt.Sprintf("https://example.com/products/hijab-voal-%s.jpg", c.Name),
			}
			if err := DB.Create(&colorVariant).Error; err != nil {
				return err
			}

			sizeVariant := models.SizeVarian{
				ColorVarianID: colorVariant.ID,
				Size:          "115x115",
				Stock:         50,
			}
			if err := DB.Create(&sizeVariant).Error; err != nil {
				return err
			}
		}

		log.Printf("Created product: %s with variants", product3.Name)
	}

	// Product 4: Pashmina Diamond Italiano
	var existingProduct4 models.Product
	if err := DB.Where("name = ?", "Pashmina Diamond Italiano").First(&existingProduct4).Error; err != nil {
		product4 := models.Product{
			Name:        "Pashmina Diamond Italiano",
			Description: "Pashmina diamond italiano, bahan premium tidak menerawang, tekstur diamond yang elegan. Ukuran 75x200cm",
			Price:       65000,
			CategoryID:  hijabCategory.ID,
			Images:      "https://example.com/products/pashmina-diamond.jpg",
			Rating:      4.8,
		}
		if err := DB.Create(&product4).Error; err != nil {
			return err
		}

		colors := []struct {
			Name  string
			Color string
		}{
			{"Black", "#000000"},
			{"Navy", "#000080"},
			{"Dark Grey", "#696969"},
			{"Caramel", "#C68E17"},
			{"Terracotta", "#E2725B"},
			{"Emerald", "#50C878"},
		}

		for _, c := range colors {
			colorVariant := models.ColorVarian{
				ProductID: product4.ID,
				Name:      c.Name,
				Color:     c.Color,
				Images:    fmt.Sprintf("https://example.com/products/pashmina-diamond-%s.jpg", c.Name),
			}
			if err := DB.Create(&colorVariant).Error; err != nil {
				return err
			}

			sizeVariant := models.SizeVarian{
				ColorVarianID: colorVariant.ID,
				Size:          "75x200",
				Stock:         45,
			}
			if err := DB.Create(&sizeVariant).Error; err != nil {
				return err
			}
		}

		log.Printf("Created product: %s with variants", product4.Name)
	}

	// Product 5: Kemeja Wanita Premium
	var existingProduct5 models.Product
	if err := DB.Where("name = ?", "Kemeja Wanita Premium").First(&existingProduct5).Error; err != nil {
		product5 := models.Product{
			Name:        "Kemeja Wanita Premium",
			Description: "Kemeja wanita berbahan katun supernova, nyaman dan tidak gerah. Cocok untuk kerja dan hangout",
			Price:       135000,
			CategoryID:  kemejaCategory.ID,
			Images:      "https://example.com/products/kemeja-wanita.jpg",
			Rating:      4.6,
		}
		if err := DB.Create(&product5).Error; err != nil {
			return err
		}

		colors := []struct {
			Name  string
			Color string
		}{
			{"White", "#FFFFFF"},
			{"Black", "#000000"},
			{"Soft Blue", "#89CFF0"},
			{"Beige", "#F5F5DC"},
			{"Pink", "#FFC0CB"},
		}

		sizes := []string{"S", "M", "L", "XL"}

		for _, c := range colors {
			colorVariant := models.ColorVarian{
				ProductID: product5.ID,
				Name:      c.Name,
				Color:     c.Color,
				Images:    fmt.Sprintf("https://example.com/products/kemeja-wanita-%s.jpg", c.Name),
			}
			if err := DB.Create(&colorVariant).Error; err != nil {
				return err
			}

			for _, size := range sizes {
				sizeVariant := models.SizeVarian{
					ColorVarianID: colorVariant.ID,
					Size:          size,
					Stock:         35,
				}
				if err := DB.Create(&sizeVariant).Error; err != nil {
					return err
				}
			}
		}

		log.Printf("Created product: %s with variants", product5.Name)
	}

	// Product 6: Dress Casual Wanita
	var existingProduct6 models.Product
	if err := DB.Where("name = ?", "Dress Casual Wanita").First(&existingProduct6).Error; err != nil {
		product6 := models.Product{
			Name:        "Dress Casual Wanita",
			Description: "Dress casual dengan model trendy, bahan katun rayon yang adem dan nyaman. Bisa untuk daily maupun hangout",
			Price:       165000,
			CategoryID:  dressCategory.ID,
			Images:      "https://example.com/products/dress-casual.jpg",
			Rating:      4.7,
		}
		if err := DB.Create(&product6).Error; err != nil {
			return err
		}

		colors := []struct {
			Name  string
			Color string
		}{
			{"Navy Floral", "#000080"},
			{"Maroon Polka", "#800000"},
			{"Sage Green", "#87AE73"},
			{"Mustard", "#FFDB58"},
			{"Coral", "#FF7F50"},
		}

		sizes := []string{"M", "L", "XL", "XXL"}

		for _, c := range colors {
			colorVariant := models.ColorVarian{
				ProductID: product6.ID,
				Name:      c.Name,
				Color:     c.Color,
				Images:    fmt.Sprintf("https://example.com/products/dress-casual-%s.jpg", c.Name),
			}
			if err := DB.Create(&colorVariant).Error; err != nil {
				return err
			}

			for _, size := range sizes {
				sizeVariant := models.SizeVarian{
					ColorVarianID: colorVariant.ID,
					Size:          size,
					Stock:         28,
				}
				if err := DB.Create(&sizeVariant).Error; err != nil {
					return err
				}
			}
		}

		log.Printf("Created product: %s with variants", product6.Name)
	}

	return nil
}

func seedPaymentMethods() error {
	paymentMethods := []models.PaymentMethod{
		{
			AccountName:   "PT E-Commerce Indonesia",
			AccountNumber: "1234567890",
			BankName:      "Bank BCA",
			BankImages:    "https://example.com/banks/bca.png",
		},
		{
			AccountName:   "PT E-Commerce Indonesia",
			AccountNumber: "9876543210",
			BankName:      "Bank Mandiri",
			BankImages:    "https://example.com/banks/mandiri.png",
		},
		{
			AccountName:   "PT E-Commerce Indonesia",
			AccountNumber: "5555666677",
			BankName:      "Bank BRI",
			BankImages:    "https://example.com/banks/bri.png",
		},
		{
			AccountName:   "PT E-Commerce Indonesia",
			AccountNumber: "1111222233",
			BankName:      "Bank BNI",
			BankImages:    "https://example.com/banks/bni.png",
		},
		{
			AccountName:   "PT E-Commerce Indonesia",
			AccountNumber: "7777888899",
			BankName:      "Bank CIMB Niaga",
			BankImages:    "https://example.com/banks/cimb.png",
		},
	}

	for _, method := range paymentMethods {
		var existingMethod models.PaymentMethod
		err := DB.Where("account_number = ? AND bank_name = ?", method.AccountNumber, method.BankName).First(&existingMethod).Error
		if err != nil {
			// Jika record not found (belum ada), create baru
			if err := DB.Create(&method).Error; err != nil {
				log.Printf("Error creating payment method %s: %v", method.BankName, err)
				return err
			}
			log.Printf("Created payment method: %s - %s", method.BankName, method.AccountNumber)
		} else {
			// Jika sudah ada, skip
			log.Printf("Payment method already exists: %s - %s", method.BankName, method.AccountNumber)
		}
	}

	return nil
}

func seedShippingMethods() error {
	shippingMethods := []models.Shipping{
		{
			Name:  "JNE Regular",
			Price: 15000,
			State: "active",
		},
		{
			Name:  "JNE Express",
			Price: 25000,
			State: "active",
		},
		{
			Name:  "J&T Regular",
			Price: 12000,
			State: "active",
		},
		{
			Name:  "J&T Express",
			Price: 20000,
			State: "active",
		},
		{
			Name:  "SiCepat Regular",
			Price: 13000,
			State: "active",
		},
		{
			Name:  "SiCepat Express",
			Price: 22000,
			State: "active",
		},
		{
			Name:  "Anteraja Regular",
			Price: 10000,
			State: "active",
		},
		{
			Name:  "Anteraja Same Day",
			Price: 30000,
			State: "active",
		},
	}

	for _, method := range shippingMethods {
		var existingMethod models.Shipping
		// Gunakan error check yang benar
		err := DB.Where("name = ?", method.Name).First(&existingMethod).Error
		if err != nil {
			// Jika record not found (belum ada), create baru
			if err := DB.Create(&method).Error; err != nil {
				log.Printf("Error creating shipping method %s: %v", method.Name, err)
				return err
			}
			log.Printf("Created shipping method: %s", method.Name)
		} else {
			// Jika sudah ada, skip
			log.Printf("Shipping method already exists: %s", method.Name)
		}
	}

	return nil
}

// ResetSeedingStatus resets the seeding status - useful for fresh migrations
func ResetSeedingStatus() error {
	// Simply truncate permissions to trigger re-seeding
	if err := DB.Exec("TRUNCATE TABLE permissions CASCADE").Error; err != nil {
		return fmt.Errorf("failed to reset seeding status: %w", err)
	}
	log.Println("Seeding status reset successfully")
	return nil
}

// ForceSeedDatabase forces seeding even if already completed
func ForceSeedDatabase(cfg *config.Config) error {
	log.Println("Force seeding database...")

	if err := seedPermissions(); err != nil {
		return err
	}

	if err := seedRoles(); err != nil {
		return err
	}

	if err := seedDefaultAdmin(cfg); err != nil {
		return err
	}

	if err := seedCategories(); err != nil {
		return err
	}

	if err := seedProducts(); err != nil {
		return err
	}

	if err := seedPaymentMethods(); err != nil {
		return err
	}

	if err := seedShippingMethods(); err != nil {
		return err
	}

	log.Println("Force seeding completed successfully")
	return nil
}
