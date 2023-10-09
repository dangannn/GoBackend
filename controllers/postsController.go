package controllers

import (
	"GoBackend/models"
	"GoBackend/services"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strconv"
)

type PostController struct {
	postsService *services.PostService
}

func NewPostController(postService *services.PostService) *PostController {
	return &PostController{
		postsService: postService,
	}
}

func (pc PostController) Create(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Error while reading create post request body", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var post models.Post

	err = json.Unmarshal(body, &post)
	if err != nil {
		log.Println("Error while unmarshaling create post request body", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response, responseErr := pc.postsService.Create(&post)
	if responseErr != nil {
		c.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (pc PostController) GetById(c *gin.Context) {
	var id = c.Param("id")
	response, responseErr := pc.postsService.GetById(id)
	if responseErr != nil {
		c.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (pc PostController) Update(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Error while reading create post request body", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	type tmpPost struct {
		Id      string
		Title   string
		Content string
	}
	var tmp tmpPost

	err = json.Unmarshal(body, &tmp)
	if err != nil {
		log.Println("Error while unmarshaling create post request body", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tmpPostId, err := strconv.Atoi(tmp.Id)
	if err != nil {
		fmt.Println(err)
		c.Abort()
		return
	}
	post := models.Post{
		Id:      uint(tmpPostId),
		Title:   tmp.Title,
		Content: tmp.Content,
	}

	response, responseErr := pc.postsService.Update(&post)
	if responseErr != nil {
		c.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (pc PostController) Delete(c *gin.Context) {
	var id = c.Param("id")
	responseErr := pc.postsService.Delete(id)
	if responseErr != nil {
		c.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	c.JSON(http.StatusOK, "post deleted")
}

func (pc PostController) GetAll(c *gin.Context) {

	response, responseErr := pc.postsService.GetAll()
	if responseErr != nil {
		c.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (pc PostController) GetPage(c *gin.Context) {
	var page = c.Param("page")
	response, responseErr := pc.postsService.GetPage(page)
	if responseErr != nil {
		c.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (pc PostController) GetApprovedComments(c *gin.Context) {
	var id = c.Param("id")
	response, responseErr := pc.postsService.GetApprovedComments(id)
	if responseErr != nil {
		c.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	c.JSON(http.StatusOK, response)
}
