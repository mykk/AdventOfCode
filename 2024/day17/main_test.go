package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const EXAMPLEDATA1 string = `Register A: 729
Register B: 0
Register C: 0

Program: 0,1,5,4,3,0`

const EXAMPLEDATA2 string = `Register A: 2024
Register B: 0
Register C: 0

Program: 0,3,5,4,3,0`

func Test(t *testing.T) {
	t.Run("TestChronospatialComputer", func(t *testing.T) {

		registers, instructions, err := ParseInput(EXAMPLEDATA1)
		assert.Nil(t, err)
		assert.Equal(t, "4,6,3,5,6,3,5,2,1,0", SolveChronospatialInstructions(registers, instructions))
	})

	t.Run("TestFindSelf", func(t *testing.T) {
		_, instructions, err := ParseInput(EXAMPLEDATA2)
		assert.Nil(t, err)
		assert.Equal(t, 117440, FindSelf(instructions))
	})
}
