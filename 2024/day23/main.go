package main

import (
	fn "AoC/functional"

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

func intersect(splice1, splice2 []string) []string {
	set1, set2 := spliceToSet(splice1), spliceToSet(splice2)

	intersected := make([]string, 0, len(set1))
	for item := range set1 {
		if set2.Contains(item) {
			intersected = append(intersected, item)
		}
	}

	return intersected
}

func CountTConnections(network map[string][]string) (count int) {
	usedGlobal := make(aoc.Set[string])

	for computer, connections := range network {
		if !strings.HasPrefix(computer, "t") {
			continue
		}

		usedCurrent := make(aoc.Set[string])
		for _, connected := range connections {
			if usedGlobal.Contains(connected) {
				continue
			}

			count += fn.CountIf(intersect(connections, network[connected]), func(interConnected string) bool {
				return !usedCurrent.Contains(interConnected) && !usedGlobal.Contains(interConnected)
			})

			usedCurrent.Add(connected)
		}
		usedGlobal.Add(computer)
	}
	return count
}

func getMaxClusterSize(network map[string][]string, cluster []string, seen aoc.Set[string]) []string {
	currentAsStr := strings.Join(fn.Sorted(cluster, func(a, b string) bool { return a < b }), "")
	if seen.Contains(currentAsStr) {
		return []string{}
	}
	seen.Add(currentAsStr)

	clusterConnections := network[cluster[0]]
	for _, clusterNode := range cluster[1:] {
		clusterConnections = intersect(clusterConnections, network[clusterNode])
	}

	if len(clusterConnections) == 0 {
		return cluster
	}

	maxCluster := cluster
	for _, connection := range clusterConnections {
		current := getMaxClusterSize(network, append(append(make([]string, 0, len(cluster)+1), cluster...), connection), seen)
		if len(current) > len(maxCluster) {
			maxCluster = current
		}
	}

	return maxCluster
}

func FindMaxCluster(network map[string][]string) string {
	maxCluster := make([]string, 0)
	seen := make(aoc.Set[string])
	for computer := range network {
		current := getMaxClusterSize(network, []string{computer}, seen)
		if len(current) > len(maxCluster) {
			maxCluster = current
		}
	}

	return strings.Join(fn.Sorted(maxCluster, func(a, b string) bool { return a < b }), ",")
}

func ParseInputData(data string) map[string][]string {
	re := regexp.MustCompile(`([a-z]{2})\-([a-z]{2})`)
	matches := re.FindAllStringSubmatch(data, -1)

	network := make(map[string][]string)
	for _, match := range matches {
		network[match[1]] = append(network[match[1]], match[2])
		network[match[2]] = append(network[match[2]], match[1])
	}
	return network
}

func main() {
	inputData, err := os.ReadFile("day23.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	network := ParseInputData(string(inputData))
	if err != nil {
		fmt.Printf("Error parsing input data: %v\n", err)
		return
	}

	fmt.Printf("Part 1: %d\n", CountTConnections(network))
	fmt.Printf("Part 2: %s\n", FindMaxCluster(network))
}
