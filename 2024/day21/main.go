package main

import (
	fn "AoC/functional"
	"strconv"

	"AoC/aoc"
	"container/heap"
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
	position   Point
	distance   int
	directions []byte
}

var NumericKeypad = [][]byte{{'7', '8', '9'}, {'4', '5', '6'}, {'1', '2', '3'}, {'#', '0', 'A'}}
var NumericValues = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A'}

var DirectionalKeypad = [][]byte{{'#', '^', 'A'}, {'<', 'v', '>'}}
var DirectionalValues = []byte{'<', '>', '^', 'v', 'A'}

var DirectionMap = map[Direction]byte{{dx: 1, dy: 0}: '>',
	{dx: -1, dy: 0}: '<',
	{dx: 0, dy: 1}:  'v',
	{dx: 0, dy: -1}: '^'}

func withinBounds(grid [][]byte, x, y int) bool {
	return y >= 0 && y < len(grid) && x >= 0 && x < len(grid[y])
}

func findPosition(grid [][]byte, item byte) Point {
	for y, row := range grid {
		for x, cell := range row {
			if item == cell {
				return Point{x: x, y: y}
			}
		}
	}
	panic("no item found")
}

func findAllPaths(grid [][]byte, start, end byte) [][]byte {
	startPosition := findPosition(grid, start)
	endPosition := findPosition(grid, end)

	finalStates := make([][]byte, 0)
	states := aoc.NewHeap[State](func(a, b State) bool { return a.distance < b.distance })

	states.PushItem(State{position: startPosition, distance: 0})

	for states.Len() != 0 {
		state := states.PopItem()

		if state.position == endPosition {
			finalStates = append(finalStates, state.directions)
			continue
		}

		if len(finalStates) != 0 && len(finalStates[0]) < state.distance {
			return finalStates
		}

		for _, direction := range Directions {
			position := Point{x: state.position.x + direction.dx, y: state.position.y + direction.dy}

			if !withinBounds(grid, position.x, position.y) || grid[position.y][position.x] == '#' {
				continue
			}

			if len(state.directions) > 1 && state.directions[0] != state.directions[1] && state.directions[1] != DirectionMap[direction] {
				continue
			}

			directions := make([]byte, len(state.directions))
			copy(directions, state.directions)
			directions = append(directions, DirectionMap[direction])

			heap.Push(states, State{position: position, distance: state.distance + 1, directions: directions})
		}
	}

	return finalStates
}

type BytePair struct {
	a, b byte
}

func makePathMap(grid [][]byte, values []byte) map[BytePair][][]byte {
	pathMap := make(map[BytePair][][]byte)
	for _, a := range values {
		if a == '#' {
			continue
		}
		for _, b := range values {
			if b == '#' {
				continue
			}
			pathMap[BytePair{a: a, b: b}] = findAllPaths(grid, a, b)
		}
	}
	return pathMap
}

func getOrDefault(code []byte, index int, defaultValue byte) byte {
	if index < 0 {
		return defaultValue
	}
	return code[index]
}

func typeInKeypadInstructions(code []byte, pathMap map[BytePair][][]byte, limit int) [][]byte {
	instructions := make([][]byte, 0)

	for i, key := range code {
		newInstructions := make([][]byte, 0)

		for i, path := range pathMap[BytePair{getOrDefault(code, i-1, 'A'), key}] {
			if i >= limit {
				break
			}
			if len(instructions) != 0 {
				for _, currentInstruction := range instructions {
					currentInstructionCopy := make([]byte, len(currentInstruction))
					copy(currentInstructionCopy, currentInstruction)
					newInstructions = append(newInstructions, append(append(currentInstructionCopy, path...), 'A'))
				}
			} else {
				newPath := make([]byte, 0, len(path))
				newPath = append(newPath, path...)
				newPath = append(newPath, 'A')
				newInstructions = append(newInstructions, newPath)
			}
		}
		instructions = newInstructions
	}

	return instructions
}

// func calculateFinalLen(directionalInstructions []byte, robotCount int) {

// }
func getCodeComplexity(code []byte, numpadPaths, directionalPaths map[BytePair][][]byte, robotCount int) int {
	directions := typeInKeypadInstructions(code, numpadPaths, 2)

	for i := 0; i < robotCount; i++ {
		newDirections := make([][]byte, 0)
		for _, direction := range directions {
			newDirections = append(newDirections, typeInKeypadInstructions(direction, directionalPaths, 1)...)
		}
		directions = newDirections
	}

	minLen := fn.Reduce(directions, len(directions[0]), func(_ int, minLen int, direction []byte) int {
		if len(direction) < minLen {
			return len(direction)
		}
		return minLen
	})

	codeNumber, _ := strconv.Atoi(string(code[:3]))
	return minLen * codeNumber
}

func GetCodeSumComplexities(codes [][]byte, robotCount int) int {
	numpadPaths := makePathMap(NumericKeypad, NumericValues)
	directionalPaths := makePathMap(DirectionalKeypad, DirectionalValues)

	return fn.Reduce(codes, 0, func(_, sum int, code []byte) int {
		return sum + getCodeComplexity(code, numpadPaths, directionalPaths, robotCount)
	})
}

func ParseInputData(input string) [][]byte {
	return fn.MustTransform(fn.GetLines(input), func(line string) []byte {
		return []byte(line)
	})
}

func main() {
	inputData, err := os.ReadFile("day21.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	codes := ParseInputData(string(inputData))
	if err != nil {
		fmt.Printf("Error parsing input data: %v\n", err)
		return
	}

	fmt.Printf("Part 1: %d\n", GetCodeSumComplexities(codes, 2))
}
