package models

import "time"

// User represents a user in the system
type User struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name" binding:"required"`
	Email       string     `json:"email" binding:"required,email" gorm:"unique"`
	PhoneNumber string     `json:"phone_number" binding:"required" gorm:"unique"`
	Password    string     `json:"password,omitempty" binding:"required,min=6"`
	Photo       string     `json:"photo,omitempty"`
	FullName    string     `json:"full_name,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	Address     string     `json:"address,omitempty"`
	Gender      string     `json:"gender,omitempty"`
	Status      string     `json:"status" gorm:"default:ACTIVE"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// UserResponse represents user data returned to client (without password)
type UserResponse struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	PhoneNumber string     `json:"phone_number"`
	Photo       string     `json:"photo,omitempty"`
	FullName    string     `json:"full_name,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	Address     string     `json:"address,omitempty"`
	Gender      string     `json:"gender,omitempty"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required,min=6"`
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required,min=6"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Name        string     `json:"name,omitempty"`
	Email       string     `json:"email,omitempty"`
	PhoneNumber string     `json:"phone_number,omitempty"`
	Photo       string     `json:"photo,omitempty"`
	FullName    string     `json:"full_name,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	Address     string     `json:"address,omitempty"`
	Gender      string     `json:"gender,omitempty"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:          u.ID,
		Name:        u.Name,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		Photo:       u.Photo,
		FullName:    u.FullName,
		DateOfBirth: u.DateOfBirth,
		Address:     u.Address,
		Gender:      u.Gender,
		Status:      u.Status,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
