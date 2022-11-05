package models

type ResourceAction struct {
	ID             int    `json:"id" gorm:"primary_key"`
	Key            string `json:"key" validate:"required,min=4"`
	Name           string `json:"name" validate:"required,min=4"`
	PermissionName string `json:"permission_name"`
	Description    string `json:"description"`
	ResourceID     int    `json:"-" gorm:"foreignKey:ID"`
}
