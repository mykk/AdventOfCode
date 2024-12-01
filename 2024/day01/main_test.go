package main

import (
   fn "AoC/functional"
   "testing"

   "github.com/stretchr/testify/assert"
)

const EXAMPLEDATA string = `3   4
4   3
2   5
1   3
3   9
3   3`

func TestParseInputData(t *testing.T) {
   t.Run("TestParseInputData", func(t *testing.T) {
      left, right, err := ParseInputData(fn.GetLines(string(EXAMPLEDATA)))
      assert.Nil(t, err)
      assert.Equal(t, []int{3, 4, 2, 1, 3, 3}, left)
      assert.Equal(t, []int{4, 3, 5, 3, 9, 3}, right)
   })
}

func TestFindDistnace(t *testing.T) {
   t.Run("TestTestFindDistnace", func(t *testing.T) {
      assert.Equal(t, 11, FindDistnace([]int{3, 4, 2, 1, 3, 3}, []int{4, 3, 5, 3, 9, 3}))
   })
}

func TestFindSimilarityCount(t *testing.T) {
   t.Run("TestFindSimilarityCount", func(t *testing.T) {
      assert.Equal(t, 31, FindSimilarityCount([]int{3, 4, 2, 1, 3, 3}, []int{4, 3, 5, 3, 9, 3}))
   })
}
