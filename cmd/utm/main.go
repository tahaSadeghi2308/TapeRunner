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
//
// Optionally pass a 4th argument: a symbol-map file (digit=original_symbol
// per line). When given, the program ALSO prints M's final tape decoded
// back into its own, human-readable alphabet - that decoded line is the
// part that actually matters, not the raw digit/marker soup.
//
//	go run ./cmd/utm machines/utm_table.tsv machines/utm_demo6_final_gate_2bit.txt machines/utm_demo6_final_gate_2bit.map
func main() {
	if len(os.Args) != 3 && len(os.Args) != 4 {
		fmt.Println("usage: go run ./cmd/utm <table_file> <input_file> [symbol_map_file]")
		fmt.Println("  table_file      : path to machines/utm_table.tsv")
		fmt.Println("  input_file      : path to a demo file containing one line: rules#w")
		fmt.Println("  symbol_map_file : optional, e.g. machines/utm_demoX.map, to decode the final tape")
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

	// The part that actually matters: M's own final tape, in M's own alphabet.
	if len(os.Args) == 4 {
		symMap, err := machines.LoadSymbolMap(os.Args[3])
		if err != nil {
			fmt.Println("error loading symbol map:", err)
			os.Exit(1)
		}
		decoded := machines.DecodeFinalTape(final.TapeContent, symMap)
		fmt.Println()
		fmt.Println("Decoded final tape of M (M's own alphabet):")
		fmt.Println(decoded)
	}
}
