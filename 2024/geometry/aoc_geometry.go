package aoc_geometry

import "errors"

type Point struct {
	X, Y int
}

type Direction struct {
	dx, dy int
}

func WithinBounds[T any](grid [][]T, pos Point) bool {
	return pos.Y >= 0 && pos.Y < len(grid) && pos.X >= 0 && pos.X < len(grid[pos.Y])
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
	if id != grid[startPoint.Y][startPoint.X] {
		panic("starting position id should match given id")
	}

	currentPoint := startPoint

	for {
		if currentPoint.X == len(grid[currentPoint.Y])-1 || grid[currentPoint.Y][currentPoint.X+1] != id {
			currentPoint.X++
			return currentPoint, SouthWalker{}
		}

		currentPoint.X++
		if currentPoint.Y > 0 && grid[currentPoint.Y-1][currentPoint.X] == id {
			return currentPoint, NorthWalker{}
		}
	}
}

func (WestWalker) Walk(grid [][]byte, id byte, startPoint Point) (endPoint Point, nextWalker PerimeterWalker) {
	if id != grid[startPoint.Y-1][startPoint.X-1] {
		panic("starting position id should match given id")
	}

	currentPoint := startPoint

	for {
		currentPoint.X--

		if currentPoint.X == 0 || grid[currentPoint.Y-1][currentPoint.X-1] != id {
			return currentPoint, NorthWalker{}
		}

		if currentPoint.X != 0 && currentPoint.Y < len(grid) && grid[currentPoint.Y][currentPoint.X-1] == id {
			return currentPoint, SouthWalker{}
		}
	}
}

func (NorthWalker) Walk(grid [][]byte, id byte, startPoint Point) (endPoint Point, nextWalker PerimeterWalker) {
	if id != grid[startPoint.Y-1][startPoint.X] {
		panic("starting position id should match given id")
	}

	currentPoint := startPoint

	for {
		currentPoint.Y--

		if currentPoint.Y == 0 || grid[currentPoint.Y-1][currentPoint.X] != id {
			return currentPoint, EastWalker{}
		}

		if currentPoint.Y > 0 && currentPoint.X > 0 && grid[currentPoint.Y-1][currentPoint.X-1] == id {
			return currentPoint, WestWalker{}
		}
	}
}

func (SouthWalker) Walk(grid [][]byte, id byte, startPoint Point) (endPoint Point, nextWalker PerimeterWalker) {
	if id != grid[startPoint.Y][startPoint.X-1] {
		panic("starting position id should match given id")
	}

	currentPoint := startPoint

	for {
		currentPoint.Y++

		if currentPoint.Y == len(grid) || grid[currentPoint.Y][currentPoint.X-1] != id {
			return currentPoint, WestWalker{}
		}

		if currentPoint.Y < len(grid) && currentPoint.X < len(grid[currentPoint.Y]) && grid[currentPoint.Y][currentPoint.X] == id {
			return currentPoint, EastWalker{}
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

func WalkPerimeter(startPoint Point, grid [][]byte) (perimeter []Point) {
	id := grid[startPoint.Y][startPoint.X]

	var currentWalker PerimeterWalker = EastWalker{}
	startPoint = walkToEastStart(id, startPoint, grid)

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
	if id != grid[point.Y][point.X] {
		panic("position id should match given id")
	}
	directions := []Direction{{-1, 0}, {1, 0}, {0, 1}, {0, -1}}

	area.Add(point)
	for _, direction := range directions {
		nextPoint := Point{point.X + direction.dx, point.Y + direction.dy}
		if area.Contains(nextPoint) || !WithinBounds(grid, nextPoint) || grid[nextPoint.Y][nextPoint.X] != id {
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
		diagonalPoint := Point{point.X + diagonalDirection.dx, point.Y + diagonalDirection.dy}

		if !WithinBounds(grid, diagonalPoint) {
			return errors.New("not a hole")
		}

		if grid[diagonalPoint.Y][diagonalPoint.X] != id && !area.Contains(diagonalPoint) {
			return errors.New("not a hole")
		}
	}

	for _, direction := range directions {
		nextPoint := Point{point.X + direction.dx, point.Y + direction.dy}
		if !WithinBounds(grid, nextPoint) {
			return errors.New("not a hole")
		}

		if grid[nextPoint.Y][nextPoint.X] != id && !area.Contains(nextPoint) {
			return errors.New("not a hole")
		}

		if !holeArea.Contains(nextPoint) && grid[nextPoint.Y][nextPoint.X] == id {
			if err := walkHole(id, nextPoint, grid, area, holeArea, ignoreHolePoints); err != nil {
				return err
			}
		}
	}

	return nil
}

func collectHoles(id byte, point Point, grid [][]byte, area Set[Point], walked Set[Point], holes []Hole, ignoreHolePoints Set[Point]) []Hole {
	if id != grid[point.Y][point.X] {
		panic("position id should match given id")
	}
	directions := []Direction{{-1, 0}, {1, 0}, {0, 1}, {0, -1}}

	walked.Add(point)
	for _, direction := range directions {
		nextPoint := Point{point.X + direction.dx, point.Y + direction.dy}
		if walked.Contains(nextPoint) || !WithinBounds(grid, nextPoint) {
			continue
		}
		if grid[nextPoint.Y][nextPoint.X] == id {
			holes = collectHoles(id, nextPoint, grid, area, walked, holes, ignoreHolePoints)
		} else if !ignoreHolePoints.Contains(nextPoint) {
			holeArea := make(Set[Point])
			if err := walkHole(grid[nextPoint.Y][nextPoint.X], nextPoint, grid, area, holeArea, ignoreHolePoints); err == nil {
				holePerimeter := WalkPerimeter(nextPoint, grid)
				holes = append(holes, Hole{id: grid[nextPoint.Y][nextPoint.X], Area: holeArea, Perimeter: holePerimeter})
			}
		}
	}
	return holes
}

func WalkArea(id byte, startPoint Point, grid [][]byte) (Set[Point], []Hole) {
	if id != grid[startPoint.Y][startPoint.X] {
		panic("starting position id should match given id")
	}
	area := walkArea(id, startPoint, grid, make(Set[Point]))
	holes := collectHoles(id, startPoint, grid, area, make(Set[Point]), []Hole{}, make(Set[Point]))
	return area, holes
}

func Walk(startPoint Point, grid [][]byte) Area {
	id := grid[startPoint.Y][startPoint.X]
	area, holes := WalkArea(id, startPoint, grid)
	perimeter := WalkPerimeter(startPoint, grid)
	return Area{id: id, Perimeter: perimeter, Area: area, Holes: holes}
}

func constructArea(point Point, grid [][]byte, usedCoordinates Set[Point]) Area {
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

func AreasFromGrid(grid [][]byte) (areas []Area) {
	usedCoordinates := make(Set[Point])
	for y, row := range grid {
		for x := range row {
			point := Point{X: x, Y: y}
			if usedCoordinates.Contains(point) {
				continue
			}
			areas = append(areas, constructArea(point, grid, usedCoordinates))
		}
	}

	return
}
