package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// y,x coordinates
type Direction [2]int

var (
	Up    = Direction{-1, 0}
	Down  = Direction{1, 0}
	Left  = Direction{0, -1}
	Right = Direction{0, 1}
)

var symbol2dir = map[rune]Direction{
	'^': Up,
	'>': Right,
	'<': Left,
	'v': Down,
}

var dir2symbol = map[Direction]byte{
	Up:    '^',
	Right: '>',
	Left:  '<',
	Down:  'v',
}

var rotatedDir = map[Direction]Direction{
	Up:    Right,
	Right: Down,
	Left:  Up,
	Down:  Left,
}

// DEBUGGING FUNCTIONS
func printGrid(grid [][]byte, gPos [2]int, ySize, xSize int) {
	redAnsi := "\033[31m"  // ANSI code for red
	resetAnsi := "\033[0m" // ANSI code to reset formatting
	clearScreen()
	var b strings.Builder
	for y := range ySize {
		for x := range xSize {
			cell := grid[y][x]
			if y == gPos[0] && x == gPos[1] {
				b.WriteString(redAnsi)
				b.WriteByte(cell)
				b.WriteString(resetAnsi)
			} else {
				b.WriteByte(cell)
			}
		}
		b.WriteByte('\n')
	}

	fmt.Println(b.String())
}

func clearScreen() {
	cmd := exec.Command("clear") // Linux and macOS
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
	fmt.Println()
}

func main() {
	visual := flag.Bool("visual", false, "Visualize the guard's path")

	flag.Parse()
	filename := flag.Args()[0]

	inputB, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	input := string(inputB)

	lines := strings.Split(input, "\n")
	ySize := len(lines) - 1
	xSize := len(lines[0])

	grid := make([][]byte, ySize)
	for y := range ySize {
		grid[y] = make([]byte, xSize)
		for x := range xSize {
			grid[y][x] = lines[y][x]
		}
	}

	// Find guard dir and position
	gDir := [2]int{}
	gPos := [2]int{}
	for y, line := range lines {
		for x, letter := range line {
			if direction, exists := symbol2dir[letter]; exists {
				gPos = [2]int{y, x}
				gDir = direction
			}
		}
	}

	// Possible obstructions
	pObstructions := 0
	// Simulate guard walk

	isGridEscapable := func(grid [][]byte, gPos [2]int, gDir Direction) bool {
		for true {
			if *visual {
				printGrid(grid, gPos, ySize, xSize)
			}

			// Find out new position
			nextPos := [2]int{gPos[0] + gDir[0], gPos[1] + gDir[1]}
			// Check if out of bounds, reached exit
			if nextPos[0] < 0 || nextPos[0] >= ySize || nextPos[1] < 0 || nextPos[1] >= xSize {
				break
			}

			// Check if the newPos is on a rock, if so, change direction
			if grid[nextPos[0]][nextPos[1]] == '#' {
				gDir = rotatedDir[gDir]
				continue
			}

			if grid[nextPos[0]][nextPos[1]] == dir2symbol[gDir] {
				return false
			}

			// Move guard to next position
			gPos = nextPos
			grid[gPos[0]][gPos[1]] = dir2symbol[gDir]
		}
		return true
	}

	// Find possible postions for obstruction
	obstructionPos := make([][2]int, 0)
	for y := range ySize {
		for x := range xSize {
			// Don't test the starting field of guard
			if y == gPos[0] && x == gPos[1] {
				continue
			}
			if grid[y][x] != '#' {
				obstructionPos = append(obstructionPos, [2]int{y, x})
			}
		}
	}

	// Check all possibile positions for obstructions for loops
	for _, pos := range obstructionPos {
		if *visual {
			fmt.Printf("Checking obstruction at %d,%d\n", pos[0], pos[1])
		}

		gridC := make([][]byte, len(grid))
		for y := range grid {
			gridC[y] = make([]byte, xSize)
			copy(gridC[y], grid[y])
		}

		gridC[pos[0]][pos[1]] = '#'
		if !isGridEscapable(gridC, gPos, gDir) {
			if *visual {
				fmt.Println("Guard is stuck in a loop")
			}
			pObstructions++
		}
	}
	fmt.Printf("Number of possible obstructions in the grid: %d\n", pObstructions)
}
