package users

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/go-Echelon/go-Echelon/pkg/delta/util"
)

// @Summary      Logout User
// @Description  Clears the refresh token cookie and removes the active session from the database.
// @Tags         users
// @Produce      json
// @Success      200  {object}  map[string]interface{} "Logged out successfully"
// @Router       /auth/session/logout [post]
func logout(c *gin.Context) {
	c.SetCookie("refreshToken", "", -1, "/", "", true, true)

	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
		return
	}

	claims, err := util.ParseRefreshToken(refreshToken)
	if err == nil {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		db := getDB(c)

		_ = db.Sessions().DeleteSessionByUserID(ctx, claims.UserID)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
