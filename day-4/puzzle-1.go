package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func countOccurances(s string) int {
	// Compile the regular expression
	re1 := regexp.MustCompile("XMAS")
	re2 := regexp.MustCompile("SAMX")

	// Find all matches
	matches1 := re1.FindAllStringIndex(s, -1)
	matches2 := re2.FindAllStringIndex(s, -1)

	// Return the count of matches
	return len(matches1) + len(matches2)
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

	// Scan horizontally
	for _, line := range lines {
		occurences += countOccurances(line)
	}

	// Scan vertically
	for x := 0; x < gridSize; x++ {
		// Build string from column
		var colBuilder strings.Builder
		for y := 0; y < gridSize; y++ {
			colBuilder.WriteByte(lines[y][x])
		}
		column := colBuilder.String()

		occurences += countOccurances(column)
	}

	// Scan diagonally
	// Lower triangle left to right
	for startRow := 0; startRow < gridSize; startRow++ {
		var b strings.Builder

		// From left to right
		for x, y := 0, startRow; x < gridSize && y < gridSize; x, y = x+1, y+1 {
			b.WriteByte(lines[y][x])
		}
		occurences += countOccurances(b.String())
	}

	// Upper triangle left to right
	for startCol := 1; startCol < gridSize; startCol++ {
		var b strings.Builder
		for x, y := startCol, 0; x < gridSize && y < gridSize; x, y = x+1, y+1 {
			b.WriteByte(lines[y][x])
		}
		occurences += countOccurances(b.String())
	}

	// Lower triangle right to left
	for startRow := 0; startRow < gridSize; startRow++ {
		var b strings.Builder

		// From left to right
		for x, y := gridSize-1, startRow; x > -1 && y < gridSize; x, y = x-1, y+1 {
			b.WriteByte(lines[y][x])
		}
		occurences += countOccurances(b.String())
	}

	// Upper triangle right to left
	for startCol := 0; startCol < gridSize-1; startCol++ {
		var b strings.Builder
		for x, y := startCol, 0; x > -1 && y < gridSize; x, y = x-1, y+1 {
			b.WriteByte(lines[y][x])
		}
		occurences += countOccurances(b.String())
	}

	fmt.Printf("Occurences: %d\n", occurences)
}
