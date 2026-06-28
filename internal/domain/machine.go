package domain

import "errors"

type Direction string

const (
	Left  Direction = "L"
	Right Direction = "R"
)

type Transition struct {
	NextState string
	Write     rune
	Move      Direction
}

type StateSymbol struct {
	State  string
	Symbol rune
}

type TMDefinition struct {
	States      []string
	InputAlpha  []rune
	TapeAlpha   []rune
	Transitions map[StateSymbol]Transition
	StartState  string
	AcceptState []string
	RejectState []string
	BlankSymbol rune
}

func NewTMDefinition(blank rune, start string) *TMDefinition {
	return &TMDefinition{
		Transitions: make(map[StateSymbol]Transition),
		BlankSymbol: blank,
		StartState:  start,
	}
}
