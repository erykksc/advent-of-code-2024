package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type TowelSet []string
type TowelSetKey string

func (s *TowelSet) Push(towel ...string) {
	*s = append(*s, towel...)
}

func (s *TowelSet) Pop() string {
	if len(*s) == 0 {
		return ""
	}
	towel := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return towel
}

func (s TowelSet) Key() TowelSetKey {
	return TowelSetKey(strings.Join(s, ","))
}

func Ways2implement(design string, towels []string, implemented map[string]int) (int, error) {
	if ways, exists := implemented[design]; exists {
		return ways, nil
	}
	if design == "" {
		return 1, nil
	}

	allWays := 0
	for _, towel := range towels {
		// fmt.Printf("Checking towel: %s\n", towel)
		if strings.HasPrefix(design, towel) {
			ways, err := Ways2implement(design[len(towel):], towels, implemented)
			if err != nil {
				continue
			}

			allWays += ways
		}
	}

	implemented[design] = allWays

	return allWays, nil
}

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	towels := strings.Split(scanner.Text(), ", ")

	scanner.Scan() // empty line between available towels and required

	implementationWays := 0
	for scanner.Scan() {
		design := scanner.Text()
		fmt.Printf("\nLooking for design: %s with towels: %v\n", design, towels)

		// towel set key -> number of ways to implement the design
		implemented := make(map[string]int)
		allWays, err := Ways2implement(design, towels, implemented)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Ways of implementing the design: %d\n", allWays)
		implementationWays += allWays
	}
	fmt.Printf("Sum of ways of implementing a design: %d\n", implementationWays)
}
