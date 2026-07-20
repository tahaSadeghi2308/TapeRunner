package machines

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tahaSadeghi2308/TapeRunner/internal/domain"
)

// ---------------------------------------------------------------------------
// Universal Turing Machine (bonus section).
//
// The interpreter itself is NOT written as ordinary Go control flow: it is a
// big, auto-generated Turing-machine transition table (utm_table.tsv), loaded
// here into a plain domain.Machine and executed by the very same generic
// engine (domain.Simulator / domain.Tape) used in parts 1 and 2. Only the
// bookkeeping needed to load that table and to build the initial tape lives
// in Go; the actual "interpretation" of M happens purely through table
// lookups performed by domain.Simulator, exactly like any other machine.
//
// Encoding convention (see report for full details):
//   - states of the simulated machine M : '0'..'9' (ordinary), 'A' (accept),
//     'X' (reject). The START state of M must always be relabeled as '0'.
//     (up to 10 ordinary states)
//   - tape symbols of M                 : digits '0'..'9'. Digit '0' is
//     ALWAYS the blank symbol of M (fixed convention). (up to 10 symbols)
//   - marker letters (mark which cell the simulated head is on):
//       symbol i (0-9) -> chr('a'+i), i.e. 'a'..'j'
//   - one rule of M is written as exactly:  q,a,q2,b,d;
//   - the whole encoded tape given to the UTM is:
//       ^  <rule><rule>...<rule>  #  <8 blank pad cells>  <w, first cell marked>
//     '^' is a sentinel the UTM never overwrites; '#' separates the encoded
//     program of M from the encoded working tape of M. The 8 blank pad cells
//     let M move left past its own input start (a genuine blank, matching
//     the semantics of the original dynamic Tape) - see report, "left-edge
//     limitation".
// ---------------------------------------------------------------------------

const (
	UTMSentinel = "^"
	UTMHash     = "#"
	UTMBlank    = "_"
	UTMAccept   = "ACCEPT"
	UTMReject   = "REJECT"
	UTMStart    = "prerun"
)

// digit -> marker letter used to mark the cell currently under M's head.
var utmMark = map[string]string{
	"0": "a", "1": "b", "2": "c", "3": "d", "4": "e",
	"5": "f", "6": "g", "7": "h", "8": "i", "9": "j",
}

// leftPadding blank ('0') cells are placed before w on the tape, so that a
// simulated machine which moves left past its own input start lands on a
// genuine blank cell instead of running into the '#' separator.
const leftPadding = 8

// LoadUTMTable reads the generated transition table (5 tab-separated
// columns: state, symbol, next_state, write_symbol, move) and turns it into
// a domain.Machine ready to be handed to domain.NewSimulator / RunUTM.
func LoadUTMTable(path string) (*domain.Machine, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	transitions := make(map[string]map[string]domain.Transition)
	statesSet := make(map[string]bool)

	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) != 5 {
			return nil, fmt.Errorf("utm table line %d: expected 5 tab-separated fields, got %d", lineNo, len(fields))
		}
		state, sym, next, write, move := fields[0], fields[1], fields[2], fields[3], fields[4]

		if _, ok := transitions[state]; !ok {
			transitions[state] = make(map[string]domain.Transition)
		}
		transitions[state][sym] = domain.Transition{
			Write: write,
			Move:  move,
			Next:  next,
		}
		statesSet[state] = true
		statesSet[next] = true
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	states := make([]string, 0, len(statesSet))
	for s := range statesSet {
		states = append(states, s)
	}

	m := &domain.Machine{
		States:       states,
		Transitions:  transitions,
		StartState:   UTMStart,
		Blank:        UTMBlank,
		AcceptStates: []string{UTMAccept},
		RejectStates: []string{UTMReject},
	}
	return m, nil
}

// BuildEncodedTape converts a friendly "rules#w" description into the actual
// initial tape string the generated UTM table expects (adds the sentinel and
// marks the first cell of w). Format of rulesAndInput:
//
//	"q,a,q2,b,d;q,a,q2,b,d;...;#w"
//
// (note the trailing ';' right before the single '#').
func BuildEncodedTape(rulesAndInput string) (string, error) {
	parts := strings.SplitN(rulesAndInput, "#", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("input must contain exactly one '#' separating the rules from w")
	}
	rules := parts[0]
	w := parts[1]

	pad := strings.Repeat("0", leftPadding)

	var markedW string
	if w == "" {
		markedW = pad + utmMark["0"]
	} else {
		firstMark, ok := utmMark[string(w[0])]
		if !ok {
			return "", fmt.Errorf("symbol %q of w is not a valid digit 0-9", string(w[0]))
		}
		markedW = pad + firstMark + w[1:]
	}

	return UTMSentinel + rules + UTMHash + markedW, nil
}

// RunUTM executes m on the given initial tape string with a configurable
// step budget. The UTM needs far more elementary steps than an ordinary
// machine (simulating a single step of M costs hundreds of UTM steps - see
// report, complexity section), so this does NOT reuse domain.Simulator.Run()
// directly (its MaxSteps constant of 500 is far too small here). The
// step-by-step semantics are otherwise identical to domain.Simulator.Run().
func RunUTM(m *domain.Machine, initialTape string, maxSteps int) []domain.StepResult {
	tape := domain.NewTape(initialTape, m.Blank)
	sim := domain.NewSimulator(m, tape)

	var history []domain.StepResult
	history = append(history, snapshot(sim, "Running"))

	steps := 0
	for steps < maxSteps {
		if m.IsAcceptState(sim.State) {
			history = append(history, snapshot(sim, "Accepted"))
			return history
		}
		if m.IsRejectState(sim.State) {
			history = append(history, snapshot(sim, "Rejected"))
			return history
		}

		readSymbol := sim.Tape.Read()
		stateTransitions, stateOk := m.Transitions[sim.State]
		transition, transOk := stateTransitions[readSymbol]

		if !stateOk || !transOk {
			sim.State = m.RejectStates[0]
			history = append(history, snapshot(sim, "Rejected"))
			return history
		}

		sim.Tape.Write(transition.Write)
		sim.Tape.Move(transition.Move)
		sim.State = transition.Next

		steps++
		history = append(history, snapshot(sim, "Running"))
	}

	history = append(history, snapshot(sim, "Timeout"))
	return history
}

func snapshot(s *domain.Simulator, status string) domain.StepResult {
	return domain.StepResult{
		CurrentState: s.State,
		TapeContent:  s.Tape.String(),
		HeadPosition: s.Tape.Head,
		Status:       status,
	}
}
