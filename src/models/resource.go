package models

type Resource struct {
	ID          int    `json:"id" gorm:"primary_key"`
	Key         string `json:"key" validate:"required,min=4"`
	Name        string `json:"name" validate:"required,min=4"`
	Description string `json:"description"`
	ProjectID   int    `json:"-" gorm:"foreignKey:ID"`
}

type ResourceCreateRequest struct {
	Key         string `json:"key" validate:"required,min=4"`
	Name        string `json:"name" validate:"required,min=4"`
	Description string `json:"description"`
}

type ResourceUpdateRequest struct {
	Key         string `json:"key" validate:"required,min=4"`
	Name        string `json:"name" validate:"required,min=4"`
	Description string `json:"description"`
}
