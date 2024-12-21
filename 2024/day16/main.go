package main

import (
	"AoC/aoc"
	fn "AoC/functional"
	"container/heap"
	"errors"

	"fmt"
	"os"
)

type Direction struct {
	dx, dy int
}

var Directions = []Direction{
	{dx: 1, dy: 0},
	{dx: -1, dy: 0},
	{dx: 0, dy: 1},
	{dx: 0, dy: -1},
}

type Point struct {
	x, y int
}

type PointAndDir struct {
	position  Point
	direction Direction
}

type State struct {
	PointAndDir
	points int
	path   aoc.Set[Point]
}

func withinBounds(grid [][]byte, x, y int) bool {
	return y >= 0 && y < len(grid) && x >= 0 && x < len(grid[y])
}

func SolveMaze(grid [][]byte, start, end Point) (winningPaths []State) {
	states := aoc.NewHeap[State](func(a, b State) bool { return a.points < b.points })
	visited := make(map[PointAndDir]int)

	states.PushItem(State{PointAndDir: PointAndDir{position: start, direction: Direction{dx: 1, dy: 0}}, points: 0, path: aoc.NewSet(start)})

	for states.Len() != 0 {
		state := states.PopItem()

		if len(winningPaths) != 0 && winningPaths[0].points < state.points {
			return
		}

		if previousPoints, found := visited[state.PointAndDir]; found && previousPoints < state.points {
			continue
		}
		visited[state.PointAndDir] = state.points

		if state.position == end {
			winningPaths = append(winningPaths, state)
			continue
		}

		reverseDir := Direction{dx: -state.direction.dx, dy: -state.direction.dy}

		for _, direction := range Directions {
			if direction == reverseDir {
				continue
			}
			position := Point{x: state.position.x + direction.dx, y: state.position.y + direction.dy}

			if previousPoints, found := visited[PointAndDir{position: position, direction: direction}]; found && previousPoints < state.points+1 {
				continue
			}

			if !withinBounds(grid, position.x, position.y) || grid[position.y][position.x] == '#' {
				continue
			}

			points := state.points
			if direction == state.direction {
				points += 1
			} else {
				points += 1001
			}

			newPath := state.path.Clone()
			newPath.Add(position)
			heap.Push(states, State{PointAndDir: PointAndDir{position: position, direction: direction}, points: points, path: newPath})
		}
	}
	return
}

func ParseInputData(data string) (grid [][]byte, start, end Point, err error) {
	grid = fn.MustTransform(fn.GetLines(data), func(line string) []byte { return []byte(line) })

	invalidPoint := Point{-1, -1}
	start, end = invalidPoint, invalidPoint
	for y, row := range grid {
		for x, cell := range row {
			if cell == 'S' {
				start = Point{x: x, y: y}
			}
			if cell == 'E' {
				end = Point{x: x, y: y}
			}
		}
	}
	if start == invalidPoint {
		return nil, Point{}, Point{}, errors.New("start point not found")
	}
	if end == invalidPoint {
		return nil, Point{}, Point{}, errors.New("end point not found")
	}
	return
}

func HotPaths(winningPaths []State) int {
	hotPaths := fn.Reduce(winningPaths, make(aoc.Set[Point]), func(_ int, points aoc.Set[Point], state State) aoc.Set[Point] {
		for point := range state.path {
			points.Add(point)
		}
		return points
	})
	return len(hotPaths)
}

func main() {
	inputData, err := os.ReadFile("day16.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	grid, start, end, err := ParseInputData(string(inputData))
	if err != nil {
		fmt.Printf("Error parsing input data: %v\n", err)
		return
	}

	winningPaths := SolveMaze(grid, start, end)
	if len(winningPaths) == 0 {
		fmt.Println("No paths found.")
		return
	}

	fmt.Printf("Part 1: %d\n", winningPaths[0].points)
	fmt.Printf("Part 2: %d\n", HotPaths(winningPaths))
}
