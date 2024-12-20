package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func FindTowels(design string, towels []string) ([]string, error) {
	if design == "" {
		return nil, nil
	}

	towelSet := []string{}
	for _, towel := range towels {
		if strings.HasPrefix(design, towel) {
			towelSet = append(towelSet, towel)
			restOfTowelSet, err := FindTowels(design[len(towel):], towels)
			if err != nil {
				towelSet = towelSet[:len(towelSet)-1]
				continue
			}
			// extend towelSet with restOfTowelSet
			towelSet = append(towelSet, restOfTowelSet...)
			return towelSet, nil
		}
	}
	return nil, errors.New("No matching towel found")
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

	possibleDesigns := 0
	for scanner.Scan() {
		design := scanner.Text()
		set, err := FindTowels(design, towels)
		if err != nil {
			fmt.Printf("%s FAIL: cannot be created with available towels\n", design)
		} else {
			fmt.Printf("%s SUCCESS: can be created with following towels: %v\n", design, set)
			possibleDesigns++
		}
	}
	fmt.Printf("Possible designs: %d\n", possibleDesigns)
}
