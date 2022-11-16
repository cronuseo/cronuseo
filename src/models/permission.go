package models

type Permissions struct {
	Permissions []Permission `json:"permissions"`
}

type Permission struct {
	RoleID  string   `json:"role_id"`
	Actions []Action `json:"actions"`
}

type Action struct {
	ActionID string `json:"action_id"`
}
