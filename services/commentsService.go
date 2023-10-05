package services

import (
	"GoBackend/models"
	"GoBackend/repositories"
	"net/http"
)

type CommentService struct {
	commentsRepository *repositories.CommentRepository
}

func NewCommentService(commentsRepository *repositories.CommentRepository) *CommentService {
	return &CommentService{
		commentsRepository: commentsRepository,
	}
}

func (ps CommentService) Create(comment *models.Comment) (*models.Comment, *models.ResponseError) {
	responseErr := validateComment(comment)
	if responseErr != nil {
		return nil, responseErr
	}
	return ps.commentsRepository.Create(comment)
}

func validateComment(post *models.Comment) *models.ResponseError {
	if post.Text == "" {
		return &models.ResponseError{
			Message: "Invalid title",
			Status:  http.StatusBadRequest,
		}
	}
	//TODO - other checks
	return nil
}
