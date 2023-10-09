package repositories

// PostRepository представляет репозиторий для работы с моделью Post.
import (
	"GoBackend/models"
	"gorm.io/gorm"
	"net/http"
)

type PostRepository struct {
	db *gorm.DB
}

// NewPostRepository создает новый экземпляр PostRepository.
func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db}
}

func (r *PostRepository) Create(post *models.Post) (*models.Post, *models.ResponseError) {
	err := r.db.Create(&post).Error
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return post, nil
}

func (r *PostRepository) GetById(id int) (*models.Post, *models.ResponseError) {
	var post *models.Post
	err := r.db.Where("id = ?", id).Find(&post).Error
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return post, nil
}

func (r *PostRepository) Delete(id int) *models.ResponseError {
	err := r.db.Where("id = ?", id).Delete(&models.Post{}).Error
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return nil
}

func (r *PostRepository) Update(post *models.Post) (*models.Post, *models.ResponseError) {
	err := r.db.Model(&models.Post{}).Where("id = ?", post.Id).Updates(&post).Error
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return post, nil
}

// GetAll возвращает все посты.
func (r *PostRepository) GetAll() ([]*models.Post, *models.ResponseError) {
	var posts []*models.Post
	err := r.db.Order("id desc").Find(&posts).Error
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return posts, nil
}

func (r *PostRepository) GetPage(page int) ([]*models.Post, *models.ResponseError) {
	var posts []*models.Post
	err := r.db.Limit(3).Offset(3 * page).Find(&posts).Error
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return posts, nil
}

func (r *PostRepository) GetApprovedComments(id string) (*[]models.Comment, *models.ResponseError) {
	var post models.Post
	err := r.db.Preload("Comments", "approved = true").Where("id = ?", id).First(&post).Error
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return &post.Comments, nil
}
