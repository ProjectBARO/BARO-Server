package models

import "gorm.io/gorm"

type User struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Email    string `json:"email" gorm:"unique"`
	Age      int    `json:"age"`
	Gender   string `json:"gender"`
	Deleted  gorm.DeletedAt
}
