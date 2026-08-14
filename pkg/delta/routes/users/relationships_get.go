package users

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// FetchMyRelationships handles GET /users/@me/relationships
func FetchMyRelationships(c *gin.Context) {
	// 1. Get authenticated userID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	db := getDB(c)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// 2. Fetch all relationships
	relationships, err := db.Relationships().GetUserRelationships(ctx, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch relationships"})
		return
	}

	// 3. Return (handle null slice)
	if relationships == nil {
		c.JSON(http.StatusOK, gin.H{"data": []interface{}{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": relationships})
}
