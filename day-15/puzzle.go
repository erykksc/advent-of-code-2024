package main

import (
	"fmt"
	"os"
	"strings"
)

type Direction struct {
	y, x int
}

var (
	Up    Direction = Direction{-1, 0}
	Down  Direction = Direction{1, 0}
	Left  Direction = Direction{0, -1}
	Right Direction = Direction{0, 1}
)

func (d Direction) String() string {
	switch d {
	case Up:
		return "Up"
	case Down:
		return "Down"
	case Left:
		return "Left"
	case Right:
		return "Right"
	}
	return "Unknown"
}

func parseGrid(gridStr string) [][]rune {
	grid := make([][]rune, 0)
	for y, line := range strings.Split(gridStr, "\n") {
		grid = append(grid, make([]rune, len(line)))
		for x, cell := range line {
			grid[y][x] = cell
		}
	}

	return grid
}

func NewWideGrid(grid [][]rune) [][]rune {
	wGrid := make([][]rune, 0, len(grid)) // height stays the same
	for y, row := range grid {
		wGrid = append(wGrid, make([]rune, 2*len(row)))
		for x, c := range row {
			switch c {
			case '#':
				wGrid[y][x*2] = '#'
				wGrid[y][x*2+1] = '#'
			case 'O':
				wGrid[y][x*2] = '['
				wGrid[y][x*2+1] = ']'
			case '.':
				wGrid[y][x*2] = '.'
				wGrid[y][x*2+1] = '.'

			case '@':
				wGrid[y][x*2] = '@'
				wGrid[y][x*2+1] = '.'
			}
		}
	}

	return wGrid
}

func parseMoves(movesStr string) []Direction {
	moves := make([]Direction, 0)
	for _, c := range movesStr {
		switch c {
		case '^':
			moves = append(moves, Up)
		case 'v':
			moves = append(moves, Down)
		case '<':
			moves = append(moves, Left)
		case '>':
			moves = append(moves, Right)
		}
	}

	return moves
}

func applyMove(grid [][]rune, move Direction) {
	// find player
	var py, px int
findPlayerLoop:
	for y := range len(grid) {
		for x := range len(grid[y]) {
			if grid[y][x] == '@' {
				py = y
				px = x
				break findPlayerLoop
			}
		}
	}

	moveObject(grid, py, px, move)

}

func canMove(grid [][]rune, y, x int, d Direction) bool {
	fmt.Printf("Checking if can move %s from (%d, %d)\n", d, y, x)
	height, width := len(grid), len(grid[0])
	if y < 0 || x < 0 || y >= height || x >= width {
		panic("Trying to move a cell that is out of bounds" + fmt.Sprintf("(%d, %d)", y, x))
	}

	c := grid[y][x]

	if c == '#' {
		return false
	}
	if c == '.' {
		return true
	}

	if d == Up || d == Down {
		if c == '[' {
			// next one and check from the side
			return canMove(grid, y+d.y, x, d) && canMove(grid, y+d.y, x+1, d)
		}
		if c == ']' {
			// next one and check from the side
			return canMove(grid, y+d.y, x, d) && canMove(grid, y+d.y, x-1, d)
		}
	}

	return canMove(grid, y+d.y, x+d.x, d)

}

func moveObject(grid [][]rune, y, x int, d Direction) {
	fmt.Printf("Moving object at (%d, %d) %s\n", y, x, d)
	height, width := len(grid), len(grid[0])
	if y < 0 || x < 0 || y >= height || x >= width {
		panic("Trying to move a cell that is out of bounds" + fmt.Sprintf("(%d, %d)", y, x))
	}

	c := grid[y][x]
	if !canMove(grid, y, x, d) {
		return
	}

	if c == '.' {
		return
	}

	moveObject(grid, y+d.y, x+d.x, d)
	grid[y+d.y][x+d.x] = c
	grid[y][x] = '.'

	// If moving a double box, move both sides
	if d == Up || d == Down {
		if c == '[' {
			moveObject(grid, y+d.y, x+1, d)
			grid[y+d.y][x+1] = ']'
			grid[y][x+1] = '.'
			return
		} else if c == ']' {
			moveObject(grid, y+d.y, x-1, d)
			grid[y+d.y][x-1] = '['
			grid[y][x-1] = '.'
			return
		}
	}

}

func sprintGrid(grid [][]rune) string {
	var b strings.Builder
	for _, row := range grid {
		for _, cell := range row {
			b.WriteRune(cell)
		}
		b.WriteRune('\n')
	}

	return b.String()
}

func sumOfGpsCoordinates(grid [][]rune) int {
	out := 0

	for y := range grid {
		for x := range grid[0] {
			c := grid[y][x]
			if c != 'O' && c != '[' {
				continue
			}

			out += 100*y + x
		}
	}
	return out
}

func main() {
	inputB, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	input := string(inputB)

	divider := strings.Index(input, "\n\n")
	if divider == -1 {
		panic("Invalid input, can't find divider between grid and move list")
	}

	fmt.Println("Input:\n", input)

	gridStr, movesStr := input[:divider], input[divider+2:]
	movesStr = strings.Trim(movesStr, "\n")

	grid := parseGrid(gridStr)
	moves := parseMoves(movesStr)
	wGrid := NewWideGrid(grid)

	for _, move := range moves {
		fmt.Printf("Move: %s\n", move)
		applyMove(grid, move)
		fmt.Println(sprintGrid(grid))
	}
	puzzle1result := sumOfGpsCoordinates(grid)
	fmt.Printf("Puzzle 1 result: %d\n\n", puzzle1result)

	fmt.Println("Start of part 2")
	fmt.Println(sprintGrid(wGrid))
	for _, move := range moves {
		fmt.Printf("Move: %s\n", move)
		applyMove(wGrid, move)
		fmt.Println(sprintGrid(wGrid))
	}
	puzzle2result := sumOfGpsCoordinates(wGrid)
	fmt.Printf("Puzzle 2 result: %d\n", puzzle2result)
}
