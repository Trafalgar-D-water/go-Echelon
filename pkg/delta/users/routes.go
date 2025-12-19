package users

import (
	"github.com/gin-gonic/gin"
	"github.com/go-Echelon/go-Echelon/pkg/core/database"
	"github.com/go-Echelon/go-Echelon/pkg/delta/users/routes"
	"github.com/go-Echelon/go-Echelon/pkg/delta/users/services"
)

func RegisterRoutes(r *gin.RouterGroup, db *database.Database) {
	userService := services.NewUserService(db)
	userGroup := r.Group("/users")
	routes.SignUp(userGroup, userService)
}
