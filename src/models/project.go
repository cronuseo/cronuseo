package models

type Project struct {
	ID    string `json:"project_id" db:"project_id"`
	Key   string `json:"project_key" validate:"required,min=4" db:"project_key"`
	Name  string `json:"name" validate:"required" db:"name"`
	OrgID string `json:"org_id" db:"org_id"`
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
