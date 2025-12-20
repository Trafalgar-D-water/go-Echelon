package users

import (
	"github.com/gin-gonic/gin"
	"github.com/go-Echelon/go-Echelon/pkg/core/database"
	"github.com/go-Echelon/go-Echelon/pkg/delta/users/routes"
	"github.com/go-Echelon/go-Echelon/pkg/delta/users/services"
)

// RegisterRoutes is the entry point for all User-related routes
// It matches the usage in main.go
func RegisterRoutes(r *gin.RouterGroup, db *database.Database) {
	// 1. Initialize Service (Directly with DB now)
	userService := services.NewUserService(db)

	// 2. Create the route group for /api/v1/users
	userGroup := r.Group("/users")

	// 3. Register individual routes files
	routes.SignUp(userGroup, userService)
	routes.Login(userGroup, userService)
	routes.GetUserByID(userGroup, userService)
}
