package main

import (
	"fmt"
	"os"
	"strings"
)

func isXmas(s string) bool {
	// Optimalization
	if s[4] != 'A' {
		return false
	}

	var diag1b strings.Builder
	diag1b.WriteByte(s[0])
	diag1b.WriteByte(s[4])
	diag1b.WriteByte(s[8])
	diag1 := diag1b.String()

	var diag2b strings.Builder
	diag2b.WriteByte(s[2])
	diag2b.WriteByte(s[4])
	diag2b.WriteByte(s[6])
	diag2 := diag2b.String()

	if diag1 != "MAS" && diag1 != "SAM" {
		return false
	}

	if diag2 != "MAS" && diag2 != "SAM" {
		return false
	}

	return true
}

func main() {
	filename := os.Args[1]

	inputBytes, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	input := string(inputBytes)
	lines := strings.Split(input, "\n")

	occurences := 0
	gridSize := len(lines) - 1
	fmt.Printf("Grid size: %d\n", gridSize)

	// Scan the grid with squares
	for x := 0; x < gridSize-2; x = x + 1 {
		for y := 0; y < gridSize-2; y++ {
			var b strings.Builder
			b.WriteByte(lines[y][x])
			b.WriteByte(lines[y][x+1])
			b.WriteByte(lines[y][x+2])
			b.WriteByte(lines[y+1][x])
			b.WriteByte(lines[y+1][x+1])
			b.WriteByte(lines[y+1][x+2])
			b.WriteByte(lines[y+2][x])
			b.WriteByte(lines[y+2][x+1])
			b.WriteByte(lines[y+2][x+2])

			if isXmas(b.String()) {
				occurences++
			}
		}
	}

	fmt.Printf("Occurences: %d\n", occurences)
}
