package models

import "gorm.io/gorm"

type User struct {
	ID       uint `gorm:"primary_key"`
	Name     string
	Nickname string
	Email    string `gorm:"unique"`
	Age      int
	Gender   string
	Deleted  gorm.DeletedAt
}
