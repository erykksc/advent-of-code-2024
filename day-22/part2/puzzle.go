package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// To mix a value into the secret number, calculate the bitwise XOR of the given value and the secret number.
// Then, the secret number becomes the result of that operation.
// (If the secret number is 42 and you were to mix 15 into the secret number, the secret number would become 37.)
func mix(secretValue *int, n int) {
	*secretValue = *secretValue ^ n
}

// To prune the secret number, calculate the value of the secret number modulo 16777216.
// Then, the secret number becomes the result of that operation.
// (If the secret number is 100000000 and you were to prune the secret number, the secret number would become 16113920.)
func prune(secretNum *int) {
	*secretNum = *secretNum % 16777216
}

func nextSecretNum(secretNum int) int {
	// Calculate the result of multiplying the secret number by 64.
	// Then, mix this result into the secret number.
	// Finally, prune the secret number.
	mix(&secretNum, secretNum*64)
	prune(&secretNum)

	// step 2
	// Calculate the result of dividing the secret number by 32.
	// Round the result down to the nearest integer.
	// Then, mix this result into the secret number.
	// Finally, prune the secret number.
	mix(&secretNum, secretNum/32)
	prune(&secretNum)

	// step 3
	// Calculate the result of multiplying the secret number by 2048.
	// Then, mix this result into the secret number.
	// Finally, prune the secret number.
	mix(&secretNum, secretNum*2048)
	prune(&secretNum)

	return secretNum
}

func buyBananas(sequence [4]int, bPrices, bPriceChanges []int) int {
	// Look for first appearence of a sequence
	for i := 3; i < len(bPriceChanges); i++ {
		var rangeOf4 [4]int
		for j := range 4 {
			rangeOf4[3-j] = bPriceChanges[i-j]
		}
		if sequence == rangeOf4 {
			return bPrices[i+1]
		}
	}
	// sequence not found
	return 0
}

func main() {
	iters := flag.Int("iters", 2000, "number of iterations")
	flag.Parse()

	inputB, err := os.ReadFile(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	input := strings.TrimSpace(string(inputB))
	lines := strings.Split(input, "\n")

	buyersCount := len(lines)

	// for each buyer, generate secretNums
	secretNums := make([][]int, buyersCount)
	for buyerIdx, initialSecretNum := range lines {
		secretNum, err := strconv.Atoi(initialSecretNum)
		if err != nil {
			panic(err)
		}
		secretNums[buyerIdx] = append(secretNums[buyerIdx], secretNum)

		for range *iters {
			secretNum = nextSecretNum(secretNum)
			secretNums[buyerIdx] = append(secretNums[buyerIdx], secretNum)
		}
	}

	// Calculate prices
	prices := make([][]int, buyersCount)
	for buyerIdx, buyersSecretNums := range secretNums {
		for _, secretNum := range buyersSecretNums {
			price := secretNum % 10
			prices[buyerIdx] = append(prices[buyerIdx], price)
		}
	}

	// Calculate price changes
	// NOTE: priceChanges is shifted by one to the left compared to prices
	// as prices[0] has no belonging change
	priceChanges := make([][]int, buyersCount)
	for buyersIdx, bPrices := range prices {
		for i := 1; i < len(bPrices); i++ {
			priceChanges[buyersIdx] = append(priceChanges[buyersIdx], bPrices[i]-bPrices[i-1])
		}
	}

	// For each buyer calculate by going through the sequences, a map of sequence -> bananas
	// Remember that only first sequence enountered matters
	seqs2bananas := make([]map[[4]int]int, buyersCount)
	for buyersIdx := range len(prices) {
		bPrices := prices[buyersIdx]
		bPriceChanges := priceChanges[buyersIdx]
		bSeq2Bananas := make(map[[4]int]int)

		for i := 3; i < len(bPriceChanges); i++ {
			var seq [4]int
			for j := range 4 {
				seq[3-j] = bPriceChanges[i-j]
			}

			// Try to map first accurance of sequence with appropriate bananas
			if _, exists := bSeq2Bananas[seq]; !exists {
				bSeq2Bananas[seq] = bPrices[i+1]
			}
		}
		seqs2bananas[buyersIdx] = bSeq2Bananas
	}

	// Iterate through keys of each array to find a sequence that maximizes the gain among all buyers
	// Create a slice of all available keys
	seq2total := make(map[[4]int]int)
	maxTotal := math.MinInt
	var maxSeq [4]int
	for buyersIdx := range buyersCount {
		for seq, bananas := range seqs2bananas[buyersIdx] {
			seq2total[seq] += bananas
			if seq2total[seq] > maxTotal {
				maxTotal = seq2total[seq]
				maxSeq = seq
			}
		}
	}

	boughtBananas := 0
	for buyersIdx := range buyersCount {
		bPrices := prices[buyersIdx]
		bPriceChanges := priceChanges[buyersIdx]
		bananas := buyBananas(maxSeq, bPrices, bPriceChanges)
		boughtBananas += bananas
		fmt.Printf("Buyer %d bought %d bananas\n", buyersIdx, bananas)
	}

	fmt.Printf("Max number of bananas to gain %d, with sequence %v\n", maxTotal, maxSeq)
	fmt.Printf("Trying this sequences yielded %d bananas\n", boughtBananas)
}
