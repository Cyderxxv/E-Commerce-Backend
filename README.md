# Literally Backend

A modern Go e-commerce backend API built with Gin framework and PostgreSQL.

## Features

- 🚀 Fast HTTP server with Gin
- 📦 RESTful API design
- 🔧 Middleware support (CORS, Logging)
- 📝 Environment configuration
- 🏗️ Clean architecture (handlers, services, models)
- 📊 JSON responses
- 🗄️ PostgreSQL database integration with GORM
- 🔐 Password hashing with bcrypt
- 🛒 Complete e-commerce functionality
- 📋 Purchase history tracking system
- 🔍 Advanced filtering and search capabilities
- 📊 Statistics and analytics

## Project Structure

```
literally-backend/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── handlers/            # HTTP handlers
│   │   ├── user_handler.go
│   │   ├── product_handler.go
│   │   ├── category_handler.go
│   │   ├── cart_handler.go
│   │   └── purchase_history_handler.go
│   ├── models/              # Data models
│   │   ├── user.go
│   │   ├── product.go
│   │   ├── category.go
│   │   ├── cart.go
│   │   └── purchase_history.go
│   ├── services/            # Business logic
│   │   ├── user_service.go
│   │   ├── product_service.go
│   │   ├── category_service.go
│   │   ├── cart_service.go
│   │   └── purchase_history_service.go
│   └── middleware/          # HTTP middleware
│       └── middleware.go
├── configs/                 # Configuration files
│   ├── database.go         # Database connection
│   └── migration.go        # Database migrations & seeding
├── pkg/
│   └── utils/               # Utility functions
│       └── utils.go
├── .env.example            # Environment variables example
├── .env                    # Environment variables
├── go.mod                  # Go module file
├── go.sum                  # Go dependencies
└── README.md               # This file
```

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd literally-backend
```

2. Install dependencies:
```bash
go mod download
```

3. Set up PostgreSQL database:
```bash
# Create database
createdb TestDB

# Or using psql
psql -U postgres
CREATE DATABASE TestDB;
\q
```

4. Copy environment variables:
```bash
copy .env.example .env
```

5. Configure database connection in `.env`:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=TestDB
DB_SSLMODE=disable
```

6. Run the application:
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080` and automatically:
- Connect to PostgreSQL database
- Run database migrations
- Seed sample data

## API Endpoints

### Health Check
- `GET /health` - Server health check

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - User login

### Profile Management
- `GET /api/v1/profile?user_id=1` - Get user profile (requires authentication)
- `PUT /api/v1/profile?user_id=1` - Update user profile (requires authentication)

### User Management (Admin)
- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/:id` - Get user by ID
- `POST /api/v1/users` - Create new user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

### Categories
- `GET /api/v1/categories` - Get all categories
- `GET /api/v1/categories/:id/products` - Get products by category

### Products
- `GET /api/v1/products` - Get all products
- `GET /api/v1/products?category_id=1` - Get products by category
- `GET /api/v1/products?featured=true` - Get featured products only
- `GET /api/v1/products?search=phone` - Search products
- `GET /api/v1/products/featured` - Get featured products
- `GET /api/v1/products/search?q=phone` - Search products
- `GET /api/v1/products/:id` - Get product by ID
- `POST /api/v1/products` - Create new product (admin)
- `PUT /api/v1/products/:id` - Update product (admin)
- `DELETE /api/v1/products/:id` - Delete product (admin)

### Shopping Cart
- `GET /api/v1/cart?user_id=1` - Get user's cart
- `POST /api/v1/cart?user_id=1` - Add item to cart
- `PUT /api/v1/cart/:id?user_id=1` - Update cart item quantity
- `DELETE /api/v1/cart/:id?user_id=1` - Remove item from cart
- `DELETE /api/v1/cart?user_id=1` - Clear all cart items

### Purchase History
- `GET /api/v1/purchase-history?user_id=1` - Get user's purchase history with filtering
- `GET /api/v1/purchase-history/stats?user_id=1` - Get purchase statistics
- `GET /api/v1/purchase-history/recent?user_id=1&limit=5` - Get recent purchases
- `GET /api/v1/purchase-history/search?user_id=1&q=Samsung` - Search purchase history
- `GET /api/v1/purchase-history/date-range?user_id=1&start_date=2025-01-01&end_date=2025-12-31` - Get purchases by date range
- `GET /api/v1/purchase-history/:id?user_id=1` - Get specific purchase details
- `GET /api/v1/purchase-history/can-review/:product_id?user_id=1` - Check if user can review product

#### Purchase History Filtering Options
- `status` - Filter by order status (DELIVERED, PROCESSING, SHIPPED, CANCELLED)
- `payment_method` - Filter by payment method (Cash, Credit Card, Bank Transfer, Installment)
- `is_installment` - Filter installment purchases (true/false)
- `start_date` - Filter from date (YYYY-MM-DD format)
- `end_date` - Filter to date (YYYY-MM-DD format)
- `product_name` - Filter by product name (partial match)
- `page` - Page number for pagination (default: 1)
- `limit` - Items per page (default: 10)

## Example API Usage

### Register User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone_number": "+1234567890",
    "password": "password123"
  }'
```

### Login User
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Get Categories
```bash
curl http://localhost:8080/api/v1/categories
```

### Get Featured Products
```bash
curl http://localhost:8080/api/v1/products/featured
```

### Search Products
```bash
curl "http://localhost:8080/api/v1/products/search?q=phone"
```

### Add to Cart
```bash
curl -X POST "http://localhost:8080/api/v1/cart?user_id=1" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 1,
    "quantity": 2
  }'
```

### Get Cart
```bash
curl "http://localhost:8080/api/v1/cart?user_id=1"
```

### Get Purchase History
```bash
curl "http://localhost:8080/api/v1/purchase-history?user_id=1"
```

### Get Purchase History with Filters
```bash
# Filter by delivered orders only
curl "http://localhost:8080/api/v1/purchase-history?user_id=1&status=DELIVERED"

# Filter by installment purchases
curl "http://localhost:8080/api/v1/purchase-history?user_id=1&is_installment=true"

# Filter by date range
curl "http://localhost:8080/api/v1/purchase-history?user_id=1&start_date=2025-07-01&end_date=2025-07-31"
```

### Get Purchase Statistics
```bash
curl "http://localhost:8080/api/v1/purchase-history/stats?user_id=1"
```

### Search Purchase History
```bash
curl "http://localhost:8080/api/v1/purchase-history/search?user_id=1&q=Samsung"
```

### Get Recent Purchases
```bash
curl "http://localhost:8080/api/v1/purchase-history/recent?user_id=1&limit=3"
```

## API Response Examples

### Purchase History Response
```json
{
  "data": [
    {
      "id": 1,
      "product_id": 1,
      "product_name": "Samsung Galaxy S25 Edge (12/256GB)",
      "product_image_url": "https://cdn.example.com/samsung-s25.jpg",
      "quantity": 1,
      "unit_price": 25650600,
      "total_price": 25650600,
      "order_status": "DELIVERED",
      "status_display": "Đã giao hàng",
      "payment_method": "Cash",
      "is_installment": false,
      "purchase_date": "2025-07-15T16:54:45.549616+07:00",
      "delivery_date": "2025-07-20T16:54:45.549616+07:00",
      "shipping_address": "Ho Chi Minh City, Vietnam",
      "days_since_purchase": 15,
      "can_review": true,
      "can_reorder": true
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 3,
    "total_items": 25,
    "per_page": 10
  },
  "message": "Purchase history retrieved successfully"
}
```

### Purchase Statistics Response
```json
{
  "data": {
    "total_orders": 5,
    "total_amount": 100731000,
    "delivered_orders": 3,
    "pending_orders": 1,
    "cancelled_orders": 0,
    "avg_order_value": 20146200
  },
  "message": "Purchase statistics retrieved successfully"
}
```

## Error Handling

The API returns consistent error responses in the following format:

```json
{
  "error": "Error description message"
}
```

### Common HTTP Status Codes
- `200 OK` - Request successful
- `400 Bad Request` - Invalid request parameters
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

### Example Error Responses

#### Invalid User ID
```json
{
  "error": "Invalid user ID"
}
```

#### Search Term Required
```json
{
  "error": "Search term is required"
}
```

#### Database Error
```json
{
  "error": "Failed to retrieve purchase history"
}
```

## Development

### Build
```bash
go build -o bin/server cmd/server/main.go
```

### Run
```bash
./bin/server
```

### Test
```bash
go test ./...
```

## Dependencies

- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [GORM](https://gorm.io/) - Object-relational mapping library for Go
- [PostgreSQL Driver](https://github.com/lib/pq) - PostgreSQL driver for Go
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) - Password hashing
- [godotenv](https://github.com/joho/godotenv) - Environment variables loader

## Database Schema

The application uses PostgreSQL with the following main tables:

### Users
- User authentication and profile information
- Password hashing with bcrypt
- Profile fields: name, email, phone, address, etc.

### Categories & Products
- Product catalog with categories
- Product details: name, description, price, stock, images
- Category-based product organization

### Shopping Cart
- User shopping cart management
- Cart items with quantities

### Purchase History
- Complete purchase tracking system
- Order status management (PROCESSING, SHIPPED, DELIVERED, CANCELLED)
- Payment method tracking (Cash, Credit Card, Bank Transfer, Installment)
- Installment purchase support
- Delivery tracking with tracking numbers
- Business logic for reviews and reorders

## Purchase History Features

### Status Management
- **PROCESSING**: Đang xử lý - Order is being processed
- **SHIPPED**: Đang vận chuyển - Order has been shipped
- **DELIVERED**: Đã giao hàng - Order delivered successfully
- **CANCELLED**: Đã hủy - Order cancelled

### Payment Methods
- **Cash**: Tiền mặt
- **Credit Card**: Thẻ tín dụng
- **Bank Transfer**: Chuyển khoản ngân hàng
- **Installment**: Trả góp

### Business Logic
- **Can Review**: Users can review products after delivery
- **Can Reorder**: Users can reorder delivered or cancelled items
- **Installment Tracking**: Monthly payment amounts and duration
- **Days Since Purchase**: Automatic calculation of purchase age

## Future Enhancements

- [ ] Authentication & Authorization (JWT tokens)
- [ ] Input validation & sanitization improvements
- [ ] Enhanced logging and monitoring
- [ ] Comprehensive unit tests
- [ ] Docker support
- [ ] API documentation (Swagger)
- [ ] Rate limiting
- [ ] Redis caching
- [ ] Order management system
- [ ] Payment gateway integration
- [ ] Email notifications
- [ ] File upload for product images
- [ ] Admin dashboard APIs
- [ ] Real-time notifications
- [ ] Inventory management

## License

MIT License
