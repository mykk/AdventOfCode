package aoc_geometry

type Point struct {
	x, y int
}

type Direction struct {
	dx, dy int
}

func WithinBounds[T any](grid [][]T, pos Point) bool {
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

func WalkPerimeter(startPoint Point, grid [][]byte) (perimeter []Point) {
	id := grid[startPoint.y][startPoint.x]

	var currentWalker PerimeterWalker = EastWalker{}
	currentPoint, currentWalker := currentWalker.Walk(grid, id, startPoint)
	perimeter = append(perimeter, startPoint)
	perimeter = append(perimeter, currentPoint)

	for currentPoint != startPoint {
		currentPoint, currentWalker = currentWalker.Walk(grid, id, currentPoint)
		perimeter = append(perimeter, currentPoint)
	}
	return
}

func Walk(startPoint Point, grid [][]byte) Area {
	id := grid[startPoint.y][startPoint.x]
	return Area{id: id, perimeter: WalkPerimeter(startPoint, grid), area: []Point{}, holes: [][]Point{}}
}

// func constructArea(x, y int, garden [][]byte, usedCoordinates Set[Point]) (area Area) {
// 	identifier := garden[y][x]
// 	//polygon = append(polygon, Point{x: x, y: y})

// 	directions := []Direction{{1, 0}, {0, -1}, {0, 1}, {1, 0}}

// 	currentPoint := Point{x: x, y: y}
// 	for _, dir := range directions {
// 		nextPosition := Point{x: currentPoint.x + dir.dx, y: currentPoint.y + dir.dy}
// 		if WithinBounds(garden, nextPosition) && garden[nextPosition.y][nextPosition.x] == identifier {

// 		}
// 	}
// 	return
// }

// func constructAreas(garden [][]byte) (areas []Area) {
// 	usedCoordinates := make(Set[Point])
// 	for y, row := range garden {
// 		for x, _ := range row {
// 			if usedCoordinates.Contains(Point{x: x, y: y}) {
// 				continue
// 			}
// 			areas = append(areas, constructArea(x, y, garden, usedCoordinates))
// 		}
// 	}

// 	return
// }

// func ParseInputData(data string) []Area {
// 	garden := fn.MustTransform(fn.GetLines(data), func(line string) []byte { return []byte(line) })
// 	return constructAreas(garden)
// }
