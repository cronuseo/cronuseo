package models

type Project struct {
	ID             int    `json:"id" gorm:"primary_key"`
	Key            string `json:"key" binding:"required,min=4"`
	Name           string `json:"name" binding:"required,min=4"`
	Description    string `json:"description"`
	OrganizationID int    `json:"-" gorm:"foreignKey:ID"`
}
