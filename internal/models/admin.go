package models

import "time"

// Admin represents an admin user in the system
type Admin struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" binding:"required" gorm:"unique;not null"`
	Password     string    `json:"password,omitempty" binding:"required,min=6" gorm:"-"`
	PasswordHash string    `json:"-" gorm:"column:password_hash;not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// AdminResponse represents admin data returned to client (without password)
type AdminResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AdminLoginRequest represents the request body for admin login
type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AdminRegisterRequest represents the request body for admin registration
type AdminRegisterRequest struct {
	Username        string `json:"username" binding:"required,min=3"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

// AdminLoginResponse represents the response body for admin login
type AdminLoginResponse struct {
	Token string        `json:"token"`
	Admin AdminResponse `json:"admin"`
}
