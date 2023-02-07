package entity

import "fmt"

type Tuple struct {
	SubjectId string `json:"subject"`
	Relation  string `json:"relation"`
	Object    string `json:"object"`
}

func (t Tuple) String() string {
	return fmt.Sprintf("%s%s%s", t.SubjectId, t.Relation, t.Object)
}

type CheckRequest struct {
	SubjectId string `json:"subject"`
	Relation  string `json:"permission"`
	Object    string `json:"object"`
}

type CheckRequestWithUser struct {
	Username   string `json:"username"`
	Permission string `json:"permission"`
	Resource   string `json:"resource"`
}

type CheckRequestWithPermissions struct {
	Username    string             `json:"username"`
	Permissions []PermissionObject `json:"permissions"`
	Resource    string             `json:"resource"`
}

type CheckRequestAll struct {
	Username  string           `json:"username"`
	Resources []ResourceObject `json:"resources"`
}

type PermissionObject struct {
	Permission string `json:"permission"`
}

type ResourceObject struct {
	Resource    string             `json:"resource"`
	Permissions []PermissionObject `json:"permissions"`
}

type CheckAllResponse struct {
	Values []CheckAllResult `json:"values"`
}

type CheckAllResult struct {
	Resource    string   `json:"resource"`
	Permissions []string `json:"permissions"`
}
