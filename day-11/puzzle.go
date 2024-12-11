package main

import (
	"flag"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// arang - current arangement of stones
func AfterBlink(arang []int) []int {
	newArang := make([]int, 0, len(arang))

	for _, stone := range arang {
		if stone == 0 {
			newArang = append(newArang, 1)
			continue
		}

		stoneDigits := int(math.Log10(float64(stone))) + 1
		if stoneDigits%2 == 0 {
			stoneDigitsS := strconv.Itoa(stone)
			leftStoneS := stoneDigitsS[:stoneDigits/2]
			rightStoneS := stoneDigitsS[stoneDigits/2:]

			leftStone, err := strconv.Atoi(leftStoneS)
			if err != nil {
				panic(err)
			}

			rightStone, err := strconv.Atoi(rightStoneS)
			if err != nil {
				panic(err)
			}
			newArang = append(newArang, leftStone, rightStone)
			continue
		}

		newArang = append(newArang, stone*2024)
	}
	return newArang
}

func main() {
	iterations := flag.Int("iterations", 25, "Number of iterations to run")
	verbose := flag.Bool("verbose", false, "Output stones after each iteration")
	filename := flag.String("input", "", "Input file with initial arrangement of stones")
	cores := flag.Int("cores", runtime.NumCPU(), "Number of cores to use")
	flag.Parse()

	inputB, err := os.ReadFile(*filename)
	if err != nil {
		panic(err)
	}
	input := strings.Trim(string(inputB), "\n")
	arangStr := strings.Split(input, " ")
	arang := make([]int, len(arangStr))

	for i, str := range arangStr {
		arang[i], err = strconv.Atoi(str)
		if err != nil {
			panic(err)
		}
	}
	log.Printf("Core count used: %d\n", *cores)
	log.Printf("Initial Arrangement:\n%v\n", arang)

	for i := range *iterations {
		// Evaluate chunks
		chunkSize := (len(arang) + *cores) / *cores
		log.Printf("Iteration: %d, chunkSize: %d\n", i+1, chunkSize)

		var wg sync.WaitGroup
		results := make([][]int, *cores)
		reqCap := 0 // Required capacity to have for joined results
		for coreIdx := range *cores {
			start := coreIdx * chunkSize
			end := start + chunkSize
			if end > len(arang) {
				end = len(arang)
			}
			if end < start {
				// Already dealt with the array/ array smaller than core count
				break
			}

			wg.Add(1)
			go func(resultIndex int, arangement []int) {
				defer wg.Done()
				results[resultIndex] = AfterBlink(arangement)
			}(coreIdx, arang[start:end])
			reqCap += end - start
		}
		wg.Wait()
		log.Printf("All chunks evaluated, number of stones: %d\n", reqCap)

		// Join results together
		newArang := make([]int, 0, reqCap)
		for _, result := range results {
			newArang = append(newArang, result...)
		}
		arang = newArang
		log.Printf("Chunk results joined\n")

		if *verbose {
			log.Printf("Arrangement after %d blinks:\n%v\n", i+1, arang)
		}
	}

	log.Printf("Total number of stones after %d iterations: %d\n", *iterations, len(arang))
}
