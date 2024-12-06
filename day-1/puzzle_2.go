package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Open the text file
	file, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	list1 := make([]int, 0)
	// Map of number -> number of occurences in list 2
	numsInList2 := make(map[int]int)

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

		numsInList2[num2]++

		lineIdx += 1
	}
	fmt.Printf("Parsed %d lines\n", lineIdx)

	// Traverse two lists while calculating the similarity score
	similiarity_score := 0
	for i := 0; i < len(list1); i++ {
		number := list1[i]
		occurencesInList2 := numsInList2[number]
		similiarity_score += number * occurencesInList2
	}

	// Output the distance
	fmt.Printf("similiarity_score: %d\n", similiarity_score)
}
