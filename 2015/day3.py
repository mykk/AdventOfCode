import file_reader

directions = {
    "<": complex(0, -1),
    ">": complex(0, 1),
    "^": complex(-1, 0),
    "v": complex(1, 0),
}


def part1(input_str: str) -> int:
    current = complex(0, 0)
    visited = {current}
    for c in input_str:
        visited.add(current := current + directions[c])
    return len(visited)


def part2(input_str: str) -> int:
    current1 = complex(0, 0)
    current2 = complex(0, 0)
    visited = {current1}
    for i in range(0, len(input_str), 2):
        visited.add(current1 := current1 + directions[input_str[i]])
        visited.add(current2 := current2 + directions[input_str[i + 1]])
    return len(visited)


def main():
    input_str: str = file_reader.getInput()

    print("part 1: ", part1(input_str))
    print("part 2: ", part2(input_str))


if __name__ == "__main__":
    main()
