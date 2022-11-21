package entity

type Permission struct {
	ID         string `json:"permission_id" db:"permission_id"`
	Key        string `json:"permission_key" db:"permission_key"`
	Name       string `json:"name" db:"name"`
	ResourceID string `json:"resource_id" db:"resource_id"`
}
