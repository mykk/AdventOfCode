package main

import (
	fn "AoC/functional"
	"testing"

	"github.com/stretchr/testify/assert"
)

const EXAMPLEDATA string = `7 6 4 2 1
1 2 7 8 9
9 7 6 2 1
1 3 2 4 5
8 6 4 4 1
1 3 6 7 9`

func TestParseInputData(t *testing.T) {
	t.Run("TestParseInputData", func(t *testing.T) {
		report, err := ParseInputData(fn.GetLines(string(EXAMPLEDATA)))
		assert.Nil(t, err)
		assert.Equal(t, []int{7, 6, 4, 2, 1}, report[0])
		assert.Equal(t, []int{1, 3, 6, 7, 9}, report[5])
	})
}

func TestFindGoodReports(t *testing.T) {
	t.Run("TestFindGoodReports", func(t *testing.T) {
		report, err := ParseInputData(fn.GetLines(string(EXAMPLEDATA)))
		assert.Nil(t, err)

		assert.Equal(t, 2, FindGoodReports(report, 0))
		assert.Equal(t, 4, FindGoodReports(report, 1))
	})
}
