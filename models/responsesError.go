package models

type ResponseError struct {
	Message string `gorm:"not null"`
	Status  int    `gorm:"not null"`
}
