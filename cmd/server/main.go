package main

import (
	"literally-backend/internal/handlers"
	"literally-backend/internal/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize Gin router
	router := gin.Default()

	// Add middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggingMiddleware())

	// Setup routes
	setupRoutes(router)

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRoutes(router *gin.Engine) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "Server is running",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Authentication routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		// Profile routes (requires authentication in production)
		v1.GET("/profile", handlers.GetProfile)
		v1.PUT("/profile", handlers.UpdateProfile)

		// User routes (admin only in production)
		users := v1.Group("/users")
		{
			users.GET("", handlers.GetUsers)
			users.GET("/:id", handlers.GetUserByID)
			users.POST("", handlers.CreateUser)
			users.PUT("/:id", handlers.UpdateUser)
			users.DELETE("/:id", handlers.DeleteUser)
		}

		// Category routes
		v1.GET("/categories", handlers.GetCategories)
		v1.GET("/categories/:id/products", handlers.GetProductsByCategory)

		// Product routes
		products := v1.Group("/products")
		{
			products.GET("", handlers.GetProducts)                  // GET /api/v1/products?category_id=1&featured=true&search=phone
			products.GET("/featured", handlers.GetFeaturedProducts) // GET /api/v1/products/featured
			products.GET("/search", handlers.SearchProducts)        // GET /api/v1/products/search?q=phone
			products.GET("/:id", handlers.GetProductByID)
			// Admin only routes
			products.POST("", handlers.CreateProduct)
			products.PUT("/:id", handlers.UpdateProduct)
			products.DELETE("/:id", handlers.DeleteProduct)
		}

		// Cart routes (requires authentication)
		cart := v1.Group("/cart")
		{
			cart.GET("", handlers.GetCart)               // GET /api/v1/cart?user_id=1
			cart.POST("", handlers.AddToCart)            // POST /api/v1/cart?user_id=1
			cart.PUT("/:id", handlers.UpdateCartItem)    // PUT /api/v1/cart/1?user_id=1
			cart.DELETE("/:id", handlers.RemoveFromCart) // DELETE /api/v1/cart/1?user_id=1
			cart.DELETE("", handlers.ClearCart)          // DELETE /api/v1/cart?user_id=1
		}

		// Order routes (requires authentication)
		// orders := v1.Group("/orders")
		// {
		//     orders.GET("", handlers.GetUserOrders)              // GET /api/v1/orders?user_id=1
		//     orders.POST("", handlers.CreateOrder)               // POST /api/v1/orders?user_id=1
		//     orders.GET("/:id", handlers.GetOrderByID)           // GET /api/v1/orders/1?user_id=1
		//     orders.PUT("/:id/status", handlers.UpdateOrderStatus) // PUT /api/v1/orders/1/status (admin only)
		// }

		// Wishlist routes (requires authentication)
		// wishlist := v1.Group("/wishlist")
		// {
		//     wishlist.GET("", handlers.GetWishlist)              // GET /api/v1/wishlist?user_id=1
		//     wishlist.POST("", handlers.AddToWishlist)           // POST /api/v1/wishlist?user_id=1
		//     wishlist.DELETE("/:id", handlers.RemoveFromWishlist) // DELETE /api/v1/wishlist/1?user_id=1
		// }
	}
}
