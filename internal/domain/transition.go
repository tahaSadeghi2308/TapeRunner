package domain

type Transition struct {
	Write string `json:"write"`
	Move  string `json:"move"` // "L" or "R"
	Next  string `json:"next"`
}
