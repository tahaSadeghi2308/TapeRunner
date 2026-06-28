package domain

import (
	"errors"
	"fmt"
)

type Simulator struct {
	Machine      *TMDefinition
	Tape         *Tape
	CurrentState string
	HeadPos      int
	StepCount    int
}

type Snapshot struct {
	Step         int
	State        string
	HeadPosition int
	TapeContent  string
}

func NewSimulator(machine *TMDefinition, input string) *Simulator {
	return &Simulator{
		Machine:      machine,
		Tape:         NewTape(machine.BlankSymbol, input),
		CurrentState: machine.StartState,
		HeadPos:      0,
		StepCount:    0,
	}
}