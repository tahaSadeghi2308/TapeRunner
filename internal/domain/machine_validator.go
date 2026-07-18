package domain

import "fmt"

func (m *Machine) Validate() error {
	if err := m.validateStates(); err != nil {
		return err
	}

	if err := m.validateAlphabets(); err != nil {
		return err
	}

	if err := m.validateTransitions(); err != nil {
		return err
	}
	return nil
}

func (m *Machine) validateStates() error {
	if len(m.States) == 0 {
		return fmt.Errorf("machine must contain at least one state")
	}

	if !contains(m.States, m.StartState) {
		return fmt.Errorf("start state %q does not exist", m.StartState)
	}

	if len(m.AcceptStates) == 0 {
		return fmt.Errorf("machine must contain at least one accept state")
	}

	for _, state := range m.AcceptStates {
		if !contains(m.States, state) {
			return fmt.Errorf("accept state %q does not exist", state)
		}
	}

	for _, state := range m.RejectStates {
		if !contains(m.States, state) {
			return fmt.Errorf("reject state %q does not exist", state)
		}
	}
	return nil
}

func (m *Machine) validateAlphabets() error {
	if len(m.InputAlphabet) == 0 {
		return fmt.Errorf("input alphabet is empty")
	}

	if len(m.TapeAlphabet) == 0 {
		return fmt.Errorf("tape alphabet is empty")
	}

	if !contains(m.TapeAlphabet, m.Blank) {
		return fmt.Errorf("blank symbol %q is not in tape alphabet", m.Blank)
	}

	for _, symbol := range m.InputAlphabet {
		if !contains(m.TapeAlphabet, symbol) {
			return fmt.Errorf(
				"input symbol %q is not in tape alphabet",
				symbol,
			)
		}
	}
	return nil
}

func (m *Machine) validateTransitions() error {
	for currentState, transitions := range m.Transitions {
		if !contains(m.States, currentState) {
			return fmt.Errorf(
				"transition references unknown state %q",
				currentState,
			)
		}

		for readSymbol, transition := range transitions {
			if !contains(m.TapeAlphabet, readSymbol) {
				return fmt.Errorf(
					"transition reads unknown symbol %q",
					readSymbol,
				)
			}

			if !contains(m.States, transition.Next) {
				return fmt.Errorf(
					"transition goes to unknown state %q",
					transition.Next,
				)
			}

			if !contains(m.TapeAlphabet, transition.Write) {
				return fmt.Errorf(
					"transition writes unknown symbol %q",
					transition.Write,
				)
			}

			switch transition.Move {
			case "L", "R":
			default:
				return fmt.Errorf(
					"transition from %q with symbol %q has invalid direction %q",
					currentState,
					readSymbol,
					transition.Move,
				)
			}
		}
	}
	return nil
}

func contains[T comparable](items []T, value T) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}
	return false
}