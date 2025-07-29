package services

import (
	"errors"
	"literally-backend/internal/models"
	"time"
)

// In-memory storage for demo purposes
// In production, you would use a database like PostgreSQL, MySQL, etc.
var users []models.User
var userIDCounter uint = 1

// GetAllUsers returns all users
func GetAllUsers() []models.User {
	return users
}

// GetUserByID returns a user by ID
func GetUserByID(id uint) (models.User, bool) {
	for _, user := range users {
		if user.ID == id {
			return user, true
		}
	}
	return models.User{}, false
}

// GetUserByEmail returns a user by email
func GetUserByEmail(email string) (models.User, bool) {
	for _, user := range users {
		if user.Email == email {
			return user, true
		}
	}
	return models.User{}, false
}

// CreateUser creates a new user
func CreateUser(req models.CreateUserRequest) (models.User, error) {
	// Check if email already exists
	for _, user := range users {
		if user.Email == req.Email {
			return models.User{}, errors.New("email already exists")
		}
		if user.PhoneNumber == req.PhoneNumber {
			return models.User{}, errors.New("phone number already exists")
		}
	}

	// Create new user
	user := models.User{
		ID:          userIDCounter,
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password, // In production, hash the password
		Status:      "ACTIVE",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	users = append(users, user)
	userIDCounter++

	return user, nil
}

// RegisterUser registers a new user
func RegisterUser(req models.RegisterRequest) (models.User, error) {
	// Check if email already exists
	for _, user := range users {
		if user.Email == req.Email {
			return models.User{}, errors.New("email already exists")
		}
		if user.PhoneNumber == req.PhoneNumber {
			return models.User{}, errors.New("phone number already exists")
		}
	}

	// Create new user
	user := models.User{
		ID:          userIDCounter,
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password, // In production, hash the password
		Status:      "ACTIVE",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	users = append(users, user)
	userIDCounter++

	return user, nil
}

// LoginUser authenticates a user
func LoginUser(req models.LoginRequest) (models.User, error) {
	for _, user := range users {
		if user.Email == req.Email && user.Password == req.Password {
			if user.Status != "ACTIVE" {
				return models.User{}, errors.New("account is not active")
			}
			return user, nil
		}
	}
	return models.User{}, errors.New("invalid email or password")
}

// UpdateUser updates an existing user
func UpdateUser(id uint, req models.UpdateUserRequest) (models.User, error) {
	for i, user := range users {
		if user.ID == id {
			// Update fields if provided
			if req.Name != "" {
				users[i].Name = req.Name
			}
			if req.Email != "" {
				// Check if new email already exists (excluding current user)
				for _, otherUser := range users {
					if otherUser.Email == req.Email && otherUser.ID != id {
						return models.User{}, errors.New("email already exists")
					}
				}
				users[i].Email = req.Email
			}
			if req.PhoneNumber != "" {
				// Check if new phone number already exists (excluding current user)
				for _, otherUser := range users {
					if otherUser.PhoneNumber == req.PhoneNumber && otherUser.ID != id {
						return models.User{}, errors.New("phone number already exists")
					}
				}
				users[i].PhoneNumber = req.PhoneNumber
			}
			if req.Photo != "" {
				users[i].Photo = req.Photo
			}
			if req.FullName != "" {
				users[i].FullName = req.FullName
			}
			if req.DateOfBirth != nil {
				users[i].DateOfBirth = req.DateOfBirth
			}
			if req.Address != "" {
				users[i].Address = req.Address
			}
			if req.Gender != "" {
				users[i].Gender = req.Gender
			}
			users[i].UpdatedAt = time.Now()

			return users[i], nil
		}
	}
	return models.User{}, errors.New("user not found")
}

// DeleteUser deletes a user by ID
func DeleteUser(id uint) error {
	for i, user := range users {
		if user.ID == id {
			// Remove user from slice
			users = append(users[:i], users[i+1:]...)
			return nil
		}
	}
	return errors.New("user not found")
}
