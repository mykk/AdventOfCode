package aoc_geometry

import (
	fn "AoC/functional"
)

type Point struct {
	X, Y int
}

type Direction struct {
	dx, dy int
}

func withinBounds[T any](grid [][]T, pos Point) bool {
	return pos.Y >= 0 && pos.Y < len(grid) && pos.X >= 0 && pos.X < len(grid[pos.Y])
}

type Hole struct {
	id        byte
	Perimeter []Point
	Area      Set[Point]
}

type Region struct {
	id               byte
	Perimeter        []Point
	InsidePerimeters [][]Point
	Area             Set[Point]
	Holes            []Hole
}

func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
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

type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(value T) {
	s[value] = struct{}{}
}

func (s Set[T]) Contains(value T) bool {
	_, exists := s[value]
	return exists
}

type PerimeterWalker interface {
	Walk(grid [][]byte, id byte, startPoint Point, outside Set[Point]) (endPoint Point, nextWalker PerimeterWalker)
}

type EastWalker struct{}
type WestWalker struct{}
type NorthWalker struct{}
type SouthWalker struct{}

func (EastWalker) Walk(grid [][]byte, id byte, startPoint Point, outside Set[Point]) (endPoint Point, nextWalker PerimeterWalker) {
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
			return currentPoint, SouthWalker{}
		}

		currentPoint.X++

		if currentPoint.Y > 0 && grid[currentPoint.Y-1][currentPoint.X] == id {
			return currentPoint, NorthWalker{}
		} else {
			outside.Add(Point{X: currentPoint.X, Y: currentPoint.Y - 1})
		}
	}
}

func (WestWalker) Walk(grid [][]byte, id byte, startPoint Point, outside Set[Point]) (endPoint Point, nextWalker PerimeterWalker) {
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
			return currentPoint, NorthWalker{}
		}

		if currentPoint.X != 0 && currentPoint.Y < len(grid) && grid[currentPoint.Y][currentPoint.X-1] == id {
			return currentPoint, SouthWalker{}
		} else {
			outside.Add(Point{X: currentPoint.X - 1, Y: currentPoint.Y})
		}
	}
}

func (NorthWalker) Walk(grid [][]byte, id byte, startPoint Point, outside Set[Point]) (endPoint Point, nextWalker PerimeterWalker) {
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
			return currentPoint, EastWalker{}
		}

		if currentPoint.Y > 0 && currentPoint.X > 0 && grid[currentPoint.Y-1][currentPoint.X-1] == id {
			return currentPoint, WestWalker{}
		} else {
			outside.Add(Point{X: currentPoint.X - 1, Y: currentPoint.Y})
		}
	}
}

func (SouthWalker) Walk(grid [][]byte, id byte, startPoint Point, outside Set[Point]) (endPoint Point, nextWalker PerimeterWalker) {
	if id != grid[startPoint.Y][startPoint.X-1] {
		panic("starting position id should match given id")
	}

	currentPoint := startPoint

	outside.Add(currentPoint)
	for {
		currentPoint.Y++

		if currentPoint.Y == len(grid) || grid[currentPoint.Y][currentPoint.X-1] != id {
			return currentPoint, WestWalker{}
		}

		if currentPoint.Y < len(grid) && currentPoint.X < len(grid[currentPoint.Y]) && grid[currentPoint.Y][currentPoint.X] == id {
			return currentPoint, EastWalker{}
		} else {
			outside.Add(currentPoint)
		}
	}
}

func walkToEastStart(id byte, startPoint Point, grid [][]byte) Point {
	area := walkArea(id, startPoint, grid, make(Set[Point]))
	for point := range area {
		if point.Y < startPoint.Y || point.Y == startPoint.Y && point.X < startPoint.X {
			startPoint = point
		}
	}
	return startPoint
}

func WalkPerimeterCollectOutside(startPoint Point, grid [][]byte, outside Set[Point]) (perimeter []Point) {
	id := grid[startPoint.Y][startPoint.X]

	var currentWalker PerimeterWalker = EastWalker{}
	startPoint = walkToEastStart(id, startPoint, grid)

	currentPoint, currentWalker := currentWalker.Walk(grid, id, startPoint, outside)
	perimeter = append(perimeter, startPoint)
	perimeter = append(perimeter, currentPoint)

	for currentPoint != startPoint {
		currentPoint, currentWalker = currentWalker.Walk(grid, id, currentPoint, outside)
		perimeter = append(perimeter, currentPoint)
	}
	return

}

func WalkPerimeter(startPoint Point, grid [][]byte) (perimeter []Point) {
	return WalkPerimeterCollectOutside(startPoint, grid, make(Set[Point]))
}

func walkArea(id byte, point Point, grid [][]byte, area Set[Point]) Set[Point] {
	if id != grid[point.Y][point.X] {
		panic("position id should match given id")
	}
	directions := []Direction{{-1, 0}, {1, 0}, {0, 1}, {0, -1}}

	area.Add(point)

	for _, direction := range directions {
		nextPoint := Point{point.X + direction.dx, point.Y + direction.dy}

		if area.Contains(nextPoint) || !withinBounds(grid, nextPoint) || grid[nextPoint.Y][nextPoint.X] != id {
			continue
		}
		area = walkArea(id, nextPoint, grid, area)
	}

	return area
}

func walkHole(id byte, point Point, grid [][]byte, holeArea Set[Point]) {
	directions := []Direction{{-1, 0}, {1, 0}, {0, 1}, {0, -1}}

	holeArea.Add(point)

	for _, direction := range directions {
		nextPoint := Point{point.X + direction.dx, point.Y + direction.dy}

		if !holeArea.Contains(nextPoint) && withinBounds(grid, nextPoint) && grid[nextPoint.Y][nextPoint.X] == id {
			walkHole(id, nextPoint, grid, holeArea)
		}
	}
}

func collectHoles(id byte, point Point, grid [][]byte, area Set[Point], walked Set[Point], holes []Hole, outside Set[Point]) []Hole {
	if id != grid[point.Y][point.X] {
		panic("position id should match given id")
	}
	directions := []Direction{{-1, 0}, {1, 0}, {0, 1}, {0, -1}}

	walked.Add(point)
	for _, direction := range directions {
		nextPoint := Point{point.X + direction.dx, point.Y + direction.dy}
		if walked.Contains(nextPoint) || !withinBounds(grid, nextPoint) {
			continue
		}
		if grid[nextPoint.Y][nextPoint.X] == id {
			holes = collectHoles(id, nextPoint, grid, area, walked, holes, outside)
		} else if !outside.Contains(nextPoint) && !fn.Any(holes, func(_ int, hole Hole) bool { return hole.Area.Contains(nextPoint) }) {
			holeArea := make(Set[Point])
			walkHole(grid[nextPoint.Y][nextPoint.X], nextPoint, grid, holeArea)
			holePerimeter := WalkPerimeter(nextPoint, grid)
			holes = append(holes, Hole{id: grid[nextPoint.Y][nextPoint.X], Area: holeArea, Perimeter: holePerimeter})
		}
	}
	return holes
}

func WalkArea(id byte, startPoint Point, grid [][]byte, outside Set[Point]) (Set[Point], []Hole) {
	if id != grid[startPoint.Y][startPoint.X] {
		panic("starting position id should match given id")
	}
	area := walkArea(id, startPoint, grid, make(Set[Point]))
	holes := collectHoles(id, startPoint, grid, area, make(Set[Point]), []Hole{}, outside)
	return area, holes
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
		newPerimeter := WalkPerimeter(hole.Perimeter[0], newGrid)
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

func Walk(startPoint Point, grid [][]byte) Region {
	id := grid[startPoint.Y][startPoint.X]

	outside := make(Set[Point])
	perimeter := WalkPerimeterCollectOutside(startPoint, grid, outside)

	area, holes := WalkArea(id, startPoint, grid, outside)
	return Region{id: id, Perimeter: perimeter, InsidePerimeters: getInsideParameters(grid, holes), Area: area, Holes: holes}
}

func constructArea(point Point, grid [][]byte, usedCoordinates Set[Point]) Region {
	area := Walk(point, grid)

	for point := range area.Area {
		usedCoordinates.Add(point)
	}

	for _, hole := range area.Holes {
		for point := range hole.Area {
			usedCoordinates.Add(point)
		}
	}
	return area
}

func RegionsFromGrid(grid [][]byte) (regions []Region) {
	usedCoordinates := make(Set[Point])
	for y, row := range grid {
		for x := range row {
			point := Point{X: x, Y: y}
			if usedCoordinates.Contains(point) {
				continue
			}
			regions = append(regions, constructArea(point, grid, usedCoordinates))
		}
	}

	return
}
