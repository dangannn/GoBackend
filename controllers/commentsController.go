package controllers

import "GoBackend/services"

type CommentController struct {
	commentsService *services.CommentService
}

func NewCommentController(commentService *services.CommentService) *CommentController {
	return &CommentController{
		commentsService: commentService,
	}
}
