package models

// @Description Organization information
type Tenant struct {
	ID              string `json:"tenant_id" db:"tenant_id"`
	Key             string `json:"tenant_key" validate:"required,min=4" db:"tenant_key"`
	Name            string `json:"name" validate:"required,min=4" db:"name"`
	OraganizationID string `json:"org_id" db:"org_id"`
}

type TenantCreateRequest struct {
	Key  string `json:"tenant_key" validate:"required,min=4"`
	Name string `json:"name" validate:"required,min=4"`
}

type TenantUpdateRequest struct {
	Name string `json:"name" validate:"required,min=4"`
}
