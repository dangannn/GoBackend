package services

import (
	"GoBackend/models"
	"GoBackend/repositories"
	"net/http"
)

type PostService struct {
	postsRepository *repositories.PostRepository
}

func NewPostService(postsRepository *repositories.PostRepository) *PostService {
	return &PostService{
		postsRepository: postsRepository,
	}
}

func (ps PostService) CreatePost(post *models.Post) (*models.Post, *models.ResponseError) {
	responseErr := validatePost(post)
	if responseErr != nil {
		return nil, responseErr
	}
	return ps.postsRepository.Create(post)
}

func validatePost(post *models.Post) *models.ResponseError {
	if post.Title == "" {
		return &models.ResponseError{
			Message: "Invalid title",
			Status:  http.StatusBadRequest,
		}
	}
	//TODO - other checks
	return nil
}
