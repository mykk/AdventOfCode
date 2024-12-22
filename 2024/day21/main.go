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

			directions := make([]byte, len(state.directions))
			copy(directions, state.directions)
			directions = append(directions, DirectionMap[direction])

			direcitonSwitch := 0
			if len(directions) > 1 {
				direcitonSwitch = fn.Reduce(directions[:len(directions)-1], 0, func(i, sum int, dir byte) int {
					if dir != directions[i+1] {
						return sum + 1
					}
					return sum
				})
			}
			if direcitonSwitch > 1 {
				continue
			}

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

func typeInInstructions(code []byte, pathMap map[BytePair][][]byte) [][]byte {
	var instructions [][]byte

	for i, key := range code {
		var newInstructions [][]byte

		for _, path := range pathMap[BytePair{getOrDefault(code, i-1, 'A'), key}] {
			if len(instructions) == 0 {
				newInstructions = append(newInstructions, append(append([]byte{}, path...), 'A'))
				continue
			}

			for _, currentInstruction := range instructions {
				newInstructions = append(newInstructions, append(append(append([]byte{}, currentInstruction...), path...), 'A'))
			}

		}
		instructions = newInstructions
	}

	return instructions
}

func getLenFromCache(instructions []byte, robotCount int, cache map[string]map[int]int) (int, bool) {
	if robotCountToLenMap, found := cache[string(instructions)]; !found {
		return 0, false
	} else {
		length, found := robotCountToLenMap[robotCount]
		return length, found
	}
}

func calculateSingleInstructionLen(instruction []byte, robotCount int, pathMap map[BytePair][][]byte, cache map[string]map[int]int) int {
	if length, found := getLenFromCache(instruction, robotCount, cache); found {
		return length
	}

	if _, found := cache[string(instruction)]; !found {
		cache[string(instruction)] = make(map[int]int)
	}

	newInstructions := typeInInstructions(instruction, pathMap)

	for _, newInstruction := range newInstructions {
		length := calculateLen(newInstruction, pathMap, robotCount-1, cache)

		if existingLen, found := cache[string(instruction)][robotCount]; !found || length < existingLen {
			cache[string(instruction)][robotCount] = length
		}
	}

	return cache[string(instruction)][robotCount]
}

func splitInstructions(instructions []byte) [][]byte {
	var instructionSplits [][]byte
	splitAt := 0

	for i, cell := range instructions {
		if cell == 'A' {
			instructionSplits = append(instructionSplits, append([]byte{}, instructions[splitAt:i+1]...))
			splitAt = i + 1
		}
	}

	return instructionSplits
}

func calculateLen(instructions []byte, pathMap map[BytePair][][]byte, robotCount int, cache map[string]map[int]int) int {
	if robotCount == 0 {
		return len(instructions)
	}

	instructionSplits := splitInstructions(instructions)

	return fn.Reduce(instructionSplits, 0, func(_, sum int, instruction []byte) int {
		return sum + calculateSingleInstructionLen(instruction, robotCount, pathMap, cache)
	})
}

func calculateCodeComplexity(code []byte, numpadPaths, directionalPaths map[BytePair][][]byte, robotCount int, cache map[string]map[int]int) int {
	instructions := typeInInstructions(code, numpadPaths)

	minLen := fn.Reduce(instructions[1:], calculateLen(instructions[0], directionalPaths, robotCount, cache), func(_, minLen int, instruction []byte) int {
		currentLen := calculateLen(instruction, directionalPaths, robotCount, cache)
		if currentLen < minLen {
			return currentLen
		}
		return minLen
	})

	codeNumber, _ := strconv.Atoi(string(code[:3]))
	return minLen * codeNumber
}

func GetCodeSumComplexities(codes [][]byte, robotCount int) int {
	numpadPaths := makePathMap(NumericKeypad, NumericValues)
	directionalPaths := makePathMap(DirectionalKeypad, DirectionalValues)
	cache := make(map[string]map[int]int)

	return fn.Reduce(codes, 0, func(_, sum int, code []byte) int {
		return sum + calculateCodeComplexity(code, numpadPaths, directionalPaths, robotCount, cache)
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
	fmt.Printf("Part 2: %d\n", GetCodeSumComplexities(codes, 25))
}
