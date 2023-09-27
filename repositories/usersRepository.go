package repositories

import (
	"GoBackend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) GetAllUsers() ([]*models.User, *models.ResponseError) {
	var users []*models.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return users, nil
}

func (r *UserRepository) GetUserById(id int) (*models.User, *models.ResponseError) {
	var user *models.User
	err := r.db.Where("id = ?", id).Find(&user).Error
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return user, nil
}

func (r *UserRepository) GetUserPosts(id int) (*[]models.Post, *models.ResponseError) {
	var user *models.User
	err := r.db.Preload("Posts").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return &user.Posts, nil
}

func (r *UserRepository) Login(loginRequest *models.LoginRequest) (*models.User, *models.ResponseError) {
	var user *models.User

	err := r.db.Where("email = ?", loginRequest.Email).Find(&user).Error

	if err != nil {
		return nil, &models.ResponseError{
			Message: "Wrong email or password 1",
			Status:  http.StatusInternalServerError,
		}
	}
	err1 := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))

	if err1 != nil {
		return nil, &models.ResponseError{
			Message: "Wrong email or password 2",
			Status:  http.StatusInternalServerError,
		}
	}

	return user, nil
}

func (r *UserRepository) CreateUser(user *models.User) (*models.User, *models.ResponseError) {
	err := r.db.Create(&user).Error
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return user, nil
}
