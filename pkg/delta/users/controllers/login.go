package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login handles the user authentication
func (ctrl *UserController) Login(c *gin.Context) {
	// 1. Parse JSON Request
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// 2. Timeout Context
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// 3. Call Service
	user, token, err := ctrl.Service.Login(ctx, req.Email, req.Password)
	if err != nil {
		if err.Error() == "invalid credentials" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed: " + err.Error()})
		return
	}

	// 4. Success Response
	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"message":    "Login successful",
		"token":      token,
		"user":       user,
	})
}
