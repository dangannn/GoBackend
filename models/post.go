package models

type Post struct {
	Id       uint      `gorm:"primaryKey"`
	Title    string    `gorm:"not null"`
	Content  string    `gorm:"not null"`
	Comments []Comment `gorm:"foreignKey:PostID"`
}
