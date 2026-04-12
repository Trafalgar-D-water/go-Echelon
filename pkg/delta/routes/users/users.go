package users

import (
	"github.com/gin-gonic/gin"
	"github.com/go-Echelon/go-Echelon/pkg/core/database"
)

const dbKey = "db"

// RegisterRoutes mounts all user-related routes onto /api/v1/users
// and /api/v1/auth/session.
func RegisterRoutes(r *gin.RouterGroup, db database.Database) {
	// Inject database into gin context for all user handlers
	injectDB := func(c *gin.Context) {
		c.Set(dbKey, db)
		c.Next()
	}

	userRoutes := r.Group("", injectDB)
	{
		userRoutes.POST("/users", create)
		userRoutes.POST("/auth/session/login", login)
		userRoutes.GET("/users/refresh", refresh)
		userRoutes.GET("/users/:id", fetch)
		userRoutes.POST("/auth/session/logout", logout)
		userRoutes.POST("/users/verify", verifyOTP) // OTP Verification Endpoint
	}
}

// getDB retrieves the database from gin context
func getDB(c *gin.Context) database.Database {
	return c.MustGet(dbKey).(database.Database)
}
