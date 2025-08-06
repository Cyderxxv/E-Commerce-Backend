package handlers

import (
	"literally-backend/internal/models"
	"literally-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetCart godoc
// @Summary Get user cart
// @Description Get the authenticated user's shopping cart items
// @Tags cart
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{} "Cart retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /cart [get]
func GetCart(c *gin.Context) {
	// In a real app, extract user ID from JWT token
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
		})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	cartItems := services.GetUserCart(uint(userID))
	total := services.GetCartTotal(uint(userID))

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"items": cartItems,
			"total": total,
		},
		"message": "Cart retrieved successfully",
	})
}

// AddToCart godoc
// @Summary Add item to cart
// @Description Add a product to the authenticated user's shopping cart
// @Tags cart
// @Accept json
// @Produce json
// @Security Bearer
// @Param item body object true "Cart item data"
// @Success 200 {object} map[string]interface{} "Item added to cart successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Product not found"
// @Router /cart [post]
func AddToCart(c *gin.Context) {
	// In a real app, extract user ID from JWT token
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
		})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var req models.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	cart, err := services.AddToCart(uint(userID), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    cart,
		"message": "Item added to cart successfully",
	})
}

// UpdateCartItem godoc
// @Summary Update cart item
// @Description Update the quantity of an item in the authenticated user's cart
// @Tags cart
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Cart item ID"
// @Param item body object true "Updated cart item data"
// @Success 200 {object} map[string]interface{} "Cart item updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Cart item not found"
// @Router /cart/{id} [put]
func UpdateCartItem(c *gin.Context) {
	// In a real app, extract user ID from JWT token
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
		})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	cartIDParam := c.Param("id")
	cartID, err := strconv.ParseUint(cartIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid cart item ID",
		})
		return
	}

	var req models.UpdateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	cart, err := services.UpdateCartItem(uint(userID), uint(cartID), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    cart,
		"message": "Cart item updated successfully",
	})
}

// RemoveFromCart godoc
// @Summary Remove item from cart
// @Description Remove an item from the authenticated user's cart
// @Tags cart
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Cart item ID"
// @Success 200 {object} map[string]interface{} "Item removed from cart successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid cart item ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Cart item not found"
// @Router /cart/{id} [delete]
func RemoveFromCart(c *gin.Context) {
	// In a real app, extract user ID from JWT token
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
		})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	cartIDParam := c.Param("id")
	cartID, err := strconv.ParseUint(cartIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid cart item ID",
		})
		return
	}

	err = services.RemoveFromCart(uint(userID), uint(cartID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item removed from cart successfully",
	})
}

// ClearCart godoc
// @Summary Clear user cart
// @Description Remove all items from the authenticated user's cart
// @Tags cart
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{} "Cart cleared successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /cart [delete]
func ClearCart(c *gin.Context) {
	// In a real app, extract user ID from JWT token
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
		})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	err = services.ClearUserCart(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to clear cart",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cart cleared successfully",
	})
}
