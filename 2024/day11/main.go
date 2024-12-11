package main

import (
	fn "AoC/functional"
	"strconv"
	"strings"

	"fmt"
	"os"
)

func countDigis(num int) (digits int) {
	for ; num != 0; num = num / 10 {
		digits += 1
	}
	return
}

func findMidSplit(digits int) int {
	split := 1
	for i := digits / 2; i != 0; i-- {
		split = split * 10
	}
	return split
}

type Pair[A, B any] struct {
	First  A
	Second B
}

type StoneBlinkPair Pair[int, int]

func countCountPlutonianPebble(stone, blinkCount int, cache map[StoneBlinkPair]int) int {
	if blinkCount == 0 {
		return 1
	}

	if value, found := cache[StoneBlinkPair{stone, blinkCount}]; found {
		return value
	}

	value := 0
	if stone == 0 {
		value = countCountPlutonianPebble(1, blinkCount-1, cache)
	} else if digits := countDigis(stone); digits%2 == 0 {
		split := findMidSplit(digits)
		value = countCountPlutonianPebble(stone%split, blinkCount-1, cache)
		value += countCountPlutonianPebble(stone/split, blinkCount-1, cache)
	} else {
		value = countCountPlutonianPebble(stone*2024, blinkCount-1, cache)
	}

	cache[StoneBlinkPair{stone, blinkCount}] = value

	return value
}

func CountPlutonianPebbles(stones []int, blinkCount int) int {
	cache := make(map[StoneBlinkPair]int)

	return fn.Reduce(stones, 0, func(_, sum, stone int) int {
		return sum + countCountPlutonianPebble(stone, blinkCount, cache)
	})
}

func ParseInputData(data string) ([]int, error) {
	return fn.Transform(strings.Fields(data), strconv.Atoi)
}

func main() {
	inputData, err := os.ReadFile("day11.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	stones, err := ParseInputData(string(inputData))
	if err != nil {
		fmt.Printf("Error parsing input data: %v\n", err)
		return
	}

	fmt.Printf("Part 1: %d\n", CountPlutonianPebbles(stones, 25))
	fmt.Printf("Part 2: %d\n", CountPlutonianPebbles(stones, 75))
}
