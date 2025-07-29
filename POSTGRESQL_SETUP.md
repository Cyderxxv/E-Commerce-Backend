# PostgreSQL Integration Guide

## Cài đặt PostgreSQL

### 1. Download và cài đặt PostgreSQL
- Tải từ: https://www.postgresql.org/download/windows/
- Chạy installer và làm theo hướng dẫn
- Nhớ password cho user `postgres`

### 2. Tạo Database
Mở PostgreSQL command line (psql) hoặc pgAdmin và chạy:

```sql
-- Tạo database
CREATE DATABASE literally_backend;

-- Tạo user (optional, có thể dùng postgres user)
CREATE USER literally_user WITH PASSWORD '123456';

-- Cấp quyền
GRANT ALL PRIVILEGES ON DATABASE literally_backend TO literally_user;
```

### 3. Cập nhật file .env
Sửa thông tin database trong file `.env`:

```bash
DB_HOST=localhost
DB_USER=postgres          # hoặc literally_user
DB_PASSWORD=123456        # password bạn đã đặt
DB_NAME=literally_backend
DB_PORT=5432
DB_SSLMODE=disable
```

### 4. Chạy Server
```bash
go run cmd/server/main.go
```

## Features đã tích hợp

### Database Features:
- ✅ Auto-migration: Tự động tạo tables khi khởi động
- ✅ Seed data: Tự động tạo dữ liệu mẫu
- ✅ Connection pooling: Quản lý kết nối hiệu quả
- ✅ Environment configuration: Cấu hình qua .env

### API Endpoints:
- ✅ Authentication: Register/Login với PostgreSQL
- ✅ User management: Profile với database
- ✅ Product management: CRUD với database  
- ✅ Cart system: Persistent cart trong database
- ✅ Error handling: Proper error responses

### Data Models:
- ✅ User: phone_number, photo, full_name (Flutter compatible)
- ✅ Product: image_url, rating, is_featured
- ✅ Cart: Persistent shopping cart
- ✅ GORM tags: Proper database mapping

## Testing the API

### 1. Health Check
```bash
GET http://localhost:8080/health
```

### 2. Register User
```bash
POST http://localhost:8080/api/v1/auth/register
Content-Type: application/json

{
    "email": "test@example.com",
    "password": "123456",
    "full_name": "Test User",
    "phone_number": "0123456789"
}
```

### 3. Login
```bash
POST http://localhost:8080/api/v1/auth/login
Content-Type: application/json

{
    "email": "test@example.com",
    "password": "123456"
}
```

### 4. Get Products
```bash
GET http://localhost:8080/api/v1/products
```

## Database Schema

Server sẽ tự động tạo các tables:
- `users`: User information
- `products`: Product catalog  
- `carts`: Shopping cart items
- `orders`: Order history
- `order_items`: Order details
- `reviews`: Product reviews

## Troubleshooting

### Connection Error:
1. Kiểm tra PostgreSQL service đang chạy
2. Kiểm tra thông tin trong .env file
3. Kiểm tra firewall/port 5432

### Migration Error:
1. Kiểm tra user có quyền tạo table
2. Kiểm tra database tồn tại
3. Xem log để biết lỗi cụ thể
