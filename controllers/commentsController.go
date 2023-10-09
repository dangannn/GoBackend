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
	"time"
)

type CommentController struct {
	commentsService *services.CommentService
}

func NewCommentController(commentService *services.CommentService) *CommentController {
	return &CommentController{
		commentsService: commentService,
	}
}

func (cc CommentController) Create(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading create comment request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	type TmpComment struct {
		Id        uint
		Text      string
		PostId    string
		Approved  bool
		AuthorId  uint
		CreatedAt time.Time
	}
	var tmp TmpComment

	err = json.Unmarshal(body, &tmp)
	if err != nil {
		log.Println(err)
		return
	}
	tmpPostId, err := strconv.Atoi(tmp.PostId)
	if err != nil {
		fmt.Println(err)
		ctx.Abort()
		return
	}
	receivedData := models.Comment{
		Id:        tmp.Id,
		Text:      tmp.Text,
		PostId:    uint(tmpPostId),
		Approved:  tmp.Approved,
		AuthorId:  tmp.AuthorId,
		CreatedAt: tmp.CreatedAt,
	}
	log.Println("измененные", receivedData)

	response, responseErr := cc.commentsService.Create(&receivedData)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
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
		log.Println("Error while unmarshalling create post request body", err)
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
