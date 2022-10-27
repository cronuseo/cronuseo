package models

type User struct {
	ID             int    `json:"id" gorm:"primary_key"`
	Username       string `json:"username" validate:"required,min=4"`
	Name           string `json:"name" validate:"required,min=4"`
	OrganizationID int    `json:"-" gorm:"foreignKey:ID"`
}

type UserWithGroup struct {
	ID             int               `json:"id"`
	Username       string            `json:"username"`
	Name           string            `json:"name"`
	OrganizationID int               `json:"org_id"`
	Groups         []GroupOnlyWithID `json:"groups"`
}

type GroupOnlyWithID struct {
	GroupID int `json:"group_id"`
}
