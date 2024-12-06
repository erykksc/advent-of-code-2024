package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

func main() {
	filename := os.Args[1]
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	str := string(content)

	pattern := `mul\(([0-9]{1,3}),([0-9]{1,3})\)`
	re := regexp.MustCompile(pattern)

	matches := re.FindAllStringSubmatch(str, -1)

	result := 0
	for _, match := range matches {
		num1, err := strconv.Atoi(match[1])
		if err != nil {
			panic(err)
		}
		num2, err := strconv.Atoi(match[2])
		if err != nil {
			panic(err)
		}

		result += num1 * num2
	}

	fmt.Printf("Result: %d\n", result)
}
