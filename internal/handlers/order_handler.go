package handlers

import (
	"literally-backend/internal/models"
	"literally-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserOrders godoc
// @Summary Get user orders
// @Description Get a list of orders for the authenticated user with pagination and optional status filtering
// @Tags orders
// @Accept json
// @Produce json
// @Security Bearer
// @Param status query string false "Filter by order status (PENDING, CONFIRMED, SHIPPED, DELIVERED, CANCELLED)"
// @Param limit query int false "Number of orders to return (default: 10)"
// @Param offset query int false "Number of orders to skip (default: 0)"
// @Success 200 {object} map[string]interface{} "Success response with orders and pagination"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /orders [get]
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

// GetOrderByID godoc
// @Summary Get order by ID
// @Description Get a specific order by ID for the authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Order ID"
// @Success 200 {object} map[string]interface{} "Success response with order details"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid order ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Router /orders/{id} [get]
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

// CreateOrder godoc
// @Summary Create new order
// @Description Create a new order from cart items or specific items for the authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Security Bearer
// @Param order body object false "Order creation request (optional - if not provided, creates order from cart)"
// @Success 201 {object} map[string]interface{} "Success response with created order"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /orders [post]
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
	order, err := services.CreateOrderFromRequest(userID.(uint), req)
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

// UpdateOrderStatus godoc
// @Summary Update order status
// @Description Update the status of an order (admin only)
// @Tags orders
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Order ID"
// @Param status body object true "Status update request"
// @Success 200 {object} map[string]interface{} "Success response"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /orders/{id}/status [put]
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

// GetOrderStats godoc
// @Summary Get order statistics
// @Description Get order statistics for the authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{} "Success response with order statistics"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /orders/stats [get]
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
