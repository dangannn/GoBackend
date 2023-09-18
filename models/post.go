package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Id       uint      `gorm:"primaryKey"`
	Title    string    `gorm:"not null"`
	Content  string    `gorm:"not null"`
	Comments []Comment `gorm:"foreignKey:PostID"`
}
