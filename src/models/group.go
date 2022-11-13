package models

type Group struct {
	ID       string   `json:"group_id" gorm:"primary_key" db:"group_id"`
	Key      string   `json:"group_key" validate:"required,min=4" db:"group_key"`
	Name     string   `json:"name" validate:"required,min=4" db:"name"`
	TenantID string   `json:"tenant_id" gorm:"foreignKey:ID" db:"tenant_id"`
	Users    []UserID `json:"users"`
}

type GroupCreateRequest struct {
	Key   string   `json:"group_key" validate:"required,min=4"`
	Name  string   `json:"name" validate:"required,min=4"`
	Users []UserID `json:"users"`
}

type GroupUpdateRequest struct {
	Name string `json:"name" validate:"required,min=4"`
}

type GroupUser struct {
	GroupID int `gorm:"foreignKey:ID"`
	UserID  int `gorm:"foreignKey:ID"`
}

// type GroupUsers struct {
// 	ID             int              `json:"id"`
// 	Key            string           `json:"key"`
// 	Name           string           `json:"name"`
// 	OrganizationID int              `json:"org_id"`
// 	Users          []UserOnlyWithID `json:"users"`
// }

type UserID struct {
	UserID string `json:"user_id"`
}

// type AddUsersToGroup struct {
// 	Users []UserOnlyWithID `json:"users"`
// }
