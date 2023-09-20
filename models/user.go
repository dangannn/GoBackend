package models

type User struct {
	Id       uint   `gorm:"primaryKey"`
	Name     string `gorm:"size:255"`
	Email    string `gorm:"unique;type:varchar(100);unique_index"`
	Password string `gorm:"not null"`
	Role     string `gorm:"not null"`
	Posts    []Post `gorm:"foreignKey:AuthorId"`
}

type LoginRequest struct {
	Email    string `gorm:"size:255"`
	Password string
}
