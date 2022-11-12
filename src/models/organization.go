package models

// @Description Organization information
type Organization struct {
	ID   string `json:"org_id" db:"org_id"`
	Key  string `json:"org_key" validate:"required,min=4" db:"org_key"`
	Name string `json:"name" validate:"required,min=4" db:"name"`
}

type OrganizationRequest struct {
	Key  string `json:"org_key" validate:"required,min=4"`
	Name string `json:"name" validate:"required,min=4"`
}

type OrganizationUpdateRequest struct {
	Name string `json:"name" validate:"required,min=4"`
}
