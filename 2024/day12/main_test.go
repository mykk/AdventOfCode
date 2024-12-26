package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const EXAMPLEDATA string = `RRRRIICCFF
RRRRIICCCF
VVRRRCCFFF
VVRCCCJFFF
VVVVCJJCFE
VVIVCCJJEE
VVIIICJJEE
MIIIIIJJEE
MIIISIJEEE
MMMISSJEEE`

func Test(t *testing.T) {
	t.Run("CountPebbles", func(t *testing.T) {

		regions := ParseInputData(EXAMPLEDATA)

		assert.Equal(t, 1930, CalculatePrice(regions))
		assert.Equal(t, 1206, CalculateDiscountPrice(regions))
	})
}
