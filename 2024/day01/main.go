package main

import (
   fn "AoC/functional"
   "errors"
   "fmt"
   "os"
   "strconv"
   "strings"
)

func AbsInt(n int) int {
   if n < 0 {
      return -n
   }
   return n
}

func ParseInputData(data []string) ([]int, []int, error) {
   leftList, rightList := []int{}, []int{}
   var err error = nil

   for _, line := range data {
      split_line := strings.Fields(line)
      if len(split_line) != 2 {
         return nil, nil, errors.New("unexpected data in data file")
      }

      leftVal, rightVal := split_line[0], split_line[1]

      if leftList, err = fn.TransformAppend(leftList, leftVal, func(value string) (int, error) { return strconv.Atoi(value) }); err != nil {
         return nil, nil, err
      }

      if rightList, err = fn.TransformAppend(rightList, rightVal, func(value string) (int, error) { return strconv.Atoi(value) }); err != nil {
         return nil, nil, err
      }
   }

   return leftList, rightList, nil
}

func FindDistnace(leftList []int, rightList []int) int {
   leftList = fn.Sorted(leftList, func(i, j int) bool { return i < j })
   rightList = fn.Sorted(rightList, func(i, j int) bool { return i < j })

   return fn.Reduce(leftList, 0, func(index, lhs, rhs int) int {
      return lhs + AbsInt(rhs-rightList[index])
   })
}

func FindSimilarityCount(leftList []int, rightList []int) int {
   return fn.Reduce(leftList, 0, func(_, lhs, rhs int) int {
      return lhs + rhs*fn.CountIf(rightList, func(val int) bool { return val == rhs })
   })
}

func main() {
   inputData, err := os.ReadFile("day1.txt")
   if err != nil {
      fmt.Printf("Error reading file: %v\n", err)
      return
   }
   data := fn.GetLines(string(inputData))

   leftList, rightList, err := ParseInputData(data)
   if err != nil {
      fmt.Printf("Error parsing: %v\n", err)
      return
   }

   fmt.Printf("part 1: %v\n", FindDistnace(leftList, rightList))
   fmt.Printf("part 2: %v\n", FindSimilarityCount(leftList, rightList))
}
