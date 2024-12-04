package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const EXAMPLEDATA_1 string = `MMMSXXMASM
MSAMXMSMSA
AMXSXMAAMM
MSAMASMSMX
XMASAMXAMM
XXAMMXXAMA
SMSMSASXSS
SAXAMASAAA
MAMMMXMMMM
MXMXAXMASX`

func TestFindXmas(t *testing.T) {
	t.Run("TestFindXmas", func(t *testing.T) {
		lines := ParseInputData(EXAMPLEDATA_1)

		assert.Equal(t, 18, FindXmas(lines))
		assert.Equal(t, 9, Find_X_MAS(lines))
	})
}
