package main

import (
	"fmt"
	"os"
	"strconv"
)

type Entry struct {
	space bool
	id    int
}

func FormatFile(blocks []Entry) int {
	j := 0
	for i := len(blocks) - 1; i != 0; i-- {
		if blocks[i].space {
			continue
		}

		for ; j < i; j++ {
			if !blocks[j].space {
				continue
			}
			blocks[i].space, blocks[j].space = true, false
			blocks[i].id, blocks[j].id = blocks[j].id, blocks[i].id
			break
		}
	}

	checksum := 0
	for i, el := range blocks {
		if el.space {
			continue
		}
		checksum += i * el.id
	}
	return checksum
}

func FormatFile2(blocks []Entry) int {
	for i := len(blocks) - 1; i >= 0; i-- {
		if blocks[i].space {
			continue
		}

		blockSize := 0
		blockId := blocks[i].id
		for ; i >= 0; i-- {
			if blocks[i].space || blockId != blocks[i].id {
				i++
				break
			}
			blockSize++
		}

		for j := 0; j < i; j++ {
			if !blocks[j].space {
				continue
			}

			spaceSize := 0
			for ; j < i; j++ {
				if !blocks[j].space {
					j--
					break
				}
				spaceSize++
				if blockSize == spaceSize {
					break
				}
			}

			if blockSize != spaceSize {
				continue
			}

			for k := 0; k < blockSize; k++ {
				blocks[i+k], blocks[j-k] = Entry{true, 0}, Entry{false, blockId}
			}
			break
		}
	}

	checksum := 0
	for i, el := range blocks {
		if el.space {
			continue
		}
		checksum += i * el.id
	}
	return checksum
}

func ParseInputData(data string) ([]Entry, error) {
	blocks := make([]Entry, 0, len(data))

	for i, el := range data {
		val, err := strconv.Atoi(string(el))
		if err != nil {
			return nil, err
		}
		for ; val > 0; val-- {
			if i%2 == 0 {
				blocks = append(blocks, Entry{false, i / 2})
			} else {
				blocks = append(blocks, Entry{true, 0})
			}
		}
	}
	return blocks, nil
}

func main() {
	inputData, err := os.ReadFile("day9.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	blocks1, err := ParseInputData(string(inputData))
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		return
	}

	fmt.Printf("Part 1: %d\n", FormatFile(blocks1))

	blocks2, err := ParseInputData(string(inputData))
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		return
	}
	fmt.Printf("Part 2: %d\n", FormatFile2(blocks2))
}
