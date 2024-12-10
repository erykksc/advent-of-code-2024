package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	rules := make(map[[2]int]struct{})

	rulesFilename := os.Args[1]
	updatesFilename := os.Args[2]

	rulesFile, err := os.Open(rulesFilename)
	if err != nil {
		panic(err)
	}
	defer rulesFile.Close()

	updatesFile, err := os.Open(updatesFilename)
	if err != nil {
		panic(err)
	}
	defer updatesFile.Close()

	// Parse rules
	rulesScanner := bufio.NewScanner(rulesFile)
	for rulesScanner.Scan() {
		var num1, num2 int
		_, err := fmt.Sscanf(rulesScanner.Text(), "%d|%d", &num1, &num2)
		if err != nil {
			panic(err)
		}
		rules[[2]int{num1, num2}] = struct{}{}
	}

	correctUpdate := 0
	correctUpdatesMidSum := 0

	updatesScanner := bufio.NewScanner(updatesFile)
	for updatesScanner.Scan() {
		// Parse update as int array
		update := updatesScanner.Text()
		pagesStrs := strings.Split(update, ",")
		pages := make([]int, len(pagesStrs))
		for i, pageStr := range pagesStrs {
			page, err := strconv.Atoi(pageStr)
			if err != nil {
				panic(err)
			}
			pages[i] = page
		}

		valid := true

	UpdateCheckLoop:
		for i, page1 := range pages {
			for j, page2 := range pages {
				if i == j {
					continue
				}

				// Check if page 1 should be before page2
				_, exists := rules[[2]int{page1, page2}]
				// Check if rule broken
				if exists && i > j {
					valid = false
					break UpdateCheckLoop
				}
			}
		}
		if valid {
			correctUpdate++
			correctUpdatesMidSum += pages[len(pages)/2]
		}
	}

	fmt.Printf("Correct updates: %d\n", correctUpdate)
	fmt.Printf("Mid sum of correct updates: %d\n", correctUpdatesMidSum)

}
