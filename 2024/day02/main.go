package main

import (
	fn "AoC/functional"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	MinDiff = 1
	MaxDiff = 3
)

func ParseInputData(data []string) ([][]int, error) {
	var report [][]int

	for _, line := range data {
		splitLine := strings.Fields(line)
		if len(splitLine) < 2 {
			return nil, fmt.Errorf("unexpected data format in line: '%s'", line)
		}

		lineReport, err := fn.Transform(splitLine, func(value string) (int, error) {
			return strconv.Atoi(value)
		})
		if err != nil {
			return nil, fmt.Errorf("failed to parse line '%s': %w", line, err)
		}

		report = append(report, lineReport)
	}

	return report, nil
}

func cloneWithoutElement(slice []int, index int) []int {
	result := make([]int, len(slice)-1)
	copy(result, slice[:index])
	copy(result[index:], slice[index+1:])
	return result
}

func isGoodReport(lineReport []int, errorTolerance uint) bool {
	direction := 1
	if lineReport[0] < lineReport[1] {
		direction = -1
	}

	for i := 0; i < len(lineReport)-1; i++ {
		diff := (lineReport[i] - lineReport[i+1]) * direction
		if diff < MinDiff || diff > MaxDiff {
			if errorTolerance == 0 {
				return false
			} else {
				for j := max(0, i-1); j < min(len(lineReport), i+2); j++ {
					modifiedReport := cloneWithoutElement(lineReport, j)
					if isGoodReport(modifiedReport, errorTolerance-1) {
						return true
					}
				}
				return false
			}
		}
	}
	return true
}

func FindGoodReports(report [][]int, errorTolerance uint) (goodReports int) {
	for _, lineReport := range report {
		if isGoodReport(lineReport, errorTolerance) {
			goodReports++
		}
	}

	return
}

func main() {
	inputData, err := os.ReadFile("day2.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	data := fn.GetLines(string(inputData))
	report, err := ParseInputData(data)
	if err != nil {
		fmt.Printf("Error parsing input data: %v\n", err)
		return
	}

	fmt.Printf("Part 1: %d\n", FindGoodReports(report, 0))
	fmt.Printf("Part 2: %d\n", FindGoodReports(report, 1))
}
