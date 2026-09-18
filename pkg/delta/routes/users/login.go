package users

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-Echelon/go-Echelon/pkg/core/models"
	"github.com/go-Echelon/go-Echelon/pkg/delta/middleware"
	"github.com/go-Echelon/go-Echelon/pkg/delta/util"
)

// LoginRequest defines the expected JSON body for authentication.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// @Summary      Login a User
// @Description  Login a User, generates a verification OTP, and sends the OTP to the user's email asynchronously.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body users.LoginRequest true "User Login Details"
// @Success      201  {object}  map[string]interface{} "User created successfully"
// @Failure      400  {object}  map[string]interface{} "Invalid input"
// @Failure      409  {object}  map[string]interface{} "Email already registered"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /auth/session/login [post]
func login(c *gin.Context) {
	logger := middleware.Logger(c).With(
		"component", "delta.users",
		"operation", "login",
	)

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

	refreshToken, err := util.GenerateRefreshToken(user.ID.Hex())

	if err != nil {
		logger.Error("failed to generate refresh token", "error", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	hash := sha256.Sum256([]byte(refreshToken))
	hashedToken := hex.EncodeToString(hash[:])

	hashedRefreshToken, err := bcrypt.GenerateFromPassword(
		[]byte(hashedToken),
		bcrypt.DefaultCost,
	)
	if err != nil {
		logger.Error("failed to hash refresh token", "error", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	existingSession, err := sessionStore.GetSessionByUserID(ctx, user.ID.Hex())

	switch {
	case err == nil:
		err = sessionStore.UpdateSession(
			ctx,
			existingSession.ID.Hex(),
			string(hashedRefreshToken),
			expiresAt,
		)
	case errors.Is(err, mongo.ErrNoDocuments):
		_, err = sessionStore.CreateSession(ctx, &models.Session{
			ID:           primitive.NewObjectID(),
			UserID:       user.ID.Hex(),
			RefreshToken: string(hashedRefreshToken),
			UserAgent:    c.Request.UserAgent(),
			IP:           c.ClientIP(),
			ExpiresAt:    expiresAt,
			CreatedAt:    time.Now(),
		})
	}

	if err != nil {
		logger.Error(
			"failed to create or update user session",
			"operation", "upsert_session",
			"error", err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	accessToken, err := util.GenerateAccessToken(user.ID.Hex())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to genrate access token",
		})
		return
	}

	c.SetCookie(
		"refreshToken",
		refreshToken,
		7*24*60*60,
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message":     "login successful",
		"accessToken": accessToken,
		"user":        user,
	})
}
