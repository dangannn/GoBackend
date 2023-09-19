package repositories

import (
	"GoBackend/models"
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
