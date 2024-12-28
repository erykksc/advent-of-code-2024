package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Position struct {
	y, x int
}

// Num pad
// +---+---+---+
// | 7 | 8 | 9 |
// +---+---+---+
// | 4 | 5 | 6 |
// +---+---+---+
// | 1 | 2 | 3 |
// +---+---+---+
//
//	| 0 | A |
//	+---+---+
var PosInNumerical = map[rune]Position{
	'7': {0, 0},
	'8': {0, 1},
	'9': {0, 2},
	'4': {1, 0},
	'5': {1, 1},
	'6': {1, 2},
	'1': {2, 0},
	'2': {2, 1},
	'3': {2, 2},
	'0': {3, 1},
	'A': {3, 2},
}

//	+---+---+
//	| ^ | A |
//
// +---+---+---+
// | < | v | > |
// +---+---+---+
var PosInDirectional = map[rune]Position{
	'^': {0, 1},
	'A': {0, 2},
	'<': {1, 0},
	'v': {1, 1},
	'>': {1, 2},
}

func solve(code string, symbolPos map[rune]Position) string {
	armPos, exists := symbolPos['A']
	if !exists {
		panic("Arm position not found")
	}
	var solution strings.Builder
	for _, symbol := range code {
		targetPos, exists := symbolPos[symbol]
		if !exists {
			panic("Target position not found, target: " + string(symbol))
		}
		for armPos != targetPos {
			if armPos.x > targetPos.x {
				armPos.x--
				solution.WriteRune('<')
			} else if armPos.x < targetPos.x {
				armPos.x++
				solution.WriteRune('>')
			} else if armPos.y < targetPos.y {
				armPos.y++
				solution.WriteRune('v')
			} else if armPos.y > targetPos.y {
				armPos.y--
				solution.WriteRune('^')
			}
		}
		solution.WriteRune('A')
	}
	return solution.String()
}

func complexity(code, solution string) int {
	numInCode, err := strconv.Atoi(code[:len(code)-1])
	if err != nil {
		panic(err)
	}

	return numInCode * len(solution)
}

func main() {
	input, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	codes := strings.TrimSpace(string(input))

	complexitiesSum := 0
	for _, code := range strings.Split(codes, "\n") {
		// numerical
		numSolution := solve(code, PosInNumerical)
		fmt.Printf("%s: %s (len: %d)\n", code, numSolution, len(numSolution))
		// directional 1
		dir1Solution := solve(numSolution, PosInDirectional)
		fmt.Printf("      %s (len: %d)\n", dir1Solution, len(dir1Solution))
		// directional 2
		dir2Solution := solve(dir1Solution, PosInDirectional)
		fmt.Printf("      %s (len: %d)\n", dir2Solution, len(dir2Solution))

		compx := complexity(code, dir2Solution)
		fmt.Printf("Complexity: %d\n", compx)
		complexitiesSum += compx
	}
	fmt.Printf("Sum of complexities: %d\n", complexitiesSum)
}
