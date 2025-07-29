package services

import (
	"errors"
	"literally-backend/internal/models"
	"time"
)

// In-memory storage for demo purposes
var products []models.Product
var categories []models.Category
var productIDCounter uint = 1
var categoryIDCounter uint = 1

// Initialize default categories
func init() {
	categories = []models.Category{
		{ID: 1, Name: "Phones", Icon: "phone", CreatedAt: time.Now()},
		{ID: 2, Name: "Laptops", Icon: "laptop", CreatedAt: time.Now()},
		{ID: 3, Name: "Tablets", Icon: "tablet", CreatedAt: time.Now()},
		{ID: 4, Name: "Smart Watches", Icon: "watch", CreatedAt: time.Now()},
		{ID: 5, Name: "Headphones", Icon: "headphone", CreatedAt: time.Now()},
		{ID: 6, Name: "Accessories", Icon: "accessory", CreatedAt: time.Now()},
	}
	categoryIDCounter = 7

	// Add some sample products
	products = []models.Product{
		{
			ID: 1, Name: "iPhone 15 Pro", Description: "Latest iPhone with A17 Pro chip",
			Price: 999.99, Stock: 50, ImageUrl: "https://example.com/iphone15.jpg",
			CategoryID: 1, Brand: "Apple", Rating: 4.8, ReviewCount: 120,
			IsFeatured: true, IsAvailable: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		},
		{
			ID: 2, Name: "MacBook Pro M3", Description: "Powerful laptop for professionals",
			Price: 1999.99, Stock: 30, ImageUrl: "https://example.com/macbook.jpg",
			CategoryID: 2, Brand: "Apple", Rating: 4.9, ReviewCount: 85,
			IsFeatured: true, IsAvailable: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		},
	}
	productIDCounter = 3
}

// GetAllProducts returns all products
func GetAllProducts() []models.Product {
	return products
}

// GetFeaturedProducts returns featured products
func GetFeaturedProducts() []models.Product {
	var featured []models.Product
	for _, product := range products {
		if product.IsFeatured && product.IsAvailable {
			featured = append(featured, product)
		}
	}
	return featured
}

// GetProductsByCategory returns products by category ID
func GetProductsByCategory(categoryID uint) []models.Product {
	var categoryProducts []models.Product
	for _, product := range products {
		if product.CategoryID == categoryID && product.IsAvailable {
			categoryProducts = append(categoryProducts, product)
		}
	}
	return categoryProducts
}

// GetProductByID returns a product by ID
func GetProductByID(id uint) (models.Product, bool) {
	for _, product := range products {
		if product.ID == id {
			return product, true
		}
	}
	return models.Product{}, false
}

// SearchProducts searches products by name or description
func SearchProducts(query string) []models.Product {
	var results []models.Product
	for _, product := range products {
		if contains(product.Name, query) || contains(product.Description, query) || contains(product.Brand, query) {
			if product.IsAvailable {
				results = append(results, product)
			}
		}
	}
	return results
}

// CreateProduct creates a new product
func CreateProduct(req models.CreateProductRequest) (models.Product, error) {
	// Validate category exists
	if !categoryExists(req.CategoryID) {
		return models.Product{}, errors.New("category not found")
	}

	// Create new product
	product := models.Product{
		ID:          productIDCounter,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageUrl:    req.ImageUrl,
		CategoryID:  req.CategoryID,
		Brand:       req.Brand,
		IsFeatured:  req.IsFeatured,
		IsAvailable: true,
		Rating:      0,
		ReviewCount: 0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	products = append(products, product)
	productIDCounter++

	return product, nil
}

// UpdateProduct updates an existing product
func UpdateProduct(id uint, req models.UpdateProductRequest) (models.Product, error) {
	for i, product := range products {
		if product.ID == id {
			// Update fields if provided
			if req.Name != "" {
				products[i].Name = req.Name
			}
			if req.Description != "" {
				products[i].Description = req.Description
			}
			if req.Price > 0 {
				products[i].Price = req.Price
			}
			if req.Stock >= 0 {
				products[i].Stock = req.Stock
			}
			if req.ImageUrl != "" {
				products[i].ImageUrl = req.ImageUrl
			}
			if req.CategoryID > 0 && categoryExists(req.CategoryID) {
				products[i].CategoryID = req.CategoryID
			}
			if req.Brand != "" {
				products[i].Brand = req.Brand
			}
			products[i].IsFeatured = req.IsFeatured
			products[i].IsAvailable = req.IsAvailable
			products[i].UpdatedAt = time.Now()

			return products[i], nil
		}
	}
	return models.Product{}, errors.New("product not found")
}

// DeleteProduct deletes a product by ID
func DeleteProduct(id uint) error {
	for i, product := range products {
		if product.ID == id {
			// Remove product from slice
			products = append(products[:i], products[i+1:]...)
			return nil
		}
	}
	return errors.New("product not found")
}

// Categories Management

// GetAllCategories returns all categories
func GetAllCategories() []models.Category {
	return categories
}

// GetCategoryByID returns a category by ID
func GetCategoryByID(id uint) (models.Category, bool) {
	for _, category := range categories {
		if category.ID == id {
			return category, true
		}
	}
	return models.Category{}, false
}

// Helper functions

func categoryExists(id uint) bool {
	_, exists := GetCategoryByID(id)
	return exists
}

func contains(str, substr string) bool {
	return len(str) >= len(substr) &&
		(str == substr ||
			(len(str) > len(substr) &&
				(str[:len(substr)] == substr ||
					str[len(str)-len(substr):] == substr ||
					findSubstring(str, substr))))
}

func findSubstring(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
