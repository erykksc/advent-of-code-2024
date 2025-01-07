package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Grid [][]rune

func main() {
	flag.Parse()
	inputB, err := os.ReadFile(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	input := strings.TrimSpace(string(inputB))
	keysLocksStrs := strings.Split(input, "\n\n")

	keys := make([][5]int, 0)
	locks := make([][5]int, 0)

	for _, s := range keysLocksStrs {
		lines := strings.Split(s, "\n")

		var out [5]int
		isKey := false
		if lines[0][0] == '.' {
			isKey = true
			for depth, line := range lines[:6] {
				for column := range 5 {
					if line[column] == '.' {
						continue
					}
					if out[column] > 0 {
						continue
					}

					out[column] = 6 - depth
				}
			}
			keys = append(keys, out)
		} else {
			// is lock
			for height, line := range lines[1:] {
				for column := range 5 {
					if line[column] == '.' {
						continue
					}
					out[column] = height + 1
				}
			}
			locks = append(locks, out)
		}
		fmt.Printf("%s\n", s)
		if isKey {
			fmt.Printf("Key: ")
		} else {
			fmt.Printf("Lock: ")
		}
		fmt.Printf("%v\n", out)
	}

	fmt.Printf("Total keys: %d\n", len(keys))
	fmt.Printf("Total locks: %d\n", len(locks))

	pairs := 0
	for _, key := range keys {
		for _, lock := range locks {
			fits := true
			for column := range 5 {
				if 5 < key[column]+lock[column] {
					fits = false
					break
				}
			}
			if fits {
				fmt.Printf("Found pair key-lock: %v %v\n", key, lock)
				pairs++
			}
		}
	}
	fmt.Printf("Unique pairs: %d\n", pairs)
}
