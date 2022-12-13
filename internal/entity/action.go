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

type ActionQueryResponse struct {
	Links   Links          `json:"_links"`
	Results []ActionResult `json:"results"`
	Limit   int            `json:"limit"`
	Size    int            `json:"size"`
	Cursor  int            `json:"cursor"`
}

type ActionResult struct {
	ID         string      `json:"action_id" db:"action_id"`
	Key        string      `json:"action_key" db:"action_key"`
	Name       string      `json:"name" db:"name"`
	ResourceID string      `json:"resource_id" db:"resource_id"`
	CreatedAt  string      `json:"created_at" db:"created_at"`
	UpdatedAt  string      `json:"updated_at" db:"updated_at"`
	Links      ActionLinks `json:"_links"`
}

type ActionLinks struct {
	Self string `json:"self,omitempty"`
}
