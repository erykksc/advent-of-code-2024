package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

type Direction [2]int

var (
	Up    = Direction{-1, 0}
	Down  = Direction{1, 0}
	Right = Direction{0, 1}
	Left  = Direction{0, -1}
)

// Returns a string representation of the region within the context of the full grid
func regionToString(grid [][]rune, region map[[2]int]int) string {
	// Convert to string
	var sb strings.Builder
	// sb.WriteString(fmt.Sprintf("Region bounds: (%d,%d) to (%d,%d)\n", minY, minX, maxY, maxX))
	for y, row := range grid {
		for x, cell := range row {
			if _, exists := region[[2]int{y, x}]; exists {
				sb.WriteRune(cell)
			} else {
				sb.WriteString(".")
			}
		}

		sb.WriteRune('\n')
	}

	return sb.String()
}

// Function that returns all positions of all nodes in a region
func findRegion(y, x int, plantType rune, grid [][]rune, region *map[[2]int]int) int {
	// Check if out of bounds
	if y < 0 || x < 0 || y >= len(grid) || x >= len(grid[0]) {
		return 0
	}
	// Check if already evaluated and in the region
	if _, exists := (*region)[[2]int{y, x}]; exists {
		return 1
	}
	// Different plant
	if grid[y][x] != plantType {
		return 0
	}

	(*region)[[2]int{y, x}] = 0
	// Save the current node
	neighbors := 0
	neighbors += findRegion(y-1, x, plantType, grid, region) // go left
	neighbors += findRegion(y, x+1, plantType, grid, region) // go right
	neighbors += findRegion(y+1, x, plantType, grid, region) // go down
	neighbors += findRegion(y, x-1, plantType, grid, region) // go left

	(*region)[[2]int{y, x}] = neighbors

	return 1
}

func countWalls(region map[[2]int]int) int {
	minY, minX := math.MaxInt, math.MaxInt
	maxY, maxX := 0, 0
	for cell := range region {
		if cell[0] > maxY {
			maxY = cell[0]
		}
		if cell[0] < minY {
			minY = cell[0]
		}
		if cell[1] > maxX {
			maxX = cell[1]
		}
		if cell[1] < minX {
			minX = cell[1]
		}
	}

	// Scan walls horizontally
	totalWalls := 0
	for y := minY - 1; y < maxY+2; y++ {
		buildingUpperWall := false
		buildingLowerWall := false
		upperWallStart := -1
		lowerWallStart := -1
		for x := minX; x <= maxX; x++ {
			// position y,x marks the position of potential wall
			_, plantAbove := region[[2]int{y - 1, x}]
			_, plantBelow := region[[2]int{y + 1, x}]
			_, insidePlant := region[[2]int{y, x}]

			shouldStopUpperWall := insidePlant || !plantAbove
			shouldStopLowerWall := insidePlant || !plantBelow
			shouldPlaceUpperWall := !insidePlant && plantAbove
			shouldPlaceLowerWall := !insidePlant && plantBelow

			if shouldStopUpperWall && buildingUpperWall {
				buildingUpperWall = false
				wallLen := x - upperWallStart
				// Prevents the case where initial cell is a wall
				if wallLen < 1 {
				} else {
					totalWalls++
				}
			}

			if shouldStopLowerWall && buildingLowerWall {
				buildingLowerWall = false
				wallLen := x - lowerWallStart
				// Prevents the case where initial cell is a wall
				if wallLen < 1 {
				} else {
					totalWalls++
				}
			}

			if shouldPlaceUpperWall && (!buildingUpperWall) {
				buildingUpperWall = true
				upperWallStart = x
			}

			if shouldPlaceLowerWall && (!buildingLowerWall) {
				buildingLowerWall = true
				lowerWallStart = x
			}
		}

		if buildingUpperWall {
			fmt.Printf("New wall segment until the end: %d:%d-%d, length:%d\n", y, upperWallStart, maxY+1, maxY+1-upperWallStart)
			totalWalls++
		}
		if buildingLowerWall {
			fmt.Printf("New wall segment until the end: %d:%d-%d, length:%d\n", y, upperWallStart, maxY+1, maxY+1-upperWallStart)
			totalWalls++
		}
	}

	fmt.Printf("Vertical scans\n")
	// Scan walls vertically
	for x := minX - 1; x < maxX+2; x++ {
		buildingLeftWall := false
		buildingRightWall := false
		leftWallStart := -1
		rightWallStart := -1
		for y := minY; y <= maxY; y++ {
			// position y,x marks the position of potential wall
			_, plantLeft := region[[2]int{y, x - 1}]
			_, plantRight := region[[2]int{y, x + 1}]
			_, insidePlant := region[[2]int{y, x}]

			shouldStopLeftWall := insidePlant || !plantLeft
			shouldStopRightWall := insidePlant || !plantRight
			shouldPlaceLeftWall := !insidePlant && plantLeft
			shouldPlaceRightWall := !insidePlant && plantRight

			// Build the left wall
			if shouldStopLeftWall && buildingLeftWall {
				buildingLeftWall = false
				wallLen := y - leftWallStart
				if wallLen < 1 {
					// Prevents the case where initial cell is a wall
				} else {
					fmt.Printf("New wall segment: %d:%d-%d, length:%d\n", x, leftWallStart, y, wallLen)
					totalWalls++
				}
			}
			if shouldPlaceLeftWall && (!buildingLeftWall) {
				buildingLeftWall = true
				leftWallStart = y
			}

			// Build the right wall
			if shouldStopRightWall && buildingRightWall {
				buildingRightWall = false
				wallLen := y - rightWallStart
				if wallLen < 1 {
					// Prevents the case where initial cell is a wall
				} else {
					fmt.Printf("New wall segment: %d:%d-%d, length:%d\n", x, rightWallStart, y, wallLen)
					totalWalls++
				}
			}
			if shouldPlaceRightWall && (!buildingRightWall) {
				buildingRightWall = true
				rightWallStart = y
			}
		}

		if buildingLeftWall {
			fmt.Printf("New wall segment until the end: %d:%d-%d, length:%d\n", x, leftWallStart, maxY+1, maxY+1-leftWallStart)
			totalWalls++
		}
		if buildingRightWall {
			fmt.Printf("New wall segment until the end: %d:%d-%d, length:%d\n", x, rightWallStart, maxY+1, maxY+1-rightWallStart)
			totalWalls++
		}
	}
	return totalWalls
}

func main() {
	// Read the input file into memory
	filename := os.Args[1]
	inputB, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	input := string(inputB)
	input = strings.Trim(input, "\n")
	lines := strings.Split(input, "\n")
	height, width := len(lines), len(lines[0])

	// Allocate memory for the entire grid
	allocatedGrid := make([]rune, height*width)
	grid := make([][]rune, height)
	for y := range height {
		grid[y] = allocatedGrid[y*width : (y+1)*width]
	}

	// Convert the input string into grid
	for y, line := range lines {
		for x, plantType := range line {
			grid[y][x] = plantType
		}
	}

	// Print input grid
	fmt.Println("Grid:")
	for _, row := range grid {
		for _, c := range row {
			fmt.Printf("%c", c)
		}
		fmt.Println()
	}

	// Find out regions
	totalPricePart1 := 0
	totalPricePart2 := 0
	cellVisited := make(map[[2]int]bool)
	for y := range height {
		for x := range width {
			plantType := grid[y][x]
			pos := [2]int{y, x}
			if cellVisited[pos] {
				continue
			}

			region := make(map[[2]int]int)
			findRegion(y, x, plantType, grid, &region)
			fmt.Printf("New region found:\n")
			fmt.Println(regionToString(grid, region))

			area := len(region)

			//calculate perimiter
			perimiter := 0
			fmt.Printf("Region of plant %c, area %d:\n", plantType, area)
			for plantPos, neighbors := range region {
				cellVisited[plantPos] = true // skip going through the same region again
				perimiter += 4 - neighbors
				fmt.Printf("(%d, %d)^%d ", plantPos[0], plantPos[1], neighbors)
			}
			pricePart1 := area * perimiter
			totalPricePart1 += pricePart1
			fmt.Printf("\nPerimiter of %d and price %d\n", perimiter, pricePart1)

			sides := countWalls(region)
			fmt.Printf("Total walls: %d\n", sides)

			totalPricePart2 += area * sides

		}
	}

	fmt.Printf("Total price for part 1: %d\n", totalPricePart1)
	fmt.Printf("Total price for part 2: %d\n", totalPricePart2)
}
