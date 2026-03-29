package users

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-Echelon/go-Echelon/pkg/delta/util"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// @Summary      Refresh Access Token
// @Description  Validates refresh token from cookie, rotates it, and returns a new access token
// @Tags         users
// @Produce      json
// @Success      200  {object}  map[string]interface{} "New access token"
// @Failure      400  {object}  map[string]interface{} "Invalid or missing refresh token"
// @Failure      401  {object}  map[string]interface{} "Token mismatch or unauthorized"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /users/refresh [get]
func refresh(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()

	refreshToken, err := c.Cookie("refreshToken")

	fmt.Println(refreshToken, "This is my refrsh token ")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "InValid Token bla bla "})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := claims["userId"].(string)

	db := getDB(c)

	sessionStore := db.Sessions()

	hash := sha256.Sum256([]byte(refreshToken))
	hashedToken := hex.EncodeToString(hash[:])

	session, err := sessionStore.GetSessionByUserID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session not found"})
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(session.RefreshToken),
		[]byte(hashedToken))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token MisMatch"})
		return
	}

	// Rotation
	newAccessToken, _ := util.GenerateAccessToken(userID)
	newRefreshToken, _ := util.GenerateRefreshToken(userID)

	// hash a new Refresh Token

	newHash := sha256.Sum256([]byte(newRefreshToken))
	newHashedToken := hex.EncodeToString(newHash[:])

	newBcrypt, _ := bcrypt.GenerateFromPassword([]byte(newHashedToken), bcrypt.DefaultCost)

	err = sessionStore.UpdateSession(ctx, session.ID.Hex(), string(newBcrypt), time.Now().Add(7*24*time.Hour))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("refreshToken", newRefreshToken, 60*60*24, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"accessToken": newAccessToken})
}
