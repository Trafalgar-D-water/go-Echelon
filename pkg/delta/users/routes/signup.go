package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-Echelon/go-Echelon/pkg/delta/users/controllers"
	"github.com/go-Echelon/go-Echelon/pkg/delta/users/services"
)

func SignUp(r *gin.RouterGroup, service *services.UserService) {
	controller := controllers.NewUserController(service)

	r.POST("/signup", controller.SignUp)
}
