package controllers

import "github.com/go-Echelon/go-Echelon/pkg/delta/users/services"

type UserController struct {
	Service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{Service: service}
}
