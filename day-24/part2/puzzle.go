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

func (oper Operation) String() string {
	switch oper {
	case AND:
		return "AND"
	case OR:
		return "OR"
	case XOR:
		return "XOR"
	}
	panic("Non existant operator")
}

type GateOperation struct {
	input1, input2, output string
	operator               Operation
}

func (gOper GateOperation) String() string {
	return fmt.Sprintf("%s %s %s -> %s", gOper.input1, gOper.operator, gOper.input2, gOper.output)
}

func BitRepresentation(prefix string, values map[string]bool) string {
	// Calculate sum of all values starting with 'z'
	var bitRepresentation string
	for i := range 100 {
		// znanme
		name := fmt.Sprintf(prefix+"%02d", i)
		output, exists := values[name]
		if !exists {
			// assumeif one is missing, there are no more z/x/y (based on prefix) numbers
			break
		}
		if output {
			bitRepresentation = "1" + bitRepresentation
		} else {
			bitRepresentation = "0" + bitRepresentation
		}
	}
	return bitRepresentation
}

func Assosiated(level int, gates []GateOperation) []GateOperation {
	assosiated := []GateOperation{}

	isNameAssosiated := func(name string, level int) bool {
		if name[0] != 'x' && name[0] != 'y' && name[0] != 'z' {
			return false
		}

		number, err := strconv.Atoi(name[1:])
		if err != nil {
			panic(err)
		}
		if number == level {
			return true
		}
		return false

	}

	for _, gate := range gates {
		if isNameAssosiated(gate.input1, level) || isNameAssosiated(gate.input2, level) || isNameAssosiated(gate.output, level) {
			assosiated = append(assosiated, gate)
		}
	}
	return assosiated
}

// Find dependencie sup to the next z
func Dependencies(name string, gates []GateOperation, initialLayer int) []string {
	// if name[0] == 'x' || name[0] == 'y' || name[0] == 'z' {
	// 	layer, err := strconv.Atoi(name[1:])
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//
	// 	if layer < initialLayer {
	// 		return []string{}
	// 	}
	// }
	//
	if name[0] == 'x' || name[0] == 'y' {
		return []string{name}
	}

	dependencies := []string{}

	for _, gate := range gates {
		if gate.output == name {
			fmt.Printf("%s %s %s -> %s\n", gate.input1, gate.operator, gate.input2, gate.output)
			dependencies = append(dependencies, name)
			dependencies = append(dependencies, Dependencies(gate.input1, gates, initialLayer)...)
			dependencies = append(dependencies, Dependencies(gate.input2, gates, initialLayer)...)
			break
		}
	}

	// remove duplicates
	uniqueDependencies := make(map[string]bool)
	for _, dep := range dependencies {
		uniqueDependencies[dep] = true
	}

	out := make([]string, 0, len(uniqueDependencies))

	for dep := range uniqueDependencies {
		out = append(out, dep)
	}

	return out
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

		if input1[0] == 'y' && input2[0] == 'x' {
			tmp := input1
			input1 = input2
			input2 = tmp
		}

		gOperation := GateOperation{input1, input2, output, operation}
		gates = append(gates, gOperation)
	}

	// Find out all the dependencies of the last z
	lastZnum := 0
	for _, gate := range gates {
		if gate.output[0] == 'z' {
			number, err := strconv.Atoi(gate.output[1:])
			if err != nil {
				panic(err)
			}

			if number > lastZnum {
				lastZnum = number
			}
		}
	}
	fmt.Printf("Last z is: z%02d\n", lastZnum)

	uncalculatedGates := make([]GateOperation, len(gates))
	copy(uncalculatedGates, gates)

	// calculate the output values of all gates
	for len(uncalculatedGates) > 0 {
		for gateIdx, gate := range uncalculatedGates {
			input1, input2 := gate.input1, gate.input2

			input1V, input1Exists := values[input1]
			input2V, input2Exists := values[input2]

			if !input1Exists || !input2Exists {
				continue
			}

			var outputV bool
			switch gate.operator {
			case AND:
				outputV = input1V && input2V
			case OR:
				outputV = input1V || input2V
			case XOR:
				outputV = input1V != input2V
			}

			values[gate.output] = outputV

			uncalculatedGates = append(uncalculatedGates[:gateIdx], uncalculatedGates[gateIdx+1:]...)
			break
		}
	}
	xBase2 := BitRepresentation("x", values)
	xBase10, err := strconv.ParseInt(xBase2, 2, 64)
	if err != nil {
		panic(err)
	}
	fmt.Printf("'x' is: %s, %d\n", xBase2, xBase10)

	yBase2 := BitRepresentation("y", values)
	yBase10, err := strconv.ParseInt(yBase2, 2, 64)
	if err != nil {
		panic(err)
	}
	fmt.Printf("'y' is: %s, %d\n", yBase2, yBase10)

	zBase2 := BitRepresentation("z", values)
	zBase10, err := strconv.ParseInt(zBase2, 2, 64)
	if err != nil {
		panic(err)
	}
	fmt.Printf("'z' is: %s, %d\n", zBase2, zBase10)

	sumBase10 := xBase10 + yBase10
	sumBase2 := strconv.FormatInt(sumBase10, 2)
	fmt.Printf("Sum is: %s, %d\n", sumBase2, sumBase10)
	var diff strings.Builder
	for i := range len(sumBase2) {
		if sumBase2[i] == zBase2[i] {
			diff.WriteRune(' ')
		} else {
			diff.WriteRune('X')
		}
	}
	fmt.Printf("Dif is: %s, %d\n", diff.String(), zBase10-sumBase10)

	// bit index -> name
	Cout := make(map[int]string)
	//find Cout[0]
	for _, gate := range gates {
		if gate.operator == AND && gate.input1 == "x00" && gate.input2 == "y00" {
			Cout[0] = gate.output
			break
		}
	}
	for num := 1; num < lastZnum; num++ {
		xName := fmt.Sprintf("x%02d", num)
		yName := fmt.Sprintf("y%02d", num)
		// find A XOR B -> half S
		fmt.Printf("Looking for %s XOR %s -> halfS\n", xName, yName)
		var halfS string
		for _, gate := range gates {
			// find A and B
			if gate.operator != XOR || gate.input1 != xName || gate.input2 != yName {
				continue
			}

			halfS = gate.output
			break
		}
		if len(halfS) == 0 {
			panic(fmt.Sprintf("Could not find halfS for %d", num))
		}
		fmt.Printf("Found halfS: %s\n", halfS)

		fmt.Printf("Looking for %s XOR %s -> fullS\n", halfS, Cout[num-1])
		var fullS string
		// find halfS XOR Cout[num-1] -> fullS
		for _, gate := range gates {
			if gate.operator != XOR {
				continue
			}
			if !((gate.input1 == halfS && gate.input2 == Cout[num-1]) || (gate.input1 == Cout[num-1] && gate.input2 == halfS)) {
				continue
			}
			// found S out
			fullS = gate.output
			break
		}
		if len(fullS) == 0 {
			panic(fmt.Sprintf("Could not find fullS for %d", num))
		}
		fmt.Printf("Found fullS: %s\n", fullS)

		// find A AND B
		fmt.Printf("Looking for %s AND %s -> AandB\n", xName, yName)
		var AandB string
		for _, gate := range gates {
			if gate.operator != AND || gate.input1 != xName || gate.input2 != yName {
				continue
			}

			AandB = gate.output
			break
		}
		if len(AandB) == 0 {
			panic(fmt.Sprintf("Could not find AandB for %d", num))
		}
		fmt.Printf("Found AandB: %s\n", AandB)

		// find halfS AND Cout[num-1]
		fmt.Printf("Looking for %s AND %s -> halfSAndCout\n", halfS, Cout[num-1])
		var halfSAndCout string
		for _, gate := range gates {
			if gate.operator != AND {
				continue
			}
			if !((gate.input1 == halfS && gate.input2 == Cout[num-1]) ||
				(gate.input1 == Cout[num-1] && gate.input2 == halfS)) {
				continue
			}

			halfSAndCout = gate.output
			break
		}
		if len(halfSAndCout) == 0 {
			panic(fmt.Sprintf("Could not find halfSAndCout for %d", num))
		}
		fmt.Printf("Found halfSAndCout: %s\n", halfSAndCout)

		// find AandB OR halfSAndCout -> Cout
		fmt.Printf("Looking for %s OR %s -> Cout\n", AandB, halfSAndCout)
		for _, gate := range gates {
			if gate.operator != OR {
				continue
			}
			if !((gate.input1 == AandB && gate.input2 == halfSAndCout) ||
				(gate.input1 == halfSAndCout && gate.input2 == AandB)) {
				continue
			}

			Cout[num] = gate.output
			break
		}
		if len(Cout[num]) == 0 {
			panic(fmt.Sprintf("Could not find Cout for %d", num))
		}
		fmt.Printf("Found Cout: %s\n", Cout[num])
	}
}
