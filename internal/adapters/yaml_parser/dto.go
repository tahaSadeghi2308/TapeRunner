package yamlparser

type yamlTransition struct {
	CurrentState string `yaml:"current_state"`
	Read         string `yaml:"read"`
	Write        string `yaml:"write"`
	Move         string `yaml:"move"`
	NextState    string `yaml:"next_state"`
}

type yamlMachine struct {
	States         []string         `yaml:"states"`
	InputAlphabet  []string         `yaml:"input_alphabet"`
	TapeAlphabet   []string         `yaml:"tape_alphabet"`
	Blank          string           `yaml:"blank"`
	StartState     string           `yaml:"start_state"`
	AcceptStates   []string         `yaml:"accept_states"`
	RejectStates   []string         `yaml:"reject_states"`
	Transitions    []yamlTransition `yaml:"transitions"`
}