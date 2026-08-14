package users

import (
	"net/http"
	"time"
	"context"

	"github.com/gin-gonic/gin"
)

type UserSearchResult struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Online        bool   `json:"online"`
}

// @Summary      Search Users
// @Description  Search for users by their username using a case-insensitive match. Requires JWT authentication.
// @Tags         users
// @Produce      json
// @Param        username query string true "The username to search for (minimum 2 characters)"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{} "Returns an array of clean UserSearchResult objects"
// @Failure      400  {object}  map[string]string "Invalid query or missing username"
// @Failure      401  {object}  map[string]string "Unauthorized Access"
// @Failure      500  {object}  map[string]string "Database error"
// @Router       /users/search [get]
func searchUser(c *gin.Context) {
	query := c.Query("username")
	if query == "" || len(query) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query string 'username' must be at least 2 characters"})
		return
	}

	db := getDB(c)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	users, err := db.Users().SearchUsers(ctx, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users"})
		return
	}

	// We create a slice of clean results to avoid leaking passwords!
	var results []UserSearchResult
	for _, u := range users {
		results = append(results, UserSearchResult{
			ID:            u.ID.Hex(),
			Username:      u.Username,
			Discriminator: u.Discriminator,
			Online:        u.Online,
		})
	}

	if len(results) == 0 {
		// Ensure an empty array is sent instead of null
		results = []UserSearchResult{}
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}
