package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	fmt.Println("flowchart TD")
	gateNum := 0
	for scanner.Scan() {
		line := scanner.Text()

		words := strings.Split(line, " ")

		input1, gate, input2, output := words[0], words[1], words[2], words[4]

		fmt.Printf("\t%s --> gate%d{%s}\n", input1, gateNum, gate)
		fmt.Printf("\t%s --> gate%d\n", input2, gateNum)
		fmt.Printf("\tgate%d --> %s\n\n", gateNum, output)

		gateNum++
	}
}
