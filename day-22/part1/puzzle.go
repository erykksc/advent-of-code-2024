package main

import (
	"flag"
	"fmt"
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

func main() {
	iters := flag.Int("iters", 2000, "number of iterations")
	flag.Parse()

	inputB, err := os.ReadFile(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	input := strings.TrimSpace(string(inputB))
	lines := strings.Split(input, "\n")

	sum := 0
	for _, num := range lines {
		secretNum, err := strconv.Atoi(num)
		if err != nil {
			panic(err)
		}
		fmt.Printf("\nInitial secret number: %d\n", secretNum)
		for range *iters {
			secretNum = nextSecretNum(secretNum)
			// fmt.Printf("Secret number after %d iterations: %d\n", iter, secretNum)
		}
		fmt.Printf("Secret number after %d iterations: %d\n", *iters, secretNum)
		sum += secretNum
	}
	fmt.Printf("\nSum of secret numbers: %d\n", sum)
}
