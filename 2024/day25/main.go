package main

import (
	fn "AoC/functional"

	"fmt"
	"os"
)

func countUnlockableByKey(locks [][]int, key []int, lockSize int) int {
	return fn.CountIf(locks, func(lock []int) bool {
		return fn.All(key, func(i, keyField int) bool { return keyField+lock[i] <= lockSize })
	})
}

func CountFittingKeys(keys [][]int, locks [][]int, lockSize int) int {
	return fn.Reduce(keys, 0, func(_, sum int, key []int) int {
		return sum + countUnlockableByKey(locks, key, lockSize)
	})
}

func parseGrid(grid []string) (processedRows int, columnCounts []int) {
	columnCounts = make([]int, len(grid[0]))

	for _, row := range grid {
		if row == "" {
			break
		}

		for i, cell := range row {
			if byte(cell) == '#' {
				columnCounts[i]++
			}
		}
		processedRows++
	}
	return
}

func ParseInputData(data string) (keys [][]int, locks [][]int, lockSize int) {
	lines := fn.GetLines(data)

	for i := 0; i < len(lines); {
		if fn.All([]byte(lines[i]), func(_ int, cell byte) bool { return cell == '#' }) {
			offset, lock := parseGrid(lines[i:])
			i += offset + 1
			lockSize = offset
			locks = append(locks, lock)
		} else {
			offset, key := parseGrid(lines[i:])
			i += offset + 1
			keys = append(keys, key)
		}
	}
	return
}

func main() {
	inputData, err := os.ReadFile("day25.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	keys, locks, lockSize := ParseInputData(string(inputData))

	fmt.Printf("Part 1: %d\n", CountFittingKeys(keys, locks, lockSize))
}
