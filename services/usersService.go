package services

import "GoBackend/repositories"

type UserService struct {
	usersRepository *repositories.UserRepository
}

func NewUserService(usersRepository *repositories.UserRepository) *UserService {
	return &UserService{
		usersRepository: usersRepository,
	}
}
