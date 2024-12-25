package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const EXAMPLEDATA string = `029A
980A
179A
456A
379A`

func Test(t *testing.T) {
	t.Run("TestKeypadUnlock", func(t *testing.T) {
		codes := ParseInputData(EXAMPLEDATA)

		assert.Equal(t, 126384, GetCodeSumComplexities(codes, 2))
	})
}
