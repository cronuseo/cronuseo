package models

type ResourceRole struct {
	ID         string    `json:"resource_role_id" db:"resource_role_id"`
	Key        string    `json:"resource_role_key" validate:"required,min=4" db:"resource_role_key"`
	Name       string    `json:"name" validate:"required,min=4" db:"name"`
	ResourceID string    `json:"resource_id" db:"resource_id"`
	Users      []UserID  `json:"users,omitempty"`
	Groups     []GroupID `json:"groups,omitempty"`
}

type ResourceRoleCreateRequest struct {
	Key    string    `json:"resource_role_key" validate:"required,min=4" db:"resource_role_key"`
	Name   string    `json:"name" validate:"required,min=4" db:"name"`
	Users  []UserID  `json:"users"`
	Groups []GroupID `json:"groups"`
}

type ResourceRoleUpdateRequest struct {
	Name string `json:"name" validate:"required,min=4" db:"name"`
}

type ResourceRolePatchRequest struct {
	Operations []ResourceRolePatchOperation `json:"operations"`
}

type ResourceRolePatchOperation struct {
	Operation string  `json:"op"`
	Path      string  `json:"path"`
	Values    []Value `json:"values"`
}

type Value struct {
	Value string `json:"value" db:"value"`
}
