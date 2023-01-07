package entity

type CreatePermissionsRequest struct {
	Permissions []Permission `json:"permissions"`
}

type Permission struct {
	Role     string `json:"role"`
	Action   string `json:"action"`
	Resource string `json:"resource"`
}
