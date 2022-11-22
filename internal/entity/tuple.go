package entity

type Tuple struct {
	SubjectId string `json:"subject"`
	Relation  string `json:"relation"`
	Object    string `json:"object"`
}
