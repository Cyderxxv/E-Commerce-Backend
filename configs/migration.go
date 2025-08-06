package configs

import (
	"literally-backend/internal/models"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func MigrateDatabase() {
	log.Println("Starting database migration...")

	// Auto migrate all models
	err := DB.AutoMigrate(
		&models.User{},
		&models.Admin{},
		&models.Category{},
		&models.Product{},
		&models.Cart{},
		&models.PaymentMethod{},
		&models.Order{},
		&models.OrderItem{},
		&models.InstallmentPlan{},
		&models.InstallmentPayment{},
		&models.Wishlist{},
		&models.Review{},
		&models.Notification{},
		&models.Transaction{},
		&models.PurchaseHistory{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed successfully!")

	// Seed default data
	seedDefaultData()
}

// seedDefaultData thêm dữ liệu mặc định
func seedDefaultData() {
	log.Println("Seeding default data...")

	// Seed admin first
	seedAdmin()

	// Seed users
	seedUsers()

	// Seed payment methods
	seedPaymentMethods()

	// Seed categories
	seedCategories()

	// Seed sample products
	seedSampleProducts()

	// Seed purchase history
	seedPurchaseHistory()

	log.Println("Default data seeded successfully!")
}

// mustHashPassword hashes password and panics if error
func mustHashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}

// seedPurchaseHistory thêm sample purchase history
func seedPurchaseHistory() {
	var count int64
	DB.Model(&models.PurchaseHistory{}).Count(&count)

	if count == 0 {
		// Sample purchase history records
		purchaseHistories := []models.PurchaseHistory{
			{
				UserID:          1,
				ProductID:       1,
				ProductName:     "Samsung Galaxy S25 Edge (12/256GB)",
				ProductImageURL: "https://cdn.viettablet.com/images/detailed/66/samsung-galaxy-s25-edge-111.jpg",
				Quantity:        1,
				UnitPrice:       25650600,
				TotalPrice:      25650600,
				OrderStatus:     "DELIVERED",
				PaymentMethod:   "Cash",
				IsInstallment:   false,
				PurchaseDate:    time.Now().AddDate(0, 0, -15),
				DeliveryDate:    &[]time.Time{time.Now().AddDate(0, 0, -10)}[0],
				ShippingAddress: "Ho Chi Minh City, Vietnam",
			},
			{
				UserID:            1,
				ProductID:         6,
				ProductName:       "Apple MacBook Air M4 13-inch (16/512GB)",
				ProductImageURL:   "https://bizweb.dktcdn.net/100/453/356/products/mbair-13inch-m4-midnight-1744562440665.jpg?v=1747827209317",
				Quantity:          1,
				UnitPrice:         29990000,
				TotalPrice:        29990000,
				OrderStatus:       "PROCESSING",
				PaymentMethod:     "Installment",
				IsInstallment:     true,
				InstallmentMonths: &[]int{12}[0],
				MonthlyPayment:    &[]float64{2499167}[0],
				PurchaseDate:      time.Now().AddDate(0, 0, -5),
				ShippingAddress:   "Ho Chi Minh City, Vietnam",
			},
			{
				UserID:          1,
				ProductID:       2,
				ProductName:     "Xiaomi 15S PRO (12/256GB)",
				ProductImageURL: "https://cdn.mobilecity.vn/mobilecity-vn/images/2025/05/w300/xiaomi-15s-pro-den-cac-bon.jpg.webp",
				Quantity:        1,
				UnitPrice:       14550200,
				TotalPrice:      14550200,
				OrderStatus:     "DELIVERED",
				PaymentMethod:   "Bank Transfer",
				IsInstallment:   false,
				PurchaseDate:    time.Now().AddDate(0, -1, -10),
				DeliveryDate:    &[]time.Time{time.Now().AddDate(0, -1, -5)}[0],
				ShippingAddress: "Ho Chi Minh City, Vietnam",
			},
			{
				UserID:          1,
				ProductID:       4,
				ProductName:     "Samsung Galaxy Z Flip6 (12/256GB)",
				ProductImageURL: "https://cdn2.cellphones.com.vn/insecure/rs:fill:0:358/q:90/plain/https://cellphones.com.vn/media/catalog/product/s/a/samsung-galaxy-z-flip-6-xanh-duong-4_2.png",
				Quantity:        1,
				UnitPrice:       20550200,
				TotalPrice:      20550200,
				OrderStatus:     "DELIVERED",
				PaymentMethod:   "Cash",
				IsInstallment:   false,
				PurchaseDate:    time.Now().AddDate(0, -2, -20),
				DeliveryDate:    &[]time.Time{time.Now().AddDate(0, -2, -15)}[0],
				ShippingAddress: "Ho Chi Minh City, Vietnam",
			},
			{
				UserID:          1,
				ProductID:       8,
				ProductName:     "Apple Watch Series 10 (46mm)",
				ProductImageURL: "https://store.storeimages.cdn-apple.com/1/as-images.apple.com/is/MXM23ref_FV99_VW_34FR+watch-case-46-aluminum-jetblack-nc-s10_VW_34FR+watch-face-46-aluminum-jetblack-s10_VW_34FR?wid=752&hei=720&bgc=fafafa&trim=1&fmt=p-jpg&qlt=80&.v=TnVrdDZWRlZzTURKbHFqOGh0dGpVRW5TeWJ6QW43NUFnQ2V4cmRFc1VnYUdWejZ5THhpKzJwRmRDYlhxN2o5aXB2QjR6TEZ4ZThxM3VqYkZobmlXM3RGNnlaeXQ4NGFKQTAzc0NGeHR2aVk0VEhOZEFKYmY1ZHNpalQ3YVhOWk9WV",
				Quantity:        1,
				UnitPrice:       9990000,
				TotalPrice:      9990000,
				OrderStatus:     "SHIPPED",
				PaymentMethod:   "Credit Card",
				IsInstallment:   false,
				PurchaseDate:    time.Now().AddDate(0, 0, -3),
				TrackingNumber:  "AWS123456789",
				ShippingAddress: "Ho Chi Minh City, Vietnam",
			},
		}

		for _, purchase := range purchaseHistories {
			DB.Create(&purchase)
		}
		log.Println("Purchase history seeded")
	}
}

func seedUsers() {
	var count int64
	DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		return // Users already exist
	}

	users := []models.User{
		{
			Name:         "Test User",
			Email:        "test@example.com",
			PhoneNumber:  "+84123456789",
			PasswordHash: mustHashPassword("123456"),
			Photo:        "https://avatar.com/test.jpg",
			FullName:     "Test User Full Name",
			Address:      "123 Test Street, Test City",
			Gender:       "Male",
			Status:       "ACTIVE",
		},
		{
			Name:         "John Doe",
			Email:        "john@example.com",
			PhoneNumber:  "+84987654321",
			PasswordHash: mustHashPassword("123456"),
			Photo:        "https://avatar.com/john.jpg",
			FullName:     "John Michael Doe",
			Address:      "456 Main Street, City",
			Gender:       "Male",
			Status:       "ACTIVE",
		},
		{
			Name:         "Jane Smith",
			Email:        "jane@example.com",
			PhoneNumber:  "+84111222333",
			PasswordHash: mustHashPassword("123456"),
			Photo:        "https://avatar.com/jane.jpg",
			FullName:     "Jane Elizabeth Smith",
			Address:      "789 Oak Avenue, City",
			Gender:       "Female",
			Status:       "ACTIVE",
		},
	}

	for _, user := range users {
		DB.Create(&user)
	}
}

// seedPaymentMethods thêm payment methods mặc định
func seedPaymentMethods() {
	var count int64
	DB.Model(&models.PaymentMethod{}).Count(&count)

	if count == 0 {
		paymentMethods := []models.PaymentMethod{
			{Name: "Cash", IsInstallment: false},
			{Name: "Credit Card", IsInstallment: false},
			{Name: "Bank Transfer", IsInstallment: false},
			{Name: "Installment", IsInstallment: true},
		}

		for _, pm := range paymentMethods {
			DB.Create(&pm)
		}
		log.Println("Payment methods seeded")
	}
}

// seedCategories thêm categories mặc định
func seedCategories() {
	var count int64
	DB.Model(&models.Category{}).Count(&count)

	if count == 0 {
		categories := []models.Category{
			{Name: "Phones", Icon: "phone"},
			{Name: "Laptops", Icon: "laptop"},
			{Name: "Tablets", Icon: "tablet"},
			{Name: "Smart Watches", Icon: "watch"},
			{Name: "Headphones", Icon: "headphone"},
			{Name: "Accessories", Icon: "accessory"},
		}

		for _, cat := range categories {
			DB.Create(&cat)
		}
		log.Println("Categories seeded")
	}
}

// seedSampleProducts thêm sample products
func seedSampleProducts() {
	var count int64
	DB.Model(&models.Product{}).Count(&count)

	if count == 0 {
		products := []models.Product{
			{
				Name:        "iPhone 15 Pro",
				Description: "Latest iPhone with A17 Pro chip and titanium design",
				Price:       999.99,
				Stock:       50,
				ImageUrl:    "https://example.com/iphone15pro.jpg",
				CategoryID:  1, // Phones
				Brand:       "Apple",
				IsFeatured:  true,
				IsAvailable: true,
				Rating:      4.8,
				ReviewCount: 120,
			},
			{
				Name:        "MacBook Pro M3",
				Description: "Powerful laptop for professionals with M3 chip",
				Price:       1999.99,
				Stock:       30,
				ImageUrl:    "https://example.com/macbook-pro-m3.jpg",
				CategoryID:  2, // Laptops
				Brand:       "Apple",
				IsFeatured:  true,
				IsAvailable: true,
				Rating:      4.9,
				ReviewCount: 85,
			},
			{
				Name:        "Samsung Galaxy S24 Ultra",
				Description: "Premium Android phone with S Pen",
				Price:       1199.99,
				Stock:       40,
				ImageUrl:    "https://example.com/samsung-s24-ultra.jpg",
				CategoryID:  1, // Phones
				Brand:       "Samsung",
				IsFeatured:  true,
				IsAvailable: true,
				Rating:      4.7,
				ReviewCount: 95,
			},
			{
				Name:        "iPad Pro 12.9",
				Description: "Professional tablet with M2 chip",
				Price:       1099.99,
				Stock:       25,
				ImageUrl:    "https://example.com/ipad-pro.jpg",
				CategoryID:  3, // Tablets
				Brand:       "Apple",
				IsFeatured:  false,
				IsAvailable: true,
				Rating:      4.6,
				ReviewCount: 67,
			},
			{
				Name:        "Apple Watch Series 9",
				Description: "Advanced smartwatch with health monitoring",
				Price:       399.99,
				Stock:       60,
				ImageUrl:    "https://example.com/apple-watch-9.jpg",
				CategoryID:  4, // Smart Watches
				Brand:       "Apple",
				IsFeatured:  false,
				IsAvailable: true,
				Rating:      4.5,
				ReviewCount: 145,
			},
			{
				Name:        "AirPods Pro 2",
				Description: "Premium wireless earbuds with noise cancellation",
				Price:       249.99,
				Stock:       80,
				ImageUrl:    "https://example.com/airpods-pro-2.jpg",
				CategoryID:  5, // Headphones
				Brand:       "Apple",
				IsFeatured:  false,
				IsAvailable: true,
				Rating:      4.4,
				ReviewCount: 234,
			},
		}

		for _, product := range products {
			DB.Create(&product)
		}
		log.Println("Sample products seeded")
	}
}

// seedAdmin creates default admin user
func seedAdmin() {
	var count int64
	DB.Model(&models.Admin{}).Count(&count)

	if count == 0 {
		admin := models.Admin{
			Username:     "admin",
			PasswordHash: mustHashPassword("admin123"),
		}

		DB.Create(&admin)
		log.Println("Default admin created: username=admin, password=admin123")
	}
}
