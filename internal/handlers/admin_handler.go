package handlers

import (
	"literally-backend/internal/models"
	"literally-backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminLogin godoc
// @Summary Admin login
// @Description Authenticate admin and return JWT token
// @Tags admin-auth
// @Accept json
// @Produce json
// @Param credentials body models.AdminLoginRequest true "Admin login credentials"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized - Invalid credentials"
// @Router /admin/auth/login [post]
func AdminLogin(c *gin.Context) {
	var req models.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	authResponse, err := services.AdminLogin(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    authResponse,
		"message": "Admin login successful",
	})
}

// AdminRegister godoc
// @Summary Admin register
// @Description Register a new admin account and return JWT token
// @Tags admin-auth
// @Accept json
// @Produce json
// @Param admin body models.AdminRegisterRequest true "Admin registration data"
// @Success 201 {object} map[string]interface{} "Registration successful"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid input"
// @Router /admin/auth/register [post]
func AdminRegister(c *gin.Context) {
	var req models.AdminRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	authResponse, err := services.AdminRegister(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    authResponse,
		"message": "Admin registered successfully",
	})
}

// GetAdminProfile godoc
// @Summary Get admin profile
// @Description Get the authenticated admin's profile information
// @Tags admin-auth
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{} "Profile retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Admin not found"
// @Router /admin/profile [get]
func GetAdminProfile(c *gin.Context) {
	// Extract admin ID from JWT token
	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Admin not authenticated",
		})
		return
	}

	admin, found := services.GetAdminByID(adminID.(uint))
	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Admin not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    admin,
		"message": "Admin profile retrieved successfully",
	})
}
