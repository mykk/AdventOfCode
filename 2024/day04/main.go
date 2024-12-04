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

func countAS(grid [][]byte, pos Position, target byte, direction Direction) int {
	newPos := Position{pos.x + direction.x, pos.y + direction.y}
	if withinBounds(grid, newPos) && grid[newPos.x][newPos.y] == target {
		if target == 'S' {
			return 1
		}
		return countAS(grid, newPos, 'S', direction)
	}
	return 0
}

func countMAS(grid [][]byte, directions []Direction, pos Position, dirTransform func(Direction) Direction) (xmasCount int) {
	for _, direction := range directions {
		newPos := Position{pos.x + direction.x, pos.y + direction.y}
		if withinBounds(grid, newPos) && grid[newPos.x][newPos.y] == 'M' {
			xmasCount += countAS(grid, newPos, 'A', dirTransform(direction))
		}
	}
	return
}

func FindXmas(grid [][]byte) (xmasCount int) {
	directions := []Direction{{0, 1}, {0, -1}, {-1, 0}, {1, 0}, {1, 1}, {-1, 1}, {1, -1}, {-1, -1}}

	for x, row := range grid {
		for y, cell := range row {
			if cell == 'X' {
				xmasCount += countMAS(grid, directions, Position{x, y}, func(dir Direction) Direction { return dir })
			}
		}
	}
	return
}

func Find_X_MAS(grid [][]byte) (xmasCount int) {
	directions := []Direction{{1, 1}, {-1, 1}, {1, -1}, {-1, -1}}
	dirTransform := func(dir Direction) Direction { return Direction{-dir.x, -dir.y} }

	for x, row := range grid {
		for y, cell := range row {
			if cell == 'A' && countMAS(grid, directions, Position{x, y}, dirTransform) == 2 {
				xmasCount += 1
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
