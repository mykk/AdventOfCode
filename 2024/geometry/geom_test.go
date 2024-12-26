package aoc_geometry

import (
	fn "AoC/functional"

	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	t.Run("PerimeterTest1", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("AAA"))

		regions := RegionsFromGrid(testData)

		assert.Equal(t, []Point{{0, 0}, {3, 0}, {3, 1}, {0, 1}, {0, 0}}, regions[0].Perimeter)
	})

	t.Run("PerimeterTest2", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte(" A "))
		testData = append(testData, []byte("AAA"))

		regions := RegionsFromGrid(testData)

		for _, region := range regions {
			if !region.Area.Contains(Point{1, 0}) {
				continue
			}
			assert.Equal(t, []Point{{1, 0}, {2, 0}, {2, 1}, {3, 1}, {3, 2}, {0, 2}, {0, 1}, {1, 1}, {1, 0}}, region.Perimeter)
		}
	})

	t.Run("PerimeterTest3", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("AA "))
		testData = append(testData, []byte(" A "))
		testData = append(testData, []byte("AAA"))

		regions := RegionsFromGrid(testData)
		for _, region := range regions {
			if !region.Area.Contains(Point{0, 0}) {
				continue
			}
			assert.Equal(t, []Point{{0, 0}, {2, 0}, {2, 2}, {3, 2}, {3, 3}, {0, 3}, {0, 2}, {1, 2}, {1, 1}, {0, 1}, {0, 0}}, region.Perimeter)
		}
	})

	t.Run("PerimeterTest4", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte(" AA"))
		testData = append(testData, []byte(" A "))
		testData = append(testData, []byte("AAA"))

		regions := RegionsFromGrid(testData)

		for _, region := range regions {
			if !region.Area.Contains(Point{1, 0}) {
				continue
			}
			assert.Equal(t, []Point{{1, 0}, {3, 0},
				{3, 1}, {2, 1},
				{2, 2}, {3, 2},
				{3, 3}, {0, 3}, {0, 2},
				{1, 2}, {1, 0}}, region.Perimeter)
		}
	})

	t.Run("PerimeterTest5", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("AAA"))
		testData = append(testData, []byte(" A "))
		testData = append(testData, []byte("AAA"))

		regions := RegionsFromGrid(testData)

		for _, region := range regions {
			if !region.Area.Contains(Point{0, 0}) {
				continue
			}
			assert.Equal(t, []Point{{0, 0}, {3, 0},
				{3, 1}, {2, 1},
				{2, 2}, {3, 2},
				{3, 3}, {0, 3}, {0, 2},
				{1, 2}, {1, 1}, {0, 1}, {0, 0}}, region.Perimeter)
		}
	})

	t.Run("PerimeterTest6", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("  A"))
		testData = append(testData, []byte("AAA"))
		testData = append(testData, []byte(" A "))
		testData = append(testData, []byte("AAA"))

		regions := RegionsFromGrid(testData)

		for _, region := range regions {
			if !region.Area.Contains(Point{2, 0}) {
				continue
			}
			expected := []Point{{2, 0}, {3, 0},
				{3, 2}, {2, 2},
				{2, 3}, {3, 3},
				{3, 4}, {0, 4}, {0, 3},
				{1, 3}, {1, 2}, {0, 2}, {0, 1},
				{2, 1}, {2, 0}}

			assert.Equal(t, expected, region.Perimeter)
		}
	})

	t.Run("PerimeterTest7", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("AAA"))
		testData = append(testData, []byte("AAA"))
		testData = append(testData, []byte("AAA"))
		testData = append(testData, []byte("AAA"))

		regions := RegionsFromGrid(testData)

		for _, region := range regions {
			if !region.Area.Contains(Point{0, 0}) {
				continue
			}
			expected := []Point{{0, 0}, {3, 0}, {3, 4}, {0, 4}, {0, 0}}

			assert.Equal(t, expected, regions[0].Perimeter)
		}
	})

	t.Run("PerimeterTest8", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("AAAA"))

		regions := RegionsFromGrid(testData)

		assert.Equal(t, []Point{{0, 0}, {4, 0}, {4, 1}, {0, 1}, {0, 0}}, regions[0].Perimeter)
	})

	t.Run("PerimeterTest9", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("A    "))
		testData = append(testData, []byte("AAAAA"))

		regions := RegionsFromGrid(testData)

		assert.Equal(t, []Point{{0, 0}, {1, 0}, {1, 1}, {5, 1}, {5, 2}, {0, 2}, {0, 0}}, regions[0].Perimeter)
	})

	t.Run("PerimeterTest10", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("A    "))
		testData = append(testData, []byte(" AAAA"))

		regions := RegionsFromGrid(testData)

		for _, region := range regions {
			if !region.Area.Contains(Point{X: 1, Y: 1}) {
				continue
			}
			assert.Equal(t, []Point{{1, 1}, {5, 1}, {5, 2}, {1, 2}, {1, 1}}, region.Perimeter)
		}
	})

	t.Run("PerimeterTest11", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("    AA"))
		testData = append(testData, []byte("AAAA  "))

		regions := RegionsFromGrid(testData)

		for _, region := range regions {
			if !region.Area.Contains(Point{X: 1, Y: 1}) {
				continue
			}
			assert.Equal(t, []Point{{0, 1}, {4, 1}, {4, 2}, {0, 2}, {0, 1}}, region.Perimeter)
		}
	})

	t.Run("PerimeterTest12_1", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("AA    AA"))
		testData = append(testData, []byte("  AAAA  "))
		testData = append(testData, []byte("AA    AA"))

		regions := RegionsFromGrid(testData)

		for _, region := range regions {
			if !region.Area.Contains(Point{X: 3, Y: 1}) {
				continue
			}
			assert.Equal(t, []Point{{2, 1}, {6, 1}, {6, 2}, {2, 2}, {2, 1}}, region.Perimeter)
		}
	})

	t.Run("PerimeterTest12_2", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("AA    AA"))
		testData = append(testData, []byte("  A   AA"))
		testData = append(testData, []byte("  AAAA  "))
		testData = append(testData, []byte("AA    AA"))

		regions := RegionsFromGrid(testData)

		for _, region := range regions {
			if !region.Area.Contains(Point{X: 2, Y: 1}) {
				continue
			}
			assert.Equal(t, []Point{{2, 1}, {3, 1}, {3, 2}, {6, 2}, {6, 3}, {2, 3}, {2, 1}}, region.Perimeter)
		}
	})

	t.Run("PerimeterTest12_3", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("AA    AA"))
		testData = append(testData, []byte("  A   AA"))
		testData = append(testData, []byte("  AAAA  "))
		testData = append(testData, []byte("AA   A  "))
		testData = append(testData, []byte("AA    AA"))

		regions := RegionsFromGrid(testData)

		for _, region := range regions {
			if !region.Area.Contains(Point{X: 2, Y: 1}) {
				continue
			}
			assert.Equal(t, []Point{{2, 1}, {3, 1}, {3, 2}, {6, 2}, {6, 4}, {5, 4}, {5, 3}, {2, 3}, {2, 1}}, region.Perimeter)
		}
	})

	t.Run("PerimeterTest12_4", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("AA    AA"))
		testData = append(testData, []byte("  A   AA"))
		testData = append(testData, []byte("  AAAA  "))
		testData = append(testData, []byte("AA   A  "))
		testData = append(testData, []byte("AA   AAA"))

		regions := RegionsFromGrid(testData)

		for _, region := range regions {
			if !region.Area.Contains(Point{X: 2, Y: 1}) {
				continue
			}
			assert.Equal(t, []Point{{2, 1}, {3, 1}, {3, 2}, {6, 2}, {6, 4}, {8, 4}, {8, 5}, {5, 5}, {5, 3}, {2, 3}, {2, 1}}, region.Perimeter)
		}
	})

	t.Run("PerimeterTest12_5", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("A"))
		testData = append(testData, []byte("A"))
		testData = append(testData, []byte("A"))

		regions := RegionsFromGrid(testData)

		assert.Equal(t, []Point{{0, 0}, {1, 0}, {1, 3}, {0, 3}, {0, 0}}, regions[0].Perimeter)
	})

	t.Run("PerimeterTest12_6", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte(" A"))
		testData = append(testData, []byte("A "))
		testData = append(testData, []byte("A "))
		testData = append(testData, []byte("A "))

		regions := RegionsFromGrid(testData)

		for _, region := range regions {
			if !region.Area.Contains(Point{X: 0, Y: 1}) {
				continue
			}
			assert.Equal(t, []Point{{0, 1}, {1, 1}, {1, 4}, {0, 4}, {0, 1}}, region.Perimeter)
		}
	})

	t.Run("PerimeterTest12_7", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte(" WW   "))
		testData = append(testData, []byte("WWWW  "))
		testData = append(testData, []byte("WWW W "))
		testData = append(testData, []byte("WWWWW "))

		regions := RegionsFromGrid(testData)

		for _, region := range regions {
			if !region.Area.Contains(Point{X: 1, Y: 0}) {
				continue
			}
			assert.Equal(t, []Point{{1, 0}, {3, 0}, {3, 1},
				{4, 1}, {4, 2}, {3, 2}, {3, 3},
				{4, 3}, {4, 2}, {5, 2}, {5, 4},
				{0, 4}, {0, 1}, {1, 1}, {1, 0}}, region.Perimeter)
		}
	})

	t.Run("PerimeterTest12_8", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("HH HH "))
		testData = append(testData, []byte("HHHH H"))
		testData = append(testData, []byte("HHHHHH"))

		regions := RegionsFromGrid(testData)

		for _, region := range regions {
			if !region.Area.Contains(Point{X: 0, Y: 0}) {
				continue
			}
			assert.Equal(t, []Point{{0, 0}, {2, 0},
				{2, 1}, {3, 1},
				{3, 0}, {5, 0},
				{5, 1}, {4, 1},
				{4, 2}, {5, 2},
				{5, 1}, {6, 1},
				{6, 3}, {0, 3}, {0, 0}}, region.Perimeter)

			assert.Equal(t, 24, region.GetInsidePerimeter()+region.GetOutsidePerimeter())
		}
	})

	t.Run("AreaTest1", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("AAAA"))
		testData = append(testData, []byte("BBCD"))
		testData = append(testData, []byte("BBCC"))
		testData = append(testData, []byte("EEEC"))

		regions := RegionsFromGrid(testData)
		assert.Equal(t, 5, len(regions))

		assert.True(t, fn.All(regions, func(_ int, region Region) bool { return len(region.Holes) == 0 }))
		assert.True(t, fn.All(regions, func(_ int, region Region) bool {
			if region.id == 'A' {
				return len(region.Area) == 4
			}

			if region.id == 'B' {
				return len(region.Area) == 4
			}

			if region.id == 'C' {
				return len(region.Area) == 4
			}

			if region.id == 'D' {
				return len(region.Area) == 1
			}

			if region.id == 'E' {
				return len(region.Area) == 3
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

		regions := RegionsFromGrid(testData)
		assert.Equal(t, 1, len(regions))
		assert.Equal(t, 4, len(regions[0].Holes))

		assert.True(t, fn.All(regions[0].Holes, func(_ int, hole Hole) bool {
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

		regions := RegionsFromGrid(testData)
		assert.Equal(t, 2, len(regions))
		assert.True(t, fn.All(regions, func(_ int, region Region) bool {
			if region.id == 'O' && len(region.Holes) != 1 {
				return false
			}

			if region.id == 'O' {
				return len(region.Holes[0].Area) == 3
			}

			if region.id == 'B' {
				return len(region.Holes) == 0
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

		regions := RegionsFromGrid(testData)
		assert.Equal(t, 3, len(regions))
		assert.True(t, fn.All(regions, func(_ int, region Region) bool {
			return len(region.Holes) == 0
		}))
	})

	t.Run("AreaTest5__", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("OOOOO"))
		testData = append(testData, []byte("OXXXO"))
		testData = append(testData, []byte("OOBBO"))
		testData = append(testData, []byte("OOOOO"))
		testData = append(testData, []byte("OOOBB"))

		regions := RegionsFromGrid(testData)
		for _, region := range regions {
			if region.id != 'O' {
				continue
			}
			assert.Equal(t, 20, region.GetOutsidePerimeter())
			assert.Equal(t, 10, region.GetInsidePerimeter())
		}
	})
}
