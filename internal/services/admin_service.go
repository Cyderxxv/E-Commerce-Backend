package services

import (
	"errors"
	"literally-backend/configs"
	"literally-backend/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AdminJWTClaims represents the JWT token claims for admin
type AdminJWTClaims struct {
	AdminID  uint   `json:"admin_id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// AdminAuthService handles JWT operations for admin
type AdminAuthService struct {
	secretKey []byte
}

// NewAdminAuthService creates a new AdminAuthService instance
func NewAdminAuthService() *AdminAuthService {
	secret := "your-super-secret-admin-jwt-key-change-this-in-production"
	return &AdminAuthService{
		secretKey: []byte(secret),
	}
}

// GenerateAdminToken generates a JWT token for admin
func (s *AdminAuthService) GenerateAdminToken(admin models.Admin) (string, error) {
	claims := AdminJWTClaims{
		AdminID:  admin.ID,
		Username: admin.Username,
		IsAdmin:  true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "literally-backend",
			Subject:   "admin-auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ValidateAdminToken validates a JWT token and returns admin claims
func (s *AdminAuthService) ValidateAdminToken(tokenString string) (*AdminJWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AdminJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AdminJWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// AdminLogin authenticates admin and returns auth response
func AdminLogin(req models.AdminLoginRequest) (models.AdminLoginResponse, error) {
	admin, err := LoginAdmin(req)
	if err != nil {
		return models.AdminLoginResponse{}, err
	}

	authService := NewAdminAuthService()
	token, err := authService.GenerateAdminToken(admin)
	if err != nil {
		return models.AdminLoginResponse{}, err
	}

	return models.AdminLoginResponse{
		Admin: models.AdminResponse{
			ID:        admin.ID,
			Username:  admin.Username,
			CreatedAt: admin.CreatedAt,
			UpdatedAt: admin.UpdatedAt,
		},
		Token: token,
	}, nil
}

// AdminRegister creates new admin and returns auth response
func AdminRegister(req models.AdminRegisterRequest) (models.AdminLoginResponse, error) {
	// Validate password confirmation
	if req.Password != req.ConfirmPassword {
		return models.AdminLoginResponse{}, errors.New("passwords do not match")
	}

	admin, err := RegisterAdmin(req)
	if err != nil {
		return models.AdminLoginResponse{}, err
	}

	authService := NewAdminAuthService()
	token, err := authService.GenerateAdminToken(admin)
	if err != nil {
		return models.AdminLoginResponse{}, err
	}

	return models.AdminLoginResponse{
		Admin: models.AdminResponse{
			ID:        admin.ID,
			Username:  admin.Username,
			CreatedAt: admin.CreatedAt,
			UpdatedAt: admin.UpdatedAt,
		},
		Token: token,
	}, nil
}

// LoginAdmin validates admin credentials
func LoginAdmin(req models.AdminLoginRequest) (models.Admin, error) {
	var admin models.Admin
	db := configs.GetDB()

	// Find admin by username
	if err := db.Where("username = ?", req.Username).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Admin{}, errors.New("invalid username or password")
		}
		return models.Admin{}, err
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(req.Password)); err != nil {
		return models.Admin{}, errors.New("invalid username or password")
	}

	return admin, nil
}

// RegisterAdmin creates a new admin account
func RegisterAdmin(req models.AdminRegisterRequest) (models.Admin, error) {
	db := configs.GetDB()

	// Check if admin already exists
	var existingAdmin models.Admin
	if err := db.Where("username = ?", req.Username).First(&existingAdmin).Error; err == nil {
		return models.Admin{}, errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.Admin{}, err
	}

	// Create new admin
	admin := models.Admin{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
	}

	if err := db.Create(&admin).Error; err != nil {
		return models.Admin{}, err
	}

	return admin, nil
}

// GetAdminByID retrieves admin by ID
func GetAdminByID(adminID uint) (models.AdminResponse, bool) {
	var admin models.Admin
	db := configs.GetDB()

	if err := db.First(&admin, adminID).Error; err != nil {
		return models.AdminResponse{}, false
	}

	return models.AdminResponse{
		ID:        admin.ID,
		Username:  admin.Username,
		CreatedAt: admin.CreatedAt,
		UpdatedAt: admin.UpdatedAt,
	}, true
}

// CreateAdmin creates a new admin (for seeding or admin management)
func CreateAdmin(username, password string) error {
	db := configs.GetDB()

	// Check if admin already exists
	var existingAdmin models.Admin
	if err := db.Where("username = ?", username).First(&existingAdmin).Error; err == nil {
		return errors.New("admin already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create admin
	admin := models.Admin{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	return db.Create(&admin).Error
}
