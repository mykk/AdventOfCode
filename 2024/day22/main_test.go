package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const EXAMPLEDATA1 string = `1
10
100
2024`

const EXAMPLEDATA2 string = `1
2
3
2024`

func Test(t *testing.T) {
	t.Run("TestSecretNumberMangling", func(t *testing.T) {
		numbers, err := ParseInputData(EXAMPLEDATA1)
		assert.Nil(t, err)

		assert.Equal(t, 37327623, MangleSecretNumbers(numbers))
	})

	t.Run("TestMaxBananas", func(t *testing.T) {
		numbers, err := ParseInputData(EXAMPLEDATA2)
		assert.Nil(t, err)

		assert.Equal(t, 23, GetTheMostBananas(numbers))
	})
}
