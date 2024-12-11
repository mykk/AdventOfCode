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

func countCountPlutonianPebble(stone, blinkCount int, cache map[int]map[int]int) int {
	if blinkCount == 0 {
		return 1
	}

	if stoneValues, found := cache[stone]; found {
		if value, found := stoneValues[blinkCount]; found {
			return value
		}
	} else {
		cache[stone] = make(map[int]int)
	}

	if stone == 0 {
		cache[stone][blinkCount] = countCountPlutonianPebble(1, blinkCount-1, cache)
	} else if digits := countDigis(stone); digits%2 == 0 {
		split := findMidSplit(digits)
		value := countCountPlutonianPebble(stone%split, blinkCount-1, cache)
		value += countCountPlutonianPebble(stone/split, blinkCount-1, cache)
		cache[stone][blinkCount] = value
	} else {
		cache[stone][blinkCount] = countCountPlutonianPebble(stone*2024, blinkCount-1, cache)
	}

	return cache[stone][blinkCount]
}

func CountPlutonianPebbles(stones []int, blinkCount int) int {
	cache := make(map[int]map[int]int)

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
