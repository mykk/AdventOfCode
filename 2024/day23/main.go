package main

import (
	"AoC/aoc"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func spliceToSet(splice []string) aoc.Set[string] {
	set := make(aoc.Set[string])
	for _, item := range splice {
		set.Add(item)
	}
	return set
}

func intersect(splice1, splice2 []string) aoc.Set[string] {
	set1, set2 := spliceToSet(splice1), spliceToSet(splice2)

	set := make(aoc.Set[string])

	for item := range set1 {
		if set2.Contains(item) {
			set.Add(item)
		}
	}

	return set
}

func CountTConnections(connections map[string][]string) (count int) {
	usedGlobal := make(aoc.Set[string])

	for computer, computerConnections := range connections {
		if !strings.HasPrefix(computer, "t") {
			continue
		}

		used := make(aoc.Set[string])
		for _, connected := range computerConnections {
			if usedGlobal.Contains(connected) {
				continue
			}

			interConnections := intersect(computerConnections, connections[connected])
			for interConnected := range interConnections {
				if !used.Contains(interConnected) && !usedGlobal.Contains(interConnected) {
					count += 1
				}
			}

			used.Add(connected)
		}
		usedGlobal.Add(computer)
	}
	return count
}

func ParseInputData(data string) map[string][]string {
	re := regexp.MustCompile(`([a-z]{2})\-([a-z]{2})`)
	matches := re.FindAllStringSubmatch(data, -1)

	connections := make(map[string][]string)
	for _, match := range matches {
		connections[match[1]] = append(connections[match[1]], match[2])
		connections[match[2]] = append(connections[match[2]], match[1])
	}
	return connections
}

func main() {
	inputData, err := os.ReadFile("day23.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	connections := ParseInputData(string(inputData))
	if err != nil {
		fmt.Printf("Error parsing input data: %v\n", err)
		return
	}

	fmt.Printf("Part 1: %d\n", CountTConnections(connections))
}
