package services

import (
	"errors"
	"literally-backend/internal/models"
	"time"
)

// In-memory storage for cart items
var carts []models.Cart
var cartIDCounter uint = 1

// GetUserCart returns all cart items for a user
func GetUserCart(userID uint) []models.CartWithProduct {
	var userCarts []models.CartWithProduct
	for _, cart := range carts {
		if cart.UserID == userID {
			product, exists := GetProductByID(cart.ProductID)
			if exists {
				userCarts = append(userCarts, models.CartWithProduct{
					Cart:    cart,
					Product: product,
				})
			}
		}
	}
	return userCarts
}

// AddToCart adds a product to user's cart
func AddToCart(userID uint, req models.AddToCartRequest) (models.Cart, error) {
	// Check if product exists
	product, exists := GetProductByID(req.ProductID)
	if !exists {
		return models.Cart{}, errors.New("product not found")
	}

	// Check if product is available
	if !product.IsAvailable {
		return models.Cart{}, errors.New("product is not available")
	}

	// Check stock
	if product.Stock < req.Quantity {
		return models.Cart{}, errors.New("insufficient stock")
	}

	// Check if item already exists in cart
	for i, cart := range carts {
		if cart.UserID == userID && cart.ProductID == req.ProductID {
			// Update quantity
			newQuantity := cart.Quantity + req.Quantity
			if product.Stock < newQuantity {
				return models.Cart{}, errors.New("insufficient stock")
			}
			carts[i].Quantity = newQuantity
			carts[i].UpdatedAt = time.Now()
			return carts[i], nil
		}
	}

	// Create new cart item
	cart := models.Cart{
		ID:        cartIDCounter,
		UserID:    userID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	carts = append(carts, cart)
	cartIDCounter++

	return cart, nil
}

// UpdateCartItem updates quantity of a cart item
func UpdateCartItem(userID uint, cartID uint, req models.UpdateCartRequest) (models.Cart, error) {
	for i, cart := range carts {
		if cart.ID == cartID && cart.UserID == userID {
			// Check stock
			product, exists := GetProductByID(cart.ProductID)
			if !exists {
				return models.Cart{}, errors.New("product not found")
			}

			if product.Stock < req.Quantity {
				return models.Cart{}, errors.New("insufficient stock")
			}

			carts[i].Quantity = req.Quantity
			carts[i].UpdatedAt = time.Now()
			return carts[i], nil
		}
	}
	return models.Cart{}, errors.New("cart item not found")
}

// RemoveFromCart removes an item from user's cart
func RemoveFromCart(userID uint, cartID uint) error {
	for i, cart := range carts {
		if cart.ID == cartID && cart.UserID == userID {
			// Remove cart item from slice
			carts = append(carts[:i], carts[i+1:]...)
			return nil
		}
	}
	return errors.New("cart item not found")
}

// ClearUserCart removes all items from user's cart
func ClearUserCart(userID uint) error {
	var newCarts []models.Cart
	for _, cart := range carts {
		if cart.UserID != userID {
			newCarts = append(newCarts, cart)
		}
	}
	carts = newCarts
	return nil
}

// GetCartTotal calculates total amount for user's cart
func GetCartTotal(userID uint) float64 {
	total := 0.0
	for _, cart := range carts {
		if cart.UserID == userID {
			product, exists := GetProductByID(cart.ProductID)
			if exists {
				total += product.Price * float64(cart.Quantity)
			}
		}
	}
	return total
}
