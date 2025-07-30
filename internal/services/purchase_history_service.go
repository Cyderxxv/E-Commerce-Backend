package services

import (
	"errors"
	"literally-backend/configs"
	"literally-backend/internal/models"
	"time"

	"gorm.io/gorm"
)

// GetUserPurchaseHistory retrieves purchase history for a specific user
func GetUserPurchaseHistory(userID uint, filter models.PurchaseHistoryFilter) ([]models.PurchaseHistoryResponse, int64, error) {
	// Set default values for pagination
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	var purchases []models.PurchaseHistory
	var total int64

	// Build query
	query := configs.DB.Where("user_id = ?", userID)

	// Apply filters
	if filter.Status != "" {
		query = query.Where("order_status = ?", filter.Status)
	}

	if filter.PaymentMethod != "" {
		query = query.Where("payment_method = ?", filter.PaymentMethod)
	}

	if filter.IsInstallment != nil {
		query = query.Where("is_installment = ?", *filter.IsInstallment)
	}

	if filter.StartDate != nil {
		query = query.Where("purchase_date >= ?", filter.StartDate)
	}

	if filter.EndDate != nil {
		query = query.Where("purchase_date <= ?", filter.EndDate)
	}

	if filter.ProductName != "" {
		query = query.Where("product_name ILIKE ?", "%"+filter.ProductName+"%")
	}

	// Count total records
	if err := query.Model(&models.PurchaseHistory{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Set pagination defaults
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	offset := (filter.Page - 1) * filter.Limit

	// Execute query with pagination
	if err := query.
		Preload("Product").
		Order("purchase_date DESC").
		Offset(offset).
		Limit(filter.Limit).
		Find(&purchases).Error; err != nil {
		return nil, 0, err
	}

	// Convert to response format
	var responses []models.PurchaseHistoryResponse
	for _, purchase := range purchases {
		responses = append(responses, purchase.ToResponse())
	}

	return responses, total, nil
}

// GetPurchaseHistoryByID retrieves a specific purchase history record
func GetPurchaseHistoryByID(userID, purchaseID uint) (models.PurchaseHistoryResponse, error) {
	var purchase models.PurchaseHistory

	if err := configs.DB.
		Preload("Product").
		Where("id = ? AND user_id = ?", purchaseID, userID).
		First(&purchase).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.PurchaseHistoryResponse{}, errors.New("purchase history not found")
		}
		return models.PurchaseHistoryResponse{}, err
	}

	return purchase.ToResponse(), nil
}

// GetUserPurchaseStats retrieves purchase statistics for a user
func GetUserPurchaseStats(userID uint) (models.PurchaseHistoryStats, error) {
	var stats models.PurchaseHistoryStats

	// Total orders and amount
	var totalAmount float64
	if err := configs.DB.Model(&models.PurchaseHistory{}).
		Where("user_id = ?", userID).
		Select("COUNT(*) as count, COALESCE(SUM(total_price), 0) as total").
		Row().Scan(&stats.TotalOrders, &totalAmount); err != nil {
		return stats, err
	}

	stats.TotalAmount = totalAmount

	// Orders by status
	type StatusCount struct {
		Status string
		Count  int
	}

	var statusCounts []StatusCount
	if err := configs.DB.Model(&models.PurchaseHistory{}).
		Where("user_id = ?", userID).
		Select("order_status as status, COUNT(*) as count").
		Group("order_status").
		Find(&statusCounts).Error; err != nil {
		return stats, err
	}

	for _, sc := range statusCounts {
		switch sc.Status {
		case "DELIVERED":
			stats.DeliveredOrders = sc.Count
		case "PENDING", "PROCESSING":
			stats.PendingOrders += sc.Count
		case "CANCELLED":
			stats.CancelledOrders = sc.Count
		}
	}

	// Calculate average order value
	if stats.TotalOrders > 0 {
		stats.AvgOrderValue = stats.TotalAmount / float64(stats.TotalOrders)
	}

	return stats, nil
}

// GetRecentPurchases retrieves user's recent purchases (last 30 days)
func GetRecentPurchases(userID uint, limit int) ([]models.PurchaseHistoryResponse, error) {
	var purchases []models.PurchaseHistory

	if limit <= 0 {
		limit = 5
	}

	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	if err := configs.DB.
		Preload("Product").
		Where("user_id = ? AND purchase_date >= ?", userID, thirtyDaysAgo).
		Order("purchase_date DESC").
		Limit(limit).
		Find(&purchases).Error; err != nil {
		return nil, err
	}

	var responses []models.PurchaseHistoryResponse
	for _, purchase := range purchases {
		responses = append(responses, purchase.ToResponse())
	}

	return responses, nil
}

// CreatePurchaseHistory creates a new purchase history record
func CreatePurchaseHistory(req models.PurchaseHistory) (models.PurchaseHistory, error) {
	if err := configs.DB.Create(&req).Error; err != nil {
		return models.PurchaseHistory{}, err
	}

	// Load relationships
	if err := configs.DB.
		Preload("Product").
		First(&req, req.ID).Error; err != nil {
		return models.PurchaseHistory{}, err
	}

	return req, nil
}

// UpdatePurchaseHistoryStatus updates the status of a purchase history record
func UpdatePurchaseHistoryStatus(purchaseID uint, status string) error {
	updates := map[string]interface{}{
		"order_status": status,
		"updated_at":   time.Now(),
	}

	// Set delivery date if status is DELIVERED
	if status == "DELIVERED" {
		updates["delivery_date"] = time.Now()
	}

	return configs.DB.Model(&models.PurchaseHistory{}).
		Where("id = ?", purchaseID).
		Updates(updates).Error
}

// GetPurchaseHistoryByOrder retrieves purchase history for a specific order
func GetPurchaseHistoryByOrder(orderID uint) ([]models.PurchaseHistory, error) {
	var purchases []models.PurchaseHistory

	if err := configs.DB.
		Preload("Product").
		Where("order_id = ?", orderID).
		Find(&purchases).Error; err != nil {
		return nil, err
	}

	return purchases, nil
}

// CanUserReviewProduct checks if user can review a product (must have purchased and delivered)
func CanUserReviewProduct(userID, productID uint) (bool, error) {
	var count int64

	if err := configs.DB.Model(&models.PurchaseHistory{}).
		Where("user_id = ? AND product_id = ? AND order_status = ?", userID, productID, "DELIVERED").
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetPopularProducts retrieves popular products based on purchase history
func GetPopularProducts(limit int) ([]models.Product, error) {
	if limit <= 0 {
		limit = 10
	}

	var products []models.Product

	// Get products ordered by purchase frequency
	if err := configs.DB.
		Table("products").
		Select("products.*, COUNT(purchase_history.product_id) as purchase_count").
		Joins("LEFT JOIN purchase_history ON products.id = purchase_history.product_id").
		Where("products.is_available = ?", true).
		Group("products.id").
		Order("purchase_count DESC").
		Limit(limit).
		Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

// GetPurchaseHistoryByDateRange retrieves purchase history within a date range
func GetPurchaseHistoryByDateRange(userID uint, startDate, endDate time.Time) ([]models.PurchaseHistoryResponse, float64, error) {
	var purchases []models.PurchaseHistory
	var totalAmount float64

	if err := configs.DB.
		Preload("Product").
		Where("user_id = ? AND purchase_date BETWEEN ? AND ?", userID, startDate, endDate).
		Order("purchase_date DESC").
		Find(&purchases).Error; err != nil {
		return nil, 0, err
	}

	var responses []models.PurchaseHistoryResponse
	for _, purchase := range purchases {
		responses = append(responses, purchase.ToResponse())
		totalAmount += purchase.TotalPrice
	}

	return responses, totalAmount, nil
}

// SearchUserPurchaseHistory searches purchase history by product name or other criteria
func SearchUserPurchaseHistory(userID uint, searchTerm string, limit int) ([]models.PurchaseHistoryResponse, error) {
	var purchases []models.PurchaseHistory

	if limit <= 0 {
		limit = 20
	}

	if err := configs.DB.
		Preload("Product").
		Where("user_id = ? AND (product_name ILIKE ? OR tracking_number ILIKE ?)",
			userID, "%"+searchTerm+"%", "%"+searchTerm+"%").
		Order("purchase_date DESC").
		Limit(limit).
		Find(&purchases).Error; err != nil {
		return nil, err
	}

	var responses []models.PurchaseHistoryResponse
	for _, purchase := range purchases {
		responses = append(responses, purchase.ToResponse())
	}

	return responses, nil
}
