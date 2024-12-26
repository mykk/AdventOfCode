package aoc_geometry

import (
	"AoC/aoc"
	fn "AoC/functional"
)

type Point struct {
	X, Y int
}

type Direction struct {
	DX, DY int
}

var Directions = []Direction{{-1, 0}, {1, 0}, {0, 1}, {0, -1}}

type Hole struct {
	id        byte
	Perimeter []Point
	Area      aoc.Set[Point]
}

type Region struct {
	id               byte
	Perimeter        []Point
	InsidePerimeters [][]Point
	Area             aoc.Set[Point]
	Holes            []Hole
}

func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func withinBounds[T any](grid [][]T, pos Point) bool {
	return pos.Y >= 0 && pos.Y < len(grid) && pos.X >= 0 && pos.X < len(grid[pos.Y])
}

func getPerimeter(polygon []Point) int {
	return fn.Reduce(polygon[:len(polygon)-1], 0, func(i, perimeter int, point Point) int {
		nextPoint := polygon[i+1]
		return perimeter + absInt(point.X-nextPoint.X) + absInt(point.Y-nextPoint.Y)
	})
}

func (region *Region) GetOutsidePerimeter() int {
	return getPerimeter(region.Perimeter)
}

func (region *Region) GetInsidePerimeter() int {
	return fn.Reduce(region.InsidePerimeters, 0, func(_, perimeter int, hole []Point) int {
		return perimeter + getPerimeter(hole)
	})
}

func (hole *Hole) GetPerimeter() int {
	return getPerimeter(hole.Perimeter)
}

type perimeterWalker interface {
	walk(grid [][]byte, startPoint Point, outside aoc.Set[Point], checkId func(id byte) bool) (endPoint Point, nextWalker perimeterWalker)
}

type eastWalker struct{}
type westWalker struct{}
type northWalker struct{}
type southWalker struct{}

func (eastWalker) walk(grid [][]byte, startPoint Point, outside aoc.Set[Point], checkId func(id byte) bool) (endPoint Point, nextWalker perimeterWalker) {
	if !checkId(grid[startPoint.Y][startPoint.X]) {
		panic("starting position id should match given id")
	}

	currentPoint := startPoint

	if currentPoint.X == 0 || !checkId(grid[currentPoint.Y][currentPoint.X-1]) {
		outside.Add(Point{X: currentPoint.X - 1, Y: currentPoint.Y})
	}
	outside.Add(Point{X: currentPoint.X, Y: currentPoint.Y - 1})

	for {
		if currentPoint.X == len(grid[currentPoint.Y])-1 || !checkId(grid[currentPoint.Y][currentPoint.X+1]) {
			currentPoint.X++

			outside.Add(Point{X: currentPoint.X, Y: currentPoint.Y})
			return currentPoint, southWalker{}
		}

		currentPoint.X++

		if currentPoint.Y > 0 && checkId(grid[currentPoint.Y-1][currentPoint.X]) {
			return currentPoint, northWalker{}
		}

		outside.Add(Point{X: currentPoint.X, Y: currentPoint.Y - 1})
	}
}

func (westWalker) walk(grid [][]byte, startPoint Point, outside aoc.Set[Point], checkId func(id byte) bool) (endPoint Point, nextWalker perimeterWalker) {
	if !checkId(grid[startPoint.Y-1][startPoint.X-1]) {
		panic("starting position id should match given id")
	}

	currentPoint := startPoint

	if currentPoint.X == len(grid[startPoint.Y-1]) || !checkId(grid[startPoint.Y-1][startPoint.X]) {
		outside.Add(Point{X: currentPoint.X, Y: startPoint.Y - 1})
	}
	outside.Add(Point{X: currentPoint.X - 1, Y: currentPoint.Y})

	for {
		currentPoint.X--

		if currentPoint.X == 0 || !checkId(grid[currentPoint.Y-1][currentPoint.X-1]) {
			outside.Add(Point{X: currentPoint.X - 1, Y: currentPoint.Y - 1})
			return currentPoint, northWalker{}
		}

		if currentPoint.X != 0 && currentPoint.Y < len(grid) && checkId(grid[currentPoint.Y][currentPoint.X-1]) {
			return currentPoint, southWalker{}
		}
		outside.Add(Point{X: currentPoint.X - 1, Y: currentPoint.Y})
	}
}

func (northWalker) walk(grid [][]byte, startPoint Point, outside aoc.Set[Point], checkId func(id byte) bool) (endPoint Point, nextWalker perimeterWalker) {
	if !checkId(grid[startPoint.Y-1][startPoint.X]) {
		panic("starting position id should match given id")
	}

	currentPoint := startPoint

	if currentPoint.Y == len(grid) || !checkId(grid[currentPoint.Y][currentPoint.X]) {
		outside.Add(Point{X: currentPoint.X, Y: currentPoint.Y + 1})
	}
	outside.Add(Point{X: currentPoint.X - 1, Y: currentPoint.Y})

	for {
		currentPoint.Y--

		if currentPoint.Y == 0 || !checkId(grid[currentPoint.Y-1][currentPoint.X]) {
			outside.Add(Point{X: currentPoint.X, Y: currentPoint.Y - 1})
			return currentPoint, eastWalker{}
		}

		if currentPoint.Y > 0 && currentPoint.X > 0 && checkId(grid[currentPoint.Y-1][currentPoint.X-1]) {
			return currentPoint, westWalker{}
		}
		outside.Add(Point{X: currentPoint.X - 1, Y: currentPoint.Y})
	}
}

func (southWalker) walk(grid [][]byte, startPoint Point, outside aoc.Set[Point], checkId func(id byte) bool) (endPoint Point, nextWalker perimeterWalker) {
	if !checkId(grid[startPoint.Y][startPoint.X-1]) {
		panic("starting position id should match given id")
	}

	currentPoint := startPoint

	outside.Add(currentPoint)

	for {
		currentPoint.Y++

		if currentPoint.Y == len(grid) || !checkId(grid[currentPoint.Y][currentPoint.X-1]) {
			return currentPoint, westWalker{}
		}

		if currentPoint.Y < len(grid) && currentPoint.X < len(grid[currentPoint.Y]) && checkId(grid[currentPoint.Y][currentPoint.X]) {
			return currentPoint, eastWalker{}
		}

		outside.Add(currentPoint)
	}
}

func walkToStart(startPoint Point, grid [][]byte, area aoc.Set[Point]) Point {
	for point := range area {
		if point.Y < startPoint.Y || point.Y == startPoint.Y && point.X < startPoint.X {
			startPoint = point
		}
	}
	return startPoint
}

func constructPerimeterAndOutside(startPoint Point, grid [][]byte, area aoc.Set[Point], checkId func(byte) bool) (perimeter []Point, outside aoc.Set[Point]) {
	outside = make(aoc.Set[Point])

	startPoint = walkToStart(startPoint, grid, area)

	var currentWalker perimeterWalker = eastWalker{}
	currentPoint, currentWalker := currentWalker.walk(grid, startPoint, outside, checkId)
	perimeter = append(perimeter, startPoint)
	perimeter = append(perimeter, currentPoint)

	for currentPoint != startPoint {
		currentPoint, currentWalker = currentWalker.walk(grid, currentPoint, outside, checkId)
		perimeter = append(perimeter, currentPoint)
	}
	return
}

func constructCurrentIdPerimeterAndOutside(startPoint Point, grid [][]byte, area aoc.Set[Point]) (perimeter []Point, outside aoc.Set[Point]) {
	id := grid[startPoint.Y][startPoint.X]
	return constructPerimeterAndOutside(startPoint, grid, area, func(currentId byte) bool { return id == currentId })
}

func constructArea(point Point, grid [][]byte, area aoc.Set[Point], checkId func(byte) bool) aoc.Set[Point] {
	if !checkId(grid[point.Y][point.X]) {
		panic("position id should match given id")
	}

	area.Add(point)

	for _, direction := range Directions {
		nextPoint := Point{point.X + direction.DX, point.Y + direction.DY}

		if area.Contains(nextPoint) || !withinBounds(grid, nextPoint) || !checkId(grid[nextPoint.Y][nextPoint.X]) {
			continue
		}
		area = constructArea(nextPoint, grid, area, checkId)
	}

	return area
}

func constructCurrentIdArea(id byte, point Point, grid [][]byte, area aoc.Set[Point]) aoc.Set[Point] {
	return constructArea(point, grid, area, func(currentId byte) bool { return currentId == id })
}

func constructHoles(id byte, point Point, grid [][]byte, area aoc.Set[Point], walked aoc.Set[Point], holes []Hole, outside aoc.Set[Point]) []Hole {
	if id != grid[point.Y][point.X] {
		panic("position id should match given id")
	}

	walked.Add(point)
	for _, direction := range Directions {
		nextPoint := Point{point.X + direction.DX, point.Y + direction.DY}
		if walked.Contains(nextPoint) || !withinBounds(grid, nextPoint) {
			continue
		}

		if grid[nextPoint.Y][nextPoint.X] == id {
			holes = constructHoles(id, nextPoint, grid, area, walked, holes, outside)
		} else if !outside.Contains(nextPoint) && fn.All(holes, func(_ int, hole Hole) bool { return !hole.Area.Contains(nextPoint) }) {
			holeArea := constructCurrentIdArea(grid[nextPoint.Y][nextPoint.X], nextPoint, grid, make(aoc.Set[Point]))
			holePerimeter, _ := constructCurrentIdPerimeterAndOutside(nextPoint, grid, holeArea)
			holes = append(holes, Hole{id: grid[nextPoint.Y][nextPoint.X], Area: holeArea, Perimeter: holePerimeter})
		}
	}
	return holes
}

func getInsideParameters(id byte, grid [][]byte, holes []Hole) (insidePerimeters [][]Point) {
	coveredAreas := make(aoc.Set[Point])

	for _, hole := range holes {
		if coveredAreas.Contains(hole.Perimeter[0]) {
			continue
		}

		joinedHoleArea := constructArea(hole.Perimeter[0], grid, make(aoc.Set[Point]), func(currentId byte) bool { return currentId != id })
		holePerimeter, _ := constructPerimeterAndOutside(hole.Perimeter[0], grid, joinedHoleArea, func(currentId byte) bool { return currentId != id })
		insidePerimeters = append(insidePerimeters, holePerimeter)

		for areaPoint := range joinedHoleArea {
			coveredAreas.Add(areaPoint)
		}
	}
	return
}

func constructRegion(startPoint Point, grid [][]byte) Region {
	id := grid[startPoint.Y][startPoint.X]

	area := constructCurrentIdArea(id, startPoint, grid, make(aoc.Set[Point]))
	perimeter, outside := constructCurrentIdPerimeterAndOutside(startPoint, grid, area)
	holes := constructHoles(id, startPoint, grid, area, make(aoc.Set[Point]), []Hole{}, outside)
	insidePerimeters := getInsideParameters(id, grid, holes)

	return Region{id: id, Perimeter: perimeter, InsidePerimeters: insidePerimeters, Area: area, Holes: holes}
}

func constructRegionAndAppendUsed(point Point, grid [][]byte, usedCoordinates aoc.Set[Point]) Region {
	region := constructRegion(point, grid)

	for point := range region.Area {
		usedCoordinates.Add(point)
	}

	for _, hole := range region.Holes {
		for point := range hole.Area {
			usedCoordinates.Add(point)
		}
	}
	return region
}

func RegionsFromGrid(grid [][]byte) (regions []Region) {
	usedCoordinates := make(aoc.Set[Point])
	for y, row := range grid {
		for x := range row {
			point := Point{X: x, Y: y}
			if usedCoordinates.Contains(point) {
				continue
			}
			regions = append(regions, constructRegionAndAppendUsed(point, grid, usedCoordinates))
		}
	}

	return
}
