package models

import "time"

type Post struct {
	Id        uint      `gorm:"primaryKey"`
	Title     string    `gorm:"not null"`
	Content   string    `gorm:"not null"`
	Comments  []Comment `gorm:"foreignKey:PostId"`
	Checks    uint      `gorm:"not null"`
	AuthorId  uint      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt time.Time
}
