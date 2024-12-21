package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const EXAMPLEDATA string = `###############
#...#...#.....#
#.#.#.#.#.###.#
#S#...#.#.#...#
#######.#.#.###
#######.#.#...#
#######.#.###.#
###..E#...#...#
###.#######.###
#...###...#...#
#.#####.#.###.#
#.#...#.#.#...#
#.#.#.#.#.#.###
#...#...#...###
###############`

func Test(t *testing.T) {
	t.Run("TestRaceConditionRace", func(t *testing.T) {
		grid, start, end, err := ParseInputData(string(EXAMPLEDATA))
		assert.Nil(t, err)

		assert.Equal(t, 30, CountCheats(grid, start, end, 4, 2))
		assert.Equal(t, 285, CountCheats(grid, start, end, 50, 20))
	})
}
