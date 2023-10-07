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

type CommentController struct {
	commentsService *services.CommentService
}

func NewCommentController(commentService *services.CommentService) *CommentController {
	return &CommentController{
		commentsService: commentService,
	}
}

func (cc CommentController) GetAllUnapproved(ctx *gin.Context) {

	response, responseErr := cc.commentsService.GetAllUnapproved()
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (cc CommentController) Moderate(ctx *gin.Context) {
	var id = ctx.Param("id")
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading moderation comment body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var comment models.Comment

	err = json.Unmarshal(body, &comment)
	log.Println("com contr", comment)
	if err != nil {
		log.Println("Error while unmarshaling create post request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	responseErr := cc.commentsService.Moderate(id, &comment)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.JSON(http.StatusOK, "comment moderated")
}
