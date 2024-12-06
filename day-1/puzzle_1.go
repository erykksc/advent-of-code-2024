package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	// Open the text file
	file, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	list1 := make([]int, 0)
	list2 := make([]int, 0)

	scanner := bufio.NewScanner(file)
	lineIdx := 0
	// Scan the text file line by line
	for scanner.Scan() {
		line := scanner.Text()

		// Parse first and second number
		var num1, num2 int
		n, err := fmt.Sscanf(line, "%d   %d", &num1, &num2)
		// If there are not two numbers, there is likely an error in the input
		if err != nil {
			panic(err)
		}
		if n != 2 {
			panic("Invalid input: couldn't parse two integers")
		}

		// Add them to two lists
		list1 = append(list1, num1)
		list2 = append(list2, num2)

		lineIdx += 1
	}

	fmt.Printf("Parsed %d lines\n", lineIdx)

	// Sort two lists
	sort.Ints(list1)
	sort.Ints(list2)

	// Traverse two lists while calculating the distance
	distance := 0
	for i := 0; i < len(list1); i++ {
		d := list1[i] - list2[i]
		if d < 0 {
			d = -d
		}
		distance += d
	}

	// Output the distance
	fmt.Printf("Total distance: %d\n", distance)
}
