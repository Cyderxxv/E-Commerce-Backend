package handlers

import (
	"literally-backend/internal/models"
	"literally-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetProducts godoc
// @Summary Get products
// @Description Get a list of products with optional filtering by category, featured status, or search query
// @Tags products
// @Accept json
// @Produce json
// @Param category_id query string false "Filter by category ID"
// @Param featured query string false "Filter featured products (true/false)"
// @Param search query string false "Search products by name or description"
// @Success 200 {object} map[string]interface{} "Products retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid category ID"
// @Router /products [get]
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

// GetFeaturedProducts godoc
// @Summary Get featured products
// @Description Get a list of featured products
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Featured products retrieved successfully"
// @Router /products/featured [get]
func GetFeaturedProducts(c *gin.Context) {
	products := services.GetFeaturedProducts()

	c.JSON(http.StatusOK, gin.H{
		"data":    products,
		"message": "Featured products retrieved successfully",
	})
}

// GetCategories godoc
// @Summary Get all categories
// @Description Get a list of all product categories
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Categories retrieved successfully"
// @Router /categories [get]
func GetCategories(c *gin.Context) {
	categories := services.GetAllCategories()

	c.JSON(http.StatusOK, gin.H{
		"data":    categories,
		"message": "Categories retrieved successfully",
	})
}

// GetCategoryByID godoc
// @Summary Get category by ID
// @Description Get a specific category by its ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]interface{} "Category retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid category ID"
// @Failure 404 {object} map[string]interface{} "Category not found"
// @Router /categories/{id} [get]
func GetCategoryByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid category ID",
		})
		return
	}

	category, exists := services.GetCategoryByID(uint(id))
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Category not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    category,
		"message": "Category retrieved successfully",
	})
}

// GetProductsByCategory godoc
// @Summary Get products by category
// @Description Get a list of products in a specific category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]interface{} "Products retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid category ID"
// @Failure 404 {object} map[string]interface{} "Category not found"
// @Router /categories/{id}/products [get]
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

// SearchProducts godoc
// @Summary Search products
// @Description Search for products by name or description
// @Tags products
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Success 200 {object} map[string]interface{} "Search results retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Search query is required"
// @Router /products/search [get]
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

// GetProductByID godoc
// @Summary Get product by ID
// @Description Get a specific product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]interface{} "Product retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid product ID"
// @Failure 404 {object} map[string]interface{} "Product not found"
// @Router /products/{id} [get]
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

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product (admin only)
// @Tags products
// @Accept json
// @Produce json
// @Security Bearer
// @Param product body object true "Product creation data"
// @Success 201 {object} map[string]interface{} "Product created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /products [post]
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

// UpdateProduct godoc
// @Summary Update product by ID
// @Description Update a specific product by its ID (admin only)
// @Tags products
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Product ID"
// @Param product body object true "Updated product data"
// @Success 200 {object} map[string]interface{} "Product updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Product not found"
// @Router /products/{id} [put]
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

// DeleteProduct godoc
// @Summary Delete product by ID
// @Description Delete a specific product by its ID (admin only)
// @Tags products
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]interface{} "Product deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid product ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Product not found"
// @Router /products/{id} [delete]
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

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new product category (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param category body models.CreateCategoryRequest true "Category creation data"
// @Success 201 {object} map[string]interface{} "Category created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /categories [post]
func CreateCategory(c *gin.Context) {
	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	category, err := services.CreateCategory(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    category,
		"message": "Category created successfully",
	})
}

// UpdateCategory godoc
// @Summary Update category by ID
// @Description Update a specific category by its ID (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Category ID"
// @Param category body models.UpdateCategoryRequest true "Updated category data"
// @Success 200 {object} map[string]interface{} "Category updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Category not found"
// @Router /categories/{id} [put]
func UpdateCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid category ID",
		})
		return
	}

	var req models.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	category, err := services.UpdateCategory(uint(id), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    category,
		"message": "Category updated successfully",
	})
}

// DeleteCategory godoc
// @Summary Delete category by ID
// @Description Delete a specific category by its ID (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]interface{} "Category deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid category ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Category not found"
// @Router /categories/{id} [delete]
func DeleteCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid category ID",
		})
		return
	}

	err = services.DeleteCategory(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category deleted successfully",
	})
}
