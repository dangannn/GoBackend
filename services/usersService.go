package services

import (
	"GoBackend/models"
	"GoBackend/repositories"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type UserService struct {
	usersRepository *repositories.UserRepository
}

func NewUserService(usersRepository *repositories.UserRepository) *UserService {
	return &UserService{
		usersRepository: usersRepository,
	}
}

func (us UserService) CreateUser(user *models.User) (*models.User, *models.ResponseError) {
	responseErr := validateUser(user)
	bcryptPassword, _ := bcrypt.GenerateFromPassword(user.HashedPassword, bcrypt.DefaultCost)
	user.HashedPassword = bcryptPassword
	if responseErr != nil {
		return nil, responseErr
	}
	return us.usersRepository.CreateUser(user)
}

func (us UserService) GetAllUsers() ([]*models.User, *models.ResponseError) {
	return us.usersRepository.GetAllUsers()
}

func validateUser(user *models.User) *models.ResponseError {
	if user.Name == "" {
		return &models.ResponseError{
			Message: "Invalid title",
			Status:  http.StatusBadRequest,
		}
	}
	//TODO - other checks
	return nil
}
