package database

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Uid      string `gorm:"unique"`
	Username string
	Token    string
}
