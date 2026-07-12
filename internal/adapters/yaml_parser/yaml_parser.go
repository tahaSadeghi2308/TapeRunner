package yamlparser

import (
	"os"

	"github.com/tahaSadeghi2308/TapeRunner/internal/domain"
	"gopkg.in/yaml.v3"
)

type YamlParser struct{}

func NewYamlParser() *YamlParser {
	return &YamlParser{}
}

func (p *YamlParser) Read(filePath string) (*domain.Machine, error) {
	if err := p.Validator(filePath); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var yml yamlMachine
	if err := yaml.Unmarshal(data, &yml); err != nil {
		return nil, err
	}

	transitions := make(map[string]map[string]domain.Transition)
	for _, t := range yml.Transitions {
		if _, ok := transitions[t.CurrentState]; !ok {
			transitions[t.CurrentState] = make(map[string]domain.Transition)
		}

		transitions[t.CurrentState][t.Read] = domain.Transition{
			Write: t.Write,
			Move:  t.Move,
			Next:  t.NextState,
		}
	}

	machine := &domain.Machine{
		States:        yml.States,
		InputAlphabet: yml.InputAlphabet,
		TapeAlphabet:  yml.TapeAlphabet,
		Transitions:   transitions,
		StartState:    yml.StartState,
		Blank:         yml.Blank,
		AcceptStates:  yml.AcceptStates,
		RejectStates:  yml.RejectStates,
	}

	if err := machine.Validate(); err != nil {
		return nil, err
	}
	return machine, nil
}
