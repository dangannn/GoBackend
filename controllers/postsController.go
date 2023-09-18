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

type PostController struct {
	postsService *services.PostService
}

func NewPostController(postService *services.PostService) *PostController {
	return &PostController{
		postsService: postService,
	}
}

func (pc PostController) CreatePost(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading create post request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var post models.Post

	err = json.Unmarshal(body, &post)
	if err != nil {
		log.Println("Error while unmarshaling create post request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response, responseErr := pc.postsService.CreatePost(&post)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (pc PostController) RetrieveAllPosts(ctx *gin.Context) {

	response, responseErr := pc.postsService.RetrieveAllPosts()
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (pc PostController) GetComments(ctx *gin.Context) {
	var id string = ctx.Param("id")
	response, responseErr := pc.postsService.GetComments(id)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
