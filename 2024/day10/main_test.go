package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const EXAMPLEDATA_1 string = `012345
123456
234567
345678
4.6789
56789.`

const EXAMPLEDATA_2 string = `89010123
78121874
87430965
96549874
45678903
32019012
01329801
10456732`

func Test1(t *testing.T) {
	t.Run("TestTrails1", func(t *testing.T) {
		grid := ParseInputData(EXAMPLEDATA_1)

		trailCount, trailScore := CountTrails(grid)
		assert.Equal(t, 2, trailCount)
		assert.Equal(t, 227, trailScore)
	})
}

func Test2(t *testing.T) {
	t.Run("TestTrails2", func(t *testing.T) {
		grid := ParseInputData(EXAMPLEDATA_2)

		trailCount, trailScore := CountTrails(grid)
		assert.Equal(t, 36, trailCount)
		assert.Equal(t, 81, trailScore)
	})
}
