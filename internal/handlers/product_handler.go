package handlers

import (
	"literally-backend/internal/models"
	"literally-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetProducts handles GET /api/v1/products
func GetProducts(c *gin.Context) {
	// Check for query parameters
	categoryID := c.Query("category_id")
	featured := c.Query("featured")
	search := c.Query("search")

	var products []models.Product

	if search != "" {
		// Search products
		products = services.SearchProducts(search)
	} else if categoryID != "" {
		// Get products by category
		catID, err := strconv.ParseUint(categoryID, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid category ID",
			})
			return
		}
		products = services.GetProductsByCategory(uint(catID))
	} else if featured == "true" {
		// Get featured products
		products = services.GetFeaturedProducts()
	} else {
		// Get all products
		products = services.GetAllProducts()
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    products,
		"message": "Products retrieved successfully",
	})
}

// GetFeaturedProducts handles GET /api/v1/products/featured
func GetFeaturedProducts(c *gin.Context) {
	products := services.GetFeaturedProducts()

	c.JSON(http.StatusOK, gin.H{
		"data":    products,
		"message": "Featured products retrieved successfully",
	})
}

// GetCategories handles GET /api/v1/categories
func GetCategories(c *gin.Context) {
	categories := services.GetAllCategories()

	c.JSON(http.StatusOK, gin.H{
		"data":    categories,
		"message": "Categories retrieved successfully",
	})
}

// GetProductsByCategory handles GET /api/v1/categories/:id/products
func GetProductsByCategory(c *gin.Context) {
	categoryIDParam := c.Param("id")
	categoryID, err := strconv.ParseUint(categoryIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid category ID",
		})
		return
	}

	// Check if category exists
	_, exists := services.GetCategoryByID(uint(categoryID))
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Category not found",
		})
		return
	}

	products := services.GetProductsByCategory(uint(categoryID))

	c.JSON(http.StatusOK, gin.H{
		"data":    products,
		"message": "Products retrieved successfully",
	})
}

// SearchProducts handles GET /api/v1/products/search
func SearchProducts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Search query is required",
		})
		return
	}

	products := services.SearchProducts(query)

	c.JSON(http.StatusOK, gin.H{
		"data":    products,
		"message": "Search results retrieved successfully",
	})
}

// GetProductByID handles GET /api/v1/products/:id
func GetProductByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	product, exists := services.GetProductByID(uint(id))
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Product not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    product,
		"message": "Product retrieved successfully",
	})
}

// CreateProduct handles POST /api/v1/products
func CreateProduct(c *gin.Context) {
	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	product, err := services.CreateProduct(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    product,
		"message": "Product created successfully",
	})
}

// UpdateProduct handles PUT /api/v1/products/:id
func UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	product, err := services.UpdateProduct(uint(id), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    product,
		"message": "Product updated successfully",
	})
}

// DeleteProduct handles DELETE /api/v1/products/:id
func DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	err = services.DeleteProduct(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}
