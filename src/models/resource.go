package models

type Resource struct {
	ID        string `json:"resource_id" db:"resource_id"`
	Key       string `json:"resource_key" validate:"required,min=4" db:"resource_key"`
	Name      string `json:"name" validate:"required,min=4" db:"name"`
	ProjectID string `json:"project_id" db:"project_id"`
}

type ResourceCreateRequest struct {
	Key  string `json:"resource_key" validate:"required,min=4"`
	Name string `json:"name" validate:"required,min=4"`
}

type ResourceUpdateRequest struct {
	Name string `json:"name" validate:"required,min=4"`
}
