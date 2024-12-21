package main

import (
	"AoC/aoc"
	fn "AoC/functional"
	"container/heap"
	"errors"
	"strconv"
	"strings"

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
	position Point
	distance int
}

func withinBounds(limit Point, point Point) bool {
	if point.x < 0 || point.y < 0 {
		return false
	}
	return point.x <= limit.x && point.y <= limit.y
}

func FindPath(fallingBits []Point, limit int, start, end Point) (int, error) {
	states := aoc.NewHeap[State](func(a, b State) bool { return a.distance < b.distance })
	visited := make(aoc.Set[Point])

	states.PushItem(State{position: start, distance: 0})

	for states.Len() != 0 {
		state := states.PopItem()

		if state.position == end {
			return state.distance, nil
		}

		if visited.Contains(state.position) {
			continue
		}
		visited.Add(state.position)

		for _, direction := range Directions {
			position := Point{x: state.position.x + direction.dx, y: state.position.y + direction.dy}

			if visited.Contains(position) {
				continue
			}

			if !withinBounds(end, position) || fn.Contains(fallingBits[:limit], position) {
				continue
			}

			heap.Push(states, State{position: position, distance: state.distance + 1})
		}
	}

	return 0, errors.New("no escape")
}

func FindNoEscapeLimit(fallingBits []Point, good int, start, end Point) (Point, error) {
	if _, err := FindPath(fallingBits, len(fallingBits), start, end); err == nil {
		return Point{}, errors.New("path is never blocked")
	}

	bad := len(fallingBits) - 1
	for bad > good+1 {
		test := (good + bad) / 2
		if _, err := FindPath(fallingBits, test+1, start, end); err != nil {
			bad = test
		} else {
			good = test
		}
	}
	return fallingBits[bad], nil
}

func ParseInputData(data string) ([]Point, error) {
	return fn.Transform(fn.GetLines(data), func(line string) (Point, error) {
		point := strings.Split(line, ",")
		if len(point) != 2 {
			return Point{}, errors.New("bad input")
		}
		x, err := strconv.Atoi(point[0])
		if err != nil {
			return Point{}, err
		}

		y, err := strconv.Atoi(point[1])
		if err != nil {
			return Point{}, err
		}

		return Point{x: x, y: y}, nil
	})
}

func main() {
	inputData, err := os.ReadFile("day18.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	fallingBits, err := ParseInputData(string(inputData))
	if err != nil {
		fmt.Printf("Error parsing input data: %v\n", err)
		return
	}

	escape, err := FindPath(fallingBits, 1024, Point{x: 0, y: 0}, Point{x: 70, y: 70})
	if err != nil {
		fmt.Printf("there is no escape")
		return
	}

	fmt.Printf("Part 1: %d\n", escape)

	noEscape, err := FindNoEscapeLimit(fallingBits, 1024, Point{x: 0, y: 0}, Point{x: 70, y: 70})
	if err != nil {
		fmt.Printf("there is always ways to escape")
		return
	}

	fmt.Printf("Part 2: %d,%d\n", noEscape.x, noEscape.y)
}
