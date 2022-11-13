package models

type Group struct {
	ID       string   `json:"group_id" gorm:"primary_key" db:"group_id"`
	Key      string   `json:"group_key" validate:"required,min=4" db:"group_key"`
	Name     string   `json:"name" validate:"required,min=4" db:"name"`
	TenantID string   `json:"tenant_id" gorm:"foreignKey:ID" db:"tenant_id"`
	Users    []UserID `json:"users,omitempty"`
}

type GroupCreateRequest struct {
	Key   string   `json:"group_key" validate:"required,min=4"`
	Name  string   `json:"name" validate:"required,min=4"`
	Users []UserID `json:"users"`
}

type GroupUpdateRequest struct {
	Name string `json:"name" validate:"required,min=4"`
}

type UserID struct {
	UserID string `json:"user_id" db:"user_id"`
}

type GroupPatchRequest struct {
	Operations []GroupPatchOperation `json:"operations"`
}

type GroupPatchOperation struct {
	Operation string   `json:"op"`
	Users     []UserID `json:"users"`
}
