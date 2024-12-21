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

type State struct {
	position          Point
	duration          int
	path              []Point
	positionDurations map[Point]int
}

func withinBounds(grid [][]byte, x, y int) bool {
	return y >= 0 && y < len(grid) && x >= 0 && x < len(grid[y])
}

func FindPath(grid [][]byte, start, end Point) State {
	states := aoc.NewHeap[State](func(a, b State) bool { return a.duration < b.duration })
	visited := make(aoc.Set[Point])

	states.PushItem(State{position: start, duration: 0, path: []Point{start}, positionDurations: map[Point]int{start: 0}})

	for states.Len() != 0 {
		state := states.PopItem()
		if visited.Contains(state.position) {
			continue
		}
		visited.Add(state.position)

		if state.position == end {
			return state
		}

		for _, direction := range Directions {
			position := Point{x: state.position.x + direction.dx, y: state.position.y + direction.dy}
			if visited.Contains(position) {
				continue
			}

			if !withinBounds(grid, position.x, position.y) || grid[position.y][position.x] == '#' {
				continue
			}

			positionDurations := make(map[Point]int, len(state.positionDurations))
			for key, value := range state.positionDurations {
				positionDurations[key] = value
			}
			positionDurations[position] = state.duration + 1

			path := make([]Point, len(state.path))
			copy(path, state.path)
			heap.Push(states, State{position: position, duration: state.duration + 1, positionDurations: positionDurations, path: append(path, position)})
		}
	}

	panic("no path found")
}

type CheatState struct {
	position Point
	duration int
}

func FindCheatPaths(grid [][]byte, finalState State, savesAtLeast, cheatDuration int) (cheats int) {
	if len(finalState.path) < savesAtLeast {
		return
	}

	for duration, point := range finalState.path[:len(finalState.path)-savesAtLeast] {
		states := aoc.NewHeap[CheatState](func(a, b CheatState) bool { return a.duration < b.duration })
		visited := make(aoc.Set[Point])

		states.PushItem(CheatState{position: point, duration: 0})

		for states.Len() != 0 {
			state := states.PopItem()
			if visited.Contains(state.position) {
				continue
			}
			visited.Add(state.position)

			if original, found := finalState.positionDurations[state.position]; found && original-(duration+state.duration) >= savesAtLeast {
				cheats++
			}
			if state.duration == cheatDuration {
				continue
			}

			for _, dir := range Directions {
				position := Point{x: state.position.x + dir.dx, y: state.position.y + dir.dy}
				if visited.Contains(position) {
					continue
				}
				states.PushItem(CheatState{position: position, duration: state.duration + 1})
			}
		}
	}
	return
}

func CountCheats(grid [][]byte, start, end Point, savesAtLeast, cheatDuration int) int {
	return FindCheatPaths(grid, FindPath(grid, start, end), savesAtLeast, cheatDuration)
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

func main() {
	inputData, err := os.ReadFile("day20.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	grid, start, end, err := ParseInputData(string(inputData))
	if err != nil {
		fmt.Printf("Error parsing input data: %v\n", err)
		return
	}

	fmt.Printf("Part 1: %d\n", CountCheats(grid, start, end, 100, 2))
	fmt.Printf("Part 2: %d\n", CountCheats(grid, start, end, 100, 20))
}
