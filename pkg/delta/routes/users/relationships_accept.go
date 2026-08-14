package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AcceptFriendRequest handles PUT /users/@me/relationships/:id
func AcceptFriendRequest(c *gin.Context) {
	// 1. Get the authenticated userID
	actorID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// 2. Get the TargetID from the URL parameter
	targetID := c.Param("id")
	if targetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "target_id is required in the URL"})
		return
	}

	db := getDB(c)
	ctx := c.Request.Context()

	// 3. Update the relationship type to 'Friend' (1)
	err := db.Relationships().AcceptRequest(ctx, actorID.(string), targetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to accept request: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request accepted successfully"})
}
