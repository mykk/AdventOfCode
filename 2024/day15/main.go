package main

import (
	fn "AoC/functional"
	"errors"
	"slices"

	"fmt"
	"os"
)

type Direction struct {
	dx, dy int
}

type Point struct {
	x, y int
}

type Box struct{}
type DoubleBoxStart struct{}
type DoubleBoxEnd struct{}
type Wall struct{}
type Empty struct{}

func moveDoubleBoxHorizontally(warehouse map[Point]interface{}, robot Point, dir Direction) Point {
	firstBox := Point{x: robot.x + dir.dx, y: robot.y}
	currentPos := firstBox

	boxCount := 0

	for {
		switch warehouse[currentPos].(type) {
		case Empty:
			{
				for boxCount != 0 {
					warehouse[currentPos] = warehouse[Point{x: currentPos.x - dir.dx, y: currentPos.y}]
					currentPos = Point{x: currentPos.x - dir.dx, y: currentPos.y}

					warehouse[currentPos] = warehouse[Point{x: currentPos.x - dir.dx, y: currentPos.y}]
					currentPos = Point{x: currentPos.x - dir.dx, y: currentPos.y}

					boxCount -= 2
				}
				warehouse[currentPos] = Empty{}
				return currentPos
			}
		case DoubleBoxStart, DoubleBoxEnd:
			{
				currentPos = Point{x: currentPos.x + dir.dx, y: currentPos.y + dir.dy}
				boxCount++
			}
		case Wall:
			return robot
		default:
			panic("unexpected warehouse object")
		}
	}
}

func getDoubleBoxStartEndInVerticalDir(warehouse map[Point]interface{}, point Point, dir Direction) (Point, Point) {
	if _, ok := warehouse[Point{x: point.x, y: point.y + dir.dy}].(DoubleBoxStart); ok {
		return Point{x: point.x, y: point.y + dir.dy}, Point{x: point.x + 1, y: point.y + dir.dy}
	}
	if _, ok := warehouse[Point{x: point.x, y: point.y + dir.dy}].(DoubleBoxEnd); ok {
		return Point{x: point.x - 1, y: point.y + dir.dy}, Point{x: point.x, y: point.y + dir.dy}
	}

	panic("wrong box type")
}

func doubleBoxMovableVertically(warehouse map[Point]interface{}, start, end Point, dir Direction) bool {
	for _, point := range []Point{start, end} {
		switch warehouse[Point{x: point.x, y: point.y + dir.dy}].(type) {
		case DoubleBoxStart, DoubleBoxEnd:
			{
				nextStart, nextEnd := getDoubleBoxStartEndInVerticalDir(warehouse, point, dir)
				if !doubleBoxMovableVertically(warehouse, nextStart, nextEnd, dir) {
					return false
				}
			}
		case Wall:
			return false
		}
	}
	return true
}

func moveDoubleBoxVertically(warehouse map[Point]interface{}, point Point, dir Direction) {
	start, end := getDoubleBoxStartEndInVerticalDir(warehouse, point, dir)

	for _, currentPoint := range []Point{start, end} {
		nextPoint := Point{x: currentPoint.x, y: currentPoint.y + dir.dy}

		switch warehouse[nextPoint].(type) {
		case Empty:
			{
				warehouse[nextPoint] = warehouse[currentPoint]
				warehouse[currentPoint] = Empty{}
			}
		case DoubleBoxStart, DoubleBoxEnd:
			{
				moveDoubleBoxVertically(warehouse, currentPoint, dir)
				warehouse[nextPoint] = warehouse[currentPoint]
				warehouse[currentPoint] = Empty{}
			}
		default:
			panic("unexpected warehouse object")
		}
	}
}

func tryMoveDoubleBoxVertically(warehouse map[Point]interface{}, robot Point, dir Direction) Point {
	start, end := getDoubleBoxStartEndInVerticalDir(warehouse, robot, dir)
	if !doubleBoxMovableVertically(warehouse, start, end, dir) {
		return robot
	}

	moveDoubleBoxVertically(warehouse, robot, dir)
	return Point{x: robot.x, y: robot.y + dir.dy}
}

func moveDoubleBox(warehouse map[Point]interface{}, robot Point, dir Direction) Point {
	left := Direction{dx: 1, dy: 0}
	right := Direction{dx: -1, dy: 0}

	if dir == left || dir == right {
		return moveDoubleBoxHorizontally(warehouse, robot, dir)
	}

	return tryMoveDoubleBoxVertically(warehouse, robot, dir)
}

func moveBox(warehouse map[Point]interface{}, robot Point, dir Direction) Point {
	firstBox := Point{x: robot.x + dir.dx, y: robot.y + dir.dy}
	currentPos := firstBox

	for {
		switch warehouse[currentPos].(type) {
		case Empty:
			{
				warehouse[firstBox] = Empty{}
				warehouse[currentPos] = Box{}
				return firstBox
			}
		case Box:
			currentPos = Point{x: currentPos.x + dir.dx, y: currentPos.y + dir.dy}
		case Wall:
			return robot
		default:
			panic("unexpected warehouse object")
		}
	}
}

func countGps(warehouse map[Point]interface{}) int {
	gps := 0
	for point, object := range warehouse {
		switch object.(type) {
		case Box, DoubleBoxStart:
			gps += point.x + point.y*100
		}
	}
	return gps
}

func MoveRobot(warehouse map[Point]interface{}, robot Point, directions []Direction) int {
	for _, dir := range directions {
		pos := Point{x: robot.x + dir.dx, y: robot.y + dir.dy}
		switch warehouse[pos].(type) {
		case Empty:
			robot = pos
		case Wall:
			{
			}
		case Box:
			robot = moveBox(warehouse, robot, dir)
		case DoubleBoxStart, DoubleBoxEnd:
			robot = moveDoubleBox(warehouse, robot, dir)
		default:
			panic("unexpected warehouse object")
		}
	}

	return countGps(warehouse)
}

func ParseInputData(data string, double bool) (warehouse map[Point]interface{}, robot Point, directions []Direction, err error) {
	lines := fn.GetLines(data)
	split := slices.Index(lines, "")
	grid := fn.MustTransform(lines[:split], func(line string) []byte { return []byte(line) })

	warehouse = make(map[Point]interface{})
	errorPoint := Point{x: -1, y: -1}
	robot = errorPoint
	for y, row := range grid {
		for x, cell := range row {
			point := Point{x: x, y: y}
			if double {
				point = Point{x: x * 2, y: y}
			}

			if cell == '@' {
				robot = point
			}
			if cell == 'O' {
				if !double {
					warehouse[point] = Box{}
				} else {
					warehouse[point] = DoubleBoxStart{}
					warehouse[Point{x: point.x + 1, y: point.y}] = DoubleBoxEnd{}
				}
			} else if cell == '#' {
				warehouse[point] = Wall{}
				if double {
					warehouse[Point{x: point.x + 1, y: point.y}] = Wall{}
				}
			} else {
				warehouse[point] = Empty{}
				if double {
					warehouse[Point{x: point.x + 1, y: point.y}] = Empty{}
				}
			}
		}
	}
	if robot == errorPoint {
		err = errors.New("robot position not found")
		return
	}

	directionMap := map[byte]Direction{
		'^': Direction{dx: 0, dy: -1},
		'v': Direction{dx: 0, dy: 1},
		'>': Direction{dx: 1, dy: 0},
		'<': Direction{dx: -1, dy: 0}}

	instructions := fn.Reduce(lines[split+1:], []byte{}, func(_ int, instr []byte, line string) []byte { return append(instr, []byte(line)...) })
	directions, err = fn.Transform(instructions, func(intruction byte) (Direction, error) {
		if dir, ok := directionMap[intruction]; ok {
			return dir, nil
		}
		return Direction{}, errors.New("uknown instruction")
	})

	return
}

func main() {
	inputData, err := os.ReadFile("day15.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	warehouse, robot, directions, err := ParseInputData(string(inputData), false)
	if err != nil {
		fmt.Printf("Error parsing input data: %v\n", err)
		return
	}
	fmt.Printf("Part 1: %d\n", MoveRobot(warehouse, robot, directions))

	warehouse, robot, directions, err = ParseInputData(string(inputData), true)
	if err != nil {
		fmt.Printf("Error parsing input data: %v\n", err)
		return
	}
	fmt.Printf("Part 2: %d\n", MoveRobot(warehouse, robot, directions))
}
