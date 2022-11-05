package models

type Group struct {
	ID             int    `json:"id" gorm:"primary_key"`
	Key            string `json:"key" validate:"required,min=4"`
	Name           string `json:"name" validate:"required,min=4"`
	OrganizationID int    `json:"-" gorm:"foreignKey:ID"`
}

type GroupUser struct {
	GroupID int `gorm:"foreignKey:ID"`
	UserID  int `gorm:"foreignKey:ID"`
}

type GroupUsers struct {
	ID             int              `json:"id"`
	Key            string           `json:"key"`
	Name           string           `json:"name"`
	OrganizationID int              `json:"org_id"`
	Users          []UserOnlyWithID `json:"users"`
}

type UserOnlyWithID struct {
	UserID int `json:"user_id"`
}
