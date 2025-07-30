package handlers

import (
	"literally-backend/internal/models"
	"literally-backend/internal/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetUserPurchaseHistory handles GET /api/v1/purchase-history
func GetUserPurchaseHistory(c *gin.Context) {
	// In a real app, get user ID from JWT token
	// For now, we'll get it from query param or use a default
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		userIDStr = "1" // Default user for testing
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Parse filters
	var filter models.PurchaseHistoryFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Set default values for pagination
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	// Parse date filters if provided
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filter.StartDate = &startDate
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filter.EndDate = &endDate
		}
	}

	purchases, total, err := services.GetUserPurchaseHistory(uint(userID), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Calculate pagination info
	totalPages := 0
	if filter.Limit > 0 {
		totalPages = int((total + int64(filter.Limit) - 1) / int64(filter.Limit))
	}

	c.JSON(http.StatusOK, gin.H{
		"data": purchases,
		"pagination": gin.H{
			"current_page": filter.Page,
			"total_pages":  totalPages,
			"total_items":  total,
			"per_page":     filter.Limit,
		},
		"message": "Purchase history retrieved successfully",
	})
}

// GetPurchaseHistoryByID handles GET /api/v1/purchase-history/:id
func GetPurchaseHistoryByID(c *gin.Context) {
	// In a real app, get user ID from JWT token
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		userIDStr = "1" // Default user for testing
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	purchaseIDStr := c.Param("id")
	purchaseID, err := strconv.ParseUint(purchaseIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid purchase ID",
		})
		return
	}

	purchase, err := services.GetPurchaseHistoryByID(uint(userID), uint(purchaseID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    purchase,
		"message": "Purchase history retrieved successfully",
	})
}

// GetUserPurchaseStats handles GET /api/v1/purchase-history/stats
func GetUserPurchaseStats(c *gin.Context) {
	// In a real app, get user ID from JWT token
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		userIDStr = "1" // Default user for testing
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	stats, err := services.GetUserPurchaseStats(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    stats,
		"message": "Purchase statistics retrieved successfully",
	})
}

// GetRecentPurchases handles GET /api/v1/purchase-history/recent
func GetRecentPurchases(c *gin.Context) {
	// In a real app, get user ID from JWT token
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		userIDStr = "1" // Default user for testing
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 5
	}

	purchases, err := services.GetRecentPurchases(uint(userID), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    purchases,
		"message": "Recent purchases retrieved successfully",
	})
}

// SearchPurchaseHistory handles GET /api/v1/purchase-history/search
func SearchPurchaseHistory(c *gin.Context) {
	// In a real app, get user ID from JWT token
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		userIDStr = "1" // Default user for testing
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	searchTerm := c.Query("q")
	if searchTerm == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Search term is required",
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}

	purchases, err := services.SearchUserPurchaseHistory(uint(userID), searchTerm, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    purchases,
		"message": "Search results retrieved successfully",
	})
}

// GetPurchaseHistoryByDateRange handles GET /api/v1/purchase-history/date-range
func GetPurchaseHistoryByDateRange(c *gin.Context) {
	// In a real app, get user ID from JWT token
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		userIDStr = "1" // Default user for testing
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Both start_date and end_date are required",
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid start_date format. Use YYYY-MM-DD",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid end_date format. Use YYYY-MM-DD",
		})
		return
	}

	purchases, totalAmount, err := services.GetPurchaseHistoryByDateRange(uint(userID), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"purchases":    purchases,
			"total_amount": totalAmount,
			"start_date":   startDate.Format("2006-01-02"),
			"end_date":     endDate.Format("2006-01-02"),
		},
		"message": "Purchase history for date range retrieved successfully",
	})
}

// CanReviewProduct handles GET /api/v1/purchase-history/can-review/:product_id
func CanReviewProduct(c *gin.Context) {
	// In a real app, get user ID from JWT token
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		userIDStr = "1" // Default user for testing
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	canReview, err := services.CanUserReviewProduct(uint(userID), uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"can_review": canReview,
			"product_id": productID,
			"user_id":    userID,
		},
		"message": "Review eligibility checked successfully",
	})
}
