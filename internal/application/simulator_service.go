package application

import (
	"github.com/tahaSadeghi2308/TapeRunner/internal/domain"
	"github.com/tahaSadeghi2308/TapeRunner/internal/ports"
)

type SimulationService struct {
	MachineSetup ports.MachineReader
}

func NewSimulationService(machineSetup ports.MachineReader) *SimulationService {
	return &SimulationService{
		MachineSetup: machineSetup,
	}
}

func (s *SimulationService) Run(machine *domain.Machine, initialTape string) []domain.StepResult {
	tape := domain.NewTape(initialTape, machine.Blank)
	simulator := domain.NewSimulator(machine, tape)

	return simulator.Run()
}
