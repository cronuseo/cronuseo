package entity

import "time"

type CreatePermissionsRequest struct {
	Permissions []Permission `json:"permissions"`
}

type Permission struct {
	Role     string `json:"role"`
	Action   string `json:"action"`
	Resource string `json:"resource"`
}

type CheckMetrics struct {
	Request   CheckRequestWithUser
	Result    bool
	Timestamp time.Time
}
