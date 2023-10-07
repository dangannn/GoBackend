package repositories

import (
	"GoBackend/models"
	"gorm.io/gorm"
	"net/http"
)

type EmailRepository struct {
	db *gorm.DB
}

func NewEmailRepository(db *gorm.DB) *EmailRepository {
	return &EmailRepository{db}
}

func (r *EmailRepository) GetDailyStats() (*models.DailyStats, *models.ResponseError) {
	var stats models.DailyStats
	err := r.db.First(&stats).Error
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return &stats, nil
}

func (r *EmailRepository) ResetDailyStats() *models.ResponseError {
	err := r.db.Model(&models.DailyStats{}).Where("id = ?", 1).Updates(map[string]interface{}{"views": 0, "new_comments": 0}).Error
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return nil
}

func (r *EmailRepository) AddView() *models.ResponseError {
	err := r.db.Model(&models.DailyStats{}).Where("id = ?", 1).Update("views", gorm.Expr("views + 1")).Error
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return nil
}
func (r *EmailRepository) AddNewComment() *models.ResponseError {
	err := r.db.Model(&models.DailyStats{}).Where("id = ?", 1).Update("new_comments", gorm.Expr("new_comments + 1")).Error
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return nil
}
