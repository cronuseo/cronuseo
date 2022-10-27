package models

type Project struct {
	ID             int    `json:"id" gorm:"primary_key"`
	Key            string `json:"key" validate:"required,min=4"`
	Name           string `json:"name" validate:"required"`
	Description    string `json:"description"`
	OrganizationID int    `json:"-" gorm:"foreignKey:ID"`
}
