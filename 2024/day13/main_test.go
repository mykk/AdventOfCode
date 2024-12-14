package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const EXAMPLEDATA string = `Button A: X+94, Y+34
Button B: X+22, Y+67
Prize: X=8400, Y=5400

Button A: X+26, Y+66
Button B: X+67, Y+21
Prize: X=12748, Y=12176

Button A: X+17, Y+86
Button B: X+84, Y+37
Prize: X=7870, Y=6450

Button A: X+69, Y+23
Button B: X+27, Y+71
Prize: X=18641, Y=10279`

func Test(t *testing.T) {
	t.Run("TestClawMachine", func(t *testing.T) {
		machines, err := ParseInputData(EXAMPLEDATA)
		assert.Nil(t, err)

		assert.Equal(t, int64(480), WinPrizes(machines, 0))
	})
}
