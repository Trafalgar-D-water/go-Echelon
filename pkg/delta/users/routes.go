package users

import (
	"github.com/gin-gonic/gin"
	"github.com/go-Echelon/go-Echelon/pkg/core/database"
	"github.com/go-Echelon/go-Echelon/pkg/delta/middleware"
	"github.com/go-Echelon/go-Echelon/pkg/delta/users/routes"
	"github.com/go-Echelon/go-Echelon/pkg/delta/users/services"
)

func RegisterRoutes(r *gin.RouterGroup, db *database.Database) {
	userService := services.NewUserService(db)

	// 2. Create the route group for /api/v1/users
	userGroup := r.Group("/users")
	routes.SignUp(userGroup, userService)
	routes.Login(userGroup, userService)
	protected := userGroup.Group("/")
	protected.Use(middleware.AuthMiddleware())
	routes.GetUserByID(protected, userService)
}
