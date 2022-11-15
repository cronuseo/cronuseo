package models

type ResourceAction struct {
	ID         string `json:"resource_action_id" db:"resource_action_id"`
	Key        string `json:"resource_action_key" validate:"required,min=4" db:"resource_action_key"`
	Name       string `json:"name" validate:"required,min=4" db:"name"`
	ResourceID string `json:"resource_id" db:"resource_id"`
}

type ResourceActionCreateRequest struct {
	Key  string `json:"resource_action_key" validate:"required,min=4" db:"resource_action_key"`
	Name string `json:"name" validate:"required,min=4" db:"name"`
}

type ResourceActionUpdateRequest struct {
	Name string `json:"name" validate:"required,min=4" db:"name"`
}
