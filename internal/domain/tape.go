package domain

type Tape struct {
	Symbols []string
	Head    int
	Blank   string
}

func NewTape(input string, blank string) *Tape {
	symbols := make([]string, len(input))
	for i, char := range input {
		symbols[i] = string(char)
	}
	if len(symbols) == 0 {
		symbols = []string{blank}
	}
	return &Tape{
		Symbols: symbols,
		Head:    0,
		Blank:   blank,
	}
}

func (t *Tape) Read() string {
	return t.Symbols[t.Head]
}

func (t *Tape) Write(symbol string) {
	t.Symbols[t.Head] = symbol
}

func (t *Tape) Move(direction string) {
	if direction == "R" {
		t.Head++
		if t.Head >= len(t.Symbols) {
			t.Symbols = append(t.Symbols, t.Blank)
		}
	} else if direction == "L" {
		t.Head--
		if t.Head < 0 {
			t.Symbols = append([]string{t.Blank}, t.Symbols...)
			t.Head = 0
		}
	}
}

func (t *Tape) String() string {
	result := ""
	for _, s := range t.Symbols {
		result += s
	}
	return result
}