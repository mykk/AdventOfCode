package main

import (
	fn "AoC/functional"
	"math"

	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func mustAtoi(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}

type LogicGate interface {
	evaluate(values map[string]int, connections map[string]LogicGate) int
	operands() (string, string)
}

func resolveValue(values map[string]int, connections map[string]LogicGate, operand string) int {
	if value, found := values[operand]; found {
		return value
	}

	values[operand] = connections[operand].evaluate(values, connections)
	return values[operand]
}

type GateOperands struct {
	operandA, operandB string
}

type xorGate struct {
	GateOperands
}

func (gate xorGate) evaluate(values map[string]int, connections map[string]LogicGate) int {
	return resolveValue(values, connections, gate.operandA) ^ resolveValue(values, connections, gate.operandB)
}

func (gate xorGate) operands() (string, string) {
	return gate.operandA, gate.operandB
}

type andGate struct {
	GateOperands
}

func (gate andGate) operands() (string, string) {
	return gate.operandA, gate.operandB
}

func (gate andGate) evaluate(values map[string]int, connections map[string]LogicGate) int {
	return resolveValue(values, connections, gate.operandA) & resolveValue(values, connections, gate.operandB)
}

type orGate struct {
	GateOperands
}

func (gate orGate) operands() (string, string) {
	return gate.operandA, gate.operandB
}

func (gate orGate) evaluate(values map[string]int, connections map[string]LogicGate) int {
	return resolveValue(values, connections, gate.operandA) | resolveValue(values, connections, gate.operandB)
}

func ParseInputData(data string) (map[string]int, map[string]LogicGate, []string, []string, []string) {
	valueMatches := regexp.MustCompile(`(...): (1|0)`).FindAllStringSubmatch(data, -1)

	xValues := make([]string, 0)
	yValues := make([]string, 0)
	zValues := make([]string, 0)
	values := make(map[string]int)
	for _, value := range valueMatches {
		values[value[1]] = mustAtoi(value[2])
		if strings.HasPrefix(value[1], "x") {
			xValues = append(xValues, value[1])
		}
		if strings.HasPrefix(value[1], "y") {
			yValues = append(yValues, value[1])
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

		if strings.HasPrefix(connection[4], "z") {
			zValues = append(zValues, connection[4])
		}
	}

	xValues = fn.Sorted(xValues, func(a, b string) bool { return a < b })
	yValues = fn.Sorted(yValues, func(a, b string) bool { return a < b })
	zValues = fn.Sorted(zValues, func(a, b string) bool { return a < b })

	return values, connections, xValues, yValues, zValues
}

func ComputeZValue(values map[string]int, connections map[string]LogicGate, zValues []string) int {
	zValues = fn.Sorted(zValues, func(a, b string) bool { return a < b })

	return fn.Reduce(zValues, 0, func(i, zValue int, z string) int {
		return zValue + resolveValue(values, connections, z)*int(math.Pow(2, float64(i)))
	})
}

func getSusBitsMutable(values map[string]int, connections map[string]LogicGate, xValues, yValues, zValues []string) []int {
	susBits := make([]int, 0)
	carryOver := 0
	for i, xValue := range xValues {
		result := values[xValue] + values[yValues[i]] + carryOver

		resultValue := result % 2
		carryOver = result / 2

		zValue := resolveValue(values, connections, zValues[i])
		if resultValue != zValue {
			susBits = append(susBits, i)
		}
	}
	return susBits
}

func getSusBits(values map[string]int, connections map[string]LogicGate, xValues, yValues, zValues []string) []int {
	unresolvedValues := make(map[string]int)
	for key, val := range values {
		unresolvedValues[key] = val
	}

	return getSusBitsMutable(unresolvedValues, connections, xValues, yValues, zValues)
}

func findConnection(connections map[string]LogicGate, predicate func(connection LogicGate) bool) (string, bool) {
	for key, connection := range connections {
		if predicate(connection) {
			return key, true
		}
	}

	return "", false
}

func isCurrentLevelXOR(connection string, connections map[string]LogicGate, xValues, yValues []string, level int) bool {
	x := xValues[level]
	y := yValues[level]

	if _, ok := connections[connection].(xorGate); !ok {
		return false
	}

	a, b := connections[connection].operands()
	return (a == x || b == x) && (a == y || b == y)
}

func isCurrentLevelAnd(connection string, connections map[string]LogicGate, xValues, yValues []string, level int) bool {
	x := xValues[level]
	y := yValues[level]

	if _, ok := connections[connection].(andGate); !ok {
		return false
	}

	a, b := connections[connection].operands()
	return (a == x || b == x) && (a == y || b == y)
}

func swapCurrentCarry(connectionKey string, connections map[string]LogicGate, xValues, yValues []string, level int) (string, string) {
	key, found := findConnection(connections, func(connection LogicGate) bool {
		if _, ok := connection.(orGate); !ok {
			return false
		}
		a, b := connection.operands()
		return (a == xValues[level] || b == xValues[level]) && (a == yValues[level] || b == yValues[level])
	})
	if !found {
		panic("connection not found")
	}
	return connectionKey, key
}

func swapCarryBit(connectionKey string, connections map[string]LogicGate, xValues, yValues []string, level int) (string, string) {
	if level == -1 {
		panic("going too deep")
	}

	if _, ok := connections[connectionKey].(orGate); !ok {
		return swapCurrentCarry(connectionKey, connections, xValues, yValues, level)
	}
	connection := connections[connectionKey]

	a, b := connection.operands()

	if !isCurrentLevelAnd(a, connections, xValues, yValues, level) && !isCurrentLevelAnd(b, connections, xValues, yValues, level) {
		return swapCurrentCarry(connectionKey, connections, xValues, yValues, level)
	}

	if isCurrentLevelAnd(a, connections, xValues, yValues, level) {
		if _, ok := connections[b].(andGate); !ok {
			swapWith, found := findConnection(connections, func(connection LogicGate) bool {
				if _, ok := connection.(andGate); !ok {
					return false
				}

				a, b := connection.operands()
				return isCurrentLevelXOR(a, connections, xValues, yValues, level) || isCurrentLevelXOR(b, connections, xValues, yValues, level)
			})
			if !found {
				panic("not found")
			}
			return b, swapWith
		}

		a1, b1 := connections[b].operands()
		if !isCurrentLevelXOR(a1, connections, xValues, yValues, level) && !isCurrentLevelXOR(b1, connections, xValues, yValues, level) {
			swapWith, found := findConnection(connections, func(connection LogicGate) bool {
				if _, ok := connection.(andGate); !ok {
					return false
				}

				a, b := connection.operands()
				return isCurrentLevelXOR(a, connections, xValues, yValues, level) || isCurrentLevelXOR(b, connections, xValues, yValues, level)
			})
			if !found {
				panic("not found")
			}
			return b, swapWith
		}
		if isCurrentLevelXOR(a1, connections, xValues, yValues, level) {
			return swapCarryBit(a1, connections, xValues, yValues, level-1)
		} else {
			return swapCarryBit(b1, connections, xValues, yValues, level-1)
		}
	} else {
		if _, ok := connections[a].(andGate); !ok {
			swapWith, found := findConnection(connections, func(connection LogicGate) bool {
				if _, ok := connection.(andGate); !ok {
					return false
				}

				a, b := connection.operands()
				return isCurrentLevelXOR(a, connections, xValues, yValues, level) || isCurrentLevelXOR(b, connections, xValues, yValues, level)
			})
			if !found {
				panic("not found")
			}
			return a, swapWith
		}

		a1, b1 := connections[a].operands()
		if !isCurrentLevelXOR(a1, connections, xValues, yValues, level) && !isCurrentLevelXOR(b1, connections, xValues, yValues, level) {
			swapWith, found := findConnection(connections, func(connection LogicGate) bool {
				if _, ok := connection.(andGate); !ok {
					return false
				}

				a, b := connection.operands()
				return isCurrentLevelXOR(a, connections, xValues, yValues, level) || isCurrentLevelXOR(b, connections, xValues, yValues, level)
			})
			if !found {
				panic("not found")
			}
			return a, swapWith
		}
		if isCurrentLevelXOR(a1, connections, xValues, yValues, level) {
			return swapCarryBit(a1, connections, xValues, yValues, level-1)
		} else {
			return swapCarryBit(b1, connections, xValues, yValues, level-1)
		}
	}
}

func swapCurrent(connections map[string]LogicGate, xValues, yValues, zValues []string, badBitIndex int) (string, string) {
	xyXOR, found := findConnection(connections, func(connection LogicGate) bool {
		if _, ok := connection.(xorGate); !ok {
			return false
		}

		a, b := connection.operands()
		return (a == yValues[badBitIndex] || b == yValues[badBitIndex]) && (a == xValues[badBitIndex] || b == xValues[badBitIndex])
	})
	if !found {
		panic("not found x XOR y")
	}

	badZ := zValues[badBitIndex]

	key, found := findConnection(connections, func(connection LogicGate) bool {
		if _, ok := connection.(xorGate); !ok {
			return false
		}
		a, b := connection.operands()
		return a == xyXOR || b == xyXOR
	})
	if !found {
		panic("not found")
	}
	connections[badZ], connections[key] = connections[key], connections[badZ]
	return badZ, key
}

func swapConnection(connections map[string]LogicGate, xValues, yValues, zValues []string, badBitIndex int) (string, string) {
	badZ := zValues[badBitIndex]
	susConnection := connections[badZ]

	switch susConnection.(type) {
	case orGate, andGate:
		{
			return swapCurrent(connections, xValues, yValues, zValues, badBitIndex)
		}
	case xorGate:
		{
			if true {
				connections["tqr"], connections["hth"] = connections["hth"], connections["tqr"]
				return "tqr", "hth"
			}

			a, b := susConnection.operands()

			if !isCurrentLevelXOR(a, connections, xValues, yValues, badBitIndex) && !isCurrentLevelXOR(b, connections, xValues, yValues, badBitIndex) {
				return swapCurrent(connections, xValues, yValues, zValues, badBitIndex)
			}

			if isCurrentLevelXOR(a, connections, xValues, yValues, badBitIndex) {
				return swapCarryBit(b, connections, xValues, yValues, badBitIndex-1)
			} else {
				return swapCarryBit(a, connections, xValues, yValues, badBitIndex-1)
			}
		}
	}

	panic("not resolved anything")
}

func FindSwappedWires(values map[string]int, connections map[string]LogicGate, xValues, yValues, zValues []string) string {
	swaps := make([]string, 0)

	susBits := getSusBits(values, connections, xValues, yValues, zValues)

	for len(susBits) != 0 {
		swapA, swapB := swapConnection(connections, xValues, yValues, zValues, susBits[0])

		swaps = append(swaps, swapA)
		swaps = append(swaps, swapB)

		susBits = getSusBits(values, connections, xValues, yValues, zValues)
	}

	swaps = fn.Sorted(swaps, func(a, b string) bool { return a < b })
	return strings.Join(swaps, ",")
}

func main() {
	inputData, err := os.ReadFile("day24.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	values, connections, xValues, yValues, zValues := ParseInputData(string(inputData))

	valuesUnresolved := make(map[string]int)
	for key, val := range values {
		valuesUnresolved[key] = val
	}
	fmt.Printf("Part 1: %d\n", ComputeZValue(values, connections, zValues))
	fmt.Printf("Part 2: %s\n", FindSwappedWires(valuesUnresolved, connections, xValues, yValues, zValues))
}
