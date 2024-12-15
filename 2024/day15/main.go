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
type Wall struct{}
type Empty struct{}

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
			{
				currentPos = Point{x: currentPos.x + dir.dx, y: currentPos.y + dir.dy}
			}
		case Wall:
			{
				return robot
			}
		}
	}
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
			{
				robot = moveBox(warehouse, robot, dir)
			}
		}
	}

	gps := 0
	for point, object := range warehouse {
		switch object.(type) {
		case Box:
			{
				gps += point.x + point.y*100
			}
		}
	}
	return gps
}

func ParseInputData(data string) (warehouse map[Point]interface{}, robot Point, directions []Direction, err error) {
	lines := fn.GetLines(data)
	split := slices.Index(lines, "")
	grid := fn.MustTransform(lines[:split], func(line string) []byte { return []byte(line) })

	gridObjectMap := map[byte]interface{}{'O': Box{}, '#': Wall{}}

	warehouse = make(map[Point]interface{})
	errorPoint := Point{x: -1, y: -1}
	robot = errorPoint
	for y, row := range grid {
		for x, cell := range row {
			if cell == '@' {
				robot = Point{x: x, y: y}
			}
			if object, ok := gridObjectMap[cell]; ok {
				warehouse[Point{x: x, y: y}] = object
			} else {
				warehouse[Point{x: x, y: y}] = Empty{}
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

	warehouse, robot, directions, err := ParseInputData(string(inputData))
	if err != nil {
		fmt.Printf("Error parsing input data: %v\n", err)
		return
	}

	fmt.Printf("Part 1: %d\n", MoveRobot(warehouse, robot, directions))
	//fmt.Printf("Part 2: %d\n", FindChristmasTree(robots, 101, 103))
}
