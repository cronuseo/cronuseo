package models

type Project struct {
	ID       string `json:"id" db:"project_id"`
	Key      string `json:"project_key" validate:"required,min=4" db:"project_key"`
	Name     string `json:"name" validate:"required" db:"name"`
	TenantID string `json:"tenant_id" db:"tenant_id"`
}

type ProjectCreateRequest struct {
	Key         string `json:"project_key" validate:"required,min=4"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type ProjectUpdateRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}
