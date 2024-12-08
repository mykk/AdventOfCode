package main

import (
	fn "AoC/functional"
	"errors"

	"fmt"
	"os"
)

type Position struct {
	x, y int
}

type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(value T) {
	s[value] = struct{}{}
}

func withinBounds(mapSize Position, pos Position) bool {
	return pos.y >= 0 && pos.y < mapSize.y && pos.x >= 0 && pos.x < mapSize.x
}

func appendIfWithinBounds(set Set[Position], mapSize Position, pos Position) error {
	if withinBounds(mapSize, pos) {
		set.Add(pos)
		return nil
	}
	return errors.New("out of bounds")
}

type AntinodeGenerator func(a1, a2 Position) func() (Position, Position, error)

func countAntinodes(antennas map[byte][]Position, mapSize Position, antinodeGenerator AntinodeGenerator) int {
	antinodes := make(Set[Position])

	for _, positions := range antennas {
		for i, a1 := range positions[:len(positions)-1] {
			for _, a2 := range positions[i+1:] {
				gen := antinodeGenerator(a1, a2)

				for antinode1, antinode2, err := gen(); err == nil; antinode1, antinode2, err = gen() {
					err1, err2 := appendIfWithinBounds(antinodes, mapSize, antinode1), appendIfWithinBounds(antinodes, mapSize, antinode2)
					if err1 != nil && err2 != nil {
						break
					}
				}
			}
		}
	}
	return len(antinodes)
}

func CountSingleAntinodes(antennas map[byte][]Position, mapSize Position) int {
	singleAntiNodeGenerator := func(a1, a2 Position) func() (Position, Position, error) {
		terminate := false
		diff := Position{x: a1.x - a2.x, y: a1.y - a2.y}

		return func() (Position, Position, error) {
			if terminate {
				return Position{}, Position{}, errors.New("no antinodes")
			}
			terminate = true
			return Position{x: a1.x + diff.x, y: a1.y + diff.y}, Position{x: a2.x - diff.x, y: a2.y - diff.y}, nil
		}
	}

	return countAntinodes(antennas, mapSize, singleAntiNodeGenerator)
}

func CountHarmonicsAntinodes(antennas map[byte][]Position, mapSize Position) int {
	harmonicsAntiNodeGenerator := func(a1, a2 Position) func() (Position, Position, error) {
		counter := -1
		diff := Position{x: a1.x - a2.x, y: a1.y - a2.y}

		return func() (Position, Position, error) {
			counter++
			return Position{x: a1.x + diff.x*counter, y: a1.y + diff.y*counter}, Position{x: a2.x - diff.x*counter, y: a2.y - diff.y*counter}, nil
		}
	}

	return countAntinodes(antennas, mapSize, harmonicsAntiNodeGenerator)
}

func ParseInputData(data string) (map[byte][]Position, Position, error) {
	lines := fn.GetLines(data)
	if len(lines) == 0 {
		return nil, Position{}, errors.New("empty map")
	}

	antennas := make(map[byte][]Position)
	for row, line := range lines {
		for column, antenna := range []byte(line) {
			if antenna == '.' {
				continue
			}

			if _, exists := antennas[antenna]; exists {
				antennas[antenna] = append(antennas[antenna], Position{x: column, y: row})
			} else {
				antennas[antenna] = []Position{{x: column, y: row}}
			}
		}
	}

	return antennas, Position{x: len(lines[0]), y: len(lines)}, nil
}

func main() {
	inputData, err := os.ReadFile("day8.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	antennas, mapSize, err := ParseInputData(string(inputData))
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		return
	}

	fmt.Printf("Part 1: %d\n", CountSingleAntinodes(antennas, mapSize))
	fmt.Printf("Part 2: %d\n", CountHarmonicsAntinodes(antennas, mapSize))
}
