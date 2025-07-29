package configs

import (
	"literally-backend/internal/models"
	"log"
)

// MigrateDatabase tự động tạo tables từ models
func MigrateDatabase() {
	log.Println("Starting database migration...")

	// Auto migrate all models
	err := DB.AutoMigrate(
		&models.User{},
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

	// Seed payment methods
	seedPaymentMethods()

	// Seed categories
	seedCategories()

	// Seed sample products
	seedSampleProducts()

	log.Println("Default data seeded successfully!")
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
