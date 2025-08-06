package main

import (
	"literally-backend/configs"
	_ "literally-backend/docs" // Import generated docs
	"literally-backend/internal/handlers"
	"literally-backend/internal/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Literally Backend API
// @version         1.0
// @description     This is a comprehensive e-commerce backend API server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database connection
	configs.ConnectDatabase()

	// Run migrations and seed data
	configs.MigrateDatabase()

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
	// Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
		// Authentication routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		// Profile routes (requires authentication)
		profile := v1.Group("/")
		profile.Use(middleware.AuthMiddleware())
		{
			profile.GET("/profile", handlers.GetProfile)
			profile.PUT("/profile", handlers.UpdateProfile)
		}

		// User routes (admin only - for now, require authentication)
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("", handlers.GetUsers)
			users.GET("/:id", handlers.GetUserByID)
			users.POST("", handlers.CreateUser)
			users.PUT("/:id", handlers.UpdateUser)
			users.DELETE("/:id", handlers.DeleteUser)
		}

		// Category routes (public)
		v1.GET("/categories", handlers.GetCategories)
		v1.GET("/categories/:id/products", handlers.GetProductsByCategory)

		// Product routes (public for read, protected for write)
		products := v1.Group("/products")
		{
			products.GET("", handlers.GetProducts)                  // GET /api/v1/products?category_id=1&featured=true&search=phone
			products.GET("/featured", handlers.GetFeaturedProducts) // GET /api/v1/products/featured
			products.GET("/search", handlers.SearchProducts)        // GET /api/v1/products/search?q=phone
			products.GET("/:id", handlers.GetProductByID)
		}

		// Admin product routes (requires authentication)
		adminProducts := v1.Group("/products")
		adminProducts.Use(middleware.AuthMiddleware())
		{
			adminProducts.POST("", handlers.CreateProduct)
			adminProducts.PUT("/:id", handlers.UpdateProduct)
			adminProducts.DELETE("/:id", handlers.DeleteProduct)
		}

		// Cart routes (requires authentication)
		cart := v1.Group("/cart")
		cart.Use(middleware.AuthMiddleware())
		{
			cart.GET("", handlers.GetCart)               // GET /api/v1/cart
			cart.POST("", handlers.AddToCart)            // POST /api/v1/cart
			cart.PUT("/:id", handlers.UpdateCartItem)    // PUT /api/v1/cart/1
			cart.DELETE("/:id", handlers.RemoveFromCart) // DELETE /api/v1/cart/1
			cart.DELETE("", handlers.ClearCart)          // DELETE /api/v1/cart
		}

		// Purchase History routes (requires authentication)
		purchaseHistory := v1.Group("/purchase-history")
		purchaseHistory.Use(middleware.AuthMiddleware())
		{
			purchaseHistory.GET("", handlers.GetUserPurchaseHistory)                   // GET /api/v1/purchase-history?page=1&limit=10&status=DELIVERED
			purchaseHistory.GET("/stats", handlers.GetUserPurchaseStats)               // GET /api/v1/purchase-history/stats
			purchaseHistory.GET("/recent", handlers.GetRecentPurchases)                // GET /api/v1/purchase-history/recent?limit=5
			purchaseHistory.GET("/search", handlers.SearchPurchaseHistory)             // GET /api/v1/purchase-history/search?q=samsung
			purchaseHistory.GET("/date-range", handlers.GetPurchaseHistoryByDateRange) // GET /api/v1/purchase-history/date-range?user_id=1&start_date=2024-01-01&end_date=2024-12-31
			purchaseHistory.GET("/:id", handlers.GetPurchaseHistoryByID)               // GET /api/v1/purchase-history/1?user_id=1
			purchaseHistory.GET("/can-review/:product_id", handlers.CanReviewProduct)  // GET /api/v1/purchase-history/can-review/1?user_id=1
		}

		// Order routes (requires authentication)
		orders := v1.Group("/orders")
		orders.Use(middleware.AuthMiddleware())
		{
			orders.GET("", handlers.GetUserOrders)                // GET /api/v1/orders
			orders.POST("", handlers.CreateOrder)                 // POST /api/v1/orders
			orders.GET("/stats", handlers.GetOrderStats)          // GET /api/v1/orders/stats
			orders.GET("/:id", handlers.GetOrderByID)             // GET /api/v1/orders/:id
			orders.PUT("/:id/status", handlers.UpdateOrderStatus) // PUT /api/v1/orders/:id/status (admin only)
		}

		// Wishlist routes (requires authentication)
		// wishlist := v1.Group("/wishlist")
		// {
		//     wishlist.GET("", handlers.GetWishlist)              // GET /api/v1/wishlist?user_id=1
		//     wishlist.POST("", handlers.AddToWishlist)           // POST /api/v1/wishlist?user_id=1
		//     wishlist.DELETE("/:id", handlers.RemoveFromWishlist) // DELETE /api/v1/wishlist/1?user_id=1
		// }
	}
}
