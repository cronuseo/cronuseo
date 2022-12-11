package entity

type Action struct {
	LogicalKey string `json:"-" db:"id"`
	ID         string `json:"action_id" db:"action_id"`
	Key        string `json:"action_key" db:"action_key"`
	Name       string `json:"name" db:"name"`
	ResourceID string `json:"resource_id" db:"resource_id"`
	CreatedAt  string `json:"created_at" db:"created_at"`
	UpdatedAt  string `json:"updated_at" db:"updated_at"`
}
