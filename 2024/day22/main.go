package main

import (
	"AoC/aoc"
	fn "AoC/functional"
	"fmt"
	"os"
	"strconv"
)

func mangleNumber(number int, onSecretNumberGenerated func(secretNumber int)) int {
	for i := 0; i < 2000; i++ {
		number = ((number * 64) ^ number) % 16777216
		number = ((number / 32) ^ number) % 16777216
		number = ((number * 2048) ^ number) % 16777216
		onSecretNumberGenerated(number)
	}
	return number
}

func MangleSecretNumbers(numbers []int) int {
	return fn.Reduce(numbers, 0, func(_, sum, current int) int { return sum + mangleNumber(current, func(int) {}) })
}

type Sequence struct {
	a, b, c, d int
}

func GetTheMostBananas(numbers []int) int {
	allSequences := make(map[Sequence]int)

	for _, current := range numbers {
		currentSequences := make(aoc.Set[Sequence])

		sequenceBuilder := make([]int, 0, 4)
		lastNumber := current % 10
		mangleNumber(current, func(secretNumber int) {
			delta := lastNumber - secretNumber%10
			lastNumber = secretNumber % 10

			if len(sequenceBuilder) != 4 {
				sequenceBuilder = append(sequenceBuilder, delta)
				return
			}

			sequenceBuilder[0], sequenceBuilder[1], sequenceBuilder[2], sequenceBuilder[3] =
				sequenceBuilder[1], sequenceBuilder[2], sequenceBuilder[3], delta

			sequence := Sequence{sequenceBuilder[0], sequenceBuilder[1], sequenceBuilder[2], sequenceBuilder[3]}
			if !currentSequences.Contains(sequence) {
				currentSequences.Add(sequence)
				if bananas, found := allSequences[sequence]; found {
					allSequences[sequence] = bananas + lastNumber
				} else {
					allSequences[sequence] = lastNumber
				}
			}
		})
	}

	maxBananas := 0
	for _, bananas := range allSequences {
		if bananas > maxBananas {
			maxBananas = bananas
		}
	}
	return maxBananas
}

func ParseInputData(data string) ([]int, error) {
	return fn.Transform(fn.GetLines(data), strconv.Atoi)
}

func main() {
	inputData, err := os.ReadFile("day22.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	numbers, err := ParseInputData(string(inputData))
	if err != nil {
		fmt.Printf("Error parsing input data: %v\n", err)
		return
	}

	fmt.Printf("Part 1: %d\n", MangleSecretNumbers(numbers))
	fmt.Printf("Part 2: %d\n", GetTheMostBananas(numbers))
}
