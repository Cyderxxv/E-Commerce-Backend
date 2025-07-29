package models

import "time"

// Review represents a product review
type Review struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	ProductID uint      `json:"product_id"`
	Rating    float64   `json:"rating" binding:"required,min=1,max=5"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ReviewWithUser represents review with user information
type ReviewWithUser struct {
	Review
	User UserResponse `json:"user"`
}

// CreateReviewRequest represents request to create a review
type CreateReviewRequest struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Rating    float64 `json:"rating" binding:"required,min=1,max=5"`
	Comment   string  `json:"comment"`
}

// UpdateReviewRequest represents request to update a review
type UpdateReviewRequest struct {
	Rating  float64 `json:"rating,omitempty"`
	Comment string  `json:"comment,omitempty"`
}

// Notification represents a user notification
type Notification struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"` // ORDER, PAYMENT, PROMOTION, etc.
	IsRead    bool      `json:"is_read" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateNotificationRequest represents request to create a notification
type CreateNotificationRequest struct {
	UserID  uint   `json:"user_id" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Message string `json:"message" binding:"required"`
	Type    string `json:"type" binding:"required"`
}

// MarkAsReadRequest represents request to mark notification as read
type MarkAsReadRequest struct {
	IsRead bool `json:"is_read"`
}

// Transaction represents a transaction
type Transaction struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	UserID          uint      `json:"user_id"`
	OrderID         uint      `json:"order_id"`
	Amount          float64   `json:"amount"`
	TransactionDate time.Time `json:"transaction_date"`
}

// TransactionWithOrder represents transaction with order details
type TransactionWithOrder struct {
	Transaction
	Order Order        `json:"order"`
	User  UserResponse `json:"user"`
}
