package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Position struct {
	y, x int
}

func (p Position) Up() Position {
	return Position{p.y + 1, p.x}
}

func (p Position) Down() Position {
	return Position{p.y - 1, p.x}
}

func (p Position) Left() Position {
	return Position{p.y, p.x - 1}
}

func (p Position) Right() Position {
	return Position{p.y, p.x + 1}
}

func FindPath(grid [][]rune) []Position {
	if len(grid) == 0 {
		panic("Empty grid")
	}
	width, height := len(grid[0]), len(grid)

	parent := make(map[Position]Position)
	visited := make(map[Position]bool)

	// BFS through grid
	queue := []Position{}
	queue = append(queue, Position{0, 0})
	for len(queue) > 0 {
		// pop first item from unvisited
		cPos := queue[0]
		queue = queue[1:]

		// if out of boundries
		if cPos.y < 0 || cPos.y >= height || cPos.x < 0 || cPos.x >= width {
			continue
		}

		if visited[cPos] {
			continue
		}

		cNode := grid[cPos.y][cPos.x]
		if cNode == '#' {
			continue
		}

		// Reached target
		target := Position{height - 1, width - 1}
		if cPos == target {
			return ReconstructPath(parent, Position{0, 0}, target)
		}

		visited[cPos] = true

		up := cPos.Up()
		down := cPos.Down()
		left := cPos.Left()
		right := cPos.Right()
		queue = append(queue, up, down, left, right)

		// Set parents of positions found from cPos
		if !visited[up] {
			parent[up] = cPos
		}
		if !visited[down] {
			parent[down] = cPos
		}
		if !visited[left] {
			parent[left] = cPos
		}
		if !visited[right] {
			parent[right] = cPos
		}
	}
	return nil

}

func ReconstructPath(parent map[Position]Position, start, target Position) []Position {
	path := make([]Position, 0)
	for target != start {
		path = append(path, target)
		target = parent[target]
	}
	path = append(path, start)
	slices.Reverse(path)
	return path
}

func printGrid(height, width int, path, bytesPos []Position) {
	onPath := make(map[Position]bool)
	for _, pos := range path {
		onPath[pos] = true
	}
	corrupted := make(map[Position]bool)
	for _, pos := range bytesPos {
		corrupted[pos] = true
	}

	var gridStr strings.Builder
	// print out grid
	for y := range height {
		for x := range width {
			pos := Position{y, x}
			if onPath[pos] {
				gridStr.WriteRune('O')
				continue
			}
			if corrupted[pos] {
				gridStr.WriteRune('#')
				continue
			}
			gridStr.WriteRune('.')
		}
		gridStr.WriteRune('\n')
	}
	fmt.Println(gridStr.String())
}

func main() {
	bytes2read := flag.Int("bytes", 0, "Number of bytes to read")
	flag.Parse()
	file, err := os.Open(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Get the grid size
	scanner.Scan()
	gridLine := strings.TrimLeft(scanner.Text(), "grid: ")
	var height, width int
	_, err = fmt.Sscanf(gridLine, "%d,%d", &height, &width)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Height: %d, Width: %d\n", height, width)
	scanner.Scan() // Scan the empty line between grid size and falling bytes coordinates

	bytesPos := make([]Position, 0)
	for scanner.Scan() {
		line := scanner.Text()
		var y, x int
		_, err := fmt.Sscanf(line, "%d,%d", &x, &y)
		if err != nil {
			panic(err)
		}
		bytesPos = append(bytesPos, Position{y, x})
		if len(bytesPos) == *bytes2read {
			break
		}
	}

	// Print grid with all bytes drawn
	// printGrid(height, width, nil, bytesPos)

	var path []Position
	for i := 0; i < len(bytesPos); i++ {
		fmt.Printf("Reading byte %d: %d,%d\n", i+1, bytesPos[i].x, bytesPos[i].y)
		// copy grid to have fresh grid for motivating
		grid := make([][]rune, height)
		for y := range height {
			grid[y] = make([]rune, width)
			for x := range width {
				grid[y][x] = '.'
			}
		}
		for _, pos := range bytesPos[:i+1] {
			grid[pos.y][pos.x] = '#'
		}

		path = FindPath(grid)
		if path == nil {
			fmt.Printf("After %d bytes: no path found, last corrupted byte: %d,%d\n", i+1, bytesPos[i].x, bytesPos[i].y)
			return
		}

		for _, pos := range path {
			grid[pos.y][pos.x] = 'O'
		}

	}
	fmt.Printf("After %d bytes: shortest route: %d\n", *bytes2read, len(path)-1) // -1 to path length as they don't count the initial position as a step
	printGrid(height, width, path, bytesPos[:*bytes2read])
}
