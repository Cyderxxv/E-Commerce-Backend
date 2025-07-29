# Literally Backend

A modern Go backend API built with Gin framework.

## Features

- 🚀 Fast HTTP server with Gin
- 📦 RESTful API design
- 🔧 Middleware support (CORS, Logging)
- 📝 Environment configuration
- 🏗️ Clean architecture (handlers, services, models)
- 📊 JSON responses
- ⚡ In-memory storage (demo)

## Project Structure

```
literally-backend/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── handlers/            # HTTP handlers
│   │   ├── user_handler.go
│   │   └── product_handler.go
│   ├── models/              # Data models
│   │   ├── user.go
│   │   └── product.go
│   ├── services/            # Business logic
│   │   ├── user_service.go
│   │   └── product_service.go
│   └── middleware/          # HTTP middleware
│       └── middleware.go
├── pkg/
│   └── utils/               # Utility functions
│       └── utils.go
├── configs/                 # Configuration files
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

3. Copy environment variables:
```bash
copy .env.example .env
```

4. Run the application:
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

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
- [godotenv](https://github.com/joho/godotenv) - Environment variables loader

## Future Enhancements

- [ ] Database integration (PostgreSQL/MySQL)
- [ ] Authentication & Authorization (JWT)
- [ ] Input validation & sanitization
- [ ] Logging improvements
- [ ] Unit tests
- [ ] Docker support
- [ ] API documentation (Swagger)
- [ ] Rate limiting
- [ ] Caching

## License

MIT License
