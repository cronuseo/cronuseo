package models

type User struct {
	ID        string `json:"user_id" db:"user_id"`
	Username  string `json:"username" validate:"required,min=4" db:"username"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	TenantID  string `json:"tenant_id" db:"tenant_id"`
}

type UserCreateRequest struct {
	ID        string `json:"user_id"`
	Username  string `json:"username" validate:"required,min=4"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UserUpdateRequest struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

type UserWithGroup struct {
	ID             int               `json:"id"`
	Username       string            `json:"username"`
	FirstName      string            `json:"firstname"`
	LastName       string            `json:"lastname"`
	OrganizationID int               `json:"org_id"`
	Groups         []GroupOnlyWithID `json:"groups"`
}

type GroupOnlyWithID struct {
	GroupID int `json:"group_id"`
}
