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

type Robot struct {
	p Point
	v Direction
}

func findQuarterAmount(robots []Robot, x, y int, predicate func(robot Robot, xSplit, ySplit int) bool) int {
	xSplit := (x - 1) / 2
	ySplit := (y - 1) / 2

	return fn.Reduce(robots, 0, func(_, sum int, robot Robot) int {
		if predicate(robot, xSplit, ySplit) {
			return sum + 1
		}
		return sum
	})
}

func calculateSafetyFactor(robots []Robot, x, y int) int {
	q1 := findQuarterAmount(robots, x, y, func(robot Robot, xSplit, ySplit int) bool { return robot.p.x < xSplit && robot.p.y < ySplit })
	q2 := findQuarterAmount(robots, x, y, func(robot Robot, xSplit, ySplit int) bool { return robot.p.x > xSplit && robot.p.y < ySplit })
	q3 := findQuarterAmount(robots, x, y, func(robot Robot, xSplit, ySplit int) bool { return robot.p.x < xSplit && robot.p.y > ySplit })
	q4 := findQuarterAmount(robots, x, y, func(robot Robot, xSplit, ySplit int) bool { return robot.p.x > xSplit && robot.p.y > ySplit })

	return q1 * q2 * q3 * q4
}

func getMovedRobots(robots []Robot, t, x, y int) []Robot {
	return fn.MustTransform(robots, func(robot Robot) Robot {
		position := Point{x: (robot.p.x + robot.v.dx*t) % x, y: (robot.p.y + robot.v.dy*t) % y}
		if position.x < 0 {
			position.x += x
		}

		if position.y < 0 {
			position.y += y
		}

		return Robot{p: position, v: robot.v}
	})
}

func MoveRobots(robots []Robot, t, x, y int) int {
	movedRobots := getMovedRobots(robots, t, x, y)
	return calculateSafetyFactor(movedRobots, x, y)
}

func christmasTreeReppresent(robots []Robot) bool {
	trunk := make([]int, 0, 20)
	for i := 1; i <= 20; i++ {
		trunk = append(trunk, i)
	}

	return fn.Any(robots, func(_ int, robot Robot) bool {
		return fn.All(trunk, func(_, i int) bool {
			return fn.Any(robots, func(_ int, other Robot) bool {
				return robot.p.x == other.p.x && robot.p.y+i == other.p.y
			})
		})
	})
}

func printTheThree(robots []Robot, x, y int) {
	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			sum := fn.Reduce(robots, 0, func(_, sum int, robot Robot) int {
				if robot.p.x == j && robot.p.y == i {
					return sum + 1
				}
				return sum
			})
			fmt.Printf("%d", sum)
		}
		fmt.Printf("\n")
	}
}

func FindChristmasTree(robots []Robot, x, y int) int {
	for i := 0; i < 1000000; i++ {
		movedRobots := getMovedRobots(robots, i, x, y)
		if christmasTreeReppresent(movedRobots) {
			printTheThree(movedRobots, x, y)
			return i
		}
	}

	panic("reached the unreachable")
}

func ParseInputData(data string) ([]Robot, error) {
	re := regexp.MustCompile(`p\=(\d+),(\d+) v=(-?\d+),(-?\d+)`)

	return fn.Transform(fn.GetLines(data), func(line string) (Robot, error) {
		matches := re.FindAllStringSubmatch(line, -1)
		if len(matches) != 1 {
			return Robot{}, errors.New("bad input")
		}
		p := Point{x: MustAtoi(matches[0][1]), y: MustAtoi(matches[0][2])}
		v := Direction{dx: MustAtoi(matches[0][3]), dy: MustAtoi(matches[0][4])}
		return Robot{p: p, v: v}, nil
	})
}

func main() {
	inputData, err := os.ReadFile("day14.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	robots, err := ParseInputData(string(inputData))
	if err != nil {
		fmt.Printf("Error parsing input data: %v\n", err)
		return
	}

	fmt.Printf("Part 1: %d\n", MoveRobots(robots, 100, 101, 103))
	fmt.Printf("Part 2: %d\n", FindChristmasTree(robots, 101, 103))
}
