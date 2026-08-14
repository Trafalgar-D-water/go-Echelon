package users

import (
	"github.com/gin-gonic/gin"
	"github.com/go-Echelon/go-Echelon/pkg/core/database"
	"github.com/go-Echelon/go-Echelon/pkg/delta/middleware"
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

	// PUBLIC ROUTES
	userRoutes := r.Group("", injectDB)
	{
		userRoutes.POST("/users", create)
		userRoutes.POST("/auth/session/login", login)
		userRoutes.GET("/users/refresh", refresh)
		userRoutes.POST("/users/verify", verifyOTP) // OTP Verification Endpoint
	}

	// PROTECTED ROUTES (Requires JWT Access Token)
	protected := r.Group("", injectDB, middleware.Auth())
	{
		protected.GET("/users/@me/relationships", FetchMyRelationships)
		protected.GET("/users/search", searchUser)
		protected.GET("/users/:id", fetch)
		protected.POST("/auth/session/logout", logout)
		
		// Relationships 
		protected.POST("/users/@me/relationships", SendRelationshipRequest)
		protected.PUT("/users/@me/relationships/:id", AcceptFriendRequest)
		// Assuming you will add these handlers soon!
		// protected.DELETE("/users/@me/relationships/:id", RemoveRelationship)
	}
}

// getDB retrieves the database from gin context
func getDB(c *gin.Context) database.Database {
	return c.MustGet(dbKey).(database.Database)
}
