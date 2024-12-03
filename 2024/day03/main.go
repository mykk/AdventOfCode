package main

import (
	fn "AoC/functional"
	"regexp"
	"strconv"

	"fmt"
	"os"
)

type Multiples struct {
	a, b int
}

func MustAtoi(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}

func ParseInputData(data string) []Multiples {
	re := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)
	matches := re.FindAllStringSubmatch(data, -1)

	return fn.MustTransform(matches, func(match []string) Multiples {
		return Multiples{MustAtoi(match[1]), MustAtoi(match[2])}
	})
}

func ParseInputDataDoDont(data string) []Multiples {
	re := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)|do\(\)|don't\(\)`)
	matches := re.FindAllStringSubmatch(data, -1)

	do := true
	return fn.MustTransform(matches, func(match []string) Multiples {
		if match[0] == "do()" {
			do = true
		} else if match[0] == "don't()" {
			do = false
		} else if do {
			return Multiples{MustAtoi(match[1]), MustAtoi(match[2])}
		}
		return Multiples{0, 0}
	})
}

func SumMultiples(multiples []Multiples) int {
	return fn.Reduce(multiples, 0, func(index int, sum int, multiples Multiples) int {
		return sum + multiples.a*multiples.b
	})
}

func main() {
	inputData, err := os.ReadFile("day3.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	fmt.Printf("Part 1: %d\n", SumMultiples(ParseInputData(string(inputData))))
	fmt.Printf("Part 2: %d\n", SumMultiples(ParseInputDataDoDont(string(inputData))))
}
