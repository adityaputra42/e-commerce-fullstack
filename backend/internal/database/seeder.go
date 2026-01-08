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
		{
			Name:         "super_admin",
			Description:  "Full system access with all permissions",
			Level:        4,
			IsSystemRole: true,
		},
		{
			Name:         "admin",
			Description:  "Administrative access to manage store",
			Level:        3,
			IsSystemRole: true,
		},
		{
			Name:         "vendor",
			Description:  "Vendor access to manage own products",
			Level:        2,
			IsSystemRole: true,
		},
		{
			Name:         "customer",
			Description:  "Customer access to browse and purchase",
			Level:        1,
			IsSystemRole: true,
		},
	}

	rolePermissions := map[string][]string{
		"super_admin": {
			"users.create", "users.read", "users.update", "users.delete",
			"roles.create", "roles.read", "roles.update", "roles.delete",
			"permissions.read",
			"profile.read", "profile.update",
			"products.create", "products.read", "products.update", "products.delete",
			"categories.create", "categories.read", "categories.update", "categories.delete",
			"orders.create", "orders.read", "orders.update", "orders.delete",
			"payment_methods.create", "payment_methods.read", "payment_methods.update", "payment_methods.delete",
			"shipping.create", "shipping.read", "shipping.update", "shipping.delete",
			"transactions.create", "transactions.read", "transactions.update",
			"payments.create", "payments.read", "payments.update", "payments.delete",
			"dashboard.read", "analytics.read", "activity_logs.read",
		},
		"admin": {
			"users.read", "users.update",
			"roles.read", "permissions.read",
			"profile.read", "profile.update",
			"products.create", "products.read", "products.update", "products.delete",
			"categories.create", "categories.read", "categories.update", "categories.delete",
			"orders.read", "orders.update",
			"payment_methods.read", "payment_methods.update",
			"shipping.read", "shipping.update",
			"transactions.read", "transactions.update",
			"payments.read", "payments.update",
			"dashboard.read", "analytics.read", "activity_logs.read",
		},
		"vendor": {
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
		"customer": {
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

		err := DB.Where("name = ?", role.Name).First(&existingRole).Error
		if err == nil {
			log.Printf("Role already exists: %s", role.Name)
			continue
		}

		if err := DB.Create(&role).Error; err != nil {
			return err
		}

		var permissions []models.Permission
		if err := DB.
			Where("name IN ?", rolePermissions[role.Name]).
			Find(&permissions).Error; err != nil {
			return err
		}

		if err := DB.Model(&role).
			Association("Permissions").
			Replace(&permissions); err != nil {
			return err
		}

		log.Printf(
			"Created role: %s (level=%d) with %d permissions",
			role.Name,
			role.Level,
			len(permissions),
		)
	}

	return nil
}

func seedDefaultAdmin(cfg *config.Config) error {
	if cfg.System.DefaultAdminEmail == "" || cfg.System.DefaultAdminPassword == "" {
		log.Println("Skipping default admin creation - credentials not provided")
		return nil
	}

	var existingAdmin models.User
	if err := DB.Where("email = ?", cfg.System.DefaultAdminEmail).
		First(&existingAdmin).Error; err == nil {
		log.Println("Default admin already exists")
		return nil
	}

	var superAdminRole models.Role
	if err := DB.Where("name = ?", "super_admin").
		First(&superAdminRole).Error; err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(cfg.System.DefaultAdminPassword)
	if err != nil {
		return err
	}

	now := time.Now()
	admin := models.User{
		Email:           cfg.System.DefaultAdminEmail,
		Username:        "superadmin",
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

	log.Printf("Created default super admin user: %s", admin.Email)
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
			Images:      "https://www.static-src.com/wcsstore/Indraprastha/images/catalog/full//92/MTA-75004802/no_brand_abaya_zahra_limitid_edition_full01_gdime8vw.jpg",
			Rating:      4.8,
		}
		if err := DB.Create(&product1).Error; err != nil {
			return err
		}

		colors := []struct {
			Name   string
			Color  string
			Images string
		}{
			{"Navy", "#000080", "https://www.static-src.com/wcsstore/Indraprastha/images/catalog/full/catalog-image/110/MTA-179396237/brd-44261_gamis-syari-set-hijab-khimar-ceruty-premium-by-andini-group-ori-jumbo-standar-seragaman-fashion-muslimah-terbaru-terlaris_full02-e1b33b18.jpg"},
			{"Maroon", "#800000", "https://img.lazcdn.com/g/p/68ed4884f2da0a5667007c7fb0ccbecc.jpg_720x720q80.jpg"},
			{"Dark Green", "#006400", "https://static.desty.app/desty-omni/20240325/80626df11feb4261a5b9dc1a756e4ca7.jpg?x-oss-process=image/format,webp"},
			{"Black", "#000000", "https://cdn1-production-images-kly.akamaized.net/U310EfKS_F_JJ0dd9-PC6MVGt0o=/500x667/smart/filters:quality(75):strip_icc()/kly-media-production/medias/5339829/original/067980900_1757131082-e0f38f98-d5fd-4d5f-bf8d-37f8bd2b2a87.jpg"},
			{"Dusty Pink", "#DCAE96", "https://img.lazcdn.com/g/ff/kf/S54750a774d5040bdb4956fa2f146ecc7o.jpg_720x720q80.jpg"},
		}

		sizes := []string{"S", "M", "L", "XL", "XXL"}

		for _, c := range colors {
			colorVariant := models.ColorVarian{
				ProductID: product1.ID,
				Name:      c.Name,
				Color:     c.Color,
				Images:    c.Images,
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
			Images:      "https://www.static-src.com/wcsstore/Indraprastha/images/catalog/full/catalog-image/106/MTA-157206948/br-m036969-01094_gamis-o-ring-o-ring-rempel-busui-friendly-cod-wanita-premium-berkualitas_full01-efdded63.jpg",
			Rating:      4.9,
		}
		if err := DB.Create(&product2).Error; err != nil {
			return err
		}

		colors := []struct {
			Name   string
			Color  string
			Images string
		}{
			{"Mint Green", "#98FF98", "https://img.lazcdn.com/g/p/1eafd95350a27c04bb378c1d4e2cf178.jpg_720x720q80.jpg"},
			{"Baby Blue", "#89CFF0", "https://ethica-collection.com/wp-content/uploads/2023/10/AYUMI-409-GREEN-WP-4-1024x1024.webp?w=1024&h=1024&q=90"},
			{"Lavender", "#E6E6FA", "https://p16-oec-va.ibyteimg.com/tos-maliva-i-o3syd03w52-us/36fffd0a36c749a499108b82b1b55bea~tplv-o3syd03w52-resize-webp:800:800.webp?dr=15584&t=555f072d&ps=933b5bde&shp=6ce186a1&shcp=e1be8f53&idc=my2&from=1826719393"},
			{"Cream", "#FFFDD0", "https://www.static-src.com/wcsstore/Indraprastha/images/catalog/full/catalog-image/110/MTA-179338844/no_brand_pussat_ashiya_one_set_gamis_syar-i_-_french_khimar_crinkle_airflow_premium_dress_syari_set_hijab_wanita_muslim_full01_b479cd14.jpg"},
			{"Soft Grey", "#D3D3D3", "https://p16-oec-sg.ibyteimg.com/tos-alisg-i-aphluv4xwc-sg/a73998fb39014bdbbfd3047977867d06~tplv-aphluv4xwc-resize-webp:800:800.webp?dr=15584&t=555f072d&ps=933b5bde&shp=6ce186a1&shcp=e1be8f53&idc=my2&from=1826719393"},
		}

		sizes := []string{"M", "L", "XL", "XXL", "JUMBO"}

		for _, c := range colors {
			colorVariant := models.ColorVarian{
				ProductID: product2.ID,
				Name:      c.Name,
				Color:     c.Color,
				Images:    c.Images,
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
			Images:      "https://www.umamascarves.co.id/wp-content/uploads/2024/04/Cover-4.jpg",
			Rating:      4.7,
		}
		if err := DB.Create(&product3).Error; err != nil {
			return err
		}

		colors := []struct {
			Name  string
			Color string
			Image string
		}{
			{"Black", "#000000", "https://www.static-src.com/wcsstore/Indraprastha/images/catalog/full//103/MTA-87495151/no_brand_jilbab_hijab_kerudung_muslimah_segiempat_voal_paris_premium_110x110cm_full03_4fc403fe.jpg"},
			{"White", "#FFFFFF", "https://img.lazcdn.com/g/ff/kf/S7df725e9786e4362bd5b61e545435c15b.jpg_720x720q80.jpg"},
			{"Navy", "#000080", "https://www.static-src.com/wcsstore/Indraprastha/images/catalog/full/catalog-image/102/MTA-174474074/authentism_kerudung-segi-empat-polos-voal-ultrafine-premium-superfine-hijab-authentism-rachita-voal_full02.jpg"},
			{"Maroon", "#800000", "https://www.static-src.com/wcsstore/Indraprastha/images/catalog/full/catalog-image/112/MTA-179085936/alya-hijab-by-naja_adzana-voal-segiempat-alyahijabbynaja_full07.jpg"},
			{"Dusty Pink", "#DCAE96", "https://www.hijabwanitacantik.com/cdn/shop/files/GalaxySky_aca3a48a-d4c6-4c94-ac05-ce6ddaec3311_grande.jpg?v=1756709659"},
			{"Milo", "#C19A6B", "https://www.static-src.com/wcsstore/Indraprastha/images/catalog/full/catalog-image/MTA-148317862/lozy_hijab_lozy_hijab_-_kirana_paris_plain_japan_milo_full01_e8k12c2d.jpg"},
			{"Army Green", "#4B5320", "https://lozy.id/cdn/shop/files/SQUARE21_9eb64f6e-a8fd-40fa-83c3-38f06670727b_800x.jpg?v=1710404451"},
			{"Chocolate", "#D2691E", "https://img.id.my-best.com/product_images/7654c5bac371cbcc17af8ae3923368c4.png?ixlib=rails-4.3.1&q=70&lossless=0&w=800&h=800&fit=clip&s=36308e5f425dbea19cd5f7b9ba694fe2"},
		}

		for _, c := range colors {
			colorVariant := models.ColorVarian{
				ProductID: product3.ID,
				Name:      c.Name,
				Color:     c.Color,
				Images:    c.Image,
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
			Images:      "https://admincerdas.s3.ap-southeast-1.amazonaws.com/20201128/15515-78e6c6bf5b353eb5e89e53ba6bfe104b.jpg",
			Rating:      4.8,
		}
		if err := DB.Create(&product4).Error; err != nil {
			return err
		}

		colors := []struct {
			Name  string
			Color string
			Image string
		}{
			{"Black", "#000000", "https://s3.belanjapasti.com/media/image/pashmina-diamond-crepe-italiano-part-iii-655524.jpg"},
			{"Navy", "#000080", "https://www.jagel.id/api/listimage/v/Hijab-Pashmina-Sabyan--Diamond-Italiano-Premium--Hijab-1-878605996c369f80.jpg"},
			{"Dark Grey", "#696969", "https://s3-ap-southeast-1.amazonaws.com/plugolive/vendor/542/product/photo_2023-04-18_10-46-12_1681790127141.jpg"},
			{"Caramel", "#C68E17", "https://admincerdas.s3.ap-southeast-1.amazonaws.com/20201128/15515-6ece68c8a34e56f0eb57b00bfa65f8bc.jpg"},
			{"Terracotta", "#E2725B", "https://www.static-src.com/wcsstore/Indraprastha/images/catalog/full//91/MTA-10981431/elzatta_pashmina_polos_elzatta_selendang_selvia_anindita_full12_gqosovae.jpg"},
			{"Emerald", "#50C878", "https://admincerdas.s3.ap-southeast-1.amazonaws.com/20201128/15515-7816e646dea786903d96eecf8ca4cddd.jpg"},
		}

		for _, c := range colors {
			colorVariant := models.ColorVarian{
				ProductID: product4.ID,
				Name:      c.Name,
				Color:     c.Color,
				Images:    c.Image,
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
			Images:      "https://www.hijup.com/magazine/wp-content/uploads/2022/12/4b8d043b-oversize-katun-premium.jpeg",
			Rating:      4.6,
		}
		if err := DB.Create(&product5).Error; err != nil {
			return err
		}

		colors := []struct {
			Name  string
			Color string
			Image string
		}{
			{"White", "#FFFFFF", "https://www.static-src.com/wcsstore/Indraprastha/images/catalog/full/catalog-image/103/MTA-180348393/wulfi_wulfi_atasan_kemeja_putih_kerah_jatuh_white_untuk_kasual_kerja_kantor_lengan_panjang_full01_qp8oo3p2.jpg"},
			{"Black", "#000000", "https://www.static-src.com/wcsstore/Indraprastha/images/catalog/full//catalog-image/97/MTA-139572024/no_brand_kemeja_laura_-_kemeja_basic_wanita_-_kemeja_polos_-_kemeja_kantor_full03_pl5s15pe.jpg"},
			{"Soft Blue", "#89CFF0", "https://www.static-src.com/wcsstore/Indraprastha/images/catalog/full//109/MTA-51676579/no_brand_kemeja_wanita_polos_lengan_panjang_toyobo_full14_jpu9s6t8.jpg"},
			{"Beige", "#F5F5DC", "https://www.static-src.com/wcsstore/Indraprastha/images/catalog/full/catalog-image/101/MTA-182255924/brd-136048_davelline-axelle-atasan-wanita-kemeja-kerah-pita-lengan-panjang-katun-premium_full04-db59df59.jpg"},
			{"Pink", "#FFC0CB", "https://www.static-src.com/wcsstore/Indraprastha/images/catalog/full/catalog-image/101/MTA-155622485/zahra_signature_lezahrasignature_kemeja_katun_zahra_-_fit_l-xl-_atasan_muslim_wanita_blouse_polos_full01_e3gcd157.jpg"},
		}

		sizes := []string{"S", "M", "L", "XL"}

		for _, c := range colors {
			colorVariant := models.ColorVarian{
				ProductID: product5.ID,
				Name:      c.Name,
				Color:     c.Color,
				Images:    c.Image,
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
			Images:      "https://parasayu.net/wp-content/uploads/2020/12/Rok-kaos-dan-outer.jpg",
			Rating:      4.7,
		}
		if err := DB.Create(&product6).Error; err != nil {
			return err
		}

		colors := []struct {
			Name  string
			Color string
			Image string
		}{
			{"Navy Floral", "#000080", "https://www.deliahijab.com/wp-content/uploads/2021/04/Busana-Muslim-Lebaran-Karmila-Dress-Delia-Hijab-Navy.jpeg"},
			{"Maroon Polka", "#800000", "https://www.hijup.com/magazine/wp-content/uploads/2024/04/1bfa5b00-dress-aksen-lilit-dan-tali.jpeg"},
			{"Sage Green", "#87AE73", "https://www.hijup.com/magazine/wp-content/uploads/2023/02/7e707a1c-ruffle-sage-green-dress-hijab.jpeg"},
			{"Mustard", "#FFDB58", "https://zizara.co.id/wp-content/uploads/2023/01/LINE_ALBUM_AliyaArkadewiAmoreAnika_230119_37.jpg"},
			{"Coral", "#FF7F50", "https://www.static-src.com/wcsstore/Indraprastha/images/catalog/full//98/MTA-59346233/no_brand_jesla_fashion_kimmy_dress_mosscrape_baju_gamis_terbaru_simple_dress_casual_baju_wanita_terbaru_dres_elegant_gamis_simple_best_seller_full01_sthlifk5.jpg"},
		}

		sizes := []string{"M", "L", "XL", "XXL"}

		for _, c := range colors {
			colorVariant := models.ColorVarian{
				ProductID: product6.ID,
				Name:      c.Name,
				Color:     c.Color,
				Images:    c.Image,
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
			BankImages:    "https://www.bca.co.id/-/media/Feature/Card/List-Card/Tentang-BCA/Brand-Assets/Logo-BCA/Logo-BCA_Biru.png",
		},
		{
			AccountName:   "PT E-Commerce Indonesia",
			AccountNumber: "9876543210",
			BankName:      "Bank Mandiri",
			BankImages:    "https://e7.pngegg.com/pngimages/24/85/png-clipart-bank-mandiri-logo-credit-card-bank-text-logo.png",
		},
		{
			AccountName:   "PT E-Commerce Indonesia",
			AccountNumber: "5555666677",
			BankName:      "Bank BRI",
			BankImages:    "https://logobase.net/wp-content/uploads/2025/09/Bank-Rakyat-Indonesia-BRI-Logo.webp",
		},
		{
			AccountName:   "PT E-Commerce Indonesia",
			AccountNumber: "1111222233",
			BankName:      "Bank BNI",
			BankImages:    "https://upload.wikimedia.org/wikipedia/commons/thumb/f/f0/Bank_Negara_Indonesia_logo_%282004%29.svg/1280px-Bank_Negara_Indonesia_logo_%282004%29.svg.png",
		},
		{
			AccountName:   "PT E-Commerce Indonesia",
			AccountNumber: "7777888899",
			BankName:      "Bank CIMB Niaga",
			BankImages:    "https://e7.pngegg.com/pngimages/98/783/png-clipart-logo-cimb-brand-font-text-loan-text-logo.png",
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
