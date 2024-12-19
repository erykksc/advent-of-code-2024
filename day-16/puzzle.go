package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

type Direction [2]int

var (
	North Direction = [2]int{-1, 0}
	South Direction = [2]int{1, 0}
	West  Direction = [2]int{0, -1}
	East  Direction = [2]int{0, 1}
)

var AllDirections = []Direction{North, South, West, East}

func (d Direction) Rune() rune {
	switch d {
	case North:
		return '^'
	case South:
		return 'v'
	case West:
		return '<'
	case East:
		return '>'
	}
	panic("Invalid direction")
}

func (d Direction) Rotations(target Direction) int {
	if d == target {
		return 0
	}
	if d[0] == -target[0] || d[1] == -target[1] {
		return 2
	}

	return 1 // Can rotate clockwise and anticlockwise
}

type Position struct {
	y, x   int
	facing Direction
}

func (p Position) Rotate() Position {
	facing := Direction{p.facing[1], -p.facing[0]}
	return Position{p.y, p.x, facing}
}

func (p Position) Forward() Position {
	return Position{p.y + p.facing[0], p.x + p.facing[1], p.facing}
}

func (p Position) DirectionTo(target Position) Direction {
	if target.y < p.y {
		return North
	} else if target.y > p.y {
		return South
	} else if target.x < p.x {
		return West
	} else if target.x > p.x {
		return East
	}

	// Same position, the default starting position is east
	return East
}

// Pop unvisited key with lowest dist
func popSmallest(unvisited map[Position]bool, dist map[Position]int) Position {
	if len(unvisited) == 0 {
		panic("No unvisited nodes")
	}
	minDist := math.MaxInt
	var minKey Position
	for key := range unvisited {
		d, exists := dist[key]
		if !exists {
			panic("Key from unvisited not in dist")
		}
		if d < minDist {
			minDist = d
			minKey = key
		}
	}

	delete(unvisited, minKey)
	return minKey
}

func markPrev(pos, sPos Position, grid [][]rune, prev map[Position][]Position) {
	fmt.Printf("Marking %v, prev: %v\n", pos, prev[pos])
	grid[pos.y][pos.x] = 'O'
	if pos == sPos {
		return
	}
	for _, prevPos := range prev[pos] {
		markPrev(prevPos, sPos, grid, prev)
	}
	return
}

func conver2grid(input string) (grid [][]rune, startPos Position, endPos [4]Position) {
	// convert input into a grid
	grid = make([][]rune, 0)
	for y, line := range strings.Split(input, "\n") {
		if len(line) == 0 {
			continue
		}
		grid = append(grid, make([]rune, len(line)))
		for x, c := range line {
			grid[y][x] = c
			if c == 'S' {
				startPos = Position{y, x, East}
			}
			if c == 'E' {
				for i, d := range AllDirections {
					endPos[i] = Position{y, x, d}
				}

			}
		}
	}
	return
}

func solvePuzzle1(input string) (dists map[Position]int) {
	// convert input into a grid
	grid := make([][]rune, 0)
	unvisited := make(map[Position]bool)
	var startPos Position
	for y, line := range strings.Split(input, "\n") {
		grid = append(grid, make([]rune, len(line)))
		for x, c := range line {
			grid[y][x] = c
			if c == '#' {
				continue
			}
			if c == 'S' {
				startPos = Position{
					y, x, East,
				}
			}

			for _, direction := range AllDirections {
				pos := Position{
					y,
					x,
					direction,
				}
				unvisited[pos] = true
			}
		}
	}

	// Run Dijsktra
	// Set all distances to infinity
	dists = make(map[Position]int)
	for key := range unvisited {
		dists[key] = math.MaxInt
	}
	dists[startPos] = 0

	prev := make(map[Position]Position)
	prev[startPos] = startPos

	// Start from starting spot
	for len(unvisited) > 0 {
		cPos := popSmallest(unvisited, dists)
		fmt.Printf("Visiting %v\n", cPos)

		// Move forward
		dist := dists[cPos] + 1
		nextPos := cPos.Forward()
		c := grid[nextPos.y][nextPos.x]
		if c != '#' && dist < dists[nextPos] {
			dists[nextPos] = dist
		}

		// Rotate
		nextPos = cPos
		for i := range 3 {
			nextPos = nextPos.Rotate()

			dist := dists[cPos] + 1000
			// Rotation of 180
			if i == 1 {
				dist += 1000
			}
			if dist < dists[nextPos] {
				dists[nextPos] = dist
				fmt.Printf("	Setting dist of rotate %v to %d\n", nextPos, dist)
			}
		}
	}
	return dists
}

func solvePuzzle2(input string) (dists map[Position]int, prev map[Position][]Position) {
	// convert input into a grid
	grid := make([][]rune, 0)
	unvisited := make(map[Position]bool)
	var startPos Position
	for y, line := range strings.Split(input, "\n") {
		grid = append(grid, make([]rune, len(line)))
		for x, c := range line {
			grid[y][x] = c
			if c == '#' {
				continue
			}
			if c == 'S' {
				startPos = Position{
					y, x, East,
				}
			}

			for _, direction := range AllDirections {
				pos := Position{
					y,
					x,
					direction,
				}
				unvisited[pos] = true
			}
		}
	}

	// Run Dijsktra
	// Set all distances to infinity
	dists = make(map[Position]int)
	for key := range unvisited {
		dists[key] = math.MaxInt
	}
	dists[startPos] = 0

	prev = make(map[Position][]Position)
	prev[startPos] = make([]Position, 1)
	prev[startPos][0] = startPos

	// Start from starting spot
	for len(unvisited) > 0 {
		cPos := popSmallest(unvisited, dists)
		fmt.Printf("Visiting %v\n", cPos)

		// Move forward
		dist := dists[cPos] + 1
		nextPos := cPos.Forward()
		c := grid[nextPos.y][nextPos.x]
		if c != '#' {
			if dist < dists[nextPos] {
				dists[nextPos] = dist
				prev[nextPos] = make([]Position, 1)
				prev[nextPos][0] = cPos
			} else if dist == dists[nextPos] {
				prev[nextPos] = append(prev[nextPos], cPos)
			}
		}

		// Rotate
		nextPos = cPos
		for i := range 3 {
			nextPos = nextPos.Rotate()

			dist := dists[cPos] + 1000
			// Rotation of 180
			if i == 1 {
				dist += 1000
			}
			if dist < dists[nextPos] {
				dists[nextPos] = dist
				prev[nextPos] = make([]Position, 1)
				prev[nextPos][0] = cPos
			} else if dist == dists[nextPos] {
				prev[nextPos] = append(prev[nextPos], cPos)
			}
		}
	}
	return
}

func main() {
	inputB, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	input := string(inputB)

	grid, sPos, ePositions := conver2grid(input)
	dists := solvePuzzle1(input)

	minDist := math.MaxInt
	for _, ePos := range ePositions {
		if dists[ePos] < minDist {
			minDist = dists[ePos]
		}
	}
	fmt.Printf("Distance from S to E: %d\n", minDist)

	_, _ = grid, sPos

	dists, prev := solvePuzzle2(input)
	for _, ePo := range ePositions {
		if dists[ePo] == minDist {
			markPrev(ePo, sPos, grid, prev)
			break
		}
	}

	marked := 0
	var b strings.Builder
	for y, row := range grid {
		b.WriteString(fmt.Sprintf("%02d: ", y))
		for _, c := range row {
			b.WriteRune(c)
			if c == 'O' {
				marked++
			}
		}
		b.WriteRune('\n')
	}
	fmt.Printf("Grid, size (height, width) %d, %d:\n%s", len(grid), len(grid[0]), b.String())
	fmt.Printf("Marked: %d\n", marked)
}
