import file_reader


def parseInput(programInput: list[str]) -> dict[complex, bool]:
    return {
        complex(i, j): c == "#"
        for i, row in enumerate(programInput)
        for j, c in enumerate(row)
    }


def getNewState(cell: complex, grid: dict[complex, bool], onOffRule: callable) -> bool:
    return onOffRule(
        grid,
        cell,
        sum(
            [
                1
                for i in range(-1, 2)
                for j in range(-1, 2)
                if grid.get(cell + complex(i, j), False) and (j != 0 or i != 0)
            ]
        ),
    )


def getNewGrid(
    grid: dict[complex, bool], count: int, onOffRule: callable
) -> dict[complex, bool]:
    if count <= 0:
        return grid
    return getNewGrid(
        {cell: getNewState(cell, grid, onOffRule) for cell in grid},
        count - 1,
        onOffRule,
    )


def getNewGridPart1(grid: dict[complex, bool], count: int) -> dict[complex, bool]:
    return getNewGrid(
        grid,
        count,
        lambda grid, cell, onCount: (grid[cell] and (onCount == 2 or onCount == 3))
        or (not grid[cell] and onCount == 3),
    )


def getNewGridPart2(grid: dict[complex, bool], count: int) -> dict[complex, bool]:
    corners = [
        complex(0, 0),
        complex(0, max([x.imag for x in grid])),
        complex(max([x.real for x in grid]), max([x.imag for x in grid])),
        complex(max([x.real for x in grid]), 0),
    ]
    return getNewGrid(
        {x: x in corners or grid[x] for x in grid},
        count,
        lambda grid, cell, onCount: (grid[cell] and (onCount == 2 or onCount == 3))
        or (not grid[cell] and onCount == 3)
        or (cell in corners),
    )


def main():
    grid = parseInput(file_reader.getInput().splitlines())

    grid1 = getNewGridPart1(grid, 100)
    print("part 1: ", sum([1 for x in grid1 if grid1[x]]))

    grid2 = getNewGridPart2(grid, 100)
    print("part 2: ", sum([1 for x in grid2 if grid2[x]]))


example = """.#.#.#
...##.
#....#
..#...
#.#..#
####..""".split(
    "\n"
)

exampleGridPart1 = getNewGridPart1(parseInput(example), 4)
assert 4 == sum([1 for x in exampleGridPart1 if exampleGridPart1[x]])

exampleGridPart2 = getNewGridPart2(parseInput(example), 5)
assert 17 == sum([1 for x in exampleGridPart2 if exampleGridPart2[x]])

if __name__ == "__main__":
    main()
