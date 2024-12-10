package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func SolvePuzzle1(tMap [][]int) int {
	trailheads := [][2]int{}
	// Find all trailheads
	for y := range len(tMap) {
		for x, height := range tMap[y] {
			if height == 0 {
				trailheads = append(trailheads, [2]int{y, x})
			}
		}
	}

	totalScore := 0

	for _, trailhead := range trailheads {
		queue := [][2]int{trailhead}
		visited := make(map[[2]int]bool)
		score := 0

		for len(queue) > 0 {
			// Pop node from queue
			node := queue[0]
			queue = queue[1:]

			// Check if visited
			if visited[node] {
				continue
			} else {
				visited[node] = true
			}

			y, x := node[0], node[1]

			if tMap[y][x] == 9 {
				score++
				continue
			}

			// Find all reachable directions
			if y > 0 && tMap[y][x]+1 == tMap[y-1][x] {
				queue = append(queue, [2]int{y - 1, x})
			}
			if y < len(tMap)-1 && tMap[y][x]+1 == tMap[y+1][x] {
				queue = append(queue, [2]int{y + 1, x})
			}
			if x > 0 && tMap[y][x]+1 == tMap[y][x-1] {
				queue = append(queue, [2]int{y, x - 1})
			}
			if x < len(tMap[y])-1 && tMap[y][x]+1 == tMap[y][x+1] {
				queue = append(queue, [2]int{y, x + 1})
			}
		}
		totalScore += score
	}
	return totalScore
}

func SolvePuzzle2(tMap [][]int) int {
	trailheads := [][2]int{}
	// Find all trailheads
	for y := range len(tMap) {
		for x, height := range tMap[y] {
			if height == 0 {
				trailheads = append(trailheads, [2]int{y, x})
			}
		}
	}

	totalRating := 0

	for _, trailhead := range trailheads {
		queue := [][2]int{trailhead}
		rating := 0

		for len(queue) > 0 {
			// Pop node from queue
			node := queue[0]
			queue = queue[1:]

			y, x := node[0], node[1]

			if tMap[y][x] == 9 {
				rating++
				continue
			}

			// Find all reachable directions
			if y > 0 && tMap[y][x]+1 == tMap[y-1][x] {
				queue = append(queue, [2]int{y - 1, x})
			}
			if y < len(tMap)-1 && tMap[y][x]+1 == tMap[y+1][x] {
				queue = append(queue, [2]int{y + 1, x})
			}
			if x > 0 && tMap[y][x]+1 == tMap[y][x-1] {
				queue = append(queue, [2]int{y, x - 1})
			}
			if x < len(tMap[y])-1 && tMap[y][x]+1 == tMap[y][x+1] {
				queue = append(queue, [2]int{y, x + 1})
			}
		}
		totalRating += rating
	}
	return totalRating
}

func main() {
	fileB, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	input := string(fileB)

	// Parse input
	lines := strings.Split(strings.Trim(input, "\n"), "\n")
	tMap := make([][]int, len(lines))
	for y, line := range lines {
		tMap[y] = make([]int, len(line))
		for x, c := range line {
			if c == '.' {
				tMap[y][x] = -1
				continue
			}
			num, err := strconv.Atoi(string(c))
			if err != nil {
				panic(err)
			}
			tMap[y][x] = num
		}
	}

	// output tMap
	for _, row := range tMap {
		for _, cell := range row {
			print(cell)
		}
		println()
	}

	totalScore := SolvePuzzle1(tMap)
	fmt.Printf("Total score: %v\n", totalScore)

	totalRating := SolvePuzzle2(tMap)
	fmt.Printf("Total rating: %v\n", totalRating)
}
