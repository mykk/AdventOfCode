package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const EXAMPLEDATA string = `2333133121414131402`

func Test(t *testing.T) {
	t.Run("TestAntinodes", func(t *testing.T) {
		entries, err := ParseInputData(EXAMPLEDATA)
		assert.Nil(t, err)

		assert.Equal(t, 1928, FormatFileSingles(entries))
		assert.Equal(t, 2858, FormatFileChunks(entries))
	})
}
