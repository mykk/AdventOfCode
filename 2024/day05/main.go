package main

import (
	fn "AoC/functional"
	"errors"
	"strconv"
	"strings"

	"fmt"
	"os"
)

func getMidElement(slice []int) int {
	return slice[len(slice)/2]
}

func validUpdate(rules map[int][]int, update []int) bool {
	for index, page := range update {
		if fn.Any(update[:index], func(_, check int) bool { return fn.Contains(rules[page], check) }) {
			return false
		}
	}

	return true
}

func fixUpdate(rules map[int][]int, update []int) []int {
	for index, page := range update {
		for i := 0; i < index; i++ {
			if fn.Contains(rules[page], update[i]) {
				newUpdate := make([]int, 0, len(update))
				newUpdate = append(newUpdate, update[:i]...)
				newUpdate = append(newUpdate, page)
				newUpdate = append(newUpdate, update[i:index]...)
				newUpdate = append(newUpdate, update[index+1:]...)

				return fixUpdate(rules, newUpdate)
			}
		}
	}

	return update
}

func CountValidUpdates(rules map[int][]int, updates [][]int) int {
	return fn.Reduce(updates, 0, func(_ int, sum int, update []int) int {
		if validUpdate(rules, update) {
			return sum + getMidElement(update)
		}
		return sum
	})
}

func CountCorrectedUpdates(rules map[int][]int, updates [][]int) int {
	return fn.Reduce(updates, 0, func(_ int, sum int, update []int) int {
		if !validUpdate(rules, update) {
			return sum + getMidElement(fixUpdate(rules, update))
		}
		return sum
	})
}

func ParseInputData(data string) (map[int][]int, [][]int, error) {
	rules := make(map[int][]int)
	updates := [][]int{}

	ruleParse := true
	lines := fn.GetLines(data)
	for _, line := range lines {
		if line == "" {
			ruleParse = false
			continue
		}

		if ruleParse {
			parts := strings.Split(line, "|")
			if len(parts) != 2 {
				return nil, nil, errors.New("unexpected file format")
			}
			val1, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, nil, err
			}

			val2, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, nil, err
			}
			rules[val1] = append(rules[val1], val2)
		} else {
			update, err := fn.Transform(strings.Split(line, ","), func(el string) (int, error) { return strconv.Atoi(el) })
			if err != nil {
				return nil, nil, err
			}
			updates = append(updates, update)
		}
	}
	return rules, updates, nil
}

func main() {
	inputData, err := os.ReadFile("day5.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	rules, updates, err := ParseInputData(string(inputData))
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
	}
	fmt.Printf("Part 1: %d\n", CountValidUpdates(rules, updates))
	fmt.Printf("Part 2: %d\n", CountCorrectedUpdates(rules, updates))
}
