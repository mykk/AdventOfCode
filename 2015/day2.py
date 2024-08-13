import file_reader
from functools import reduce


def part1_calc(w: int, h: int, length: int) -> int:
    return (
        2 * length * w
        + 2 * w * h
        + 2 * h * length
        + min([w * h, h * length, length * w])
    )


def part1(dimensions: list[tuple[int, int, int]]) -> int:
    return reduce(lambda x, dims: x + part1_calc(*dims), dimensions, 0)


def part2_calc(w: int, h: int, length: int) -> int:
    return sum(sorted([w, h, length])[0:2]) * 2 + w * h * length


def part2(dimensions: list[tuple[int, int, int]]) -> int:
    return reduce(lambda x, dims: x + part2_calc(*dims), dimensions, 0)


def main():
    input_str: str = file_reader.getInput()
    dimensions = [
        (int(w), int(h), int(l))
        for w, h, l in [x.split("x") for x in input_str.splitlines()]
    ]

    print("part 1: ", part1(dimensions))
    print("part 2: ", part2(dimensions))


if __name__ == "__main__":
    main()
