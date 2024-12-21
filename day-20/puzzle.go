package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Route struct {
	start, exit Position
}

type Direction struct {
	y, x int
}

var (
	Up    = Direction{-1, 0}
	Down  = Direction{1, 0}
	Left  = Direction{0, -1}
	Right = Direction{0, 1}
)

var AllDirections = []Direction{Up, Down, Left, Right}

type Position struct {
	y, x int
}

func (pos Position) Move(dir Direction) Position {
	return Position{pos.y + dir.y, pos.x + dir.x}
}

func (pos Position) ManhattanDistance(target Position) int {
	xDist := pos.x - target.x
	if xDist < 0 {
		xDist = -xDist
	}
	yDist := pos.y - target.y
	if yDist < 0 {
		yDist = -yDist
	}
	return xDist + yDist
}

func (pos Position) WithinDistance(maxDist int) []Position {
	if maxDist < 0 {
		panic("Radius < 0: " + strconv.Itoa(maxDist))
	}
	if maxDist == 0 {
		return []Position{pos}
	}

	result := make([]Position, 0)
	// scan entire square excluding those outside of radius
	candidate := Position{pos.y - maxDist, pos.x - maxDist}
	endPos := Position{pos.y + maxDist, pos.x + maxDist}
	for candidate != endPos {
		dist := pos.ManhattanDistance(candidate)
		if dist <= maxDist {
			result = append(result, candidate)
		}
		// update candidate
		candidate.x++
		if candidate.x > endPos.x {
			candidate.x = 0
			candidate.y++
		}
	}

	return result
}

type Grid2D [][]rune

func (grid Grid2D) Get(pos Position) rune {
	if !grid.IsInBounds(pos) {
		panic("Position out of bounds")
	}
	return grid[pos.y][pos.x]
}

func (grid Grid2D) Set(pos Position, char rune) {
	if !grid.IsInBounds(pos) {
		panic("Position out of bounds")
	}
	grid[pos.y][pos.x] = char
}

func (grid Grid2D) IsInBounds(pos Position) bool {
	height, width := len(grid), len(grid[0])
	return pos.x >= 0 && pos.x < width && pos.y >= 0 && pos.y < height
}

func (grid Grid2D) String() string {
	var b strings.Builder
	for y := range grid {
		for x := range grid[y] {
			c := grid[y][x]
			b.WriteRune(c)
		}
		b.WriteRune('\n')
	}
	return b.String()
}

func main() {
	cheatDuration := flag.Uint("cheat", 2, "Duration of the cheat in picoseconds")
	minCheatSave := flag.Uint("saved", 100, "At least how many picoseconds saved to count towards summary")
	flag.Parse()
	file, err := os.Open(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	grid := make(Grid2D, 0)

	var startPos, endPos Position
	height := 0
	for scanner.Scan() {
		line := scanner.Text()
		grid = append(grid, make([]rune, len(line)))
		for x, c := range line {
			grid[height][x] = c
			if c == 'S' {
				startPos.y = height
				startPos.x = x
			} else if c == 'E' {
				endPos.y = height
				endPos.x = x
			}
		}
		height++
	}

	fmt.Printf(grid.String())
	fmt.Printf("StartPos: %v EndPos:%v\n", startPos, endPos)

	dists := make(map[Position]int)
	// measure distance from end to start
	dists[endPos] = 0
	queue := []Position{endPos}
	for len(queue) > 0 {
		pos := queue[0]
		queue = queue[1:]

		for _, dir := range AllDirections {
			adj := pos.Move(dir)
			if !grid.IsInBounds(adj) {
				continue
			}
			if _, exists := dists[adj]; exists {
				continue
			}
			if grid.Get(adj) == '#' {
				continue
			}
			queue = append(queue, adj)
			dists[adj] = dists[pos] + 1
		}
	}
	fmt.Printf("No cheating route takes %d picoseconds\n", dists[startPos])

	cheatRoutes := make(map[Route]int) // (start,end) -> savedTime
	visited := make(map[Position]bool)
	queue = []Position{startPos}
	for len(queue) > 0 {
		pos := queue[0]
		queue = queue[1:]
		if visited[pos] {
			continue
		}
		visited[pos] = true

		// Try to cheat
		for _, cheatExitPos := range pos.WithinDistance(int(*cheatDuration)) {
			if !grid.IsInBounds(cheatExitPos) {
				continue
			}
			if grid.Get(cheatExitPos) == '#' {
				continue
			}

			savedTime := dists[pos] - dists[cheatExitPos] - pos.ManhattanDistance(cheatExitPos)
			if savedTime < int(*minCheatSave) {
				continue
			}
			route := Route{pos, cheatExitPos}
			bestSavedTime, exists := cheatRoutes[route]
			if !exists || savedTime < bestSavedTime {
				cheatRoutes[route] = savedTime
			}
		}

		// Advance in the race track
		for _, dir := range AllDirections {
			adj := pos.Move(dir)
			if visited[adj] {
				continue
			}
			if !grid.IsInBounds(adj) {
				continue
			}
			if grid.Get(adj) == '#' {
				continue
			}
			queue = append(queue, adj)
		}
	}

	fmt.Printf("%d PICOSECONDS CHEATS:\n", *cheatDuration)
	savedTimeRoutes := make(map[int]int) // saved time -> number of routes
	for _, savedTime := range cheatRoutes {
		savedTimeRoutes[savedTime]++
	}
	var savedTimes []int
	for savedTime := range savedTimeRoutes {
		savedTimes = append(savedTimes, savedTime)
	}
	slices.Sort(savedTimes)
	savedSummary := 0
	for _, savedTime := range savedTimes {
		numberOfRoutes := savedTimeRoutes[savedTime]
		fmt.Printf("There are %4d cheats that save %4d picoseconds\n", numberOfRoutes, savedTime)
		savedSummary += numberOfRoutes

	}
	fmt.Printf("No cheating route takes %d picoseconds\n", dists[startPos])
	fmt.Printf("SUMMARY:%d picosecond cheats that save at least %d picoseconds: %d \n", *cheatDuration, *minCheatSave, savedSummary)
}
