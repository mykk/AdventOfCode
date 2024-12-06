package main

import (
	fn "AoC/functional"
	"errors"

	"fmt"
	"os"
)

type Direction struct {
	dx, dy int
}

type Position struct {
	x, y int
}

type State struct {
	pos Position
	dir Direction
}

type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(value T) {
	s[value] = struct{}{}
}

func (s Set[T]) Contains(value T) bool {
	_, exists := s[value]
	return exists
}

func withinBounds(grid [][]byte, pos Position) bool {
	return pos.y >= 0 && pos.y < len(grid) && pos.x >= 0 && pos.x < len(grid[pos.y])
}

func getNextDirection(dir Direction) Direction {
	var turn = map[Direction]Direction{
		{dx: 0, dy: 1}:  {dx: -1, dy: 0},
		{dx: 1, dy: 0}:  {dx: 0, dy: 1},
		{dx: 0, dy: -1}: {dx: 1, dy: 0},
		{dx: -1, dy: 0}: {dx: 0, dy: -1},
	}

	return turn[dir]
}

func getNextPostion(pos Position, direction Direction, grid [][]byte, obsticle *Position) (Position, Direction, error) {
	nextPostion := Position{pos.x + direction.dx, pos.y + direction.dy}

	if !withinBounds(grid, nextPostion) {
		return Position{}, Direction{}, errors.New("out of bounds")
	}
	if grid[nextPostion.y][nextPostion.x] == '#' || (obsticle != nil && nextPostion == *obsticle) {
		return getNextPostion(pos, getNextDirection(direction), grid, obsticle)
	}
	return nextPostion, direction, nil
}

func TrackGuard(guardPostion Position, grid [][]byte, obsticle *Position) (Set[Position], error) {
	positionSet := make(Set[Position])
	stateSet := make(Set[State])
	direction := Direction{0, -1}

	positionSet.Add(guardPostion)

	position, direction, err := getNextPostion(guardPostion, direction, grid, obsticle)

	for err == nil {
		positionSet.Add(position)
		if stateSet.Contains(State{position, direction}) {
			return nil, errors.New("guard stuck in loop")
		}
		stateSet.Add(State{Position{position.x, position.y}, Direction{direction.dx, direction.dy}})

		position, direction, err = getNextPostion(position, direction, grid, obsticle)
	}

	return positionSet, nil
}

func TrapGuard(guardPosition Position, positionSet Set[Position], grid [][]byte) (count int) {
	for pos := range positionSet {
		if pos == guardPosition {
			continue
		}
		if _, err := TrackGuard(guardPosition, grid, &pos); err != nil {
			count++
		}
	}

	return
}

func ParseInputData(data string) (Position, [][]byte, error) {
	lines := fn.GetLines(data)

	var guardPosition *Position = nil
	for y, line := range lines {
		for x, cell := range line {
			if cell == '^' {
				guardPosition = &Position{x, y}
				break
			}
		}
	}

	if guardPosition == nil {
		return Position{}, nil, errors.New("guard position not found")
	}

	grid := fn.MustTransform(lines, func(line string) []byte { return []byte(line) })
	return *guardPosition, grid, nil
}

func main() {
	inputData, err := os.ReadFile("day6.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	guardPosition, grid, err := ParseInputData(string(inputData))
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		return
	}
	positions, err := TrackGuard(guardPosition, grid, nil)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("Part 1: %d\n", len(positions))
	fmt.Printf("Part 2: %d\n", TrapGuard(guardPosition, positions, grid))
}
