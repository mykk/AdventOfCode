package main

import (
	fn "AoC/functional"
	"errors"
	"regexp"
	"strconv"

	"fmt"
	"os"
)

func MustAtoi(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}

type Direction struct {
	dx, dy int
}

type Point struct {
	x, y int
}

type ClawMachine struct {
	buttonA Direction
	buttonB Direction
	prize   Point
}

func winPrize(clawMachine ClawMachine, goalOffset int) int {
	bestCombo := 0
	currentGoal := Point{x: clawMachine.prize.x + goalOffset + clawMachine.buttonB.dx, y: clawMachine.prize.y + goalOffset + clawMachine.buttonB.dy}

	for bPress := 0; currentGoal.x > 0 && currentGoal.y > 0; bPress++ {
		currentGoal = Point{currentGoal.x - clawMachine.buttonB.dx, currentGoal.y - clawMachine.buttonB.dy}

		if currentGoal.x%clawMachine.buttonA.dx != 0 || currentGoal.y%clawMachine.buttonA.dy != 0 {
			continue
		}

		if currentGoal.x/clawMachine.buttonA.dx != currentGoal.y/clawMachine.buttonA.dy {
			continue
		}

		aPress := currentGoal.x / clawMachine.buttonA.dx
		current := aPress*3 + bPress
		if bestCombo == 0 || current < bestCombo {
			bestCombo = current
		}
	}

	return bestCombo
}

func WinPrizes(clawMachines []ClawMachine, goalOffset int) int {
	return fn.Reduce(clawMachines, 0, func(_, sum int, machine ClawMachine) int { return sum + winPrize(machine, goalOffset) })
}

func ParseInputData(data string) ([]ClawMachine, error) {
	reA := regexp.MustCompile(`Button A: X\+(\d+), Y\+(\d+)`)
	reB := regexp.MustCompile(`Button B: X\+(\d+), Y\+(\d+)`)
	rePrize := regexp.MustCompile(`Prize: X\=(\d+), Y\=(\d+)`)

	inputLines := fn.GetLines(data)
	machines := make([]ClawMachine, 0, len(inputLines)/3)

	for i := 0; i < len(inputLines); i++ {
		if i+3 > len(inputLines) {
			return nil, errors.New("bad input")
		}
		matches := reA.FindAllStringSubmatch(inputLines[i], -1)
		if len(matches) != 1 && len(matches[0]) != 3 {
			return nil, errors.New("bad input")
		}
		buttonA := Direction{dx: MustAtoi(matches[0][1]), dy: MustAtoi(matches[0][2])}
		i++

		matches = reB.FindAllStringSubmatch(inputLines[i], -1)
		if len(matches) != 1 && len(matches[0]) != 3 {
			return nil, errors.New("bad input")
		}
		buttonB := Direction{dx: MustAtoi(matches[0][1]), dy: MustAtoi(matches[0][2])}
		i++

		matches = rePrize.FindAllStringSubmatch(inputLines[i], -1)
		if len(matches) != 1 && len(matches[0]) != 3 {
			return nil, errors.New("bad input")
		}
		prize := Point{x: MustAtoi(matches[0][1]), y: MustAtoi(matches[0][2])}
		i++

		machines = append(machines, ClawMachine{buttonA: buttonA, buttonB: buttonB, prize: prize})
	}

	return machines, nil
}

func main() {
	inputData, err := os.ReadFile("day13.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	machines, err := ParseInputData(string(inputData))
	if err != nil {
		fmt.Printf("Error parsing input data: %v\n", err)
		return
	}

	fmt.Printf("Part 1: %d\n", WinPrizes(machines, 0))
	fmt.Printf("Part 2: %d\n", WinPrizes(machines, 10000000000000))
}
