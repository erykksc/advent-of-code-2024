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

	invalidCount := 0
	midSum := 0

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

		originallyValid := true

		for true {
			valid := true
		UpdateCheckLoop:
			for i, pageI := range pages {
				for j, pageJ := range pages {
					if i == j {
						continue
					}

					// Check if page 1 should be before page2
					_, exists := rules[[2]int{pageI, pageJ}]
					// Check if rule broken
					if exists && i > j {
						originallyValid = false
						valid = false
						pages[i] = pageJ
						pages[j] = pageI
						break UpdateCheckLoop
					}
				}
			}
			// Should only calculate the mid sum of updates that were originally invalid
			if originallyValid {
				break
			} else if valid {
				// Finally found a valid ordering invalidCount++
				invalidCount++
				midSum += pages[len(pages)/2]
				break
			}
		}
	}

	fmt.Printf("Originally invalid updates: %d\n", invalidCount)
	fmt.Printf("Mid sum of correct updates: %d\n", midSum)

}
