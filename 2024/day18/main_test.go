package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const EXAMPLEDATA string = `5,4
4,2
4,5
3,0
2,1
6,3
2,4
1,5
0,6
3,3
2,6
5,1
1,2
5,5
2,5
6,5
1,4
0,4
6,4
1,1
6,1
1,0
0,5
1,6
2,0`

func Test(t *testing.T) {
	t.Run("TestMazeRace", func(t *testing.T) {
		fallingBits, err := ParseInputData(EXAMPLEDATA)
		assert.Nil(t, err)

		escape, err := FindPath(fallingBits, 12, Point{x: 0, y: 0}, Point{x: 6, y: 6})
		assert.Nil(t, err)
		assert.Equal(t, 22, escape)

		noEscape, err := FindNoEscapeLimit(fallingBits, 12, Point{x: 0, y: 0}, Point{x: 6, y: 6})
		assert.Nil(t, err)
		assert.Equal(t, Point{x: 6, y: 1}, noEscape)

	})
}
