package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title    string    `gorm:"not null"`
	Content  string    `gorm:"not null"`
	Comments []Comment `gorm:"foreignKey:PostID"`
}
