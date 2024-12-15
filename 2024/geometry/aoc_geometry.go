package aoc_geometry

import "errors"

type Point struct {
	x, y int
}

type Direction struct {
	dx, dy int
}

func WithinBounds[T any](grid [][]T, pos Point) bool {
	return pos.y >= 0 && pos.y < len(grid) && pos.x >= 0 && pos.x < len(grid[pos.y])
}

type Hole struct {
	id        byte
	Perimeter []Point
	Area      Set[Point]
}

type Area struct {
	id        byte
	Perimeter []Point
	Area      Set[Point]
	Holes     []Hole
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
	Walk(grid [][]byte, id byte, startPoint Point) (endPoint Point, nextWalker PerimeterWalker)
}

type EastWalker struct{}
type WestWalker struct{}
type NorthWalker struct{}
type SouthWalker struct{}

func (EastWalker) Walk(grid [][]byte, id byte, startPoint Point) (endPoint Point, nextWalker PerimeterWalker) {
	if id != grid[startPoint.y][startPoint.x] {
		panic("starting position id should match given id")
	}

	currentPoint := startPoint

	for {
		if currentPoint.y > 0 && grid[currentPoint.y-1][currentPoint.x] == id {
			return currentPoint, NorthWalker{}
		}

		if len(grid[currentPoint.y])-1 > currentPoint.x && grid[currentPoint.y][currentPoint.x+1] == id {
			currentPoint.x++
			continue
		}

		currentPoint.x++
		return currentPoint, SouthWalker{}
	}
}

func (WestWalker) Walk(grid [][]byte, id byte, startPoint Point) (endPoint Point, nextWalker PerimeterWalker) {
	if id != grid[startPoint.y-1][startPoint.x-1] {
		panic("starting position id should match given id")
	}

	currentPoint := startPoint

	for {
		if currentPoint.x > 0 && len(grid)-1 > currentPoint.y && grid[currentPoint.y][currentPoint.x-1] == id {
			return currentPoint, SouthWalker{}
		}

		if currentPoint.x > 0 && grid[currentPoint.y-1][currentPoint.x-1] == id {
			currentPoint.x--
			continue
		}

		return currentPoint, NorthWalker{}
	}
}

func (NorthWalker) Walk(grid [][]byte, id byte, startPoint Point) (endPoint Point, nextWalker PerimeterWalker) {
	if id != grid[startPoint.y-1][startPoint.x] {
		panic("starting position id should match given id")
	}

	currentPoint := startPoint

	for {
		if currentPoint.y > 0 && currentPoint.x > 0 && grid[currentPoint.y-1][currentPoint.x-1] == id {
			return currentPoint, WestWalker{}
		}

		if currentPoint.y > 0 && grid[currentPoint.y-1][currentPoint.x] == id {
			currentPoint.y--
			continue
		}
		return currentPoint, EastWalker{}
	}
}

func (SouthWalker) Walk(grid [][]byte, id byte, startPoint Point) (endPoint Point, nextWalker PerimeterWalker) {
	if id != grid[startPoint.y][startPoint.x-1] {
		panic("starting position id should match given id")
	}

	currentPoint := startPoint

	for {
		if len(grid) > currentPoint.y && len(grid[currentPoint.y]) > currentPoint.x && grid[currentPoint.y][currentPoint.x] == id {
			return currentPoint, EastWalker{}
		}

		if currentPoint.x > 0 && len(grid) > currentPoint.y && grid[currentPoint.y][currentPoint.x-1] == id {
			currentPoint.y++
			continue
		}
		return currentPoint, WestWalker{}
	}
}

func walkToEastStart(startPoint Point, grid [][]byte) Point {
	currentPoint := startPoint
	id := grid[startPoint.y][startPoint.x]

	for currentPoint.x > 0 {
		if currentPoint.y > 0 && grid[currentPoint.y-1][currentPoint.x-1] == id {
			return currentPoint
		}

		if grid[currentPoint.y][currentPoint.x-1] != id {
			return currentPoint
		}
		currentPoint.x--
	}
	return currentPoint
}

func WalkPerimeter(startPoint Point, grid [][]byte) (perimeter []Point) {
	id := grid[startPoint.y][startPoint.x]

	var currentWalker PerimeterWalker = EastWalker{}
	startPoint = walkToEastStart(startPoint, grid)

	currentPoint, currentWalker := currentWalker.Walk(grid, id, startPoint)
	perimeter = append(perimeter, startPoint)
	perimeter = append(perimeter, currentPoint)

	for currentPoint != startPoint {
		currentPoint, currentWalker = currentWalker.Walk(grid, id, currentPoint)
		perimeter = append(perimeter, currentPoint)
	}
	return
}

func walkArea(id byte, point Point, grid [][]byte, area Set[Point]) Set[Point] {
	if id != grid[point.y][point.x-1] {
		panic("position id should match given id")
	}
	directions := []Direction{{-1, 0}, {1, 0}, {0, 1}, {0, -1}}

	area.Add(point)
	for _, direction := range directions {
		nextPoint := Point{point.x + direction.dx, point.y + direction.dy}
		if area.Contains(nextPoint) || !WithinBounds(grid, nextPoint) || grid[nextPoint.y][nextPoint.x] != id {
			continue
		}
		walkArea(id, nextPoint, grid, area)
	}

	return area
}

func walkHole(id byte, point Point, grid [][]byte, area Set[Point], holeArea, ignoreHolePoints Set[Point]) error {
	directions := []Direction{{-1, 0}, {1, 0}, {0, 1}, {0, -1}}
	diagonalDirections := []Direction{{-1, -1}, {1, 1}, {-1, 1}, {1, -1}}

	holeArea.Add(point)
	ignoreHolePoints.Add(point)

	for _, diagonalDirection := range diagonalDirections {
		diagonalPoint := Point{point.x + diagonalDirection.dx, point.y + diagonalDirection.dy}
		if grid[diagonalPoint.y][diagonalPoint.x] != id && !area.Contains(diagonalPoint) {
			return errors.New("not a hole")
		}
	}

	for _, direction := range directions {
		nextPoint := Point{point.x + direction.dx, point.y + direction.dy}
		if !WithinBounds(grid, nextPoint) {
			continue
		}

		if grid[nextPoint.y][nextPoint.x] != id && !area.Contains(nextPoint) {
			return errors.New("not a hole")
		}

		if err := walkHole(id, nextPoint, grid, area, holeArea, ignoreHolePoints); err != nil {
			return err
		}
	}

	return nil
}

func collectHoles(id byte, point Point, grid [][]byte, area Set[Point], walked Set[Point], holes []Hole, ignoreHolePoints Set[Point]) []Hole {
	if id != grid[point.y][point.x-1] {
		panic("position id should match given id")
	}
	directions := []Direction{{-1, 0}, {1, 0}, {0, 1}, {0, -1}}

	walked.Add(point)
	for _, direction := range directions {
		nextPoint := Point{point.x + direction.dx, point.y + direction.dy}
		if walked.Contains(nextPoint) || !WithinBounds(grid, nextPoint) {
			continue
		}
		if grid[nextPoint.y][nextPoint.x] == id {
			collectHoles(id, nextPoint, grid, area, walked, holes, ignoreHolePoints)
		} else if !ignoreHolePoints.Contains(nextPoint) {
			holeArea := make(Set[Point])
			if err := walkHole(grid[nextPoint.y][nextPoint.x], nextPoint, grid, area, holeArea, ignoreHolePoints); err == nil {
				holePerimeter := WalkPerimeter(nextPoint, grid)
				holes = append(holes, Hole{id: grid[nextPoint.y][nextPoint.x], Area: holeArea, Perimeter: holePerimeter})
			}
		}
	}

	return holes
}

func WalkArea(id byte, startPoint Point, grid [][]byte) (Set[Point], []Hole) {
	if id != grid[startPoint.y][startPoint.x-1] {
		panic("starting position id should match given id")
	}
	area := walkArea(id, startPoint, grid, make(Set[Point]))
	holes := collectHoles(id, startPoint, grid, area, make(Set[Point]), []Hole{}, make(Set[Point]))
	return area, holes
}

func Walk(startPoint Point, grid [][]byte) Area {
	id := grid[startPoint.y][startPoint.x]
	perimeter := WalkPerimeter(startPoint, grid)
	area, holes := WalkArea(id, startPoint, grid)
	return Area{id: id, Perimeter: perimeter, Area: area, Holes: holes}
}
