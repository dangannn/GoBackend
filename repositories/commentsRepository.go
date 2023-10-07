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
func (r *CommentRepository) Delete(id int) *models.ResponseError {
	err := r.db.Where("id = ?", id).Delete(&models.Comment{}).Error
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return nil
}

func (r *CommentRepository) Moderate(id int, comment *models.Comment) *models.ResponseError {
	err := r.db.Model(&models.Comment{}).Where("id = ?", id).Update("approved", comment.Approved).Error
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return nil
}

func (r *CommentRepository) GetAllUnapproved() ([]*models.Comment, *models.ResponseError) {
	var comments []*models.Comment
	err := r.db.Order("id desc").Where("approved = ?", false).Find(&comments).Error
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return comments, nil
}
