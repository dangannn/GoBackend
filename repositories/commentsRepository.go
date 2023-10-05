package repositories

import (
	"GoBackend/models"
	"gorm.io/gorm"
	"net/http"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db}
}

func (r *CommentRepository) Create(comment *models.Comment) (*models.Comment, *models.ResponseError) {
	err := r.db.Create(&comment).Error
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return comment, nil
}
