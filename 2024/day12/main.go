package main

import (
	fn "AoC/functional"
	geom "AoC/geometry"
	"fmt"
	"os"
)

func CalculatePrice(regions []geom.Region) int {
	return fn.Reduce(regions, 0, func(_, price int, region geom.Region) int {
		perimeter := region.GetInsidePerimeter() + region.GetOutsidePerimeter()

		return price + len(region.Area)*perimeter + fn.Reduce(region.Holes, 0, func(_, holePrice int, hole geom.Hole) int {
			return holePrice + len(hole.Area)*hole.GetPerimeter()
		})
	})
}

func CalculateDiscountPrice(regions []geom.Region) int {
	return fn.Reduce(regions, 0, func(_, price int, region geom.Region) int {
		perimeter := len(region.Perimeter) - 1 + fn.Reduce(region.InsidePerimeters, 0, func(_, sum int, points []geom.Point) int { return sum + len(points) - 1 })

		return price + len(region.Area)*perimeter + fn.Reduce(region.Holes, 0, func(_, holePrice int, hole geom.Hole) int {
			return holePrice + len(hole.Area)*(len(hole.Perimeter)-1)
		})
	})
}

func ParseInputData(data string) []geom.Region {
	garden := fn.MustTransform(fn.GetLines(data), func(line string) []byte { return []byte(line) })
	return geom.RegionsFromGrid(garden)
}

func main() {
	inputData, err := os.ReadFile("day12.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	regions := ParseInputData(string(inputData))

	fmt.Printf("Part 1: %d\n", CalculatePrice(regions))
	fmt.Printf("Part 2: %d\n", CalculateDiscountPrice(regions))
}
