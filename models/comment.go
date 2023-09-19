package models

type Comment struct {
	Id     uint   `gorm:"primaryKey"`
	Text   string `gorm:"not null"`
	PostID uint   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
