package users

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-Echelon/go-Echelon/internal/email"
	"github.com/go-Echelon/go-Echelon/pkg/core/models"
	"github.com/go-Echelon/go-Echelon/pkg/delta/middleware"
	"github.com/go-Echelon/go-Echelon/pkg/delta/util"
)

// SignUpRequest defines the expected JSON body for registration.
type SignUpRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	DOB      string `json:"dob" binding:"required"`
}

// @Summary      Register a new user
// @Description  Creates a new user account, generates a verification OTP, and sends the OTP to the user's email asynchronously.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body users.SignUpRequest true "User Registration Details"
// @Success      201  {object}  map[string]interface{} "User created successfully"
// @Failure      400  {object}  map[string]interface{} "Invalid input"
// @Failure      409  {object}  map[string]interface{} "Email already registered"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /users [post]
func create(c *gin.Context) {
	logger := middleware.Logger(c).With(
		"component", "delta.users",
		"operation", "create_user",
	)

	logger.Info("registration requested")

	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("registration validation failed", "error", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input",
		})
		return
	}

	dob, err := time.Parse("2006-01-02", req.DOB)
	if err != nil {
		logger.Warn("registration date of birth validation failed", "error", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid date of birth",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	db := getDB(c)
	userStore := db.Users()
	sessionStore := db.Sessions()

	count, err := userStore.CountByEmail(ctx, req.Email)
	if err != nil {
		logger.Error(
			"failed to check whether email is registered",
			"operation", "count_user_by_email",
			"error", err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	if count > 0 {
		logger.Warn("registration rejected because email is already registered")

		c.JSON(http.StatusConflict, gin.H{
			"error": "this email is already registered",
		})
		return
	}

	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		logger.Error("failed to hash user password", "error", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	otp, err := util.GenerateVerificationOTP()
	if err != nil {
		logger.Error("failed to generate verification OTP", "error", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	now := time.Now().UTC()
	otpExpiresAt := now.Add(util.VerificationOTPLifetime)

	user := &models.User{
		ID:           primitive.NewObjectID(),
		Username:     req.Username,
		Email:        req.Email,
		Password:     string(hashedBytes),
		DOB:          dob,
		IsVerified:   false,
		OTP:          otp,
		OTPExpiresAt: &otpExpiresAt,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := userStore.CreateUser(ctx, user); err != nil {
		logger.Error(
			"failed to create user",
			"operation", "create_user",
			"error", err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	logger = logger.With("user_id", user.ID.Hex())

	if err := email.SendOTP(user.Email, otp); err != nil {
		logger.Error(
			"failed to send verification email",
			"operation", "send_otp",
			"error", err,
		)

		c.JSON(http.StatusAccepted, gin.H{
			"message": "account created; request a new verification code if it does not arrive",
		})
		return
	}

	accessToken, err := util.GenerateAccessToken(user.ID.Hex())
	if err != nil {
		logger.Error("failed to generate access token", "error", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	refreshToken, err := util.GenerateRefreshToken(user.ID.Hex())
	if err != nil {
		logger.Error("failed to generate refresh token", "error", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
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

	session := &models.Session{
		ID:           primitive.NewObjectID(),
		UserID:       user.ID.Hex(),
		RefreshToken: string(hashedRefreshToken),
		UserAgent:    c.Request.UserAgent(),
		IP:           c.ClientIP(),
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		CreatedAt:    time.Now(),
	}

	if _, err := sessionStore.CreateSession(ctx, session); err != nil {
		logger.Error(
			"failed to create user session",
			"operation", "create_session",
			"error", err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	logger.Info("user registration completed")

	c.SetCookie(
		"refreshToken",
		refreshToken,
		7*24*60*60,
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusCreated, gin.H{
		"message":     "user created successfully; please check your email for the OTP",
		"accessToken": accessToken,
	})
}
