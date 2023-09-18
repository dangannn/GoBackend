package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Text   string `gorm:"not null"`
	PostID uint   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
