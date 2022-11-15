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
	Key  string `json:"resource_role_key" validate:"required,min=4" db:"resource_role_key"`
	Name string `json:"name" validate:"required,min=4" db:"name"`
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

type ResourceRoleToGroup struct {
	GroupID        int `gorm:"foreignKey:ID"`
	ResourceRoleID int `gorm:"foreignKey:ID"`
}

type ResourceRoleToUser struct {
	UserID         int `gorm:"foreignKey:ID"`
	ResourceRoleID int `gorm:"foreignKey:ID"`
}

type ResourceRoleToResourceAction struct {
	ResourceID       int `gorm:"foreignKey:ID"`
	ResourceActionID int `gorm:"foreignKey:ID"`
	ResourceRoleID   int `gorm:"foreignKey:ID"`
}

type ResourceRoleToResourceActionKey struct {
	Resource       string `json:"resource"`
	ResourceAction string `json:"resourceAction"`
	ResourceRole   string `json:"resourceRole"`
}

type ResourceRoleWithGroupsUsers struct {
	ID              int                    `json:"id"`
	Key             string                 `json:"key"`
	Name            string                 `json:"name"`
	ResourceID      int                    `json:"res_id"`
	Users           []UserID               `json:"users"`
	Groups          []GroupOnlyWithID      `json:"groups"`
	ResourceActions []ResourceActionWithID `json:"actions"`
}

type ResourceActionWithID struct {
	ResourceActionID int `json:"resact_id"`
}
