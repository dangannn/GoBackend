package models

import "time"

type Comment struct {
	Id        uint   `gorm:"primaryKey"`
	Text      string `gorm:"not null"`
	PostId    uint   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Approved  bool   `gorm:"not null,default:false"`
	AuthorId  uint   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt time.Time
}
