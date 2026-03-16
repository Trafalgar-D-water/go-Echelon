package users

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-Echelon/go-Echelon/internal/email"
	"github.com/go-Echelon/go-Echelon/pkg/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
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
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	dob, err := time.Parse("2006-01-02", req.DOB)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	db := getDB(c)
	coll := db.Users()
	// Check if user already exists
	count, err := coll.CountDocuments(ctx, bson.M{"email": req.Email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "This email is already registered"})
		return
	}

	// Hash password
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost) // put in the util function
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Generate 6-digit OTP
	rand.Seed(time.Now().UnixNano())
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	// Build user document
	user := &models.User{
		ID:         primitive.NewObjectID(),
		Username:   req.Username,
		Email:      req.Email,
		Password:   string(hashedBytes),
		DOB:        dob,
		IsVerified: false,
		OTP:        otp,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	if _, err = coll.InsertOne(ctx, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	go func() {
		if err := email.SendOTP(user.Email, otp); err != nil {
			log.Printf("Failed to send OTP to %s: %v\n", user.Email, err)
		}
	}()

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully. Please check your email for the OTP to verify your account.",
		"user":    user,
	})
}
