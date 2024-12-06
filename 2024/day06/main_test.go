package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const EXAMPLEDATA string = `....#.....
.........#
..........
..#.......
.......#..
..........
.#..^.....
........#.
#.........
......#...`

func TestUpdates(t *testing.T) {
	t.Run("TestGuardMovement", func(t *testing.T) {
		guardPosition, grid, err := ParseInputData(EXAMPLEDATA)
		assert.Nil(t, err)

		positions, err := TrackGuard(guardPosition, grid, nil)
		assert.Nil(t, err)

		assert.Equal(t, 41, len(positions))
		assert.Equal(t, 6, TrapGuard(guardPosition, positions, grid))
	})
}
