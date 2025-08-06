package handlers

import (
	"literally-backend/internal/models"
	"literally-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary Register a new user
// @Description Register a new user account
// @Tags authentication
// @Accept json
// @Produce json
// @Param user body object true "User registration data"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid input"
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	authResponse, err := services.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    authResponse,
		"message": "User registered successfully",
	})
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags authentication
// @Accept json
// @Produce json
// @Param credentials body object true "User login credentials"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized - Invalid credentials"
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	authResponse, err := services.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    authResponse,
		"message": "Login successful",
	})
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get the authenticated user's profile information
// @Tags profile
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{} "Profile retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /profile [get]
func GetProfile(c *gin.Context) {
	// Extract user ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	user, found := services.GetUserByID(userID.(uint))
	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    user.ToResponse(),
		"message": "Profile retrieved successfully",
	})
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the authenticated user's profile information
// @Tags profile
// @Accept json
// @Produce json
// @Security Bearer
// @Param profile body object true "Updated profile data"
// @Success 200 {object} map[string]interface{} "Profile updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /profile [put]
func UpdateProfile(c *gin.Context) {
	// Extract user ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := services.UpdateUser(userID.(uint), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    user.ToResponse(),
		"message": "Profile updated successfully",
	})
}

// RefreshToken godoc
// @Summary Refresh JWT token
// @Description Refresh the JWT token for authenticated user
// @Tags authentication
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{} "Token refreshed successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /auth/refresh [post]
func RefreshToken(c *gin.Context) {
	// Extract user info from current JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Get user from database to ensure user still exists
	user, found := services.GetUserByID(userID.(uint))
	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// For now, return a placeholder until we implement proper token refresh
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"user":  user.ToResponse(),
			"token": "refreshed_token_placeholder", // TODO: implement real refresh
		},
		"message": "Token refresh successful (placeholder implementation)",
	})
}

// GetUsers godoc
// @Summary Get all users
// @Description Get a list of all users (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{} "Users retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /users [get]
func GetUsers(c *gin.Context) {
	users := services.GetAllUsers()

	var userResponses []models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    userResponses,
		"message": "Users retrieved successfully",
	})
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get a specific user by their ID (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "User retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid user ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /users/{id} [get]
func GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	user, exists := services.GetUserByID(uint(id))
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    user.ToResponse(),
		"message": "User retrieved successfully",
	})
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user account (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param user body object true "User creation data"
// @Success 201 {object} map[string]interface{} "User created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := services.CreateUser(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    user.ToResponse(),
		"message": "User created successfully",
	})
}

// UpdateUser godoc
// @Summary Update user by ID
// @Description Update a specific user by their ID (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "User ID"
// @Param user body object true "Updated user data"
// @Success 200 {object} map[string]interface{} "User updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := services.UpdateUser(uint(id), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    user.ToResponse(),
		"message": "User updated successfully",
	})
}

// DeleteUser godoc
// @Summary Delete user by ID
// @Description Delete a specific user by their ID (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "User deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid user ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	err = services.DeleteUser(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

// UpdateUserStatus godoc
// @Summary Update user status
// @Description Update user status (ACTIVE/INACTIVE/SUSPENDED) by ID (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "User ID"
// @Param status body models.UpdateUserStatusRequest true "User status update data"
// @Success 200 {object} map[string]interface{} "User status updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /users/{id}/status [put]
func UpdateUserStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var req models.UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = services.UpdateUserStatus(uint(id), req.Status)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User status updated successfully",
	})
}
