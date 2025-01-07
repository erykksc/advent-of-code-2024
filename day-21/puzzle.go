package main

import (
	"fmt"
	"os"
	"slices"
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
var PosInNumpad = map[rune]Position{
	'7': {0, 0}, '8': {0, 1}, '9': {0, 2},
	'4': {1, 0}, '5': {1, 1}, '6': {1, 2},
	'1': {2, 0}, '2': {2, 1}, '3': {2, 2},
	' ': {3, 0}, '0': {3, 1}, 'A': {3, 2},
}

//	+---+---+
//	| ^ | A |
//
// +---+---+---+
// | < | v | > |
// +---+---+---+
var PosInDirpad = map[rune]Position{
	' ': {0, 0}, '^': {0, 1}, 'A': {0, 2},
	'<': {1, 0}, 'v': {1, 1}, '>': {1, 2},
}

// Graph[start][target] = path
type Graph map[rune]map[rune]string

func CreateGraph(posInPad map[rune]Position) Graph {
	graph := make(Graph)

	for startR, start := range posInPad {
		graph[startR] = make(map[rune]string)
		for targetR, target := range posInPad {
			path := strings.Repeat("<", max(0, start.x-target.x)) +
				strings.Repeat("v", max(0, target.y-start.y)) +
				strings.Repeat("^", max(0, start.y-target.y)) +
				strings.Repeat(">", max(0, target.x-start.x))

			// check if the path goes through an empty space
			// if so, reverse the order (go vertically instead of vertically first and vice versa)
			empty := posInPad[' ']

			if (empty.y == start.y && empty.x == target.x) ||
				(empty.y == target.y && empty.x == start.x) {
				pathR := []rune(path)
				slices.Reverse(pathR)
				path = string(pathR)
			}

			graph[startR][targetR] = path
		}
	}

	return graph
}

var numpadGraph = CreateGraph(PosInNumpad)
var dirpadGraph = CreateGraph(PosInDirpad)

func (graph Graph) Use(armAt rune, code []rune) (path []rune) {
	var pathB strings.Builder

	startIdx := -1
	start := armAt
	for startIdx < len(code)-1 {
		target := code[startIdx+1]
		pathB.WriteString(graph[start][target])
		pathB.WriteRune('A')

		startIdx++
		start = code[startIdx]
	}
	return []rune(pathB.String())
}

func (graph Graph) UseTimes(code []rune, times int) (path []rune) {
	for range times {
		code = graph.Use('A', code)
	}
	return code
}

// Divides x by y and rounds up
func DivCeil(x, y int) int {
	if x%y == 0 {
		return x / y
	}
	return x/y + 1
}

func SplitIntoChunks(s []rune, size int) [][]rune {
	chunks := make([][]rune, 0, DivCeil(len(s), size))
	for i := 0; i < len(s); i += size {
		end := i + size
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}

// Same as use times, but splits the code into multiple parts if it is too long
func (graph Graph) UseTimesOptimized(armAt rune, code []rune, times int) (pathlen int) {
	chunkSize := 50000000
	// fmt.Printf("armAt %c, code %s, times %d\n", armAt, string(code), times)
	if times == 0 {
		return len(code)
	}
	if len(code) <= chunkSize {
		code = graph.Use(armAt, code)
		return graph.UseTimesOptimized('A', code, times-1)
	}

	// split code into chunks of 100000000
	chunks := SplitIntoChunks(code, chunkSize)

	total := 0
	armAtAfterChunk := armAt
	for i, chunk := range chunks {
		fmt.Printf("%sProcessing %d/%d chunks\n", strings.Repeat("\t", times), i+1, len(chunks))
		total += graph.UseTimesOptimized(armAtAfterChunk, chunk, times)
		armAtAfterChunk = chunk[len(chunk)-1]
	}

	return total
}

func TestAllPaths(graph Graph) {
	for startR, targets := range graph {
		for target, path := range targets {
			fmt.Printf("From %c to %c: %s\n", startR, target, path)
		}
	}
}

func part1(input string) {
	codes := strings.Split(input, "\n")

	complexitiesSum := 0
	for i, codeS := range codes {
		code := []rune(codeS)
		solution := numpadGraph.Use('A', code)
		solution = dirpadGraph.Use('A', solution)
		solution = dirpadGraph.Use('A', solution)

		fmt.Printf("%2d. Solution for code %s: %s\n", i, string(code), string(solution))

		numInCode, err := strconv.Atoi(codeS[:len(codeS)-1])
		if err != nil {
			panic(err)
		}
		fmt.Printf("Complexity: %d * %d = %d\n\n", len(solution), numInCode, numInCode*len(solution))
		complexitiesSum += numInCode * len(solution)
	}
	fmt.Printf("Complexities sum: %d\n", complexitiesSum)
}

func part2(input string) {
	codes := strings.Split(input, "\n")

	complexitiesSum := 0
	for _, codeS := range codes {
		code := []rune(codeS)
		numpadSolution := numpadGraph.Use('A', code)

		// process character by character, to avoid saving the string into memory
		// solution := dirpadGraph.UseTimes(numpadSolution, 19)
		// solutionLen := len(solution)
		optimized := dirpadGraph.UseTimesOptimized('A', numpadSolution, 25)
		// fmt.Printf("len(solution): %d, optimized: %d\n", len(solution), optimized)
		solutionLen := optimized

		numInCode, err := strconv.Atoi(codeS[:len(codeS)-1])
		if err != nil {
			panic(err)
		}
		fmt.Printf("Code complexity: %d * %d = %d\n\n", solutionLen, numInCode, numInCode*solutionLen)
		complexitiesSum += numInCode * solutionLen
	}
	fmt.Printf("Complexities sum: %d\n", complexitiesSum)
}

func part2Pipeline(input string) {
}

func main() {
	inputB, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	input := strings.TrimSpace(string(inputB))

	fmt.Println("-----------------Part 1-----------------")
	part1(input)
	fmt.Println("-----------------Part 2-----------------")
	part2(input)
}
