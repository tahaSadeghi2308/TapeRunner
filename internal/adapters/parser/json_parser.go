package parser

import (
	"encoding/json"
	"fmt"

	"github.com/tahaSadeghi2308/TapeRunner/internal/domain"
)

func ParseMachine(data []byte) (*domain.Machine, error) {
	var m domain.Machine
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to parse machine JSON: %w", err)
	}

	if m.StartState == "" {
		return nil, fmt.Errorf("StartState is required")
	}
	if len(m.Transitions) == 0 {
		return nil, fmt.Errorf("Transitions are required")
	}

	return &m, nil
}