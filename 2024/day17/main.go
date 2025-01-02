package main

import (
	fn "AoC/functional"
	"errors"
	"fmt"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func mustAtoi(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}

func getValue(v int, registers map[byte]int) int {
	if v <= 3 {
		return v
	}
	return registers[byte('A'+v%4)]
}

func createInstructionMap(output func(int)) map[int]func(int, map[byte]int, int) int {
	instructionMap := make(map[int]func(int, map[byte]int, int) int)
	instructionMap[0] = func(i int, reg map[byte]int, v int) int {
		reg['A'] = int(float64(reg['A']) / math.Pow(2, float64(getValue(v, reg))))
		return i
	}

	instructionMap[1] = func(i int, reg map[byte]int, v int) int {
		reg['B'] = reg['B'] ^ v
		return i
	}

	instructionMap[2] = func(i int, reg map[byte]int, v int) int {
		reg['B'] = ((getValue(v, reg) % 8) + 8) % 8
		return i
	}

	instructionMap[3] = func(i int, reg map[byte]int, v int) int {
		if reg['A'] == 0 {
			return i
		}
		return v - 2
	}

	instructionMap[4] = func(i int, reg map[byte]int, _ int) int {
		reg['B'] = reg['B'] ^ reg['C']
		return i
	}

	instructionMap[5] = func(i int, reg map[byte]int, v int) int {
		output(getValue(v, reg) % 8)
		return i
	}

	instructionMap[6] = func(i int, reg map[byte]int, v int) int {
		reg['B'] = int(float64(reg['A']) / math.Pow(2, float64(getValue(v, reg))))
		return i
	}

	instructionMap[7] = func(i int, reg map[byte]int, v int) int {
		reg['C'] = int(float64(reg['A']) / math.Pow(2, float64(getValue(v, reg))))
		return i
	}

	return instructionMap
}

func solveInstructions(registers map[byte]int, instructions []int, instructionMap map[int]func(int, map[byte]int, int) int) {
	for i := 0; i < len(instructions); i += 2 {
		i = instructionMap[instructions[i]](i, registers, instructions[i+1])
	}
}

func SolveChronospatialInstructions(registers map[byte]int, instructions []int) string {
	output := make([]int, 0)
	instructionMap := createInstructionMap(func(value int) { output = append(output, value) })

	solveInstructions(registers, instructions, instructionMap)

	temp := strings.Join(fn.MustTransform(output, strconv.Itoa), ",")
	return temp
}

func reversed(s []int) []int {
	r := make([]int, 0, len(s))
	r = append(r, s...)
	slices.Reverse(r)
	return r
}

func getOrDefault(i int, m map[int]int) int {
	if val, ok := m[i]; ok {
		return val
	}
	return 0
}

func FindSelf(instructions []int) int {
	reversed := reversed(instructions)
	registers := make(map[byte]int)

	outputMonitor := 0
	instructionMap := createInstructionMap(func(value int) {
		outputMonitor = value
		registers['A'] = 0
	})

	startAt := make(map[int]int)
	regA := 0
	for i := 0; i < len(reversed); i++ {
		j := getOrDefault(i, startAt)
		for ; j < 8; j++ {
			registers['A'] = regA + j
			registers['B'] = 0
			registers['C'] = 0

			solveInstructions(registers, instructions, instructionMap)

			if outputMonitor == reversed[i] {
				startAt[i] = j + 1
				for k := range startAt {
					if k > i {
						startAt[k] = 0
					}
				}
				regA = (regA + j) * 8
				break
			}
		}

		if j > 7 {
			regA = regA/8 - (startAt[i-1] - 1)
			i = i - 2
		}
	}

	return regA / 8
}

func ParseInput(data string) (registers map[byte]int, instructions []int, err error) {
	reA := regexp.MustCompile(`Register A: (\d+)`)
	reB := regexp.MustCompile(`Register B: (\d+)`)
	reC := regexp.MustCompile(`Register C: (\d+)`)
	reProgram := regexp.MustCompile(`Program: (.*)`)

	registers = make(map[byte]int)

	parseRegisters := func(re *regexp.Regexp, reg byte) error {
		matches := re.FindStringSubmatch(data)
		if len(matches) != 2 {
			return errors.New("bad input")
		}
		registers[reg] = mustAtoi(matches[1])
		return nil
	}

	if err = parseRegisters(reA, 'A'); err != nil {
		return
	}
	if err = parseRegisters(reB, 'B'); err != nil {
		return
	}
	if err = parseRegisters(reC, 'C'); err != nil {
		return
	}

	matches := reProgram.FindStringSubmatch(data)
	if len(matches) != 2 {
		err = errors.New("bad input")
		return
	}
	instructions, err = fn.Transform(strings.Split(matches[1], ","), strconv.Atoi)

	return
}

func main() {
	inputData, err := os.ReadFile("day17.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	registers, instructions, err := ParseInput(string(inputData))
	if err != nil {
		fmt.Printf("Error parsing input data: %v\n", err)
		return
	}

	fmt.Printf("Part 1: %s\n", SolveChronospatialInstructions(registers, instructions))
	fmt.Printf("Part 2: %d\n", FindSelf(instructions))
}
