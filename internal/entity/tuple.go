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

type Relation struct {
	Relation string `json:"permission"`
}
