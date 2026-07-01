package domain

const MaxSteps = 500 
type Simulator struct {
	Machine *Machine
	Tape    *Tape
	State   string
}

func NewSimulator(m *Machine, initialTape *Tape) *Simulator {
	return &Simulator{
		Machine: m,
		Tape:    initialTape,
		State:   m.StartState,
	}
}

func (s *Simulator) Run() []StepResult {
	var history []StepResult
	steps := 0

	history = append(history, s.snapshot("Running"))

	for steps < MaxSteps {
		if s.Machine.IsAcceptState(s.State) {
			history = append(history, s.snapshot("Accepted"))
			break
		}
		if s.Machine.IsRejectState(s.State) {
			history = append(history, s.snapshot("Rejected"))
			break
		}

		readSymbol := s.Tape.Read()
		
		stateTransitions, stateOk := s.Machine.Transitions[s.State]
		transition, transOk := stateTransitions[readSymbol]

		if !stateOk || !transOk {
			s.State = s.Machine.RejectStates[0] 
			history = append(history, s.snapshot("Rejected"))
			break
		}

		s.Tape.Write(transition.Write)
		s.Tape.Move(transition.Move)
		s.State = transition.Next

		history = append(history, s.snapshot("Running"))
		steps++
	}

	if steps >= MaxSteps {
		history = append(history, s.snapshot("Timeout"))
	}

	return history
}

func (s *Simulator) snapshot(status string) StepResult {
	return StepResult{
		CurrentState: s.State,
		TapeContent:  s.Tape.String(),
		HeadPosition: s.Tape.Head,
		Status:       status,
	}
}