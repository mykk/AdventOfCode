package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const EXAMPLEDATA string = `r, wr, b, g, bwu, rb, gb, br

brwrr
bggr
gbbr
rrbgbr
ubwu
bwurrg
brgr
bbrgwb`

func Test(t *testing.T) {
	t.Run("TestMazeRace", func(t *testing.T) {
		towels, patterns, err := ParseInputData(EXAMPLEDATA)
		assert.Nil(t, err)
		assert.Equal(t, 6, CountPossiblePatterns(towels, patterns))
		assert.Equal(t, 16, CountDifferentCombinations(towels, patterns))
	})
}
