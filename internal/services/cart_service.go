package services

import (
	"errors"
	"literally-backend/configs"
	"literally-backend/internal/models"

	"gorm.io/gorm"
)

// GetUserCart returns all cart items for a user with product details
func GetUserCart(userID uint) []models.CartWithProduct {
	var carts []models.Cart
	configs.DB.Where("user_id = ?", userID).Find(&carts)

	var cartWithProducts []models.CartWithProduct
	for _, cart := range carts {
		var product models.Product
		if err := configs.DB.First(&product, cart.ProductID).Error; err == nil {
			cartWithProducts = append(cartWithProducts, models.CartWithProduct{
				Cart:    cart,
				Product: product,
			})
		}
	}

	return cartWithProducts
}

// AddToCart adds a product to user's cart
func AddToCart(userID uint, req models.AddToCartRequest) (models.Cart, error) {
	// Check if product exists and is available
	var product models.Product
	if err := configs.DB.First(&product, req.ProductID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Cart{}, errors.New("product not found")
		}
		return models.Cart{}, err
	}

	if !product.IsAvailable {
		return models.Cart{}, errors.New("product is not available")
	}

	// Check stock
	if product.Stock < req.Quantity {
		return models.Cart{}, errors.New("insufficient stock")
	}

	// Check if item already exists in cart
	var existingCart models.Cart
	if err := configs.DB.Where("user_id = ? AND product_id = ?", userID, req.ProductID).First(&existingCart).Error; err == nil {
		// Update existing cart item
		newQuantity := existingCart.Quantity + req.Quantity
		if product.Stock < newQuantity {
			return models.Cart{}, errors.New("insufficient stock")
		}

		existingCart.Quantity = newQuantity
		if err := configs.DB.Save(&existingCart).Error; err != nil {
			return models.Cart{}, err
		}
		return existingCart, nil
	}

	// Create new cart item
	cart := models.Cart{
		UserID:    userID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	if err := configs.DB.Create(&cart).Error; err != nil {
		return models.Cart{}, err
	}

	return cart, nil
}

// UpdateCartItem updates quantity of a cart item
func UpdateCartItem(userID uint, cartID uint, req models.UpdateCartRequest) (models.Cart, error) {
	var cart models.Cart

	// Find the cart item
	if err := configs.DB.Where("id = ? AND user_id = ?", cartID, userID).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Cart{}, errors.New("cart item not found")
		}
		return models.Cart{}, err
	}

	// Check stock
	var product models.Product
	if err := configs.DB.First(&product, cart.ProductID).Error; err != nil {
		return models.Cart{}, errors.New("product not found")
	}

	if product.Stock < req.Quantity {
		return models.Cart{}, errors.New("insufficient stock")
	}

	// Update quantity
	cart.Quantity = req.Quantity
	if err := configs.DB.Save(&cart).Error; err != nil {
		return models.Cart{}, err
	}

	return cart, nil
}

// RemoveFromCart removes an item from user's cart
func RemoveFromCart(userID uint, cartID uint) error {
	result := configs.DB.Where("id = ? AND user_id = ?", cartID, userID).Delete(&models.Cart{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("cart item not found")
	}

	return nil
}

// ClearUserCart removes all items from user's cart
func ClearUserCart(userID uint) error {
	return configs.DB.Where("user_id = ?", userID).Delete(&models.Cart{}).Error
}

// GetCartTotal calculates total amount for user's cart
func GetCartTotal(userID uint) float64 {
	var carts []models.Cart
	configs.DB.Where("user_id = ?", userID).Find(&carts)

	total := 0.0
	for _, cart := range carts {
		var product models.Product
		if err := configs.DB.First(&product, cart.ProductID).Error; err == nil {
			total += product.Price * float64(cart.Quantity)
		}
	}

	return total
}
