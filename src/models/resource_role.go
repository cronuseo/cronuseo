package models

import (
	"gorm.io/gorm"
)

type ResourceRole struct {
	gorm.Model
	ID          int    `json:"id" gorm:"primary_key"`
	Key         string `json:"key" binding:"required,min=4"`
	Name        string `json:"name" binding:"required,min=4"`
	Description string `json:"description"`
	ResourceID  int    `gorm:"foreignKey:ID"`
}
