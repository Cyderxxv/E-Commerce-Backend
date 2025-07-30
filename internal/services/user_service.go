package services

import (
	"errors"
	"literally-backend/configs"
	"literally-backend/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// GetAllUsers returns all users
func GetAllUsers() []models.User {
	var users []models.User
	configs.DB.Find(&users)
	return users
}

// GetUserByID returns a user by ID
func GetUserByID(id uint) (models.User, bool) {
	var user models.User
	result := configs.DB.First(&user, id)
	if result.Error != nil {
		return models.User{}, false
	}
	return user, true
}

// GetUserByEmail returns a user by email
func GetUserByEmail(email string) (models.User, bool) {
	var user models.User
	result := configs.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return models.User{}, false
	}
	return user, true
}

// CreateUser creates a new user
func CreateUser(req models.CreateUserRequest) (models.User, error) {
	// Check if email already exists
	var existingUser models.User
	if err := configs.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return models.User{}, errors.New("email already exists")
	}

	// Check if phone number already exists
	if err := configs.DB.Where("phone_number = ?", req.PhoneNumber).First(&existingUser).Error; err == nil {
		return models.User{}, errors.New("phone number already exists")
	}

	// Create new user
	user := models.User{
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password, // In production, hash the password
		Status:      "ACTIVE",
	}

	if err := configs.DB.Create(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}

// RegisterUser registers a new user
func RegisterUser(req models.RegisterRequest) (models.User, error) {
	// Check if email already exists
	var existingUser models.User
	if err := configs.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return models.User{}, errors.New("email already exists")
	}

	// Check if phone number already exists
	if err := configs.DB.Where("phone_number = ?", req.PhoneNumber).First(&existingUser).Error; err == nil {
		return models.User{}, errors.New("phone number already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	// Create new user
	user := models.User{
		Name:         req.Name,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		PasswordHash: string(hashedPassword),
		Status:       "ACTIVE",
	}

	if err := configs.DB.Create(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}

// LoginUser authenticates a user
func LoginUser(req models.LoginRequest) (models.User, error) {
	var user models.User
	result := configs.DB.Where("email = ?", req.Email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("invalid email or password")
		}
		return models.User{}, result.Error
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return models.User{}, errors.New("invalid email or password")
	}

	if user.Status != "ACTIVE" {
		return models.User{}, errors.New("account is not active")
	}

	return user, nil
}

// UpdateUser updates an existing user
func UpdateUser(id uint, req models.UpdateUserRequest) (models.User, error) {
	var user models.User

	// Find the user
	if err := configs.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}

	// Check if new email already exists (excluding current user)
	if req.Email != "" && req.Email != user.Email {
		var existingUser models.User
		if err := configs.DB.Where("email = ? AND id != ?", req.Email, id).First(&existingUser).Error; err == nil {
			return models.User{}, errors.New("email already exists")
		}
	}

	// Check if new phone number already exists (excluding current user)
	if req.PhoneNumber != "" && req.PhoneNumber != user.PhoneNumber {
		var existingUser models.User
		if err := configs.DB.Where("phone_number = ? AND id != ?", req.PhoneNumber, id).First(&existingUser).Error; err == nil {
			return models.User{}, errors.New("phone number already exists")
		}
	}

	// Update fields
	updates := make(map[string]interface{})

	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.PhoneNumber != "" {
		updates["phone_number"] = req.PhoneNumber
	}
	if req.Photo != "" {
		updates["photo"] = req.Photo
	}
	if req.FullName != "" {
		updates["full_name"] = req.FullName
	}
	if req.DateOfBirth != nil {
		updates["date_of_birth"] = req.DateOfBirth
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}
	if req.Gender != "" {
		updates["gender"] = req.Gender
	}

	// Update the user
	if err := configs.DB.Model(&user).Updates(updates).Error; err != nil {
		return models.User{}, err
	}

	// Fetch the updated user
	configs.DB.First(&user, id)

	return user, nil
}

// DeleteUser deletes a user by ID
func DeleteUser(id uint) error {
	result := configs.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
