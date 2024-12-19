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

func main() {
	flag.Parse()

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	// Scan register values
	var regA, regB, regC int64
	scanner.Scan()
	_, err = fmt.Sscanf(scanner.Text(), "Register A: %d", &regA)
	if err != nil {
		panic(err)
	}
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
	fmt.Printf("Registers before program:\nA:%d, B:%d, C:%d\n", regA, regB, regC)

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

	fmt.Printf("Program: %v\n", program)

	c := NewComputer(big.NewInt(regA), big.NewInt(regB), big.NewInt(regC), program)
	for c.NextInstruction() {
		opCode := c.Program[c.iPtr]
		operand := c.Program[c.iPtr+1]
		fmt.Printf("After opcode %d operand %d combo %d: ", opCode, operand, c.Combo(operand))

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

		fmt.Printf("A:%d, B:%d, C:%d, Output:%v\n", c.A, c.B, c.C, c.Output)
	}

	outputStrs := make([]string, len(c.Output))
	for i, num := range c.Output {
		outputStrs[i] = strconv.Itoa(num)
	}
	fmt.Printf("Output:\n%s\n", strings.Join(outputStrs, ","))
}
