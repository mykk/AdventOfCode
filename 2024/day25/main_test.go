package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const EXAMPLEDATA string = `#####
.####
.####
.####
.#.#.
.#...
.....

#####
##.##
.#.##
...##
...#.
...#.
.....

.....
#....
#....
#...#
#.#.#
#.###
#####

.....
.....
#.#..
###..
###.#
###.#
#####

.....
.....
.....
#....
#.#..
#.#.#
#####`

func Test(t *testing.T) {
	t.Run("CountUnlocks", func(t *testing.T) {
		keys, locks, lockSize := ParseInputData(EXAMPLEDATA)

		assert.Equal(t, 3, CountFittingKeys(keys, locks, lockSize))
	})
}
