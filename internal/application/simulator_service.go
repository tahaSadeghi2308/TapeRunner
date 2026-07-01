package application

import (
	"github.com/tahaSadeghi2308/TapeRunner/internal/domain"
)

type SimulationService struct{}

func NewSimulationService() *SimulationService {
	return &SimulationService{}
}

func (s *SimulationService) Run(machine *domain.Machine, initialTape string) []domain.StepResult {
	tape := domain.NewTape(initialTape, "_")
	simulator := domain.NewSimulator(machine, tape)
	
	return simulator.Run()
}