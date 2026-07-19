package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/tahaSadeghi2308/TapeRunner/internal/machines"
)

// Run with:
//
//	go run ./cmd/utm machines/utm_table.tsv machines/utm_demo1.txt
func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage: go run ./cmd/utm <table_file> <input_file>")
		fmt.Println("  table_file : path to machines/utm_table.tsv")
		fmt.Println("  input_file : path to a demo file containing one line: rules#w")
		os.Exit(1)
	}

	tablePath := os.Args[1]
	inputPath := os.Args[2]

	m, err := machines.LoadUTMTable(tablePath)
	if err != nil {
		fmt.Println("error loading UTM table:", err)
		os.Exit(1)
	}

	raw, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Println("error reading input file:", err)
		os.Exit(1)
	}

	tapeStr, err := machines.BuildEncodedTape(strings.TrimSpace(string(raw)))
	if err != nil {
		fmt.Println("error building encoded tape:", err)
		os.Exit(1)
	}

	const maxSteps = 500000
	history := machines.RunUTM(m, tapeStr, maxSteps)

	for i, step := range history {
		fmt.Printf("Step %4d | State: %-10s | Head: %-4d | Tape: %s\n",
			i, step.CurrentState, step.HeadPosition, step.TapeContent)
	}

	final := history[len(history)-1]
	fmt.Println()
	fmt.Println("Final status:", final.Status)
}
