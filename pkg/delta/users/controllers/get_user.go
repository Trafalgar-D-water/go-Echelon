package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetUserByID handles fetching a single user
func (ctrl *UserController) GetUserByID(c *gin.Context) {
	// 1. Get ID from URL
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// 2. Timeout Context
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// 3. Call Service
	user, err := ctrl.Service.GetUserByID(ctx, id)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user: " + err.Error()})
		return
	}

	// 4. Success Response
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
