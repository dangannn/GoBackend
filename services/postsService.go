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

func (ps PostService) Create(post *models.Post) (*models.Post, *models.ResponseError) {
	responseErr := validatePost(post)
	if responseErr != nil {
		return nil, responseErr
	}
	return ps.postsRepository.Create(post)
}

func (ps PostService) Delete(id string) *models.ResponseError {
	idNum, err := strconv.Atoi(id)
	if err != nil {
		return &models.ResponseError{
			Message: "Invalid id",
			Status:  http.StatusBadRequest,
		}
	}
	return ps.postsRepository.Delete(idNum)
}

func (ps PostService) GetById(id string) (*models.Post, *models.ResponseError) {
	idNum, err := strconv.Atoi(id)
	if err != nil {
		return nil, &models.ResponseError{
			Message: "Invalid id",
			Status:  http.StatusBadRequest,
		}
	}
	return ps.postsRepository.GetById(idNum)
}

func (ps PostService) Update(post *models.Post) (*models.Post, *models.ResponseError) {
	//responseErr := validatePost(post)
	//if responseErr != nil {
	//	return nil, responseErr
	//}
	return ps.postsRepository.Update(post)
}

func (ps PostService) GetAll() ([]*models.Post, *models.ResponseError) {
	return ps.postsRepository.GetAll()
}

func (ps PostService) GetPage(page string) ([]*models.Post, *models.ResponseError) {
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		return nil, &models.ResponseError{
			Message: "Invalid page",
			Status:  http.StatusBadRequest,
		}
	}
	return ps.postsRepository.GetPage(pageNum)
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
