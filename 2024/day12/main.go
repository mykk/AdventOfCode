package main

import (
	fn "AoC/functional"
	geom "AoC/geometry"
	"fmt"
	"os"
)

func AbsInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func getPerimeter(polygon []geom.Point) int {
	return fn.Reduce(polygon[:len(polygon)-1], 0, func(i, perimeter int, point geom.Point) int {
		nextPoint := polygon[i+1]
		return perimeter + AbsInt(point.X-nextPoint.X) + AbsInt(point.Y-nextPoint.Y)
	})
}

func getTotalPerimeter(area geom.Area) int {
	return getPerimeter(area.Perimeter) + fn.Reduce(area.Holes, 0, func(_, perimeter int, hole geom.Hole) int { return perimeter + getPerimeter(hole.Perimeter) })
}

func CalculatePrice(areas []geom.Area) int {
	return fn.Reduce(areas, 0, func(_, price int, area geom.Area) int {
		perimeter := getTotalPerimeter(area)
		return price + len(area.Area)*perimeter + fn.Reduce(area.Holes, 0, func(_, holePrice int, hole geom.Hole) int {
			return holePrice + len(hole.Area)*getPerimeter(hole.Perimeter)
		})
	})
}

func ParseInputData(data string) []geom.Area {
	garden := fn.MustTransform(fn.GetLines(data), func(line string) []byte { return []byte(line) })
	return geom.AreasFromGrid(garden)
}

func main() {
	inputData, err := os.ReadFile("day12.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	areas := ParseInputData(string(inputData))
	if err != nil {
		fmt.Printf("Error parsing input data: %v\n", err)
		return
	}

	fmt.Printf("Part 1: %d\n", CalculatePrice(areas))
}
