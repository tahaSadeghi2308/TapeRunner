package domain

type Machine struct {
	States        []string
	InputAlphabet []string
	TapeAlphabet  []string
	Transitions   map[string]map[string]Transition 
	StartState    string
	AcceptStates  []string
	RejectStates  []string
}

func (m *Machine) IsAcceptState(state string) bool {
	for _, s := range m.AcceptStates {
		if s == state {
			return true
		}
	}
	return false
}

func (m *Machine) IsRejectState(state string) bool {
	for _, s := range m.RejectStates {
		if s == state {
			return true
		}
	}
	return false
}