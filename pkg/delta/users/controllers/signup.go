package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-Echelon/go-Echelon/pkg/delta/users/services"
)

// UserController holds the service needed for logic
type UserController struct {
	Service *services.UserService
}

// NewUserController creates a new controller instance
func NewUserController(service *services.UserService) *UserController {
	return &UserController{Service: service}
}

// SignUpRequest defines what we expect from the frontend
type SignUpRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// SignUp handles the user registration request
// Flow:
// 1. Parse JSON body
// 2. Call Service to create user
// 3. Return response
func (ctrl *UserController) SignUp(c *gin.Context) {
	// 1. Parse JSON Request
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// 2. set a timeout for the operation (5 seconds)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// 3. Call Service Layer (Business Logic)
	user, err := ctrl.Service.SignUp(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		// Handle errors
		if err.Error() == "user already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": "This email is already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong: " + err.Error()})
		return
	}

	// 4. Success Response
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}
