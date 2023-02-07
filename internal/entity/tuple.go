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
	Username    string     `json:"username"`
	Permissions []Relation `json:"permissions"`
	Resource    string     `json:"resource"`
}

type CheckRequestAll struct {
	SubjectId string   `json:"subject"`
	Objects   []Object `json:"resources"`
}

type Relation struct {
	Relation string `json:"permission"`
}

type Object struct {
	Object    string     `json:"resource"`
	Relations []Relation `json:"permissions"`
}

type CheckAllResponse struct {
	Values []CheckAllResult `json:"values"`
}

type CheckAllResult struct {
	Resource    string   `json:"resource"`
	Permissions []string `json:"permissions"`
}
