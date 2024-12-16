package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const EXAMPLEDATA string = `#################
#...#...#...#..E#
#.#.#.#.#.#.#.#.#
#.#.#.#...#...#.#
#.#.#.#.###.#.#.#
#...#.#.#.....#.#
#.#.#.#.#.#####.#
#.#...#.#.#.....#
#.#.#####.#.###.#
#.#.#.......#...#
#.#.###.#####.###
#.#.#...#.....#.#
#.#.#.#####.###.#
#.#.#.........#.#
#.#.#.#########.#
#S#.............#
#################`

func Test(t *testing.T) {
	t.Run("TestMazeRace", func(t *testing.T) {
		grid, start, end, err := ParseInputData(EXAMPLEDATA)
		assert.Nil(t, err)
		winningPaths := SolveMaze(grid, start, end)

		assert.NotEqual(t, 0, len(winningPaths))

		assert.Equal(t, 11048, winningPaths[0].points)
		assert.Equal(t, 64, HotPaths(winningPaths))
	})
}
