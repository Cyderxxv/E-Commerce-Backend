package models

import "time"

// Product represents a product in the system
type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Price       float64   `json:"price" binding:"required,min=0"`
	Stock       int       `json:"stock" binding:"required,min=0"`
	ImageUrl    string    `json:"image_url"`
	Rating      float64   `json:"rating" gorm:"default:0"`
	ReviewCount int       `json:"review_count" gorm:"default:0"`
	CategoryID  uint      `json:"category_id"`
	Brand       string    `json:"brand"`
	IsFeatured  bool      `json:"is_featured" gorm:"default:false"`
	IsAvailable bool      `json:"is_available" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Category represents a product category
type Category struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" binding:"required"`
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"created_at"`
}

// ProductWithCategory represents product with category information
type ProductWithCategory struct {
	Product
	Category Category `json:"category"`
}

// CreateProductRequest represents the request body for creating a product
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,min=0"`
	Stock       int     `json:"stock" binding:"required,min=0"`
	ImageUrl    string  `json:"image_url"`
	CategoryID  uint    `json:"category_id"`
	Brand       string  `json:"brand"`
	IsFeatured  bool    `json:"is_featured"`
}

// UpdateProductRequest represents the request body for updating a product
type UpdateProductRequest struct {
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Stock       int     `json:"stock,omitempty"`
	ImageUrl    string  `json:"image_url,omitempty"`
	CategoryID  uint    `json:"category_id,omitempty"`
	Brand       string  `json:"brand,omitempty"`
	IsFeatured  bool    `json:"is_featured,omitempty"`
	IsAvailable bool    `json:"is_available,omitempty"`
}

// Cart represents a cart item
type Cart struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	ProductID uint      `json:"product_id"`
	Quantity  int       `json:"quantity" binding:"required,min=1"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CartWithProduct represents cart item with product details
type CartWithProduct struct {
	Cart
	Product Product `json:"product"`
}

// AddToCartRequest represents request to add item to cart
type AddToCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

// UpdateCartRequest represents request to update cart item
type UpdateCartRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

// Wishlist represents a wishlist item
type Wishlist struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	ProductID uint      `json:"product_id"`
	AddedAt   time.Time `json:"added_at"`
}

// WishlistWithProduct represents wishlist item with product details
type WishlistWithProduct struct {
	Wishlist
	Product Product `json:"product"`
}
