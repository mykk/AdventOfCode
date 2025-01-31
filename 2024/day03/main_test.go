package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const EXAMPLEDATA_1 string = "xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))"
const EXAMPLEDATA_2 string = "xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))"

func TestFindGoodReports(t *testing.T) {
	t.Run("TestSumMultiples", func(t *testing.T) {
		multiples := ParseInputData(EXAMPLEDATA_1)
		assert.Equal(t, 161, SumMultiples(multiples))

		multiples = ParseInputDataDoDont(EXAMPLEDATA_2)
		assert.Equal(t, 48, SumMultiples(multiples))
	})
}
