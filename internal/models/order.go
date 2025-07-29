package models

import "time"

// Order represents an order in the system
type Order struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	UserID          uint      `json:"user_id"`
	TotalAmount     float64   `json:"total_amount"`
	Status          string    `json:"status"`
	PaymentMethodID uint      `json:"payment_method_id"`
	IsInstallment   bool      `json:"is_installment" gorm:"default:false"`
	ShippingAddress string    `json:"shipping_address"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	OrderID   uint      `json:"order_id"`
	ProductID uint      `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

// OrderWithItems represents order with its items
type OrderWithItems struct {
	Order
	Items   []OrderItemWithProduct `json:"items"`
	User    UserResponse           `json:"user,omitempty"`
	Payment PaymentMethod          `json:"payment_method,omitempty"`
}

// OrderItemWithProduct represents order item with product details
type OrderItemWithProduct struct {
	OrderItem
	Product Product `json:"product"`
}

// PaymentMethod represents a payment method
type PaymentMethod struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name"`
	IsInstallment bool      `json:"is_installment" gorm:"default:false"`
	CreatedAt     time.Time `json:"created_at"`
}

// CreateOrderRequest represents request to create an order
type CreateOrderRequest struct {
	PaymentMethodID uint                     `json:"payment_method_id" binding:"required"`
	IsInstallment   bool                     `json:"is_installment"`
	ShippingAddress string                   `json:"shipping_address" binding:"required"`
	Items           []CreateOrderItemRequest `json:"items" binding:"required,min=1"`
}

// CreateOrderItemRequest represents request to create an order item
type CreateOrderItemRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

// UpdateOrderStatusRequest represents request to update order status
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// InstallmentPlan represents an installment plan
type InstallmentPlan struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	OrderID        uint      `json:"order_id"`
	UserID         uint      `json:"user_id"`
	TotalAmount    float64   `json:"total_amount"`
	TotalMonths    int       `json:"total_months"`
	MonthlyPayment float64   `json:"monthly_payment"`
	PaidMonths     int       `json:"paid_months" gorm:"default:0"`
	PaidAmount     float64   `json:"paid_amount" gorm:"default:0"`
	Status         string    `json:"status" gorm:"default:ACTIVE"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// InstallmentPayment represents an installment payment
type InstallmentPayment struct {
	ID                uint       `json:"id" gorm:"primaryKey"`
	InstallmentPlanID uint       `json:"installment_plan_id"`
	MonthNumber       int        `json:"month_number"`
	DueDate           time.Time  `json:"due_date"`
	Amount            float64    `json:"amount"`
	PaidDate          *time.Time `json:"paid_date,omitempty"`
	Status            string     `json:"status" gorm:"default:PENDING"`
	CreatedAt         time.Time  `json:"created_at"`
}

// CreateInstallmentPlanRequest represents request to create installment plan
type CreateInstallmentPlanRequest struct {
	OrderID     uint `json:"order_id" binding:"required"`
	TotalMonths int  `json:"total_months" binding:"required,min=3,max=36"`
}
