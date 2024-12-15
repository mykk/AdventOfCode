package aoc_geometry

import (
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

		assert.Equal(t, []Point{{0, 1}, {1, 1}, {1, 0}, {2, 0}, {2, 1}, {3, 1}, {3, 2}, {0, 2}, {0, 1}}, perimeter)
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

		assert.Equal(t, []Point{{1, 1}, {5, 1}, {5, 2}, {0, 2}, {0, 0}, {1, 0}, {1, 1}}, perimeter)
	})

	t.Run("PerimeterTest10", func(t *testing.T) {
		testData := [][]byte{}
		testData = append(testData, []byte("A    "))
		testData = append(testData, []byte(" AAAA"))

		perimeter := WalkPerimeter(Point{3, 1}, testData)

		assert.Equal(t, []Point{{1, 1}, {5, 1}, {5, 2}, {1, 2}, {1, 1}}, perimeter)
	})
}
