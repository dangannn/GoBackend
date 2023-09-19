package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id             uint   `gorm:"primaryKey"`
	Name           string `gorm:"size:255"`
	Email          string `gorm:"type:varchar(100);unique_index"`
	HashedPassword []byte
	Role           string `gorm:"not null"`
}

func (user *User) SetNewPassword(passwordString string) {
	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(passwordString), bcrypt.DefaultCost)
	user.HashedPassword = bcryptPassword
}
