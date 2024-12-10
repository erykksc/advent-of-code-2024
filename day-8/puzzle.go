package main

import (
	"fmt"
	"os"
	"strings"
)

func Solve(mapStr string) (resultMap string, antinodesCount int) {
	xDim, yDim := 0, 0
	lines := strings.Split(mapStr, "\n")
	// remove empty lines from lines
	for i := 0; i < len(lines); i++ {
		if lines[i] == "" {
			lines = append(lines[:i], lines[i+1:]...)
		}
	}

	antennas := make(map[rune][][2]int, 0)
	for y, line := range lines {
		if line == "" {
			continue
		}

		// Get the dimensions
		if y == 0 {
			xDim = len(line)
		}
		yDim++
		for x, c := range line {
			if c == '.' {
				continue
			}
			antennas[c] = append(antennas[c], [2]int{y, x})
		}
	}

	// Calculate the positions of antinodes
	antinodes := make(map[rune][][2]int, 0)
	for frequency, fAntennas := range antennas {
		// Compare each antenna to other antennas ot find the anti nodes
		for i, antenna1 := range fAntennas {
			for j, antenna2 := range fAntennas {
				if i == j {
					continue
				}
				yDiff := antenna1[0] - antenna2[0]
				xDiff := antenna1[1] - antenna2[1]
				antinode1 := [2]int{antenna1[0] + yDiff, antenna1[1] + xDiff}
				antinode2 := [2]int{antenna2[0] - yDiff, antenna2[1] - xDiff}
				antinodes[frequency] = append(antinodes[frequency], antinode1, antinode2)
			}
		}
	}

	// Copy the input to output
	output := make([][]rune, len(lines))
	for i, line := range lines {
		output[i] = []rune(line)
	}

	// create a set with all antinodes
	antinodesSet := make(map[[2]int]bool)
	for _, fAntinodes := range antinodes {
		for _, antinode := range fAntinodes {
			y, x := antinode[0], antinode[1]
			if y < 0 || y >= yDim || x < 0 || x >= xDim {
				continue
			}
			antinodesSet[antinode] = true
		}
	}

	antinodesCount = len(antinodesSet)

	// Mark all antinodes on the mapStr that aren't outside of boundries
	for antinode, _ := range antinodesSet {
		y, x := antinode[0], antinode[1]
		// Out of map bounds
		if y < 0 || y >= yDim || x < 0 || x >= xDim {
			continue
		}
		if output[y][x] != '.' {
			continue
		}
		output[y][x] = '#'
	}

	// convert output to string
	var result strings.Builder
	for _, line := range output {
		_, err := result.WriteString(string(line))
		if err != nil {
			panic(err)
		}
		_, err = result.WriteRune('\n')
		if err != nil {
			panic(err)
		}
	}

	return result.String(), antinodesCount
}

func main() {
	inputB, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	input := string(inputB)

	solution, count := Solve(input)
	fmt.Println(solution)
	fmt.Printf("Antinodes count: %d\n", count)
}
