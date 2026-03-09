package users

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// VerifyRequest expects the email and the OTP entered by the user
type VerifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

// verifyOTP handles POST /api/v1/users/verify
func verifyOTP(c *gin.Context) {
	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	db := getDB(c)
	coll := db.Users()

	// 1. Find a user matching BOTH the email and the exact OTP
	filter := bson.M{
		"email": req.Email,
		"otp":   req.OTP,
	}

	// 2. We want to update `is_verified` to true, and immediately clear the OTP so it can't be reused
	update := bson.M{
		"$set": bson.M{
			"is_verified": true,
			"updated_at":  time.Now().UTC(),
		},
		"$unset": bson.M{
			"otp": "", // Removes the OTP field from the database
		},
	}

	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process verification"})
		return
	}

	if result.MatchedCount == 0 {
		// If nothing matched, either the email is wrong, the user is already verified, or the OTP is wrong
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or OTP code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account successfully verified! You can now log in.",
	})
}
