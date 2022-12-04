package entity

type Tuple struct {
	SubjectId string `json:"subject"`
	Relation  string `json:"relation"`
	Object    string `json:"object"`
}

type CheckRequest struct {
	SubjectId string `json:"subject"`
	Relation  string `json:"permission"`
	Object    string `json:"object"`
}

type CheckRequestWithPermissions struct {
	SubjectId string     `json:"subject"`
	Relations []Relation `json:"permissions"`
	Object    string     `json:"object"`
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
