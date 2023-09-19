package services

import (
	"GoBackend/models"
	"GoBackend/repositories"
)

type UserService struct {
	usersRepository *repositories.UserRepository
}

func NewUserService(usersRepository *repositories.UserRepository) *UserService {
	return &UserService{
		usersRepository: usersRepository,
	}
}

func (us UserService) GetAllUsers() ([]*models.User, *models.ResponseError) {
	return us.usersRepository.GetAllUsers()
}
