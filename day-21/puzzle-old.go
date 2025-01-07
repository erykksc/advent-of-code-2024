// package main
//
// import (
// 	"fmt"
// 	"os"
// 	"strconv"
// 	"strings"
// )
//
// type Position struct {
// 	y, x int
// }
//
// // Num pad
// // +---+---+---+
// // | 7 | 8 | 9 |
// // +---+---+---+
// // | 4 | 5 | 6 |
// // +---+---+---+
// // | 1 | 2 | 3 |
// // +---+---+---+
// //
// //	| 0 | A |
// //	+---+---+
// var PosInNumerical = map[rune]Position{
// 	'7': {0, 0}, '8': {0, 1}, '9': {0, 2},
// 	'4': {1, 0}, '5': {1, 1}, '6': {1, 2},
// 	'1': {2, 0}, '2': {2, 1}, '3': {2, 2},
// 	' ': {3, 0}, '0': {3, 1}, 'A': {3, 2},
// }
//
// //	+---+---+
// //	| ^ | A |
// //
// // +---+---+---+
// // | < | v | > |
// // +---+---+---+
// var PosInDirectional = map[rune]Position{
// 	' ': {0, 0}, '^': {0, 1}, 'A': {0, 2},
// 	'<': {1, 0},
// 	'v': {1, 1},
// 	'>': {1, 2},
// }
//
// type Graph map[rune]map[rune]string
//
// func CreateGraph(keypad map[rune]Position, emptyField Position) Graph {
// 	graph := make(Graph)
//
// 	for a, key1 := range keypad {
// 		for b, key2 := range keypad {
// 			path := strings.Repeat("<", key2.x-key1.x) +
// 				strings.Repeat("v", key2.y-key1.y) +
// 				strings.Repeat("^", key2.x-key1.x) +
// 				strings.Repeat(">", key1.x-key2.x)
//
// 		}
// 	}
//
// 	return graph
//
// }
//
// func routes(start, target, emptyField Position) []string {
// 	solutions := []string{}
// 	var solution strings.Builder
// 	arm := start
// 	// go horizontally, then vertically
// 	for arm != target && arm != emptyField {
// 		if arm.x > target.x {
// 			arm.x--
// 			solution.WriteRune('<')
// 		} else if arm.x < target.x {
// 			arm.x++
// 			solution.WriteRune('>')
// 		} else if arm.y < target.y {
// 			arm.y++
// 			solution.WriteRune('v')
// 		} else if arm.y > target.y {
// 			arm.y--
// 			solution.WriteRune('^')
// 		}
// 	}
// 	if arm != emptyField {
// 		solution.WriteRune('A')
// 		solutions = append(solutions, solution.String())
// 	}
//
// 	// go vertically, then horizontally
// 	solution.Reset()
// 	arm = start
// 	for arm != target {
// 		if arm.y < target.y {
// 			arm.y++
// 			solution.WriteRune('v')
// 		} else if arm.y > target.y {
// 			arm.y--
// 			solution.WriteRune('^')
// 		} else if arm.x > target.x {
// 			arm.x--
// 			solution.WriteRune('<')
// 		} else if arm.x < target.x {
// 			arm.x++
// 			solution.WriteRune('>')
// 		}
// 	}
// 	solution.WriteRune('A')
//
// 	// Add if found new solution
// 	if len(solutions) == 0 {
// 		solutions = append(solutions, solution.String())
// 	} else if solution.String() != solutions[0] {
// 		solutions = append(solutions, solution.String())
// 	}
//
// 	return solutions
// }
//
// func Parts(prevParts [][]string, symbolPos map[rune]Position) [][]string {
// 	armPos, exists := symbolPos['A']
// 	if !exists {
// 		panic("Arm position not found")
// 	}
//
// 	parts := make([][]string, len(prevParts))
//
// 	for i, symbol := range prevParts {
// 		targetPos, exists := symbolPos[symbol]
// 		if !exists {
// 			panic("Target position not found, target: " + string(symbol))
// 		}
//
// 		armRoutes := routes(armPos, targetPos, symbolPos[' '])
// 		parts[i] = armRoutes
// 		// fmt.Printf("Route from: %v -> %v = %s\n", armPos, targetPos, armRoute)
// 		armPos = targetPos
// 	}
//
// 	return parts
// }
//
// // Returns all Solutions made from combinations of all paths
// func Solutions(parts [][]string) []string {
// 	if len(parts) == 0 {
// 		return []string{}
// 	}
// 	if len(parts) == 1 {
// 		return parts[0]
// 	}
// 	solutions := []string{}
// 	part := parts[0]
// 	for _, variant := range part {
// 		fSolutions := Solutions(parts[1:])
// 		for _, fSol := range fSolutions {
// 			solutions = append(solutions, variant+fSol)
// 		}
// 	}
// 	return solutions
// }
//
// func complexity(code, solution string) int {
// 	numInCode, err := strconv.Atoi(code[:len(code)-1])
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	return numInCode * len(solution)
// }
//
// func mainOLD() {
// 	input, err := os.ReadFile(os.Args[1])
// 	if err != nil {
// 		panic(err)
// 	}
// 	codes := strings.TrimSpace(string(input))
//
// 	complexitiesSum := 0
// 	for _, code := range strings.Split(codes, "\n") {
// 		fmt.Printf("Solving code %s\n", code)
// 		// numerical
// 		parts := Parts(code, PosInNumerical)
// 		numSolutions := Solutions(parts)
// 		for i, solution := range numSolutions {
// 			fmt.Printf("Solution %d: %s\n", i, solution)
// 		}
// 		parts = Parts(code, PosInDirectional)
// 		dirSolutions := Solutions(parts)
// 		for i, solution := range dirSolutions {
// 			fmt.Printf("Solution %d: %s\n", i, solution)
// 		}
//
// 		parts = Parts(code, PosInDirectional)
// 		dir2Solutions := Solutions(parts)
// 		for i, solution := range dir2Solutions {
// 			fmt.Printf("Solution %d: %s\n", i, solution)
// 		}
//
// 		fmt.Println()
//
// 		// Combine parts in all possible ways
//
// 		// fmt.Printf("%s: %s (len: %d)\n", code, parts, len(parts))
// 		// // directional 1
// 		// dir1Solution := solutions(parts, PosInDirectional)
// 		// fmt.Printf("      %s (len: %d)\n", dir1Solution, len(dir1Solution))
// 		// // directional 2
// 		// dir2Solution := solutions(dir1Solution, PosInDirectional)
// 		// fmt.Printf("      %s (len: %d)\n", dir2Solution, len(dir2Solution))
// 		//
// 		// numInCode, err := strconv.Atoi(code[:len(code)-1])
// 		// if err != nil {
// 		// 	panic(err)
// 		// }
// 		// fmt.Printf("Complexity: %d * %d = %d\n\n", len(dir2Solution), numInCode, numInCode*len(dir2Solution))
// 		// complexitiesSum += numInCode * len(dir2Solution)
// 	}
// 	fmt.Printf("Sum of complexities: %d\n", complexitiesSum)
// }
