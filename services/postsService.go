package services

import (
	"GoBackend/models"
	"GoBackend/repositories"
	"net/http"
	"strconv"
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
func (ps PostService) GetAllPosts() ([]*models.Post, *models.ResponseError) {
	return ps.postsRepository.GetAllPosts()
}

func (ps PostService) GetPostPage(page string) ([]*models.Post, *models.ResponseError) {
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		return nil, &models.ResponseError{
			Message: "Invalid page",
			Status:  http.StatusBadRequest,
		}
	}
	return ps.postsRepository.GetPostPage(pageNum)
}

func (ps PostService) GetComments(id string) (*[]models.Comment, *models.ResponseError) {
	return ps.postsRepository.GetComments(id)
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
