package main

import (
	fn "AoC/functional"

	"fmt"
	"os"
	"strconv"
)

type Entry struct {
	space bool
	id    int
}

func nextSpaceLocation(entries []Entry, initial, limit int) int {
	for i := initial; i < limit; i++ {
		if entries[i].space {
			return i
		}
	}
	return limit
}

func formatFileMut(entries []Entry, getChunkSize func([]Entry, int) int) int {
	spaceLocation := 0

	for i := len(entries) - 1; i >= 0; i-- {
		if entries[i].space {
			continue
		}

		chunkSize := getChunkSize(entries, i)
		i -= chunkSize - 1

		spaceLocation := nextSpaceLocation(entries, spaceLocation, i)

		spaceSize := 0
		for j := spaceLocation; j < i; j++ {
			if !entries[j].space {
				spaceSize = 0
				continue
			}
			spaceSize++
			if chunkSize == spaceSize {
				for k := 0; k < chunkSize; k++ {
					entries[i+k], entries[j-k] = Entry{true, 0}, Entry{false, entries[i+k].id}
				}
				break
			}
		}
	}

	return fn.Reduce(entries, 0, func(i, checksum int, entry Entry) int {
		return checksum + i*entry.id
	})
}

func formatFile(entries []Entry, getChunkSize func([]Entry, int) int) int {
	entiesCopy := make([]Entry, 0, len(entries))
	entiesCopy = append(entiesCopy, entries...)

	return formatFileMut(entiesCopy, getChunkSize)
}

func FormatFileSingles(entries []Entry) int {
	return formatFile(entries, func([]Entry, int) int { return 1 })
}

func FormatFileChunks(entries []Entry) int {
	getChunkSize := func(entries []Entry, i int) (entriesize int) {
		blockId := entries[i].id
		for ; i >= 0; i-- {
			if entries[i].space || blockId != entries[i].id {
				return
			}
			entriesize++
		}
		return
	}
	return formatFile(entries, getChunkSize)
}

func ParseInputData(data string) ([]Entry, error) {
	entries := make([]Entry, 0, len(data))

	for i, el := range data {
		val, err := strconv.Atoi(string(el))
		if err != nil {
			return nil, err
		}
		for ; val > 0; val-- {
			if i%2 == 0 {
				entries = append(entries, Entry{false, i / 2})
			} else {
				entries = append(entries, Entry{true, 0})
			}
		}
	}
	return entries, nil
}

func main() {
	inputData, err := os.ReadFile("day9.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	entries, err := ParseInputData(string(inputData))
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		return
	}

	fmt.Printf("Part 1: %d\n", FormatFileSingles(entries))
	fmt.Printf("Part 2: %d\n", FormatFileChunks(entries))
}
