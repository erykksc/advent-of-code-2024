package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Operation int

const (
	AND Operation = iota
	OR
	XOR
)

type GateOperation struct {
	input1, input2, output string
	operaton               Operation
}

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	values := make(map[string]bool)
	// Read initial cables output
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			// break between initial values and logic gates
			break
		}
		line := scanner.Text()
		name, output, found := strings.Cut(line, ": ")
		if !found {
			panic("Invalid input")
		}

		value, err := strconv.ParseBool(output)
		if err != nil {
			panic(err)
		}

		values[name] = value
	}

	gates := make([]GateOperation, 0)

	// Read logic gates
	for scanner.Scan() {
		line := scanner.Text()

		words := strings.Split(line, " ")

		input1, operationS, input2, output := words[0], words[1], words[2], words[4]

		var operation Operation
		switch operationS {
		case "AND":
			operation = AND
		case "OR":
			operation = OR
		case "XOR":
			operation = XOR
		default:
			panic("Invalid operation: " + operationS)

		}

		gOperation := GateOperation{input1, input2, output, operation}
		gates = append(gates, gOperation)
	}

	for len(gates) > 0 {
		for gateIdx, gate := range gates {
			input1, input2 := gate.input1, gate.input2

			input1V, input1Exists := values[input1]
			input2V, input2Exists := values[input2]

			if !input1Exists || !input2Exists {
				continue
			}

			var outputV bool
			switch gate.operaton {
			case AND:
				outputV = input1V && input2V
			case OR:
				outputV = input1V || input2V
			case XOR:
				outputV = input1V != input2V
			}

			fmt.Printf("Found out value of %s: %v\n", gate.output, outputV)
			values[gate.output] = outputV

			gates = append(gates[:gateIdx], gates[gateIdx+1:]...)
			break
		}
	}

	// Calculate sum of all values starting with 'z'
	var bitRepresentation string
	for i := range 100 {
		// znanme
		name := fmt.Sprintf("z%02d", i)
		output, exists := values[name]
		if !exists {
			// assumeif one is missing, there are no more z numbers
			fmt.Printf("Highest z is: %s\n", name)
			break
		}
		if output {
			bitRepresentation = "1" + bitRepresentation
		} else {
			bitRepresentation = "0" + bitRepresentation
		}
	}

	fmt.Printf("Sum of all values starting with 'z' is: %s\n", bitRepresentation)
	base10, err := strconv.ParseInt(bitRepresentation, 2, 64)
	if err != nil {
		panic(err)
	}
	fmt.Printf("In base 10 it is: %d\n", base10)
}
