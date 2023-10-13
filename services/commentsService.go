package services

import (
	"GoBackend/models"
	"GoBackend/repositories"
	"net/http"
	"strconv"
)

type CommentService struct {
	commentsRepository *repositories.CommentRepository
}

func NewCommentService(commentsRepository *repositories.CommentRepository) *CommentService {
	return &CommentService{
		commentsRepository: commentsRepository,
	}
}

func (cs CommentService) Create(comment *models.Comment) (*models.Comment, *models.ResponseError) {
	return cs.commentsRepository.Create(comment)
}

func (cs CommentService) Moderate(id string, comment *models.Comment) *models.ResponseError {
	idNum, err := strconv.Atoi(id)
	if err != nil {
		return &models.ResponseError{
			Message: "Invalid id",
			Status:  http.StatusBadRequest,
		}
	}
	if comment.Approved {
		return cs.commentsRepository.Moderate(idNum, comment)
	}
	return cs.commentsRepository.Delete(idNum)
}

func (cs CommentService) ModerateWs(id uint, comment *models.Comment) *models.ResponseError {
	idNum := int(id)

	if comment.Approved {
		return cs.commentsRepository.Moderate(idNum, comment)
	}
	return cs.commentsRepository.Delete(idNum)
}

func (cs CommentService) GetAllUnapproved() ([]*models.Comment, *models.ResponseError) {
	return cs.commentsRepository.GetAllUnapproved()
}
