package users

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// VerifyRequest expects the email and the OTP entered by the user
type VerifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

// @Summary      Verify OTP Code
// @Description  Validates the 6-digit OTP code sent to the user's email and marks the account as verified.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body users.VerifyRequest true "OTP Verification Details"
// @Success      200  {object}  map[string]interface{} "Account verified successfully"
// @Failure      400  {object}  map[string]interface{} "Invalid input"
// @Failure      401  {object}  map[string]interface{} "Invalid email or OTP code"
// @Failure      500  {object}  map[string]interface{} "Failed to process verification"
// @Router       /users/verify [post]
func verifyOTP(c *gin.Context) {
	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	db := getDB(c)
	userStore := db.Users()

	success, err := userStore.VerifyUserOTP(ctx, req.Email, req.OTP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process verification"})
		return
	}

	if !success {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or OTP code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account successfully verified! You can now log in.",
	})
}
