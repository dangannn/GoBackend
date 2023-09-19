package controllers

import (
	"GoBackend/models"
	"GoBackend/services"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
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

func (uc UsersController) CreateUser(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading create user request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var user models.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Println("Error while unmarshaling create post request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response, responseErr := uc.usersService.CreateUser(&user)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
