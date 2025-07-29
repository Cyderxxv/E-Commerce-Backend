package services

import (
	"literally-backend/internal/models"
	"time"
)

// MockJWTService simulates JWT token generation
// In production, use a proper JWT library like github.com/golang-jwt/jwt/v5
type AuthService struct{}

// GenerateToken generates a mock JWT token
func (s *AuthService) GenerateToken(user models.User) string {
	// In production, implement proper JWT token generation
	// For now, return a mock token
	return "mock_jwt_token_" + user.Email + "_" + time.Now().Format("20060102150405")
}

// ValidateToken validates a mock JWT token
func (s *AuthService) ValidateToken(token string) (uint, bool) {
	// In production, implement proper JWT token validation
	// For now, extract user ID from mock token (simplified)
	if len(token) > 15 && token[:15] == "mock_jwt_token_" {
		// Mock validation - in real scenario, decode JWT and extract user ID
		return 1, true // Return mock user ID
	}
	return 0, false
}

// Login authenticates user and returns auth response
func Login(req models.LoginRequest) (models.AuthResponse, error) {
	user, err := LoginUser(req)
	if err != nil {
		return models.AuthResponse{}, err
	}

	authService := &AuthService{}
	token := authService.GenerateToken(user)

	return models.AuthResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

// Register creates new user and returns auth response
func Register(req models.RegisterRequest) (models.AuthResponse, error) {
	user, err := RegisterUser(req)
	if err != nil {
		return models.AuthResponse{}, err
	}

	authService := &AuthService{}
	token := authService.GenerateToken(user)

	return models.AuthResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}
