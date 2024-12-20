package main

import (
	fn "AoC/functional"

	"fmt"
	"os"
)

type Direction struct {
	x, y int
}

type Position struct {
	x, y int
}

func withinBounds(grid [][]byte, pos Position) bool {
	return pos.x >= 0 && pos.x < len(grid) && pos.y >= 0 && pos.y < len(grid[pos.x])
}

func match(grid [][]byte, pos Position, target byte) bool {
	return withinBounds(grid, pos) && grid[pos.x][pos.y] == target
}

func followPattern(grid [][]byte, pos Position, pattern []byte, dir Direction) int {
	if fn.All(pattern, func(i int, cell byte) bool { return match(grid, Position{pos.x + dir.x*i, pos.y + dir.y*i}, cell) }) {
		return 1
	}
	return 0
}

func FindXmas(grid [][]byte) (xmasCount int) {
	directions := []Direction{{0, 1}, {0, -1}, {-1, 0}, {1, 0}, {1, 1}, {-1, 1}, {1, -1}, {-1, -1}}

	for x, row := range grid {
		for y := range row {
			for _, direction := range directions {
				xmasCount += followPattern(grid, Position{x, y}, []byte{'X', 'M', 'A', 'S'}, direction)
			}
		}
	}
	return
}

func Find_X_MAS(grid [][]byte) (xmasCount int) {
	directions := []Direction{{1, 1}, {-1, 1}, {1, -1}, {-1, -1}}

	for x, row := range grid {
		for y, cell := range row {
			if cell == 'A' {
				currentCount := 0
				for _, dir := range directions {
					currentCount += followPattern(grid, Position{x + dir.x, y + dir.y}, []byte{'M', 'A', 'S'}, Direction{-dir.x, -dir.y})
				}
				if currentCount == 2 {
					xmasCount++
				}
			}
		}
	}
	return
}

func ParseInputData(data string) [][]byte {
	return fn.MustTransform(fn.GetLines(data), func(line string) []byte { return []byte(line) })
}

func main() {
	inputData, err := os.ReadFile("day4.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	grid := ParseInputData(string(inputData))
	fmt.Printf("Part 1: %d\n", FindXmas(grid))
	fmt.Printf("Part 2: %d\n", Find_X_MAS(grid))
}
