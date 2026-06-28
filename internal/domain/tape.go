package domain

type Tape struct {
	content map[int]rune
	Blank   rune
}

func NewTape(blank rune, input string) *Tape {
	t := &Tape{
		content: make(map[int]rune),
		Blank:   blank,
	}
	for i, char := range input {
		t.content[i] = char
	}
	return t
}