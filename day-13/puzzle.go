package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

type Vector2D struct {
	X, Y int
}

// Add two vectors
func (v Vector2D) Add(v2 Vector2D) Vector2D {
	return Vector2D{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
	}
}

// Subtract two vectors
func (v Vector2D) Sub(v2 Vector2D) Vector2D {
	return Vector2D{
		X: v.X - v2.X,
		Y: v.Y - v2.Y,
	}
}

// Multiply vector by a scalar
func (v Vector2D) Multiply(scalar int) Vector2D {
	return Vector2D{
		X: v.X * scalar,
		Y: v.Y * scalar,
	}
}

func (v1 Vector2D) Equals(v2 Vector2D) bool {
	return v1.X == v2.X && v1.Y == v2.Y
}

func IsPastPrize(claw, prize Vector2D) bool {
	if claw.X > prize.X {
		return true
	}

	if claw.Y > prize.Y {
		return true
	}

	return false
}

type Config struct {
	Buttons []Vector2D
	Prize   Vector2D
}

func parseConfig(configStr string) (*Config, error) {
	lines := strings.Split(configStr, "\n")
	if len(lines) != 4 {
		return nil, fmt.Errorf("invalid config, need exactly 3 lines got %d", len(lines))
	}

	config := Config{
		Buttons: make([]Vector2D, 2),
	}

	// parse Button A
	fmt.Sscanf(lines[0], "Button A: X+%d, Y+%d", &config.Buttons[0].X, &config.Buttons[0].Y)
	// parse Button B
	fmt.Sscanf(lines[1], "Button B: X+%d, Y+%d", &config.Buttons[1].X, &config.Buttons[1].Y)
	// parse prize
	fmt.Sscanf(lines[2], "Prize: X=%d, Y=%d", &config.Prize.X, &config.Prize.Y)

	return &config, nil
}

func solveConfig(c Config) (tokens int) {
	buttonA, buttonB, prize := c.Buttons[0], c.Buttons[1], c.Prize

	determinant := buttonA.X*buttonB.Y - buttonB.X*buttonA.Y

	if determinant == 0 {
		// infinite solutions or none
		// Assume none
		return 0
	}

	aPresses := (prize.X*buttonB.Y - buttonB.X*prize.Y) / determinant
	bPresses := (buttonA.X*prize.Y - prize.X*buttonA.Y) / determinant

	calcPosition := func(aPresses, bPresses int) Vector2D {
		return buttonB.Multiply(bPresses).Add(buttonA.Multiply(aPresses))
	}

	diff := prize.Sub(calcPosition(aPresses, bPresses))
	fmt.Printf("Diff: %+v\n", diff)
	if diff.X == 0 && diff.Y == 0 {
		return 3*aPresses + bPresses
	} else {
		return 0
	}
}

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	var configS strings.Builder
	linesCollected := 0
	configs := make([]Config, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		configS.WriteString(line)
		configS.WriteRune('\n')
		linesCollected++

		if linesCollected == 3 {
			newConfig, err := parseConfig(configS.String())
			if err != nil {
				panic(err)
			}
			configs = append(configs, *newConfig)
			configS.Reset()
			linesCollected = 0
		}
	}

	fmt.Println("Solution for PART 1:")
	totalTokens := 0
	for _, config := range configs {
		fmt.Printf("Handling config: %+v\n", config)
		buttonA, buttonB, prize := config.Buttons[0], config.Buttons[1], config.Prize

		cheapestPrice := math.MaxInt
		for aPresses := range 100 {
			for bPresses := range 100 {
				endPosition := buttonA.Multiply(aPresses).Add(buttonB.Multiply(bPresses))

				// Prize not reached
				if !endPosition.Equals(prize) {
					continue
				}

				price := aPresses*3 + bPresses
				if price < cheapestPrice {
					cheapestPrice = price
				}
			}
		}
		if cheapestPrice == math.MaxInt {
			fmt.Println("No solution found")
			continue
		}
		fmt.Printf("Cheapest price: %d\n", cheapestPrice)
		totalTokens += cheapestPrice
	}
	fmt.Printf("Total tokens: %d\n", totalTokens)

	fmt.Println("Solution for PART 2:")
	totalTokensPart2 := 0
	for _, config := range configs {
		fmt.Printf("Handling config: %+v\n", config)

		// Part 2
		config.Prize.X += 10000000000000
		config.Prize.Y += 10000000000000
		cheapestPrice := solveConfig(config)
		fmt.Printf("Cheapest price for part 2: %d\n", cheapestPrice)
		totalTokensPart2 += cheapestPrice
	}
	fmt.Printf("Total tokens for part 2: %d\n", totalTokensPart2)
}
