package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const EXAMPLEDATA string = `125 17`

func Test(t *testing.T) {
	t.Run("CountPebbles", func(t *testing.T) {
		stones, err := ParseInputData(EXAMPLEDATA)
		assert.Nil(t, err)

		assert.Equal(t, 55312, CountPlutonianPebbles(stones, 25))
	})
}
