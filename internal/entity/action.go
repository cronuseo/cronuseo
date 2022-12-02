package entity

type Action struct {
	ID         string `json:"action_id" db:"action_id"`
	Key        string `json:"action_key" db:"action_key"`
	Name       string `json:"name" db:"name"`
	ResourceID string `json:"resource_id" db:"resource_id"`
}
