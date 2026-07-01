package domain

type StepResult struct {
	CurrentState string
	TapeContent  string
	HeadPosition int
	Status       string
}
