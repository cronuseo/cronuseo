package models

import (
	"gorm.io/gorm"
)

type Group struct {
	gorm.Model
	ID             int    `json:"id" gorm:"primary_key"`
	Key            string `json:"key" binding:"required,min=4"`
	Name           string `json:"name" binding:"required,min=4"`
	OrganizationID int    `gorm:"foreignKey:ID"`
}
