package main

import (
	fn "AoC/functional"
	"errors"
	"strconv"
	"strings"

	"fmt"
	"os"
)

type Entry struct {
	target int
	data   []int
}

type Operation func(int, int) int

func calculateSingleCalibration(target int, current int, data []int, operations []Operation) int {
	if len(data) == 0 && target == current {
		return target
	} else if len(data) == 0 {
		return 0
	}

	for _, operator := range operations {
		if target == calculateSingleCalibration(target, operator(current, data[0]), data[1:], operations) {
			return target
		}
	}

	return 0
}

func calculateTotalCalibration(entries []Entry, additionalOperations []Operation) int {
	operations := []Operation{
		func(lhs, rhs int) int { return lhs + rhs },
		func(lhs, rhs int) int { return lhs * rhs }}
	operations = append(operations, additionalOperations...)

	return fn.Reduce(entries, 0, func(_, currentCalibration int, entry Entry) int {
		if len(entry.data) == 0 {
			return currentCalibration
		}
		return currentCalibration + calculateSingleCalibration(entry.target, entry.data[0], entry.data[1:], operations)
	})
}

func CalculateSimpleCalibration(entries []Entry) int {
	return calculateTotalCalibration(entries, []Operation{})
}

func CalculateFancyCalibration(entries []Entry) int {
	combine := func(lhs, rhs int) int {
		digits := 10
		for temp := rhs / 10; temp > 0; temp /= 10 {
			digits *= 10
		}
		return lhs*digits + rhs
	}

	return calculateTotalCalibration(entries, []Operation{combine})
}

func ParseInputData(data string) ([]Entry, error) {
	return fn.Transform(fn.GetLines(data), func(line string) (Entry, error) {
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			return Entry{}, errors.New("bad input data")
		}

		target, err := strconv.Atoi(parts[0])
		if err != nil {
			return Entry{}, err
		}

		data, err := fn.Transform(strings.Fields(parts[1]), strconv.Atoi)
		if err != nil {
			return Entry{}, err
		}

		return Entry{target, data}, nil
	})
}

func main() {
	inputData, err := os.ReadFile("day7.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	entries, err := ParseInputData(string(inputData))
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		return
	}
	fmt.Printf("Part 1: %d\n", CalculateSimpleCalibration(entries))
	fmt.Printf("Part 2: %d\n", CalculateFancyCalibration(entries))
}
