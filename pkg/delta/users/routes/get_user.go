package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-Echelon/go-Echelon/pkg/delta/users/controllers"
	"github.com/go-Echelon/go-Echelon/pkg/delta/users/services"
)

// GetUserByID registers the route at GET /:id
func GetUserByID(r *gin.RouterGroup, service *services.UserService) {
	controller := controllers.NewUserController(service)
	r.GET("/:id", controller.GetUserByID)
}
