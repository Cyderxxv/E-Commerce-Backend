package services

import (
	"errors"
	"fmt"
	"time"

	"literally-backend/configs"
	"literally-backend/internal/models"

	"gorm.io/gorm"
)

type OrderService struct {
	db *gorm.DB
}

var orderService *OrderService

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{db: db}
}

func InitOrderService() {
	orderService = NewOrderService(configs.DB)
}

// Export functions for global use
func GetUserOrders(userID uint, page, limit int, status string) ([]models.Order, int64, error) {
	if orderService == nil {
		InitOrderService()
	}
	return orderService.GetUserOrders(userID, page, limit, status)
}

func GetOrderByID(orderID, userID uint) (*models.Order, error) {
	if orderService == nil {
		InitOrderService()
	}
	return orderService.GetOrderByID(orderID, userID)
}

func UpdateOrderStatus(orderID uint, status string) error {
	if orderService == nil {
		InitOrderService()
	}
	return orderService.UpdateOrderStatus(orderID, status)
}

func CreateOrder(userID uint, shippingAddress string) (*models.Order, error) {
	if orderService == nil {
		InitOrderService()
	}
	return orderService.CreateOrder(userID, shippingAddress)
}

func CreateOrderFromRequest(userID uint, req models.CreateOrderRequest) (*models.Order, error) {
	if orderService == nil {
		InitOrderService()
	}
	return orderService.CreateOrderFromRequest(userID, req)
}

func GetOrderStats(userID uint) (map[string]interface{}, error) {
	if orderService == nil {
		InitOrderService()
	}
	return orderService.GetOrderStats(userID)
}

func (s *OrderService) GetUserOrders(userID uint, page, limit int, status string) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := s.db.Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Model(&models.Order{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Preload("OrderItems").
		Preload("OrderItems.Product").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (s *OrderService) GetOrderByID(orderID, userID uint) (*models.Order, error) {
	var order models.Order
	if err := s.db.Where("id = ? AND user_id = ?", orderID, userID).
		Preload("OrderItems").
		Preload("OrderItems.Product").
		First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order not found")
		}
		return nil, err
	}
	return &order, nil
}

func (s *OrderService) UpdateOrderStatus(orderID uint, status string) error {
	result := s.db.Model(&models.Order{}).
		Where("id = ?", orderID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

func (s *OrderService) CreateOrder(userID uint, shippingAddress string) (*models.Order, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var cartItems []models.Cart
	if err := tx.Where("user_id = ?", userID).
		Find(&cartItems).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if len(cartItems) == 0 {
		tx.Rollback()
		return nil, fmt.Errorf("cart is empty")
	}

	// Validate stock availability and calculate total amount
	var totalAmount float64
	for _, item := range cartItems {
		var product models.Product
		if err := tx.Where("id = ?", item.ProductID).First(&product).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// Check if there's enough stock
		if product.Stock < item.Quantity {
			tx.Rollback()
			return nil, fmt.Errorf("insufficient stock for product %s. Available: %d, Requested: %d",
				product.Name, product.Stock, item.Quantity)
		}

		totalAmount += product.Price * float64(item.Quantity)
	}

	order := models.Order{
		UserID:          userID,
		TotalAmount:     totalAmount,
		Status:          "pending",
		ShippingAddress: shippingAddress,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create order items and update product stock
	for _, cartItem := range cartItems {
		var product models.Product
		if err := tx.Where("id = ?", cartItem.ProductID).First(&product).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// Create order item
		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
			Price:     product.Price,
		}
		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// Update product stock (decrease by purchased quantity)
		newStock := product.Stock - cartItem.Quantity
		if err := tx.Model(&product).Update("stock", newStock).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update stock for product %d: %v", cartItem.ProductID, err)
		}

		// If stock reaches 0, mark product as unavailable
		if newStock == 0 {
			if err := tx.Model(&product).Update("is_available", false).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to update availability for product %d: %v", cartItem.ProductID, err)
			}
		}
	}

	// Clear cart after successful order creation
	if err := tx.Where("user_id = ?", userID).Delete(&models.Cart{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	if err := s.db.Where("id = ?", order.ID).
		Preload("OrderItems").
		Preload("OrderItems.Product").
		First(&order).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *OrderService) GetOrderStats(userID uint) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Đếm số orders theo status
	var statusCounts []struct {
		Status string
		Count  int64
	}

	if err := s.db.Model(&models.Order{}).
		Select("status, count(*) as count").
		Where("user_id = ?", userID).
		Group("status").
		Scan(&statusCounts).Error; err != nil {
		return nil, err
	}

	stats["status_counts"] = statusCounts

	// Tổng số tiền đã chi
	var totalSpent float64
	if err := s.db.Model(&models.Order{}).
		Select("COALESCE(SUM(total_amount), 0)").
		Where("user_id = ? AND status = ?", userID, "completed").
		Scan(&totalSpent).Error; err != nil {
		return nil, err
	}

	stats["total_spent"] = totalSpent

	// Order gần nhất
	var lastOrder models.Order
	if err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&lastOrder).Error; err == nil {
		stats["last_order_date"] = lastOrder.CreatedAt
	}

	return stats, nil
}

func (s *OrderService) CreateOrderFromRequest(userID uint, req models.CreateOrderRequest) (*models.Order, error) {
	// Begin transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Calculate total amount and validate stock availability
	var totalAmount float64
	for _, item := range req.Items {
		var product models.Product
		if err := tx.Where("id = ?", item.ProductID).First(&product).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("product not found: %d", item.ProductID)
		}

		// Check if there's enough stock
		if product.Stock < item.Quantity {
			tx.Rollback()
			return nil, fmt.Errorf("insufficient stock for product %s. Available: %d, Requested: %d",
				product.Name, product.Stock, item.Quantity)
		}

		totalAmount += product.Price * float64(item.Quantity)
	}

	// Create order
	order := models.Order{
		UserID:          userID,
		TotalAmount:     totalAmount,
		Status:          "pending",
		PaymentMethodID: req.PaymentMethodID,
		IsInstallment:   req.IsInstallment,
		ShippingAddress: req.ShippingAddress,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create order items and update product stock
	for _, reqItem := range req.Items {
		var product models.Product
		if err := tx.Where("id = ?", reqItem.ProductID).First(&product).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// Create order item
		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: reqItem.ProductID,
			Quantity:  reqItem.Quantity,
			Price:     product.Price,
		}
		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// Update product stock (decrease by purchased quantity)
		newStock := product.Stock - reqItem.Quantity
		if err := tx.Model(&product).Update("stock", newStock).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update stock for product %d: %v", reqItem.ProductID, err)
		}

		// If stock reaches 0, mark product as unavailable
		if newStock == 0 {
			if err := tx.Model(&product).Update("is_available", false).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to update availability for product %d: %v", reqItem.ProductID, err)
			}
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Load complete order with items
	if err := s.db.Where("id = ?", order.ID).
		Preload("OrderItems").
		Preload("OrderItems.Product").
		First(&order).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

// Admin Order Management Functions

// GetAllOrders returns all orders for admin (with pagination and filtering)
func GetAllOrders(page, limit int, status string) ([]models.Order, int64, error) {
	if orderService == nil {
		InitOrderService()
	}
	return orderService.GetAllOrders(page, limit, status)
}

// GetOrderByIDAdmin returns order by ID for admin (no user restriction)
func GetOrderByIDAdmin(orderID uint) (*models.Order, error) {
	if orderService == nil {
		InitOrderService()
	}
	return orderService.GetOrderByIDAdmin(orderID)
}

// GetOrderStats returns order statistics for admin
func GetAdminOrderStats() (map[string]interface{}, error) {
	if orderService == nil {
		InitOrderService()
	}
	return orderService.GetAdminOrderStats()
}

func (s *OrderService) GetAllOrders(page, limit int, status string) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := s.db.Model(&models.Order{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Preload("OrderItems").
		Preload("OrderItems.Product").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (s *OrderService) GetOrderByIDAdmin(orderID uint) (*models.Order, error) {
	var order models.Order
	if err := s.db.Where("id = ?", orderID).
		Preload("OrderItems").
		Preload("OrderItems.Product").
		First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	return &order, nil
}

func (s *OrderService) GetAdminOrderStats() (map[string]interface{}, error) {
	var totalOrders int64
	var totalRevenue float64
	var pendingOrders int64
	var confirmedOrders int64
	var shippedOrders int64
	var deliveredOrders int64
	var cancelledOrders int64

	// Count total orders
	if err := s.db.Model(&models.Order{}).Count(&totalOrders).Error; err != nil {
		return nil, err
	}

	// Calculate total revenue (from delivered orders)
	if err := s.db.Model(&models.Order{}).
		Where("status = ?", "DELIVERED").
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&totalRevenue).Error; err != nil {
		return nil, err
	}

	// Count orders by status
	if err := s.db.Model(&models.Order{}).Where("status = ?", "PENDING").Count(&pendingOrders).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&models.Order{}).Where("status = ?", "CONFIRMED").Count(&confirmedOrders).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&models.Order{}).Where("status = ?", "SHIPPED").Count(&shippedOrders).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&models.Order{}).Where("status = ?", "DELIVERED").Count(&deliveredOrders).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&models.Order{}).Where("status = ?", "CANCELLED").Count(&cancelledOrders).Error; err != nil {
		return nil, err
	}

	// Get orders from last 30 days
	var recentOrders int64
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	if err := s.db.Model(&models.Order{}).
		Where("created_at >= ?", thirtyDaysAgo).
		Count(&recentOrders).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_orders":     totalOrders,
		"total_revenue":    totalRevenue,
		"pending_orders":   pendingOrders,
		"confirmed_orders": confirmedOrders,
		"shipped_orders":   shippedOrders,
		"delivered_orders": deliveredOrders,
		"cancelled_orders": cancelledOrders,
		"recent_orders":    recentOrders,
	}, nil
}
