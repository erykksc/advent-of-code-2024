package main

import (
	"flag"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

// arang - current arangement of stones
func AfterBlink(freq map[int]int) (newFreq map[int]int, total int) {
	newFreq = make(map[int]int)
	total = 0
	for stone, occurences := range freq {
		total += occurences // every stone is at least replaced by one
		if stone == 0 {
			newFreq[1] += occurences
		} else if stoneDigits := int(math.Log10(float64(stone))) + 1; stoneDigits%2 == 0 {
			stoneDigitsS := strconv.Itoa(stone)
			leftStoneS := stoneDigitsS[:stoneDigits/2]
			rightStoneS := stoneDigitsS[stoneDigits/2:]

			leftStone, err := strconv.Atoi(leftStoneS)
			if err != nil {
				panic(err)
			}
			newFreq[leftStone] += occurences

			rightStone, err := strconv.Atoi(rightStoneS)
			if err != nil {
				panic(err)
			}
			newFreq[rightStone] += occurences

			total += occurences // as there are two numbers added
		} else {
			newFreq[stone*2024] += occurences
		}
	}

	return newFreq, total
}

func main() {
	iterations := flag.Int("iterations", 25, "Number of iterations to run")
	verbose := flag.Bool("verbose", false, "Output stones after each iteration")
	filename := flag.String("input", "", "Input file with initial arrangement of stones")
	flag.Parse()

	inputB, err := os.ReadFile(*filename)
	if err != nil {
		panic(err)
	}
	input := strings.Trim(string(inputB), "\n")
	arangStr := strings.Split(input, " ")

	// Create initial frequencies
	freq := make(map[int]int)
	for _, str := range arangStr {
		stone, err := strconv.Atoi(str)
		if err != nil {
			panic(err)
		}
		freq[stone]++

	}
	log.Printf("Initial Arrangement:\n%v\n", freq)

	totalStones := 0
	for i := range *iterations {
		freq, totalStones = AfterBlink(freq)
		if *verbose {
			log.Printf("Arrangement after %d iterations:\n%v\n", i, freq)
		}
	}

	log.Printf("Total number of stones after %d iterations: %d\n", *iterations, totalStones)
}
