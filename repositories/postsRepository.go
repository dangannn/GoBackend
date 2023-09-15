package repositories

// PostRepository представляет репозиторий для работы с моделью Post.
import (
	"GoBackend/models"
	"gorm.io/gorm"
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
	r.db.Create(post)
	return post, nil
}

// FindByID находит пост по его идентификатору.
func (r *PostRepository) FindByID(id uint) (*models.Post, error) {
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

// FindAll возвращает все посты.
func (r *PostRepository) FindAll() ([]models.Post, error) {
	var posts []models.Post
	err := r.db.Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}
