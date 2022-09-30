package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       int    `json:"id" gorm:"primary_key"`
	Username string `json:"username"`
	Name     string `json:"name"`
}
