package main

import (
	fn "AoC/functional"

	"fmt"
	"os"
)

type Point struct {
	x, y int
}

type Direction struct {
	dx, dy int
}

func withinBounds(grid [][]byte, pos Point) bool {
	return pos.y >= 0 && pos.y < len(grid) && pos.x >= 0 && pos.x < len(grid[pos.y])
}

type Area struct {
	id        byte
	perimeter []Point
	area      []Point
	holes     [][]Point
}

type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(value T) {
	s[value] = struct{}{}
}

func (s Set[T]) Contains(value T) bool {
	_, exists := s[value]
	return exists
}

type IAreaWalker interface {
	Walk(grid [][]byte, id byte, startPoint Point) (endPoint Point, nextWalker IAreaWalker)
}

func Walk(point Point, grid [][]byte) Area {
	id := grid[point.y][point.x]
	area := Area{id: id, perimeter: []Point{}, area: []Point{}, holes: [][]Point{}}

	var currentWalker IAreaWalker = EastWalker{}
	startPoint := point
	for {
		point, currentWalker = currentWalker.Walk(grid, id, point)
		area.perimeter = append(area.perimeter, point)
		if point == startPoint {
			break
		}
	}

	return area
}

type EastWalker struct {
}

type WestWalker struct {
}

type NorthWalker struct {
}

type SouthWalker struct {
}

func (EastWalker) Walk(grid [][]byte, id byte, startPoint Point) (endPoint Point, nextWalker IAreaWalker) {
	currentPoint := startPoint

	for {
		if currentPoint.y > 0 && grid[currentPoint.y-1][currentPoint.x] == id {
			return currentPoint, NorthWalker{}
		}

		currentPoint.x++

		if len(grid[currentPoint.y])-1 > currentPoint.x && grid[currentPoint.y][currentPoint.x+1] == id {
			continue
		}

		return currentPoint, SouthWalker{}
	}
}

func (NorthWalker) Walk(grid [][]byte, id byte, startPoint Point) (endPoint Point, nextWalker IAreaWalker) {
	currentPoint := startPoint

	for {
		if currentPoint.x > 0 && grid[currentPoint.y][currentPoint.x-1] == id {
			return currentPoint, WestWalker{}
		}

		currentPoint.y--
		if currentPoint.y > 0 && grid[currentPoint.y-1][currentPoint.x] == id {
			continue
		}
		return currentPoint, EastWalker{}
	}
}

func (SouthWalker) Walk(grid [][]byte, id byte, startPoint Point) (endPoint Point, nextWalker IAreaWalker) {
	currentPoint := startPoint

	for {
		if currentPoint.x > 0 && grid[currentPoint.y][currentPoint.x-1] == id {
			return currentPoint, EastWalker{}
		}

		currentPoint.y++

		if currentPoint.y > 0 && grid[currentPoint.y-1][currentPoint.x] == id {
			continue
		}
		return currentPoint, WestWalker{}
	}
}

func constructArea(x, y int, garden [][]byte, usedCoordinates Set[Point]) (area Area) {
	identifier := garden[y][x]
	polygon = append(polygon, Point{x: x, y: y})

	directions := []Direction{{1, 0}, {0, -1}, {0, 1}, {1, 0}}

	currentPoint := Point{x: x, y: y}
	for _, dir := range directions {
		nextPosition := Point{x: currentPoint.x + dir.dx, y: currentPoint.y + dir.dy}
		if withinBounds(garden, nextPosition) && garden[nextPosition.y][nextPosition.x] == identifier {

		}
	}
	return
}

func constructAreas(garden [][]byte) (areas []Area) {
	usedCoordinates := make(Set[Point])
	for y, row := range garden {
		for x, _ := range row {
			if usedCoordinates.Contains(Point{x: x, y: y}) {
				continue
			}
			areas = append(areas, constructArea(x, y, garden, usedCoordinates))
		}
	}

	return
}

func ParseInputData(data string) []Area {
	garden := fn.MustTransform(fn.GetLines(data), func(line string) []byte { return []byte(line) })
	return constructAreas(garden)
}

func main() {
	inputData, err := os.ReadFile("day12.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	areas = ParseInputData(string(inputData))
	if err != nil {
		fmt.Printf("Error parsing input data: %v\n", err)
		return
	}

	//fmt.Printf("Part 1: %d\n", CountPlutonianPebbles(stones, 25))
	//fmt.Printf("Part 2: %d\n", CountPlutonianPebbles(stones, 75))
}
