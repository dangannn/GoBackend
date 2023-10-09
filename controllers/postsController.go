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

func (pc PostController) Create(ctx *gin.Context) {
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

	response, responseErr := pc.postsService.Create(&post)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (pc PostController) GetById(ctx *gin.Context) {
	var id = ctx.Param("id")
	response, responseErr := pc.postsService.GetById(id)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (pc PostController) Update(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading create post request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
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
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tmpPostId, err := strconv.Atoi(tmp.Id)
	if err != nil {
		fmt.Println(err)
		ctx.Abort()
		return
	}
	post := models.Post{
		Id:      uint(tmpPostId),
		Title:   tmp.Title,
		Content: tmp.Content,
	}

	response, responseErr := pc.postsService.Update(&post)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (pc PostController) Delete(ctx *gin.Context) {
	var id = ctx.Param("id")
	responseErr := pc.postsService.Delete(id)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.JSON(http.StatusOK, "post deleted")
}

func (pc PostController) GetAll(ctx *gin.Context) {

	response, responseErr := pc.postsService.GetAll()
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (pc PostController) GetPage(ctx *gin.Context) {
	var page string = ctx.Param("page")
	response, responseErr := pc.postsService.GetPage(page)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (pc PostController) GetApprovedComments(ctx *gin.Context) {
	var id string = ctx.Param("id")
	response, responseErr := pc.postsService.GetApprovedComments(id)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
