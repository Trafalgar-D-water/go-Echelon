package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-Echelon/go-Echelon/pkg/delta/users/controllers"
	"github.com/go-Echelon/go-Echelon/pkg/delta/users/services"
)

// Login registers the login route at POST /login
func Login(r *gin.RouterGroup, service *services.UserService) {
	controller := controllers.NewUserController(service)
	r.POST("/login", controller.Login)
}
