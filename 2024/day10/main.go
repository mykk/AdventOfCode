package main

import (
	fn "AoC/functional"

	"fmt"
	"os"
)

type Position struct {
	x, y int
}

type Direction struct {
	dx, dy int
}

type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(value T) {
	s[value] = struct{}{}
}

func withinBounds(grid [][]byte, x, y int) bool {
	return y >= 0 && y < len(grid) && x >= 0 && x < len(grid[y])
}

var directions = []Direction{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

func trackTrail(grid [][]byte, x, y int, nextPeak byte, visitPeak func(x, y int)) {
	if nextPeak == '9'+1 {
		visitPeak(x, y)
		return
	}

	for _, direction := range directions {
		nextX, nextY := x+direction.dx, y+direction.dy
		if withinBounds(grid, nextX, nextY) && grid[nextY][nextX] == nextPeak {
			trackTrail(grid, nextX, nextY, nextPeak+1, visitPeak)
		}
	}
}

func CountTrails(grid [][]byte) (trailCount int, trailScore int) {
	for y, row := range grid {
		for x, cell := range row {
			if cell == '0' {
				positions := make(Set[Position])
				trackTrail(grid, x, y, '1', func(x, y int) {
					trailScore += 1
					positions.Add(Position{x: x, y: y})
				})
				trailCount += len(positions)
			}
		}
	}

	return trailCount, trailScore
}

func ParseInputData(data string) [][]byte {
	return fn.MustTransform(fn.GetLines(data), func(line string) []byte { return []byte(line) })
}

func main() {
	inputData, err := os.ReadFile("day10.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	grid := ParseInputData(string(inputData))
	trailCount, trailScore := CountTrails(grid)
	fmt.Printf("Part 1: %d\n", trailCount)
	fmt.Printf("Part 2: %d\n", trailScore)
}
