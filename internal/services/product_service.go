package services

import (
	"errors"
	"literally-backend/configs"
	"literally-backend/internal/models"
	"strings"

	"gorm.io/gorm"
)

// GetAllProducts returns all products
func GetAllProducts() []models.Product {
	var products []models.Product
	configs.DB.Where("is_available = ?", true).Find(&products)
	return products
}

// GetFeaturedProducts returns featured products
func GetFeaturedProducts() []models.Product {
	var products []models.Product
	configs.DB.Where("is_featured = ? AND is_available = ?", true, true).Find(&products)
	return products
}

// GetProductsByCategory returns products by category ID
func GetProductsByCategory(categoryID uint) []models.Product {
	var products []models.Product
	configs.DB.Where("category_id = ? AND is_available = ?", categoryID, true).Find(&products)
	return products
}

// GetProductByID returns a product by ID
func GetProductByID(id uint) (models.Product, bool) {
	var product models.Product
	result := configs.DB.First(&product, id)
	if result.Error != nil {
		return models.Product{}, false
	}
	return product, true
}

// SearchProducts searches products by name, description, or brand
func SearchProducts(query string) []models.Product {
	var products []models.Product
	searchTerm := "%" + strings.ToLower(query) + "%"

	configs.DB.Where(
		"(LOWER(name) LIKE ? OR LOWER(description) LIKE ? OR LOWER(brand) LIKE ?) AND is_available = ?",
		searchTerm, searchTerm, searchTerm, true,
	).Find(&products)

	return products
}

// CreateProduct creates a new product
func CreateProduct(req models.CreateProductRequest) (models.Product, error) {
	// Validate category exists
	if !categoryExists(req.CategoryID) {
		return models.Product{}, errors.New("category not found")
	}

	// Create new product
	product := models.Product{
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
	}

	if err := configs.DB.Create(&product).Error; err != nil {
		return models.Product{}, err
	}

	return product, nil
}

// UpdateProduct updates an existing product
func UpdateProduct(id uint, req models.UpdateProductRequest) (models.Product, error) {
	var product models.Product

	// Find the product
	if err := configs.DB.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Product{}, errors.New("product not found")
		}
		return models.Product{}, err
	}

	// Validate category if provided
	if req.CategoryID > 0 && !categoryExists(req.CategoryID) {
		return models.Product{}, errors.New("category not found")
	}

	// Update fields
	updates := make(map[string]interface{})

	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Price > 0 {
		updates["price"] = req.Price
	}
	if req.Stock >= 0 {
		updates["stock"] = req.Stock
	}
	if req.ImageUrl != "" {
		updates["image_url"] = req.ImageUrl
	}
	if req.CategoryID > 0 {
		updates["category_id"] = req.CategoryID
	}
	if req.Brand != "" {
		updates["brand"] = req.Brand
	}
	updates["is_featured"] = req.IsFeatured
	updates["is_available"] = req.IsAvailable

	// Update the product
	if err := configs.DB.Model(&product).Updates(updates).Error; err != nil {
		return models.Product{}, err
	}

	// Fetch the updated product
	configs.DB.First(&product, id)

	return product, nil
}

// DeleteProduct deletes a product by ID (soft delete - set is_available to false)
func DeleteProduct(id uint) error {
	result := configs.DB.Model(&models.Product{}).Where("id = ?", id).Update("is_available", false)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("product not found")
	}

	return nil
}

// Categories Management

// GetAllCategories returns all categories
func GetAllCategories() []models.Category {
	var categories []models.Category
	configs.DB.Find(&categories)
	return categories
}

// GetCategoryByID returns a category by ID
func GetCategoryByID(id uint) (models.Category, bool) {
	var category models.Category
	result := configs.DB.First(&category, id)
	if result.Error != nil {
		return models.Category{}, false
	}
	return category, true
}

// CreateCategory creates a new category
func CreateCategory(req models.CreateCategoryRequest) (models.Category, error) {
	// Check if category with same name already exists
	var existingCategory models.Category
	if err := configs.DB.Where("name = ?", req.Name).First(&existingCategory).Error; err == nil {
		return models.Category{}, errors.New("category with this name already exists")
	}

	// Create new category
	category := models.Category{
		Name: req.Name,
		Icon: req.Icon,
	}

	if err := configs.DB.Create(&category).Error; err != nil {
		return models.Category{}, err
	}

	return category, nil
}

// UpdateCategory updates an existing category
func UpdateCategory(id uint, req models.UpdateCategoryRequest) (models.Category, error) {
	var category models.Category

	// Find the category
	if err := configs.DB.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Category{}, errors.New("category not found")
		}
		return models.Category{}, err
	}

	// Check if name is being updated and if it already exists
	if req.Name != "" && req.Name != category.Name {
		var existingCategory models.Category
		if err := configs.DB.Where("name = ? AND id != ?", req.Name, id).First(&existingCategory).Error; err == nil {
			return models.Category{}, errors.New("category with this name already exists")
		}
	}

	// Update fields
	updates := make(map[string]interface{})

	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}

	// Update the category
	if err := configs.DB.Model(&category).Updates(updates).Error; err != nil {
		return models.Category{}, err
	}

	// Fetch the updated category
	configs.DB.First(&category, id)

	return category, nil
}

// DeleteCategory deletes a category by ID
func DeleteCategory(id uint) error {
	var category models.Category

	// Find the category
	if err := configs.DB.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return err
	}

	// Check if there are products using this category
	var productCount int64
	configs.DB.Model(&models.Product{}).Where("category_id = ?", id).Count(&productCount)
	if productCount > 0 {
		return errors.New("cannot delete category that has products associated with it")
	}

	// Delete the category
	if err := configs.DB.Delete(&category).Error; err != nil {
		return err
	}

	return nil
}

// Helper functions

func categoryExists(id uint) bool {
	_, exists := GetCategoryByID(id)
	return exists
}
