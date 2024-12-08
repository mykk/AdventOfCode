package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const EXAMPLEDATA string = `............
........0...
.....0......
.......0....
....0.......
......A.....
............
............
........A...
.........A..
............
............`

func Test(t *testing.T) {
	t.Run("TestAntinodes", func(t *testing.T) {
		antennas, mapSize, err := ParseInputData(string(EXAMPLEDATA))
		assert.Nil(t, err)

		assert.Equal(t, 14, CountSingleAntinodes(antennas, mapSize))
		assert.Equal(t, 34, CountHarmonicsAntinodes(antennas, mapSize))
	})
}
