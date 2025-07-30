package models

import "time"

// PurchaseHistory represents a user's purchase history record
type PurchaseHistory struct {
	ID                uint       `json:"id" gorm:"primaryKey"`
	UserID            uint       `json:"user_id" gorm:"not null"`
	OrderID           *uint      `json:"order_id,omitempty"`
	ProductID         uint       `json:"product_id" gorm:"not null"`
	ProductName       string     `json:"product_name" gorm:"not null"`
	ProductImageURL   string     `json:"product_image_url,omitempty"`
	Quantity          int        `json:"quantity" gorm:"not null"`
	UnitPrice         float64    `json:"unit_price" gorm:"type:decimal(10,2);not null"`
	TotalPrice        float64    `json:"total_price" gorm:"type:decimal(10,2);not null"`
	OrderStatus       string     `json:"order_status" gorm:"not null"`
	PaymentMethod     string     `json:"payment_method,omitempty"`
	IsInstallment     bool       `json:"is_installment" gorm:"default:false"`
	InstallmentMonths *int       `json:"installment_months,omitempty"`
	MonthlyPayment    *float64   `json:"monthly_payment,omitempty" gorm:"type:decimal(10,2)"`
	PurchaseDate      time.Time  `json:"purchase_date" gorm:"not null"`
	DeliveryDate      *time.Time `json:"delivery_date,omitempty"`
	TrackingNumber    string     `json:"tracking_number,omitempty"`
	ShippingAddress   string     `json:"shipping_address,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	// Relationships
	User    User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Product Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

// PurchaseHistoryResponse represents purchase history data returned to client
type PurchaseHistoryResponse struct {
	ID                uint       `json:"id"`
	ProductID         uint       `json:"product_id"`
	ProductName       string     `json:"product_name"`
	ProductImageURL   string     `json:"product_image_url"`
	Quantity          int        `json:"quantity"`
	UnitPrice         float64    `json:"unit_price"`
	TotalPrice        float64    `json:"total_price"`
	OrderStatus       string     `json:"order_status"`
	StatusDisplay     string     `json:"status_display"`
	PaymentMethod     string     `json:"payment_method"`
	IsInstallment     bool       `json:"is_installment"`
	InstallmentMonths *int       `json:"installment_months,omitempty"`
	MonthlyPayment    *float64   `json:"monthly_payment,omitempty"`
	PurchaseDate      time.Time  `json:"purchase_date"`
	DeliveryDate      *time.Time `json:"delivery_date,omitempty"`
	TrackingNumber    string     `json:"tracking_number,omitempty"`
	ShippingAddress   string     `json:"shipping_address,omitempty"`
	DaysSincePurchase int        `json:"days_since_purchase"`
	CanReview         bool       `json:"can_review"`
	CanReorder        bool       `json:"can_reorder"`
}

// PurchaseHistoryFilter represents filters for purchase history queries
type PurchaseHistoryFilter struct {
	Status        string     `json:"status,omitempty" form:"status"`
	PaymentMethod string     `json:"payment_method,omitempty" form:"payment_method"`
	IsInstallment *bool      `json:"is_installment,omitempty" form:"is_installment"`
	StartDate     *time.Time `json:"start_date,omitempty" form:"start_date"`
	EndDate       *time.Time `json:"end_date,omitempty" form:"end_date"`
	ProductName   string     `json:"product_name,omitempty" form:"product_name"`
	Page          int        `json:"page" form:"page"`
	Limit         int        `json:"limit" form:"limit"`
}

// ToResponse converts PurchaseHistory to PurchaseHistoryResponse
func (ph *PurchaseHistory) ToResponse() PurchaseHistoryResponse {
	statusDisplay := ph.getStatusDisplay()
	daysSince := int(time.Since(ph.PurchaseDate).Hours() / 24)
	canReview := ph.OrderStatus == "DELIVERED"
	canReorder := ph.OrderStatus == "DELIVERED" || ph.OrderStatus == "CANCELLED"

	return PurchaseHistoryResponse{
		ID:                ph.ID,
		ProductID:         ph.ProductID,
		ProductName:       ph.ProductName,
		ProductImageURL:   ph.ProductImageURL,
		Quantity:          ph.Quantity,
		UnitPrice:         ph.UnitPrice,
		TotalPrice:        ph.TotalPrice,
		OrderStatus:       ph.OrderStatus,
		StatusDisplay:     statusDisplay,
		PaymentMethod:     ph.PaymentMethod,
		IsInstallment:     ph.IsInstallment,
		InstallmentMonths: ph.InstallmentMonths,
		MonthlyPayment:    ph.MonthlyPayment,
		PurchaseDate:      ph.PurchaseDate,
		DeliveryDate:      ph.DeliveryDate,
		TrackingNumber:    ph.TrackingNumber,
		ShippingAddress:   ph.ShippingAddress,
		DaysSincePurchase: daysSince,
		CanReview:         canReview,
		CanReorder:        canReorder,
	}
}

// getStatusDisplay returns human-readable status
func (ph *PurchaseHistory) getStatusDisplay() string {
	switch ph.OrderStatus {
	case "DELIVERED":
		return "Đã giao hàng"
	case "CANCELLED":
		return "Đã hủy"
	case "SHIPPED":
		return "Đang vận chuyển"
	case "PROCESSING":
		return "Đang xử lý"
	case "PENDING":
		return "Chờ xác nhận"
	default:
		return "Không xác định"
	}
}

// PurchaseHistoryStats represents purchase statistics
type PurchaseHistoryStats struct {
	TotalOrders     int     `json:"total_orders"`
	TotalAmount     float64 `json:"total_amount"`
	DeliveredOrders int     `json:"delivered_orders"`
	PendingOrders   int     `json:"pending_orders"`
	CancelledOrders int     `json:"cancelled_orders"`
	AvgOrderValue   float64 `json:"avg_order_value"`
}
