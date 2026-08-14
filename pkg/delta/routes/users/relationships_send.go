package users

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendRelationshipRequest(c *gin.Context) {
	actorID, exists := c.Get("userId")
	actorIDStr, ok := actorID.(string)
	fmt.Printf("actorID: %s\n", actorIDStr)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid userID type"})
		return
	}

	fmt.Printf("actorID: %s\n", actorIDStr)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req struct {
		TargetID string `json:"target_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	db := getDB(c)
	relStore := db.Relationships()

	err := relStore.SendRequest(c, actorID.(string), req.TargetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Request sent successfully"})
}
