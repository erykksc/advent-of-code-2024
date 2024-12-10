package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseEquations(line string) (int, []int) {
	parts := strings.Split(line, ": ")
	testValue, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}
	numsStr := strings.Split(parts[1], " ")
	nums := make([]int, len(numsStr))
	for i, n := range numsStr {
		nums[i], err = strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
	}
	return testValue, nums
}

func isEquationValid(nums []int, total, goal int, operators []string) bool {
	if total > goal {
		return false
	}
	if len(nums) == 0 {
		return total == goal
	}

	firstNum := nums[0]
	remainingNums := make([]int, len(nums)-1)
	_ = copy(remainingNums, nums[1:])

	if isEquationValid(remainingNums, total+firstNum, goal, append(operators, "+")) {
		return true
	}

	if isEquationValid(remainingNums, total*firstNum, goal, append(operators, "*")) {
		return true
	}

	// Part 2
	totalS := strconv.Itoa(total)
	firstNumS := strconv.Itoa(firstNum)
	concatenatedNum, err := strconv.Atoi(totalS + firstNumS)
	if err != nil {
		panic(err)
	}
	if isEquationValid(remainingNums, concatenatedNum, goal, append(operators, "||")) {
		return true
	}

	return false
}

func main() {
	filename := os.Args[1]

	inputB, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	input := string(inputB)

	_ = input

	lines := strings.Split(input, "\n")

	result := 0
	validEquations := 0
	totalEquations := 0

	for _, line := range lines[:len(lines)-1] {
		// Last line in most input files is empty
		if len(line) == 0 {
			continue
		}
		testValue, nums := parseEquations(line)
		if isEquationValid(nums[1:], nums[0], testValue, []string{}) {
			// fmt.Printf("Valid:\ttestValue: %d, nums: %v\n", testValue, nums)
			result += testValue
			validEquations += 1
		}
		totalEquations += 1
	}

	fmt.Printf("Result: %d\n", result)
	fmt.Printf("Valid equations: %d/%d, %.2f%%\n", validEquations, totalEquations, float64(validEquations)/float64(totalEquations)*100)
}
