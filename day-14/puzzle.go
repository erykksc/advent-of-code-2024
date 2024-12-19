package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Robot struct {
	px, py int
	vx, vy int
}

const WIDTH, HEIGHT = 101, 103
const ITERATIONS = 10000

func PerformStep(robots []Robot) {
	for i := range robots {
		robots[i].px += robots[i].vx
		robots[i].px %= WIDTH
		if robots[i].px < 0 {
			robots[i].px = WIDTH + robots[i].px
		}

		robots[i].py += robots[i].vy % HEIGHT
		robots[i].py %= HEIGHT
		if robots[i].py < 0 {
			robots[i].py = HEIGHT + robots[i].py
		}
	}
}

func RobotPositions(robots []Robot) map[[2]int]int {
	// how many robots at position
	positions := make(map[[2]int]int)
	for _, r := range robots {
		pos := [2]int{r.px, r.py}
		positions[pos] += 1
	}
	return positions
}

func RepresentRobots(robots []Robot) string {
	positions := RobotPositions(robots)

	var b strings.Builder
	for y := range HEIGHT {
		for x := range WIDTH {
			robotsAtPos := positions[[2]int{x, y}]
			if robotsAtPos == 0 {
				b.WriteRune('.')
			} else {
				b.WriteString(strconv.Itoa(robotsAtPos))
			}
		}
		b.WriteRune('\n')
	}
	return b.String()
}

func CalcDangerLvl(robots []Robot) int {
	// In quadrants
	midHeight := (HEIGHT / 2)
	midWidth := (WIDTH / 2)
	// fmt.Printf("Mid height: %d, Mid width %d\n", midHeight, midWidth)

	positions := RobotPositions(robots)
	dangerLvl := 0

	// First quadrant
	for y := 0; y < midHeight; y++ {
		for x := 0; x < midWidth; x++ {
			count := positions[[2]int{x, y}]
			dangerLvl += count
		}
	}
	// fmt.Printf("Robots in first quadrant: %d\n", dangerLvl)

	// Second quadrant
	inQuadrant := 0
	for y := 0; y < midHeight; y++ {
		for x := midWidth + 1; x < WIDTH; x++ {
			count := positions[[2]int{x, y}]
			inQuadrant += count
		}
	}
	dangerLvl *= inQuadrant
	// fmt.Printf("Robots in second quadrant: %d\n", inQuadrant)

	// Third quadrant
	inQuadrant = 0
	for y := midHeight + 1; y < HEIGHT; y++ {
		for x := 0; x < midWidth; x++ {
			count := positions[[2]int{x, y}]
			inQuadrant += count
		}
	}
	dangerLvl *= inQuadrant
	// fmt.Printf("Robots in third quadrant: %d\n", inQuadrant)

	// Fourth quadrant
	inQuadrant = 0
	for y := midHeight + 1; y < HEIGHT; y++ {
		for x := midWidth + 1; x < WIDTH; x++ {
			count := positions[[2]int{x, y}]
			inQuadrant += count
		}
	}
	dangerLvl *= inQuadrant
	// fmt.Printf("Robots in fourth quadrant: %d\n", inQuadrant)

	return dangerLvl
}

func main() {
	verbose := flag.Bool("v", false, "verbose")
	flag.Parse()
	file, err := os.Open(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	robots := make([]Robot, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var rob Robot
		_, err := fmt.Sscanf(line, "p=%d,%d v=%d,%d", &rob.px, &rob.py, &rob.vx, &rob.vy)
		if err != nil {
			panic(err)
		}
		robots = append(robots, rob)
	}

	fmt.Println("Initial state:")
	fmt.Println(RepresentRobots(robots))
	if *verbose {
		for r := range robots {
			fmt.Printf("Robot %d: %d, %d\n", r, robots[r].px, robots[r].py)
		}
	}

	lowestDanger := math.MaxInt
	lowestDangerIter := 0
	lowestDangerRobots := make([]Robot, len(robots))
	for i := range ITERATIONS {
		PerformStep(robots)
		dangerLvl := CalcDangerLvl(robots)
		if dangerLvl < lowestDanger {
			lowestDanger = dangerLvl
			lowestDangerIter = i + 1
			copy(lowestDangerRobots, robots)
		}
		if *verbose {
			fmt.Printf("After %d seconds:\n", i+1)
			fmt.Println(RepresentRobots(robots))
			// Print robots positions
			// for r := range robots {
			// 	fmt.Printf("Robot %d: %d, %d\n", r, robots[r].px, robots[r].py)
			// }
			fmt.Printf("Danger: %d\n", dangerLvl)
		}
	}
	fmt.Printf("After %d seconds:\n", ITERATIONS)
	fmt.Println(RepresentRobots(robots))
	fmt.Printf("Danger: %d\n", CalcDangerLvl(robots))

	// Part 2, print lowest danger
	fmt.Printf("Lowest danger after %d seconds:\n", lowestDangerIter)
	fmt.Println(RepresentRobots(lowestDangerRobots))

}
