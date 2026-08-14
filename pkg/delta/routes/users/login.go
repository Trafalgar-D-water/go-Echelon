package users

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-Echelon/go-Echelon/pkg/core/models"
	"github.com/go-Echelon/go-Echelon/pkg/delta/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest defines the expected JSON body for authentication.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// @Summary      Login a User
// @Description  Authenticates a user, generates access and refresh tokens, and secures the session.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body users.LoginRequest true "User Login Details"
// @Success      200  {object}  map[string]interface{} "Login successful"
// @Failure      400  {object}  map[string]interface{} "Invalid input"
// @Failure      401  {object}  map[string]interface{} "Invalid email or password"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /auth/session/login [post]
func login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	db := getDB(c)
	userStore := db.Users()
	sessionStore := db.Sessions()

	// Find user by email
	user, err := userStore.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	// Verify password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// 1. Generate Access Token
	accessToken, err := util.GenerateAccessToken(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	// 2. Generate Refresh Token
	refreshToken, err := util.GenerateRefreshToken(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// 3. Hash Refresh Token for storage
	hash := sha256.Sum256([]byte(refreshToken))
	hashedToken := hex.EncodeToString(hash[:])

	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(hashedToken), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash refresh token"})
		return
	}

	// 4. Create Session in Database
	session := &models.Session{
		ID:           primitive.NewObjectID(),
		UserID:       user.ID.Hex(),
		RefreshToken: string(hashedRefreshToken),
		UserAgent:    c.Request.UserAgent(),
		IP:           c.ClientIP(),
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		CreatedAt:    time.Now(),
	}

	_, err = sessionStore.CreateSession(ctx, session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	// 5. Set Refresh Token in HttpOnly Cookie
	c.SetCookie(
		"refreshToken",
		refreshToken,
		7*24*60*60, // 7 days in seconds
		"/",
		"",
		true, // Secure (use true if SSL is enabled)
		true, // HttpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"message":     "Login successful",
		"accessToken": accessToken,
		"user": gin.H{
			"id":       user.ID.Hex(),
			"username": user.Username,
			"email":    user.Email,
		},
	})
}
