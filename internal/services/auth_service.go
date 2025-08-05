package services

import (
"errors"
"literally-backend/internal/models"
"os"
"time"

"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the JWT token claims
type JWTClaims struct {
UserID uint   `json:"user_id"`
Email  string `json:"email"`
jwt.RegisteredClaims
}

// AuthService handles JWT operations
type AuthService struct {
secretKey []byte
}

// NewAuthService creates a new AuthService instance
func NewAuthService() *AuthService {
secret := os.Getenv("JWT_SECRET")
if secret == "" {
secret = "your-super-secret-jwt-key-change-this-in-production"
}
return &AuthService{
secretKey: []byte(secret),
}
}

// GenerateToken generates a real JWT token
func (s *AuthService) GenerateToken(user models.User) (string, error) {
claims := JWTClaims{
UserID: user.ID,
Email:  user.Email,
RegisteredClaims: jwt.RegisteredClaims{
ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
IssuedAt:  jwt.NewNumericDate(time.Now()),
NotBefore: jwt.NewNumericDate(time.Now()),
Issuer:    "literally-backend",
Subject:   "user-auth",
},
}

token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
return token.SignedString(s.secretKey)
}

// ValidateToken validates a JWT token and returns user ID
func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
return nil, errors.New("invalid signing method")
}
return s.secretKey, nil
})

if err != nil {
return nil, err
}

if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
return claims, nil
}

return nil, errors.New("invalid token")
}

// RefreshToken generates a new token with extended expiry
func (s *AuthService) RefreshToken(tokenString string) (string, error) {
claims, err := s.ValidateToken(tokenString)
if err != nil {
return "", err
}

user, found := GetUserByID(claims.UserID)
if !found {
return "", errors.New("user not found")
}

return s.GenerateToken(user)
}

// Login authenticates user and returns auth response with real JWT
func Login(req models.LoginRequest) (models.AuthResponse, error) {
user, err := LoginUser(req)
if err != nil {
return models.AuthResponse{}, err
}

authService := NewAuthService()
token, err := authService.GenerateToken(user)
if err != nil {
return models.AuthResponse{}, err
}

return models.AuthResponse{
User:  user.ToResponse(),
Token: token,
}, nil
}

// Register creates new user and returns auth response with real JWT
func Register(req models.RegisterRequest) (models.AuthResponse, error) {
user, err := RegisterUser(req)
if err != nil {
return models.AuthResponse{}, err
}

authService := NewAuthService()
token, err := authService.GenerateToken(user)
if err != nil {
return models.AuthResponse{}, err
}

return models.AuthResponse{
User:  user.ToResponse(),
Token: token,
}, nil
}
