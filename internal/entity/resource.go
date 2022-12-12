package entity

type Resource struct {
	LogicalKey string `json:"-" db:"id"`
	ID         string `json:"resource_id" db:"resource_id"`
	Key        string `json:"resource_key" db:"resource_key"`
	Name       string `json:"name" db:"name"`
	OrgID      string `json:"org_id" db:"org_id"`
	CreatedAt  string `json:"created_at" db:"created_at"`
	UpdatedAt  string `json:"updated_at" db:"updated_at"`
}

type ResourceQueryResponse struct {
	Links   Links            `json:"_links"`
	Results []ResourceResult `json:"results"`
	Limit   int              `json:"limit"`
	Size    int              `json:"size"`
	Cursor  int              `json:"cursor"`
}

type Links struct {
	Next string `json:"next,omitempty"`
	Self string `json:"self,omitempty"`
	Prev string `json:"prev,omitempty"`
}

type ResourceResult struct {
	ID        string        `json:"resource_id" db:"resource_id"`
	Key       string        `json:"resource_key" db:"resource_key"`
	Name      string        `json:"name" db:"name"`
	OrgID     string        `json:"org_id" db:"org_id"`
	CreatedAt string        `json:"created_at" db:"created_at"`
	UpdatedAt string        `json:"updated_at" db:"updated_at"`
	Links     ResourceLinks `json:"_links"`
}

type ResourceLinks struct {
	Self string `json:"self,omitempty"`
}
