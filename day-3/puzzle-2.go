package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func multiply(str string) int {
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
	return result
}

func main() {
	filename := os.Args[1]
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	instruction := string(content)
	totalResult := 0
	currentPos := 0
	do := true
	for currentPos < len(instruction) {
		if do {
			pos := strings.Index(instruction[currentPos:], "don't()")
			// Read until the end of the file
			if pos == -1 {
				totalResult += multiply(instruction[currentPos:])
				break
			}

			totalResult += multiply(instruction[currentPos : currentPos+pos])
			currentPos += pos + len("don't()")
			do = false
		} else {
			pos := strings.Index(instruction[currentPos:], "do()")
			// Read until the end of the file, finish interpreting the string
			if pos == -1 {
				break
			}

			currentPos += pos + len("do()")
			do = true
		}
	}

	fmt.Printf("Result: %d\n", totalResult)
}
