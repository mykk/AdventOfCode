package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const EXAMPLEDATA string = `47|53
97|13
97|61
97|47
75|29
61|13
75|53
29|13
97|29
53|29
61|53
97|53
61|29
47|13
75|47
97|75
47|61
75|61
47|29
75|13
53|13

75,47,61,53,29
97,61,53,29,13
75,29,13
75,97,47,61,53
61,13,29
97,13,75,29,47`

func TestUpdates(t *testing.T) {
	t.Run("TestUpdates", func(t *testing.T) {
		rules, updates, err := ParseInputData(EXAMPLEDATA)
		assert.Nil(t, err)

		assert.Equal(t, 143, CountValidUpdates(rules, updates))
		assert.Equal(t, 123, CountCorrectedUpdates(rules, updates))
	})
}
