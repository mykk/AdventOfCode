package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const EXAMPLEDATA string = `190: 10 19
3267: 81 40 27
83: 17 5
156: 15 6
7290: 6 8 6 15
161011: 16 10 13
192: 17 8 14
21037: 9 7 18 13
292: 11 6 16 20`

func Test(t *testing.T) {
	t.Run("TestCalibration", func(t *testing.T) {
		entries, err := ParseInputData(EXAMPLEDATA)
		assert.Nil(t, err)

		assert.Equal(t, 3749, CalculateSimpleCalibration(entries))
		assert.Equal(t, 11387, CalculateFancyCalibration(entries))
	})
}
