package aoc_geometry

import (
	fn "AoC/functional"

	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	t.Run("PerimeterTest1", func(t *testing.T) {
		const PERIMETER_DATA string = `AAA`

		perimeter := WalkPerimeter(Point{0, 0}, [][]byte{[]byte(PERIMETER_DATA)})

		assert.Equal(t, []Point{{0, 0}, {3, 0}, {3, 1}, {0, 1}, {0, 0}}, perimeter)
	})

	t.Run("PerimeterTest2", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte(" A "))
		testData = append(testData, []byte("AAA"))

		perimeter := WalkPerimeter(Point{0, 1}, testData)

		assert.Equal(t, []Point{{1, 0}, {2, 0}, {2, 1}, {3, 1}, {3, 2}, {0, 2}, {0, 1}, {1, 1}, {1, 0}}, perimeter)
	})

	t.Run("PerimeterTest3", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("AA "))
		testData = append(testData, []byte(" A "))
		testData = append(testData, []byte("AAA"))

		perimeter := WalkPerimeter(Point{0, 0}, testData)

		assert.Equal(t, []Point{{0, 0}, {2, 0}, {2, 2}, {3, 2}, {3, 3}, {0, 3}, {0, 2}, {1, 2}, {1, 1}, {0, 1}, {0, 0}}, perimeter)
	})

	t.Run("PerimeterTest4", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte(" AA"))
		testData = append(testData, []byte(" A "))
		testData = append(testData, []byte("AAA"))

		perimeter := WalkPerimeter(Point{1, 0}, testData)

		assert.Equal(t, []Point{{1, 0}, {3, 0},
			{3, 1}, {2, 1},
			{2, 2}, {3, 2},
			{3, 3}, {0, 3}, {0, 2},
			{1, 2}, {1, 0}}, perimeter)
	})

	t.Run("PerimeterTest5", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("AAA"))
		testData = append(testData, []byte(" A "))
		testData = append(testData, []byte("AAA"))

		perimeter := WalkPerimeter(Point{0, 0}, testData)

		assert.Equal(t, []Point{{0, 0}, {3, 0},
			{3, 1}, {2, 1},
			{2, 2}, {3, 2},
			{3, 3}, {0, 3}, {0, 2},
			{1, 2}, {1, 1}, {0, 1}, {0, 0}}, perimeter)
	})

	t.Run("PerimeterTest6", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("  A"))
		testData = append(testData, []byte("AAA"))
		testData = append(testData, []byte(" A "))
		testData = append(testData, []byte("AAA"))

		perimeter := WalkPerimeter(Point{2, 0}, testData)

		expected := []Point{{2, 0}, {3, 0},
			{3, 2}, {2, 2},
			{2, 3}, {3, 3},
			{3, 4}, {0, 4}, {0, 3},
			{1, 3}, {1, 2}, {0, 2}, {0, 1},
			{2, 1}, {2, 0}}

		assert.Equal(t, expected, perimeter)
	})

	t.Run("PerimeterTest7", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("AAA"))
		testData = append(testData, []byte("AAA"))
		testData = append(testData, []byte("AAA"))
		testData = append(testData, []byte("AAA"))

		perimeter := WalkPerimeter(Point{0, 0}, testData)

		expected := []Point{{0, 0}, {3, 0}, {3, 4}, {0, 4}, {0, 0}}

		assert.Equal(t, expected, perimeter)
	})

	t.Run("PerimeterTest8", func(t *testing.T) {
		const PERIMETER_DATA string = `AAAA`

		perimeter := WalkPerimeter(Point{3, 0}, [][]byte{[]byte(PERIMETER_DATA)})

		assert.Equal(t, []Point{{0, 0}, {4, 0}, {4, 1}, {0, 1}, {0, 0}}, perimeter)
	})

	t.Run("PerimeterTest9", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("A    "))
		testData = append(testData, []byte("AAAAA"))

		perimeter := WalkPerimeter(Point{3, 1}, testData)

		assert.Equal(t, []Point{{0, 0}, {1, 0}, {1, 1}, {5, 1}, {5, 2}, {0, 2}, {0, 0}}, perimeter)
	})

	t.Run("PerimeterTest10", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("A    "))
		testData = append(testData, []byte(" AAAA"))

		perimeter := WalkPerimeter(Point{3, 1}, testData)

		assert.Equal(t, []Point{{1, 1}, {5, 1}, {5, 2}, {1, 2}, {1, 1}}, perimeter)
	})

	t.Run("PerimeterTest11", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("    AA"))
		testData = append(testData, []byte("AAAA  "))

		perimeter := WalkPerimeter(Point{1, 1}, testData)

		assert.Equal(t, []Point{{0, 1}, {4, 1}, {4, 2}, {0, 2}, {0, 1}}, perimeter)
	})

	t.Run("PerimeterTest12_1", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("AA    AA"))
		testData = append(testData, []byte("  AAAA  "))
		testData = append(testData, []byte("AA    AA"))

		perimeter := WalkPerimeter(Point{2, 1}, testData)

		assert.Equal(t, []Point{{2, 1}, {6, 1}, {6, 2}, {2, 2}, {2, 1}}, perimeter)
	})

	t.Run("PerimeterTest12_2", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("AA    AA"))
		testData = append(testData, []byte("  A   AA"))
		testData = append(testData, []byte("  AAAA  "))
		testData = append(testData, []byte("AA    AA"))

		perimeter := WalkPerimeter(Point{2, 1}, testData)

		assert.Equal(t, []Point{{2, 1}, {3, 1}, {3, 2}, {6, 2}, {6, 3}, {2, 3}, {2, 1}}, perimeter)
	})

	t.Run("PerimeterTest12_3", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("AA    AA"))
		testData = append(testData, []byte("  A   AA"))
		testData = append(testData, []byte("  AAAA  "))
		testData = append(testData, []byte("AA   A  "))
		testData = append(testData, []byte("AA    AA"))

		perimeter := WalkPerimeter(Point{2, 1}, testData)

		assert.Equal(t, []Point{{2, 1}, {3, 1}, {3, 2}, {6, 2}, {6, 4}, {5, 4}, {5, 3}, {2, 3}, {2, 1}}, perimeter)
	})

	t.Run("PerimeterTest12_4", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("AA    AA"))
		testData = append(testData, []byte("  A   AA"))
		testData = append(testData, []byte("  AAAA  "))
		testData = append(testData, []byte("AA   A  "))
		testData = append(testData, []byte("AA   AAA"))

		perimeter := WalkPerimeter(Point{2, 1}, testData)

		assert.Equal(t, []Point{{2, 1}, {3, 1}, {3, 2}, {6, 2}, {6, 4}, {8, 4}, {8, 5}, {5, 5}, {5, 3}, {2, 3}, {2, 1}}, perimeter)
	})

	t.Run("PerimeterTest12_5", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("A"))
		testData = append(testData, []byte("A"))
		testData = append(testData, []byte("A"))

		perimeter := WalkPerimeter(Point{0, 2}, testData)

		assert.Equal(t, []Point{{0, 0}, {1, 0}, {1, 3}, {0, 3}, {0, 0}}, perimeter)
	})

	t.Run("PerimeterTest12_6", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte(" A"))
		testData = append(testData, []byte("A "))
		testData = append(testData, []byte("A "))
		testData = append(testData, []byte("A "))

		perimeter := WalkPerimeter(Point{0, 2}, testData)

		assert.Equal(t, []Point{{0, 1}, {1, 1}, {1, 4}, {0, 4}, {0, 1}}, perimeter)
	})

	t.Run("PerimeterTest12_7", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte(" WW   "))
		testData = append(testData, []byte("WWWW  "))
		testData = append(testData, []byte("WWW W "))
		testData = append(testData, []byte("WWWWW "))

		perimeter := WalkPerimeter(Point{0, 2}, testData)

		assert.Equal(t, []Point{{1, 0}, {3, 0}, {3, 1}, {4, 1}, {4, 2}, {3, 2}, {3, 3}, {4, 3}, {4, 2}, {5, 2}, {5, 4}, {0, 4}, {0, 1}, {1, 1}, {1, 0}}, perimeter)
	})

	t.Run("PerimeterTest12_8", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("HH HH "))
		testData = append(testData, []byte("HHHH H"))
		testData = append(testData, []byte("HHHHHH"))

		areas := AreasFromGrid(testData)

		assert.Equal(t, []Point{{0, 0}, {2, 0}, {2, 1}, {3, 1}, {3, 0}, {5, 0}, {5, 1}, {4, 1}, {4, 2}, {5, 2}, {5, 1}, {6, 1}, {6, 3}, {0, 3}, {0, 0}}, areas[0].Perimeter)

		perimeter := GetTotalPerimeter(areas[0])
		assert.Equal(t, 24, perimeter)
	})

	t.Run("AreaTest1", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("AAAA"))
		testData = append(testData, []byte("BBCD"))
		testData = append(testData, []byte("BBCC"))
		testData = append(testData, []byte("EEEC"))

		areas := AreasFromGrid(testData)
		assert.Equal(t, 5, len(areas))

		assert.True(t, fn.All(areas, func(_ int, area Area) bool { return len(area.Holes) == 0 }))
		assert.True(t, fn.All(areas, func(_ int, area Area) bool {
			if area.id == 'A' {
				return len(area.Area) == 4
			}

			if area.id == 'B' {
				return len(area.Area) == 4
			}

			if area.id == 'C' {
				return len(area.Area) == 4
			}

			if area.id == 'D' {
				return len(area.Area) == 1
			}

			if area.id == 'E' {
				return len(area.Area) == 3
			}

			return false
		}))
	})

	t.Run("AreaTest2", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("OOOOO"))
		testData = append(testData, []byte("OXOXO"))
		testData = append(testData, []byte("OOOOO"))
		testData = append(testData, []byte("OXOXO"))
		testData = append(testData, []byte("OOOOO"))

		areas := AreasFromGrid(testData)
		assert.Equal(t, 1, len(areas))
		assert.Equal(t, 4, len(areas[0].Holes))

		assert.True(t, fn.All(areas[0].Holes, func(_ int, hole Hole) bool {
			return len(hole.Area) == 1 && len(hole.Perimeter) == 5
		}))
	})

	t.Run("AreaTest3", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("OOOOO"))
		testData = append(testData, []byte("OXXXO"))
		testData = append(testData, []byte("OOOOO"))
		testData = append(testData, []byte("OOOOO"))
		testData = append(testData, []byte("OOOBB"))

		areas := AreasFromGrid(testData)
		assert.Equal(t, 2, len(areas))
		assert.True(t, fn.All(areas, func(_ int, area Area) bool {
			if area.id == 'O' && len(area.Holes) != 1 {
				return false
			}

			if area.id == 'O' {
				return len(area.Holes[0].Area) == 3
			}

			if area.id == 'B' {
				return len(area.Holes) == 0
			}

			return false
		}))
	})

	t.Run("AreaTest4", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("OOOOO"))
		testData = append(testData, []byte("OXXXX"))
		testData = append(testData, []byte("OOOOO"))
		testData = append(testData, []byte("OOOOO"))
		testData = append(testData, []byte("OOOBB"))

		areas := AreasFromGrid(testData)
		assert.Equal(t, 3, len(areas))
		assert.True(t, fn.All(areas, func(_ int, area Area) bool {
			return len(area.Holes) == 0
		}))
	})

	t.Run("AreaTest5__", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("OOOOO"))
		testData = append(testData, []byte("OXXXO"))
		testData = append(testData, []byte("OOBBO"))
		testData = append(testData, []byte("OOOOO"))
		testData = append(testData, []byte("OOOBB"))

		areas := AreasFromGrid(testData)
		for _, area := range areas {
			if area.id != 'O' {
				continue
			}
			assert.Equal(t, area.QuickPerimeter, GetTotalPerimeter(area))
		}
	})
}
