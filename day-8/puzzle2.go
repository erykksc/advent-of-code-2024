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
	antinodes := make(map[[2]int]bool)
	for _, fAntennas := range antennas {
		// Compare each antenna to other antennas ot find the anti nodes
		for i, antenna1 := range fAntennas {
			for j, antenna2 := range fAntennas {
				if i == j {
					continue
				}
				antinodes[antenna1] = true
				antinodes[antenna2] = true

				yDiff := antenna1[0] - antenna2[0]
				xDiff := antenna1[1] - antenna2[1]

				// go to the direction of antenna1 from antenna2
				cNode := [2]int{antenna1[0], antenna1[1]}
				for cNode[0] > -1 && cNode[0] < yDim && cNode[1] > -1 && cNode[1] < xDim {
					cNode[0] += yDiff
					cNode[1] += xDiff
					// Check if out of bounds
					if !(cNode[0] > -1 && cNode[0] < yDim && cNode[1] > -1 && cNode[1] < xDim) {
						break
					}
					antinodes[cNode] = true
				}
				// go to the direction of antenna2 from antenna1
				cNode[0], cNode[1] = antenna2[0], antenna2[1]
				for cNode[0] > -1 && cNode[0] < yDim && cNode[1] > -1 && cNode[1] < xDim {
					cNode[0] -= yDiff
					cNode[1] -= xDiff
					// Check if out of bounds
					if !(cNode[0] > -1 && cNode[0] < yDim && cNode[1] > -1 && cNode[1] < xDim) {
						break
					}
					antinodes[cNode] = true
				}
			}
		}
	}

	// Copy the input to output
	output := make([][]rune, len(lines))
	for i, line := range lines {
		output[i] = []rune(line)
	}

	// Mark all antinodes on the mapStr that aren't outside of boundries
	for antinode, _ := range antinodes {
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

	return result.String(), len(antinodes)
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
