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

type UserController struct {
	usersService *services.UserService
}

func NewUserController(usersService *services.UserService) *UserController {
	return &UserController{
		usersService: usersService,
	}
}

func (uc UserController) GetAllUsers(c *gin.Context) {

	response, responseErr := uc.usersService.GetAllUsers()
	if responseErr != nil {
		c.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (uc UserController) CreateUser(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Error while reading create user request body", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var user models.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Println("Error while unmarshaling create post request body", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response, responseErr := uc.usersService.CreateUser(&user)
	if responseErr != nil {
		c.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (uc UserController) GetUserById(c *gin.Context) {
	var id = c.Param("id")
	response, responseErr := uc.usersService.GetUserById(id)
	if responseErr != nil {
		c.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	c.JSON(http.StatusOK, response)
}
func (uc UserController) GetUserPosts(c *gin.Context) {
	var id = c.Param("id")
	response, responseErr := uc.usersService.GetUserPosts(id)
	if responseErr != nil {
		c.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (uc UserController) LoginUser(c *gin.Context) {
	log.Println(c.Request)
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Error while reading login user request body", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var loginRequest models.LoginRequest

	err = json.Unmarshal(body, &loginRequest)
	if err != nil {
		log.Println("Error while unmarshaling create post request body", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response, responseErr := uc.usersService.Login(&loginRequest)
	if responseErr != nil {
		c.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	c.JSON(http.StatusOK, response)
}
