package ports

import "github.com/tahaSadeghi2308/TapeRunner/internal/domain"

type MachineReader interface {
	Read(filePath string) (*domain.Machine, error)
	Validator(filePath string) error
}
