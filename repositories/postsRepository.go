package repositories

// PostRepository представляет репозиторий для работы с моделью Post.
import (
	"GoBackend/models"
	"github.com/google/uuid"
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

// Create создает новый пост.
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

// FindByID находит пост по его идентификатору.
func (r *PostRepository) FindByID(id uuid.UUID) (*models.Post, error) {
	var post models.Post
	err := r.db.Where("id = ?", id).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// Update обновляет существующий пост.
func (r *PostRepository) Update(post *models.Post) error {
	return r.db.Save(post).Error
}

// DeleteByID удаляет пост по его идентификатору.
func (r *PostRepository) DeleteByID(id uint) error {
	return r.db.Where("id = ?", id).Delete(&models.Post{}).Error
}

// RetrieveAllPosts возвращает все посты.
func (r *PostRepository) RetrieveAllPosts() ([]*models.Post, *models.ResponseError) {
	var posts []*models.Post
	err := r.db.Find(&posts).Error
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return posts, nil
}

func (r *PostRepository) GetPostPage(page int) ([]*models.Post, *models.ResponseError) {
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

func (r *PostRepository) GetComments(id string) (*[]models.Comment, *models.ResponseError) {
	var post models.Post
	err := r.db.Preload("Comments").Where("id = ?", id).First(&post).Error
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return &post.Comments, nil
}
