package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func IsReportSafe(levels []int) bool {
	var lastLevel int
	var isIncreasing bool
	for i, level := range levels {
		// Distance between the last level must be between 1 and 3
		if i == 0 {
			lastLevel = level
		} else {
			dist := lastLevel - level
			if dist < 0 {
				dist = -dist
			}
			if (dist < 1) || (3 < dist) {
				return false
			}
		}

		// On second iteration one can decide if it is decreasing or increasing
		if i == 1 {
			if lastLevel < level {
				isIncreasing = true
			} else {
				isIncreasing = false
			}
		} else if i > 1 {
			if isIncreasing && lastLevel > level {
				return false
			} else if !isIncreasing && lastLevel < level {
				return false
			}
		}
		lastLevel = level
	}

	return true
}

func main() {
	// Parse first argument as a file name
	filename := os.Args[1]
	fmt.Printf("Parsing file: %s\n", filename)
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	safeReports := 0
	safeReportsWithDumpener := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		report := scanner.Text()

		// Parse levels
		fields := strings.Fields(report)
		levels := make([]int, len(fields))
		for i, field := range fields {
			level, err := strconv.Atoi(field)
			if err != nil {
				panic(err)
			}
			levels[i] = level
		}

		// Check if unmodified report is safe
		if IsReportSafe(levels) {
			safeReports++
			safeReportsWithDumpener++
			continue
		}

		// Check if report without a single level is safe
		for i := range levels {
			// Create modifiedLevels
			modifiedLevels := make([]int, 0)
			modifiedLevels = append(modifiedLevels, levels[:i]...)
			modifiedLevels = append(modifiedLevels, levels[i+1:]...)

			if IsReportSafe(modifiedLevels) {
				safeReportsWithDumpener++
				break
			}
		}
	}

	fmt.Printf("Safe reports: %d\n", safeReports)
	fmt.Printf("Safe reports with problem dumpener: %d\n", safeReportsWithDumpener)
}
