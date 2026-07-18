package yamlparser

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func (p *YamlParser) Validator(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var machine yamlMachine
	if err := yaml.Unmarshal(data, &machine); err != nil {
		return err
	}

	if len(machine.States) == 0 {
		return fmt.Errorf("states section is required")
	}

	if len(machine.InputAlphabet) == 0 {
		return fmt.Errorf("input_alphabet section is required")
	}

	if len(machine.TapeAlphabet) == 0 {
		return fmt.Errorf("tape_alphabet section is required")
	}

	if machine.StartState == "" {
		return fmt.Errorf("start_state is required")
	}

	if machine.Blank == "" {
		return fmt.Errorf("blank symbol is required")
	}

	if len(machine.AcceptStates) == 0 {
		return fmt.Errorf("accept_states section is required")
	}

	if len(machine.Transitions) == 0 {
		return fmt.Errorf("at least one transition is required")
	}

	for i, t := range machine.Transitions {
		if t.CurrentState == "" {
			return fmt.Errorf("transition %d: current_state is required", i)
		}

		if t.Read == "" {
			return fmt.Errorf("transition %d: read is required", i)
		}

		if t.Write == "" {
			return fmt.Errorf("transition %d: write is required", i)
		}

		if t.NextState == "" {
			return fmt.Errorf("transition %d: next_state is required", i)
		}

		switch t.Move {
		case "L", "R":
		default:
			return fmt.Errorf(
				"transition %d: move must be either 'L' or 'R'",
				i,
			)
		}
	}
	return nil
}
