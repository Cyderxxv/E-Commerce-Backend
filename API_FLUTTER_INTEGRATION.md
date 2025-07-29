# API Documentation for Flutter App Integration

## Base URL
```
http://localhost:8080/api/v1
```

## Database Schema Compatibility

This backend is designed to work with the database schema from your Flutter app:

### Database Tables Supported:
- ✅ users (with all fields: phone_number, photo, full_name, date_of_birth, address, gender, status)
- ✅ products (with rating, review_count, image_url, brand, is_featured, is_available)
- ✅ categories (with icon field)
- ✅ carts (user cart management)
- 🔄 orders (models created, handlers pending)
- 🔄 wishlists (models created, handlers pending)
- 🔄 reviews (models created, handlers pending)
- 🔄 notifications (models created, handlers pending)

## Authentication Flow

### 1. Register
**POST** `/auth/register`
```json
{
  "name": "John Doe",
  "email": "john@example.com", 
  "phone_number": "+1234567890",
  "password": "password123"
}
```

**Response:**
```json
{
  "data": {
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "phone_number": "+1234567890",
      "status": "ACTIVE",
      "created_at": "2025-07-29T15:13:12Z"
    },
    "token": "mock_jwt_token_..."
  },
  "message": "User registered successfully"
}
```

### 2. Login
**POST** `/auth/login`
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

## Product Management

### 1. Get Categories
**GET** `/categories`
```json
{
  "data": [
    {
      "id": 1,
      "name": "Phones",
      "icon": "phone",
      "created_at": "2025-07-29T15:13:12Z"
    }
  ]
}
```

### 2. Get Products
**GET** `/products`
**GET** `/products?category_id=1`
**GET** `/products?featured=true`
**GET** `/products?search=phone`

```json
{
  "data": [
    {
      "id": 1,
      "name": "iPhone 15 Pro",
      "description": "Latest iPhone with A17 Pro chip",
      "price": 999.99,
      "stock": 50,
      "image_url": "https://example.com/iphone15.jpg",
      "rating": 4.8,
      "review_count": 120,
      "category_id": 1,
      "brand": "Apple",
      "is_featured": true,
      "is_available": true,
      "created_at": "2025-07-29T15:13:12Z"
    }
  ]
}
```

### 3. Get Featured Products
**GET** `/products/featured`

### 4. Search Products
**GET** `/products/search?q=phone`

## Shopping Cart

### 1. Get Cart
**GET** `/cart?user_id=1`
```json
{
  "data": {
    "items": [
      {
        "id": 1,
        "user_id": 1,
        "product_id": 1,
        "quantity": 2,
        "product": {
          "id": 1,
          "name": "iPhone 15 Pro",
          "price": 999.99,
          "image_url": "https://example.com/iphone15.jpg"
        }
      }
    ],
    "total": 1999.98
  }
}
```

### 2. Add to Cart
**POST** `/cart?user_id=1`
```json
{
  "product_id": 1,
  "quantity": 2
}
```

### 3. Update Cart Item
**PUT** `/cart/1?user_id=1`
```json
{
  "quantity": 3
}
```

### 4. Remove from Cart
**DELETE** `/cart/1?user_id=1`

### 5. Clear Cart
**DELETE** `/cart?user_id=1`

## Profile Management

### 1. Get Profile
**GET** `/profile?user_id=1`

### 2. Update Profile
**PUT** `/profile?user_id=1`
```json
{
  "name": "John Updated",
  "phone_number": "+9876543210",
  "photo": "https://example.com/photo.jpg",
  "full_name": "John Doe Updated",
  "date_of_birth": "1990-01-01T00:00:00Z",
  "address": "123 Main St, City, Country",
  "gender": "Male"
}
```

## Error Responses

All endpoints return errors in this format:
```json
{
  "error": "Error message description"
}
```

Common HTTP status codes:
- 200: Success
- 201: Created successfully
- 400: Bad request (validation error)
- 401: Unauthorized
- 404: Not found
- 500: Internal server error

## Notes for Flutter Integration

1. **Authentication**: Currently using mock JWT tokens. In production, implement proper JWT authentication.

2. **User ID**: Currently passed as query parameter. In production, extract from JWT token.

3. **File Upload**: For product images and user photos, implement file upload endpoints.

4. **Pagination**: Add pagination for products list when database grows.

5. **Real-time Updates**: Consider WebSocket for real-time cart updates.

6. **Database**: Replace in-memory storage with PostgreSQL/MySQL database.

## Upcoming Endpoints

These are planned for future implementation:

### Orders
- `GET /orders?user_id=1` - Get user orders
- `POST /orders?user_id=1` - Create order from cart
- `GET /orders/:id?user_id=1` - Get order details

### Wishlist
- `GET /wishlist?user_id=1` - Get user wishlist
- `POST /wishlist?user_id=1` - Add to wishlist
- `DELETE /wishlist/:id?user_id=1` - Remove from wishlist

### Reviews
- `GET /products/:id/reviews` - Get product reviews
- `POST /products/:id/reviews?user_id=1` - Add review
- `PUT /reviews/:id?user_id=1` - Update review

### Notifications
- `GET /notifications?user_id=1` - Get user notifications
- `PUT /notifications/:id/read?user_id=1` - Mark as read
