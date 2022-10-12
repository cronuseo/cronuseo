package models

type User struct {
	ID             int    `json:"id" gorm:"primary_key"`
	Username       string `json:"username"`
	Name           string `json:"name"`
	OrganizationID int    `json:"-" gorm:"foreignKey:ID"`
}
