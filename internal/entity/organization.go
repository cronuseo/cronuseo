package entity

type Organization struct {
	LogicalKey string `json:"-" db:"id"`
	ID         string `json:"org_id" db:"org_id"`
	Key        string `json:"org_key" db:"org_key"`
	Name       string `json:"name" db:"name"`
	API_KEY    string `json:"org_api_key" db:"org_api_key"`
	CreatedAt  string `json:"created_at" db:"created_at"`
	UpdatedAt  string `json:"updated_at" db:"updated_at"`
}
