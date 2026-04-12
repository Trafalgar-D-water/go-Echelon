package users

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// @Summary      Logout User
// @Description  Clears the refresh token cookie and removes the active session from the database.
// @Tags         users
// @Produce      json
// @Success      200  {object}  map[string]interface{} "Logged out successfully"
// @Router       /auth/session/logout [post]
func logout(c *gin.Context) {
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
		return
	}

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	if err == nil && token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			if userID, exists := claims["userId"].(string); exists {
				ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
				defer cancel()
				db := getDB(c)
				_ = db.Sessions().DeleteSessionByUserID(ctx, userID)
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
