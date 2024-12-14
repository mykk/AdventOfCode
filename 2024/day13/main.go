package main

import (
	fn "AoC/functional"
	"errors"
	"regexp"
	"strconv"

	"fmt"
	"os"
)

func MustAtoi(s string) int64 {
	val, _ := strconv.Atoi(s)
	return int64(val)
}

type Direction struct {
	dx, dy int64
}

type Point struct {
	x, y int64
}

type ClawMachine struct {
	a     Direction
	b     Direction
	prize Point
}

func winPrize(clawMachine ClawMachine, goalOffset int64) int64 {
	goal := Point{x: clawMachine.prize.x + goalOffset, y: clawMachine.prize.y + goalOffset}

	xa := clawMachine.a.dx
	xb := clawMachine.b.dx

	ya := clawMachine.a.dy
	yb := clawMachine.b.dy

	if (goal.x*yb-goal.y*xb)%(xa*yb-ya*xb) != 0 {
		return 0
	}
	a := (goal.x*yb - goal.y*xb) / (xa*yb - ya*xb)

	if (goal.x-a*xa)%xb != 0 {
		return 0
	}
	b := (goal.x - a*xa) / xb

	return a*3 + b
}

func WinPrizes(clawMachines []ClawMachine, goalOffset int64) int64 {
	return fn.Reduce(clawMachines, 0, func(_ int, sum int64, machine ClawMachine) int64 { return sum + winPrize(machine, goalOffset) })
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
		a := Direction{dx: MustAtoi(matches[0][1]), dy: MustAtoi(matches[0][2])}
		i++

		matches = reB.FindAllStringSubmatch(inputLines[i], -1)
		if len(matches) != 1 && len(matches[0]) != 3 {
			return nil, errors.New("bad input")
		}
		b := Direction{dx: MustAtoi(matches[0][1]), dy: MustAtoi(matches[0][2])}
		i++

		matches = rePrize.FindAllStringSubmatch(inputLines[i], -1)
		if len(matches) != 1 && len(matches[0]) != 3 {
			return nil, errors.New("bad input")
		}
		prize := Point{x: MustAtoi(matches[0][1]), y: MustAtoi(matches[0][2])}
		i++

		machines = append(machines, ClawMachine{a: a, b: b, prize: prize})
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
