package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
)

func formatBinary(n int64) string {
	// Convert to binary string
	bin := fmt.Sprintf("%b", n)

	// Pad with zeros to make length multiple of 3
	padding := (3 - len(bin)%3) % 3
	bin = strings.Repeat("0", padding) + bin

	// Split into groups of 3
	var result strings.Builder
	for i := 0; i < len(bin); i += 3 {
		if i > 0 {
			result.WriteString(" ")
		}
		result.WriteString(bin[i : i+3])
	}
	return result.String()
}

type Computer struct {
	A, B, C *big.Int
	iPtr    int
	Program []int
	Output  []int
}

func NewComputer(A, B, C *big.Int, program []int) *Computer {
	return &Computer{
		A:       A,
		B:       B,
		C:       C,
		iPtr:    0,
		Program: program,
		Output:  []int{},
	}
}

func (c *Computer) Combo(operand int) *big.Int {
	switch operand {
	case 0:
		return big.NewInt(0)
	case 1:
		return big.NewInt(1)
	case 2:
		return big.NewInt(2)
	case 3:
		return big.NewInt(3)
	case 4:
		return c.A
	case 5:
		return c.B
	case 6:
		return c.C
	case 7:
		panic("Reserved operand 7, shouldn't appear in valid program")
	}
	panic(fmt.Sprintf("Unknown operand: %d", operand))
}

func (c Computer) NextInstruction() bool {
	return c.iPtr < len(c.Program)
}

type Instruction func(operand int)

func (c *Computer) adv(operand int) {
	var denumerator big.Int
	denumerator.Exp(big.NewInt(2), c.Combo(operand), nil)

	c.A.Div(c.A, &denumerator)
}

// bitwise XOR
func (c *Computer) bxl(operand int) {
	c.B.Xor(c.B, big.NewInt(int64(operand)))
}

// combo operand modulo 8
func (c *Computer) bst(operand int) {
	c.B.Mod(c.Combo(operand), big.NewInt(8))
}

func (c *Computer) jnz(operand int) {
	if c.A.Cmp(big.NewInt(0)) == 0 {
		c.iPtr += 2
		return
	}

	c.iPtr = operand
}

// bitwise XOR of register B and C
func (c *Computer) bxc(_ int) {
	c.B.Xor(c.B, c.C)
}

func (c *Computer) out(operand int) {
	var result big.Int
	result.Mod(c.Combo(operand), big.NewInt(8))
	normal := int(result.Int64())
	c.Output = append(c.Output, normal)
}

func (c *Computer) bdv(operand int) {
	var denumerator big.Int
	denumerator.Exp(big.NewInt(2), c.Combo(operand), nil)

	c.B.Div(c.A, &denumerator)
}

func (c *Computer) cdv(operand int) {
	var denumerator big.Int
	denumerator.Exp(big.NewInt(2), c.Combo(operand), nil)

	c.C.Div(c.A, &denumerator)
}

func RunComputer(regA, regB, regC *big.Int, program []int) []int {
	c := NewComputer(regA, regB, regC, program)
	for c.NextInstruction() {
		opCode := c.Program[c.iPtr]
		operand := c.Program[c.iPtr+1]

		switch opCode {
		case 0:
			c.adv(operand)
		case 1:
			c.bxl(operand)
		case 2:
			c.bst(operand)
		case 3:
			c.jnz(operand)
		case 4:
			c.bxc(operand)
		case 5:
			c.out(operand)
		case 6:
			c.bdv(operand)
		case 7:
			c.cdv(operand)
		default:
			panic(fmt.Sprintf("Unknown opcode: %d", opCode))
		}

		// if not jumping move to next instruction
		if opCode != 3 {
			c.iPtr += 2
		}

	}
	return c.Output
}

// intial call should set programIdx to len(program)-1
// if return value is -1, couldn't find a match
func findMatch(currentRegA, programIdx int, program []int) int {
	// fmt.Printf("findMatch(%d, %d, %v)\n", currentRegA, programIdx, program)
	if programIdx < 0 {
		panic("programIdx should be >= 0")
	}
	for i := range 8 {
		regA := currentRegA + (i << (programIdx * 3))
		output := RunComputer(big.NewInt(int64(regA)), big.NewInt(0), big.NewInt(0), program)
		// fmt.Printf("tested %d, %s, output %v\n", regA, formatBinary(int64(regA)), output)
		if len(output) < len(program) {
			continue
		}
		// time.Sleep(500 * time.Millisecond)
		// fmt.Printf("%v, %s\r", output, formatBinary(int64(regA)))
		if output[programIdx] == program[programIdx] {
			fmt.Printf("Found match %d == %d for programIdx %d with %s\n", output[programIdx], program[programIdx], programIdx, formatBinary(int64(regA)))
			if programIdx == 0 {
				return regA
			}
			// find the next idx
			match := findMatch(regA, programIdx-1, program)
			if match != -1 {
				return match
			}
			fmt.Printf("Backtracking in programIdx %d\n", programIdx)
		}
	}
	return -1
}

func main() {
	flag.Parse()

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	// Scan register values
	var regB, regC int64
	scanner.Scan() // scan regA
	scanner.Scan()
	_, err = fmt.Sscanf(scanner.Text(), "Register B: %d", &regB)
	if err != nil {
		panic(err)
	}
	scanner.Scan()
	_, err = fmt.Sscanf(scanner.Text(), "Register C: %d", &regC)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Registers before program:\nA:%d, B:%d, C:%d\n", 0, regB, regC)

	// Scan an empty line in between
	scanner.Scan()

	// Scan program
	scanner.Scan()
	numsStrs := strings.TrimPrefix(scanner.Text(), "Program: ")
	programStrs := strings.Split(numsStrs, ",")
	program := make([]int, len(programStrs))
	for i, s := range programStrs {
		num, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		program[i] = num
	}

	fmt.Printf("Program (len %d): \n%v\n", len(program), program)

	foundRegA := findMatch(0, len(program)-1, program)
	if foundRegA == -1 {
		panic("Couldn't find a match")
	}
	fmt.Printf("RegA value that makes the output match program: %d\n", foundRegA)
	output := RunComputer(big.NewInt(int64(foundRegA)), big.NewInt(0), big.NewInt(0), program)
	fmt.Printf("Output of foundRegA and program: \n%v\n%v\n", output, program)
}
