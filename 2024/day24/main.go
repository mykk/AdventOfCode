package main

import (
	fn "AoC/functional"
	"math"

	"AoC/aoc"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func MustAtoi(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}

type LogicGate interface {
	evaluate(values map[string]int, connections map[string]LogicGate) int
}

type GateOperands struct {
	operandA, operandB string
}

func resolveValue(values map[string]int, connections map[string]LogicGate, operand string) int {
	if value, found := values[operand]; found {
		return value
	}

	values[operand] = connections[operand].evaluate(values, connections)
	return values[operand]
}

type xorGate struct {
	GateOperands
}

func (gate xorGate) evaluate(values map[string]int, connections map[string]LogicGate) int {
	return resolveValue(values, connections, gate.operandA) ^ resolveValue(values, connections, gate.operandB)
}

type andGate struct {
	GateOperands
}

func (gate andGate) evaluate(values map[string]int, connections map[string]LogicGate) int {
	return resolveValue(values, connections, gate.operandA) & resolveValue(values, connections, gate.operandB)
}

type orGate struct {
	GateOperands
}

func (gate orGate) evaluate(values map[string]int, connections map[string]LogicGate) int {
	return resolveValue(values, connections, gate.operandA) | resolveValue(values, connections, gate.operandB)
}

func setToSplice(set aoc.Set[string]) []string {
	splice := make([]string, 0, len(set))
	for item := range set {
		splice = append(splice, item)
	}
	return splice
}

func ParseInputData(data string) (map[string]int, map[string]LogicGate, []string, []string, []string) {
	valueMatches := regexp.MustCompile(`(...): (1|0)`).FindAllStringSubmatch(data, -1)

	xValues := make(aoc.Set[string])
	yValues := make(aoc.Set[string])
	zValues := make(aoc.Set[string])
	values := make(map[string]int)
	for _, value := range valueMatches {
		values[value[1]] = MustAtoi(value[2])
		if strings.HasPrefix(value[1], "x") {
			xValues.Add(value[1])
		}
		if strings.HasPrefix(value[1], "y") {
			yValues.Add(value[1])
		}
		if strings.HasPrefix(value[1], "z") {
			zValues.Add(value[1])
		}

	}

	connectionMatches := regexp.MustCompile(`(...) (OR|XOR|AND) (...) \-> (...)`).FindAllStringSubmatch(data, -1)

	connections := make(map[string]LogicGate)
	for _, connection := range connectionMatches {
		switch connection[2] {
		case "OR":
			connections[connection[4]] = orGate{GateOperands: GateOperands{operandA: connection[1], operandB: connection[3]}}
		case "AND":
			connections[connection[4]] = andGate{GateOperands: GateOperands{operandA: connection[1], operandB: connection[3]}}
		case "XOR":
			connections[connection[4]] = xorGate{GateOperands: GateOperands{operandA: connection[1], operandB: connection[3]}}
		}
		if strings.HasPrefix(connection[4], "x") {
			xValues.Add(connection[4])
		}
		if strings.HasPrefix(connection[4], "y") {
			yValues.Add(connection[4])
		}

		if strings.HasPrefix(connection[4], "z") {
			zValues.Add(connection[4])
		}
	}

	return values, connections, setToSplice(xValues), setToSplice(yValues), setToSplice(zValues)
}

func ComputeZValue(values map[string]int, connections map[string]LogicGate, zValues []string) int {
	zValues = fn.Sorted(zValues, func(a, b string) bool { return a < b })

	return fn.Reduce(zValues, 0, func(i, zValue int, z string) int {
		return zValue + resolveValue(values, connections, z)*int(math.Pow(2, float64(i)))
	})
}

func FindSwappedWires(values map[string]int, connections map[string]LogicGate, xValues, yValues, zValues []string) string {
	xValues = fn.Sorted(xValues, func(a, b string) bool { return a < b })
	yValues = fn.Sorted(yValues, func(a, b string) bool { return a < b })
	zValues = fn.Sorted(zValues, func(a, b string) bool { return a < b })

	fmt.Printf(" ")
	for i := len(xValues) - 1; i != -1; i-- {
		fmt.Printf("%d", values[xValues[i]])
	}

	fmt.Printf("\n ")

	for i := len(yValues) - 1; i != -1; i-- {
		fmt.Printf("%d", values[yValues[i]])
	}
	fmt.Printf("\n")

	for i := len(zValues) - 1; i != -1; i-- {
		zValue := resolveValue(values, connections, zValues[i])
		fmt.Printf("%d", zValue)
	}
	fmt.Printf("\n")

	carryOver := 0
	for i, xValue := range xValues {
		result := values[xValue] + values[yValues[i]] + carryOver

		resultValue := result % 2
		carryOver = result / 2

		zValue := resolveValue(values, connections, zValues[i])
		if resultValue != zValue {
			fmt.Printf("SusBit: %d\n", i+1)
		}
	}

	return ""
}

func main() {
	inputData, err := os.ReadFile("day24.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	values, connections, xValues, yValues, zValues := ParseInputData(string(inputData))

	fmt.Printf("Part 1: %d\n", ComputeZValue(values, connections, zValues))
	fmt.Printf("Part 2: %s\n", FindSwappedWires(values, connections, xValues, yValues, zValues))
}
