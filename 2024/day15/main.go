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

func moveDoubleBoxVertically(warehouse map[Point]interface{}, robot Point, dir Direction) Point {
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
			panic("unexpected object type in warehouse")
		}
	}
}

func getDoubleBoxStartEndInHorizontalDir(warehouse map[Point]interface{}, point Point, dir Direction) (Point, Point) {
	if _, ok := warehouse[Point{x: point.x, y: point.y + dir.dy}].(DoubleBoxStart); ok {
		return Point{x: point.x, y: point.y + dir.dy}, Point{x: point.x + 1, y: point.y + dir.dy}
	}
	if _, ok := warehouse[Point{x: point.x, y: point.y + dir.dy}].(DoubleBoxEnd); ok {
		return Point{x: point.x - 1, y: point.y + dir.dy}, Point{x: point.x, y: point.y + dir.dy}
	}

	panic("wrong box type!")
}

func doubleBoxMovableHorizontally(warehouse map[Point]interface{}, start, end Point, dir Direction) bool {
	for _, point := range []Point{start, end} {
		switch warehouse[Point{x: point.x, y: point.y + dir.dy}].(type) {
		case DoubleBoxStart, DoubleBoxEnd:
			{
				nextStart, nextEnd := getDoubleBoxStartEndInHorizontalDir(warehouse, point, dir)
				if !doubleBoxMovableHorizontally(warehouse, nextStart, nextEnd, dir) {
					return false
				}
			}
		case Wall:
			return false
		}
	}
	return true
}

func moveDoubleBoxHorizontally(warehouse map[Point]interface{}, point Point, dir Direction) {
	start, end := getDoubleBoxStartEndInHorizontalDir(warehouse, point, dir)

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
				moveDoubleBoxHorizontally(warehouse, currentPoint, dir)
				warehouse[nextPoint] = warehouse[currentPoint]
				warehouse[currentPoint] = Empty{}
			}
		default:
			panic("unexpected object")
		}
	}
}

func printMap(warehouse map[Point]interface{}) {
	for y := 0; y < 100; y++ {
		for x := 0; x < 200; x++ {
			switch warehouse[Point{x: x, y: y}].(type) {
			case Empty:
				{
					fmt.Print(".")
				}
			case DoubleBoxStart:
				{
					fmt.Print("[")
				}
			case DoubleBoxEnd:
				{
					fmt.Print("]")

				}
			case Wall:
				{
					fmt.Print("#")
				}
			default:
				break
			}
		}
		fmt.Println()
		if _, ok := warehouse[Point{x: 0, y: y}]; !ok {
			break
		}
	}
}
func tryMoveDoubleBoxHorizontally(warehouse map[Point]interface{}, robot Point, dir Direction) Point {
	start, end := getDoubleBoxStartEndInHorizontalDir(warehouse, robot, dir)
	if !doubleBoxMovableHorizontally(warehouse, start, end, dir) {
		return robot
	}

	moveDoubleBoxHorizontally(warehouse, robot, dir)
	printMap(warehouse)
	return Point{x: robot.x, y: robot.y + dir.dy}
}

func moveDoubleBox(warehouse map[Point]interface{}, robot Point, dir Direction) Point {
	left := Direction{dx: 1, dy: 0}
	right := Direction{dx: -1, dy: 0}

	if dir == left || dir == right {
		result := moveDoubleBoxVertically(warehouse, robot, dir)
		printMap(warehouse)
		return result
	}

	return tryMoveDoubleBoxHorizontally(warehouse, robot, dir)
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
			panic("unexpected warehouse unit type")
		}
	}
}

func countGps(warehouse map[Point]interface{}) int {
	gps := 0
	for point, object := range warehouse {
		if _, ok := object.(Box); ok {
			gps += point.x + point.y*100
		}
		if _, ok := object.(DoubleBoxStart); ok {
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
			panic("unexpected warehouse unit type")
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
	printMap(warehouse)

	if err != nil {
		fmt.Printf("Error parsing input data: %v\n", err)
		return
	}
	fmt.Printf("Part 2: %d\n", MoveRobot(warehouse, robot, directions))
}
