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
	walk(grid [][]byte, id byte, startPoint Point, outside aoc.Set[Point]) (endPoint Point, nextWalker perimeterWalker)
}

type eastWalker struct{}
type westWalker struct{}
type northWalker struct{}
type southWalker struct{}

func (eastWalker) walk(grid [][]byte, id byte, startPoint Point, outside aoc.Set[Point]) (endPoint Point, nextWalker perimeterWalker) {
	if id != grid[startPoint.Y][startPoint.X] {
		panic("starting position id should match given id")
	}

	currentPoint := startPoint

	if currentPoint.X == 0 || grid[currentPoint.Y][currentPoint.X-1] != id {
		outside.Add(Point{X: currentPoint.X - 1, Y: currentPoint.Y})
	}
	outside.Add(Point{X: currentPoint.X, Y: currentPoint.Y - 1})

	for {
		if currentPoint.X == len(grid[currentPoint.Y])-1 || grid[currentPoint.Y][currentPoint.X+1] != id {
			currentPoint.X++

			outside.Add(Point{X: currentPoint.X, Y: currentPoint.Y})
			return currentPoint, southWalker{}
		}

		currentPoint.X++

		if currentPoint.Y > 0 && grid[currentPoint.Y-1][currentPoint.X] == id {
			return currentPoint, northWalker{}
		} else {
			outside.Add(Point{X: currentPoint.X, Y: currentPoint.Y - 1})
		}
	}
}

func (westWalker) walk(grid [][]byte, id byte, startPoint Point, outside aoc.Set[Point]) (endPoint Point, nextWalker perimeterWalker) {
	if id != grid[startPoint.Y-1][startPoint.X-1] {
		panic("starting position id should match given id")
	}

	currentPoint := startPoint

	if currentPoint.X == len(grid[startPoint.Y-1]) || grid[startPoint.Y-1][startPoint.X] != id {
		outside.Add(Point{X: currentPoint.X, Y: startPoint.Y - 1})
	}

	outside.Add(Point{X: currentPoint.X - 1, Y: currentPoint.Y})

	for {
		currentPoint.X--

		if currentPoint.X == 0 || grid[currentPoint.Y-1][currentPoint.X-1] != id {
			outside.Add(Point{X: currentPoint.X - 1, Y: currentPoint.Y - 1})
			return currentPoint, northWalker{}
		}

		if currentPoint.X != 0 && currentPoint.Y < len(grid) && grid[currentPoint.Y][currentPoint.X-1] == id {
			return currentPoint, southWalker{}
		} else {
			outside.Add(Point{X: currentPoint.X - 1, Y: currentPoint.Y})
		}
	}
}

func (northWalker) walk(grid [][]byte, id byte, startPoint Point, outside aoc.Set[Point]) (endPoint Point, nextWalker perimeterWalker) {
	if id != grid[startPoint.Y-1][startPoint.X] {
		panic("starting position id should match given id")
	}

	currentPoint := startPoint

	outside.Add(Point{X: currentPoint.X - 1, Y: currentPoint.Y})
	if currentPoint.Y == len(grid) || grid[currentPoint.Y][currentPoint.X] != id {
		outside.Add(Point{X: currentPoint.X, Y: currentPoint.Y + 1})
	}

	for {
		currentPoint.Y--

		if currentPoint.Y == 0 || grid[currentPoint.Y-1][currentPoint.X] != id {
			outside.Add(Point{X: currentPoint.X, Y: currentPoint.Y - 1})
			return currentPoint, eastWalker{}
		}

		if currentPoint.Y > 0 && currentPoint.X > 0 && grid[currentPoint.Y-1][currentPoint.X-1] == id {
			return currentPoint, westWalker{}
		} else {
			outside.Add(Point{X: currentPoint.X - 1, Y: currentPoint.Y})
		}
	}
}

func (southWalker) walk(grid [][]byte, id byte, startPoint Point, outside aoc.Set[Point]) (endPoint Point, nextWalker perimeterWalker) {
	if id != grid[startPoint.Y][startPoint.X-1] {
		panic("starting position id should match given id")
	}

	currentPoint := startPoint

	outside.Add(currentPoint)
	for {
		currentPoint.Y++

		if currentPoint.Y == len(grid) || grid[currentPoint.Y][currentPoint.X-1] != id {
			return currentPoint, westWalker{}
		}

		if currentPoint.Y < len(grid) && currentPoint.X < len(grid[currentPoint.Y]) && grid[currentPoint.Y][currentPoint.X] == id {
			return currentPoint, eastWalker{}
		} else {
			outside.Add(currentPoint)
		}
	}
}

func walkToEastStart(id byte, startPoint Point, grid [][]byte, area aoc.Set[Point]) Point {
	for point := range area {
		if point.Y < startPoint.Y || point.Y == startPoint.Y && point.X < startPoint.X {
			startPoint = point
		}
	}
	return startPoint
}

func WalkPerimeterCollectOutside(startPoint Point, grid [][]byte, area aoc.Set[Point], outside aoc.Set[Point]) (perimeter []Point) {
	id := grid[startPoint.Y][startPoint.X]

	var currentWalker perimeterWalker = eastWalker{}
	startPoint = walkToEastStart(id, startPoint, grid, area)

	currentPoint, currentWalker := currentWalker.walk(grid, id, startPoint, outside)
	perimeter = append(perimeter, startPoint)
	perimeter = append(perimeter, currentPoint)

	for currentPoint != startPoint {
		currentPoint, currentWalker = currentWalker.walk(grid, id, currentPoint, outside)
		perimeter = append(perimeter, currentPoint)
	}
	return

}

func constructArea(id byte, point Point, grid [][]byte, area aoc.Set[Point]) aoc.Set[Point] {
	if id != grid[point.Y][point.X] {
		panic("position id should match given id")
	}
	directions := []Direction{{-1, 0}, {1, 0}, {0, 1}, {0, -1}}

	area.Add(point)

	for _, direction := range directions {
		nextPoint := Point{point.X + direction.DX, point.Y + direction.DY}

		if area.Contains(nextPoint) || !withinBounds(grid, nextPoint) || grid[nextPoint.Y][nextPoint.X] != id {
			continue
		}
		area = constructArea(id, nextPoint, grid, area)
	}

	return area
}

func constructHoles(id byte, point Point, grid [][]byte, area aoc.Set[Point], walked aoc.Set[Point], holes []Hole, outside aoc.Set[Point]) []Hole {
	if id != grid[point.Y][point.X] {
		panic("position id should match given id")
	}
	directions := []Direction{{-1, 0}, {1, 0}, {0, 1}, {0, -1}}

	walked.Add(point)
	for _, direction := range directions {
		nextPoint := Point{point.X + direction.DX, point.Y + direction.DY}
		if walked.Contains(nextPoint) || !withinBounds(grid, nextPoint) {
			continue
		}
		if grid[nextPoint.Y][nextPoint.X] == id {
			holes = constructHoles(id, nextPoint, grid, area, walked, holes, outside)
		} else if !outside.Contains(nextPoint) && !fn.Any(holes, func(_ int, hole Hole) bool { return hole.Area.Contains(nextPoint) }) {
			holeArea := constructArea(grid[nextPoint.Y][nextPoint.X], nextPoint, grid, make(aoc.Set[Point]))
			holePerimeter := WalkPerimeterCollectOutside(nextPoint, grid, holeArea, make(aoc.Set[Point]))
			holes = append(holes, Hole{id: grid[nextPoint.Y][nextPoint.X], Area: holeArea, Perimeter: holePerimeter})
		}
	}
	return holes
}

func getInsideParameters(grid [][]byte, holes []Hole) [][]Point {
	newGrid := fn.MustTransform(grid, func(row []byte) []byte { return fn.MustTransform(row, func(byte) byte { return '.' }) })

	for _, hole := range holes {
		for point := range hole.Area {
			newGrid[point.Y][point.X] = 'X'
		}
	}

	perimeters := make([][]Point, 0)

	for _, hole := range holes {
		area := constructArea('X', hole.Perimeter[0], newGrid, make(aoc.Set[Point]))
		newPerimeter := WalkPerimeterCollectOutside(hole.Perimeter[0], newGrid, area, make(aoc.Set[Point]))
		if !fn.Any(perimeters, func(_ int, perimeter []Point) bool {
			if len(newPerimeter) != len(perimeter) {
				return false
			}
			return fn.All(perimeter, func(i int, point Point) bool { return newPerimeter[i] == point })
		}) {
			perimeters = append(perimeters, newPerimeter)
		}
	}

	return perimeters
}

func constructRegion(startPoint Point, grid [][]byte) Region {
	id := grid[startPoint.Y][startPoint.X]

	area := constructArea(id, startPoint, grid, make(aoc.Set[Point]))

	outside := make(aoc.Set[Point])
	perimeter := WalkPerimeterCollectOutside(startPoint, grid, area, outside)

	holes := constructHoles(id, startPoint, grid, area, make(aoc.Set[Point]), []Hole{}, outside)

	return Region{id: id, Perimeter: perimeter, InsidePerimeters: getInsideParameters(grid, holes), Area: area, Holes: holes}
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
