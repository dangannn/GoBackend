package controllers

import (
	"GoBackend/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UsersController struct {
	usersService *services.UserService
}

func NewUsersController(usersService *services.UserService) *UsersController {
	return &UsersController{
		usersService: usersService,
	}
}

func (uc UsersController) GetAllUsers(ctx *gin.Context) {

	response, responseErr := uc.usersService.GetAllUsers()
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
