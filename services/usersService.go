package services

import (
	"GoBackend/models"
	"GoBackend/repositories"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type UserService struct {
	usersRepository *repositories.UserRepository
}

func NewUserService(usersRepository *repositories.UserRepository) *UserService {
	return &UserService{
		usersRepository: usersRepository,
	}
}

func (us UserService) Create(user *models.User) (*models.User, *models.ResponseError) {
	//hashing password
	HashPassword(user)

	return us.usersRepository.Create(user)
}

func (us UserService) GetAll() ([]*models.User, *models.ResponseError) {
	return us.usersRepository.GetAll()
}

func (us UserService) GetById(id string) (*models.User, *models.ResponseError) {
	idNum, err := strconv.Atoi(id)
	if err != nil {
		return nil, &models.ResponseError{
			Message: "Invalid id",
			Status:  http.StatusBadRequest,
		}
	}
	return us.usersRepository.GetById(idNum)
}

func (us UserService) GetUserPosts(id string) (*[]models.Post, *models.ResponseError) {
	idNum, err := strconv.Atoi(id)
	if err != nil {
		return nil, &models.ResponseError{
			Message: "Invalid id",
			Status:  http.StatusBadRequest,
		}
	}
	return us.usersRepository.GetUserPosts(idNum)
}

func (us UserService) Login(loginRequest *models.LoginRequest) (*string, *models.ResponseError) {
	log.Println(loginRequest)
	user, err1 := us.usersRepository.Login(loginRequest)
	if err1 != nil {
		return nil, &models.ResponseError{
			Message: "Invalid user",
			Status:  http.StatusBadRequest,
		}
	}
	var (
		key []byte
		t   *jwt.Token
		s   string
	)
	key = []byte("secrete-key")
	t = jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub":  loginRequest.Email,
			"exp":  time.Now().Add(time.Hour * 72).Unix(),
			"role": user.Role,
			"id":   user.Id,
		})
	s, err2 := t.SignedString(key)
	if err2 != nil {
		return nil, &models.ResponseError{
			Message: "Error while signing token",
			Status:  http.StatusBadRequest,
		}
	}
	return &s, nil
}

func HashPassword(user *models.User) *models.ResponseError {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return &models.ResponseError{
			Message: "Error while hashing password",
			Status:  http.StatusBadRequest,
		}
	}
	user.Password = string(bytes)
	return nil
}
