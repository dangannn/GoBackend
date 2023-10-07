package models

import "gorm.io/gorm"

type Email struct {
	sender    string
	receivers string
	subject   string
	body      string
}

type DailyStats struct {
	gorm.Model
	NewComments int
	Views       int
}
