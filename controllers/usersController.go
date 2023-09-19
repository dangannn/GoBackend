package controllers

import "GoBackend/services"

type UsersController struct {
	usersService *services.UserService
}

func NewUsersController(usersService *services.UserService) *UsersController {
	return &UsersController{
		usersService: usersService,
	}
}
