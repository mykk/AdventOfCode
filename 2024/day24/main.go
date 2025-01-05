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

func findFirstBadBitMutable(values map[string]int, connections map[string]LogicGate, xValues, yValues, zValues []string) (bit int, found bool) {
	carryOver := 0
	for i, xValue := range xValues {
		result := values[xValue] + values[yValues[i]] + carryOver

		resultValue := result % 2
		carryOver = result / 2

		zValue := resolveValue(values, connections, zValues[i])
		if resultValue != zValue {
			return i, true
		}
	}
	return 0, false
}

func findFirstBadBit(values map[string]int, connections map[string]LogicGate, xValues, yValues, zValues []string) (bit int, found bool) {
	unresolvedValues := make(map[string]int)
	for key, val := range values {
		unresolvedValues[key] = val
	}

	return findFirstBadBitMutable(unresolvedValues, connections, xValues, yValues, zValues)
}

func findConnection(connections map[string]LogicGate, predicate func(connection LogicGate) bool) (string, bool) {
	for key, connection := range connections {
		if predicate(connection) {
			return key, true
		}
	}

	return "", false
}

func isCurrentBitXOR(connection string, connections map[string]LogicGate, xValues, yValues []string, bit int) bool {
	if _, ok := connections[connection].(xorGate); !ok {
		return false
	}

	a, b := connections[connection].operands()
	return (a == xValues[bit] || b == xValues[bit]) && (a == yValues[bit] || b == yValues[bit])
}

func isCurrentBitAND(connection string, connections map[string]LogicGate, xValues, yValues []string, bit int) bool {
	x := xValues[bit]
	y := yValues[bit]

	if _, ok := connections[connection].(andGate); !ok {
		return false
	}

	a, b := connections[connection].operands()
	return (a == x || b == x) && (a == y || b == y)
}

func fixCarryConnection(connectionKey string, connections map[string]LogicGate, xValues, yValues []string, bit int) (string, string) {
	xANDy, found := findCurrentBitAND(connections, xValues, yValues, bit)
	if !found {
		panic("x xor y connection not found")
	}

	goodConnection, found := findConnection(connections, func(connection LogicGate) bool {
		if _, ok := connection.(orGate); !ok {
			return false
		}
		a, b := connection.operands()
		return a == xANDy || b == xANDy
	})
	if !found {
		panic("or connection with x and y not found")
	}

	connections[connectionKey], connections[goodConnection] = connections[goodConnection], connections[connectionKey]
	return connectionKey, goodConnection
}

func findCurrentBitAND(connections map[string]LogicGate, xValues, yValues []string, index int) (string, bool) {
	return findConnection(connections, func(connection LogicGate) bool {
		if _, ok := connection.(andGate); !ok {
			return false
		}

		a, b := connection.operands()
		return (a == yValues[index] || b == yValues[index]) && (a == xValues[index] || b == xValues[index])
	})
}

func fixCurrentBitAND(swapWith string, connections map[string]LogicGate, xValues, yValues []string, index int) (string, string) {
	currentBitAnd, found := findCurrentBitAND(connections, xValues, yValues, index)
	if !found {
		panic("did not find xor connection of current bit x y")
	}
	connections[swapWith], connections[currentBitAnd] = connections[currentBitAnd], connections[swapWith]
	return swapWith, currentBitAnd
}

func fixCarryPart2Connection(connectionKey string, connections map[string]LogicGate, xValues, yValues []string, bit int) (string, string) {
	xXORy, found := findCurrentBitXOR(connections, xValues, yValues, bit)
	if !found {
		panic("x xor y connection not found")
	}

	goodConnection, found := findConnection(connections, func(connection LogicGate) bool {
		if _, ok := connection.(andGate); !ok {
			return false
		}
		a, b := connection.operands()
		return a == xXORy || b == xXORy
	})
	if !found {
		panic("and connection with x xor y not found")
	}

	connections[connectionKey], connections[goodConnection] = connections[goodConnection], connections[connectionKey]
	return connectionKey, goodConnection
}

func fixCarryBitPart2(connectionKey string, connections map[string]LogicGate, xValues, yValues []string, bit int) (string, string) {
	if _, ok := connections[connectionKey].(andGate); !ok {
		return fixCarryPart2Connection(connectionKey, connections, xValues, yValues, bit)
	}

	a, b := connections[connectionKey].operands()

	if isCurrentBitXOR(a, connections, xValues, yValues, bit) {
		return fixCarryBit(b, connections, xValues, yValues, bit-1)
	}

	if isCurrentBitXOR(b, connections, xValues, yValues, bit) {
		return fixCarryBit(a, connections, xValues, yValues, bit-1)
	}

	if isCorrectCarryBitConnection(a, connections, xValues, yValues, bit-1) {
		return fixCurrentBitXor(b, connections, xValues, yValues, bit)
	}

	if isCorrectCarryBitConnection(b, connections, xValues, yValues, bit-1) {
		return fixCurrentBitXor(a, connections, xValues, yValues, bit)
	}

	return fixCarryPart2Connection(connectionKey, connections, xValues, yValues, bit)
}

func fixCarryBit(connectionKey string, connections map[string]LogicGate, xValues, yValues []string, bit int) (string, string) {
	if bit == 0 {
		panic("nothing to fix here")
	}

	if _, ok := connections[connectionKey].(orGate); !ok {
		return fixCarryConnection(connectionKey, connections, xValues, yValues, bit)
	}
	connection := connections[connectionKey]

	a, b := connection.operands()

	if isCurrentBitAND(a, connections, xValues, yValues, bit) {
		return fixCarryBitPart2(b, connections, xValues, yValues, bit)
	}

	if isCurrentBitAND(b, connections, xValues, yValues, bit) {
		return fixCarryBitPart2(a, connections, xValues, yValues, bit)
	}

	if isCorrectCarryBitPart2(connections[a], connections, xValues, yValues, bit) {
		return fixCurrentBitAND(b, connections, xValues, yValues, bit)
	}

	if isCorrectCarryBitPart2(connections[b], connections, xValues, yValues, bit) {
		return fixCurrentBitAND(a, connections, xValues, yValues, bit)
	}

	return fixCarryConnection(connectionKey, connections, xValues, yValues, bit)
}

func findCurrentBitXOR(connections map[string]LogicGate, xValues, yValues []string, index int) (string, bool) {
	return findConnection(connections, func(connection LogicGate) bool {
		if _, ok := connection.(xorGate); !ok {
			return false
		}

		a, b := connection.operands()
		return (a == yValues[index] || b == yValues[index]) && (a == xValues[index] || b == xValues[index])
	})
}

func fixCurrentBitXor(swapWith string, connections map[string]LogicGate, xValues, yValues []string, index int) (string, string) {
	currentBitXor, found := findCurrentBitXOR(connections, xValues, yValues, index)
	if !found {
		panic("did not find xor connection of current bit x y")
	}
	connections[swapWith], connections[currentBitXor] = connections[currentBitXor], connections[swapWith]
	return swapWith, currentBitXor
}

func isCorrectCarryBitPart2(connection LogicGate, connections map[string]LogicGate, xValues, yValues []string, index int) bool {
	if _, ok := connection.(andGate); !ok {
		return false
	}

	a, b := connection.operands()

	if isCurrentBitXOR(a, connections, xValues, yValues, index) {
		return isCorrectCarryBitConnection(b, connections, xValues, yValues, index-1)
	}

	if isCurrentBitXOR(b, connections, xValues, yValues, index) {
		return isCorrectCarryBitConnection(a, connections, xValues, yValues, index-1)
	}

	return false
}

func isCorrectCarryBitConnection(connection string, connections map[string]LogicGate, xValues, yValues []string, index int) bool {
	if index == 0 {
		return isCurrentBitAND(connection, connections, xValues, yValues, index)
	}

	if _, ok := connections[connection].(orGate); !ok {
		return false
	}

	a, b := connections[connection].operands()

	if isCurrentBitAND(a, connections, xValues, yValues, index) {
		return isCorrectCarryBitPart2(connections[b], connections, xValues, yValues, index)
	}

	if isCurrentBitAND(b, connections, xValues, yValues, index) {
		return isCorrectCarryBitPart2(connections[a], connections, xValues, yValues, index)
	}

	return false
}

func fixBadZConnection(connections map[string]LogicGate, xValues, yValues, zValues []string, badBitIndex int) (string, string) {
	xXORy, found := findCurrentBitXOR(connections, xValues, yValues, badBitIndex)
	if !found {
		panic("x xor y connection not found")
	}

	goodZ, found := findConnection(connections, func(connection LogicGate) bool {
		if _, ok := connection.(xorGate); !ok {
			return false
		}
		a, b := connection.operands()
		return a == xXORy || b == xXORy
	})
	if !found {
		panic("xor connection with x xor y not found")
	}

	badZ := zValues[badBitIndex]
	connections[badZ], connections[goodZ] = connections[goodZ], connections[badZ]
	return badZ, goodZ
}

func fixConnection(connections map[string]LogicGate, xValues, yValues, zValues []string, badBitIndex int) (string, string) {
	badZ := zValues[badBitIndex]
	badConnection := connections[badZ]

	if _, ok := badConnection.(xorGate); !ok {
		return fixBadZConnection(connections, xValues, yValues, zValues, badBitIndex)
	}

	a, b := badConnection.operands()

	if isCurrentBitXOR(a, connections, xValues, yValues, badBitIndex) {
		return fixCarryBit(b, connections, xValues, yValues, badBitIndex-1)
	}

	if isCurrentBitXOR(b, connections, xValues, yValues, badBitIndex) {
		return fixCarryBit(a, connections, xValues, yValues, badBitIndex-1)
	}

	if isCorrectCarryBitConnection(a, connections, xValues, yValues, badBitIndex-1) {
		return fixCurrentBitXor(b, connections, xValues, yValues, badBitIndex)
	}

	if isCorrectCarryBitConnection(b, connections, xValues, yValues, badBitIndex-1) {
		return fixCurrentBitXor(a, connections, xValues, yValues, badBitIndex)
	}

	return fixBadZConnection(connections, xValues, yValues, zValues, badBitIndex)
}

func FindSwappedWires(values map[string]int, connections map[string]LogicGate, xValues, yValues, zValues []string) string {
	swaps := make([]string, 0)

	for {
		if bit, found := findFirstBadBit(values, connections, xValues, yValues, zValues); found {
			swapA, swapB := fixConnection(connections, xValues, yValues, zValues, bit)
			swaps = append(swaps, swapA, swapB)
		} else {
			return strings.Join(fn.Sorted(swaps, func(a, b string) bool { return a < b }), ",")
		}
	}
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
