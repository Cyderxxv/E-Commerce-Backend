package handlers

import (
	"literally-backend/internal/models"
	"literally-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserOrders handles GET /api/v1/orders (requires authentication)
func GetUserOrders(c *gin.Context) {
	// Get user ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Get query parameters
	status := c.Query("status")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Calculate page from offset
	page := offset/limit + 1

	// Get orders
	orders, total, err := services.GetUserOrders(userID.(uint), page, limit, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve orders",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"orders": orders,
			"pagination": gin.H{
				"total":  total,
				"limit":  limit,
				"offset": offset,
			},
		},
		"message": "Orders retrieved successfully",
	})
}

// GetOrderByID handles GET /api/v1/orders/:id (requires authentication)
func GetOrderByID(c *gin.Context) {
	// Get user ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Get order ID from URL parameter
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	// Get order
	order, err := services.GetOrderByID(uint(orderID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Order not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    order,
		"message": "Order retrieved successfully",
	})
}

// CreateOrder handles POST /api/v1/orders (requires authentication)
func CreateOrder(c *gin.Context) {
	// Get user ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Parse request body
	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	// Create order
	order, err := services.CreateOrder(userID.(uint), req.ShippingAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create order",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    order,
		"message": "Order created successfully",
	})
}

// UpdateOrderStatus handles PUT /api/v1/orders/:id/status (admin only)
func UpdateOrderStatus(c *gin.Context) {
	// Get order ID from URL parameter
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	// Parse request body
	var req models.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	// Update order status
	err = services.UpdateOrderStatus(uint(orderID), req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update order status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order status updated successfully",
	})
}

// GetOrderStats handles GET /api/v1/orders/stats (requires authentication)
func GetOrderStats(c *gin.Context) {
	// Get user ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Get order statistics
	stats, err := services.GetOrderStats(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve order statistics",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    stats,
		"message": "Order statistics retrieved successfully",
	})
}
