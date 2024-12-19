package main

import (
	fn "AoC/functional"
	"errors"
	"strings"

	"fmt"
	"os"
)

func matchPattern(pattern string, towels []string, cache map[string]int, patternMatching func(sum, current int) int) int {
	if len(pattern) == 0 {
		return 1
	}

	if patternCount, found := cache[pattern]; found {
		return patternCount
	}

	patternCount := 0
	for _, towel := range towels {
		if after, found := strings.CutPrefix(pattern, towel); found {
			patternCount = patternMatching(patternCount, matchPattern(after, towels, cache, patternMatching))
		}
	}
	cache[pattern] = patternCount
	return patternCount
}

func countPatterns(towels, patterns []string, patternMatching func(sum, current int) int) int {
	cache := make(map[string]int)
	return fn.Reduce(patterns, 0, func(_, count int, pattern string) int {
		return count + matchPattern(pattern, towels, cache, patternMatching)
	})
}

func CountPossiblePatterns(towels, patterns []string) int {
	return countPatterns(towels, patterns, func(sum, current int) int { return sum | current })
}

func CountDifferentCombinations(towels, patterns []string) int {
	return countPatterns(towels, patterns, func(sum, current int) int { return sum + current })
}

func ParseInputData(data string) ([]string, []string, error) {
	lines := fn.GetLines(data)
	if len(lines) < 3 {
		return nil, nil, errors.New("unexpected file format")
	}
	return strings.Split(lines[0], ", "), lines[2:], nil
}

func main() {
	inputData, err := os.ReadFile("day19.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	towels, patterns, err := ParseInputData(string(inputData))
	if err != nil {
		fmt.Printf("Error parsing input data: %v\n", err)
		return
	}

	fmt.Printf("Part 1: %d\n", CountPossiblePatterns(towels, patterns))
	fmt.Printf("Part 2: %d\n", CountDifferentCombinations(towels, patterns))
}
