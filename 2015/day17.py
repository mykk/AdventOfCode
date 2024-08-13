import file_reader


def parseInput(programInput: list[str]) -> list[int]:
    return [int(x) for x in programInput]


def fillEggnog(
    containers: list[int],
    remainingCapacity: int,
    current: list[int] = None,
    capacities: dict[int, int] = None,
) -> int:
    if current is None:
        current = []

    if remainingCapacity == 0:
        if capacities is not None:
            containerCount = len(current)
            capacities[containerCount] = capacities.get(containerCount, 0) + 1
        return 1
    elif remainingCapacity < 0:
        return 0

    temp = containers.copy()
    return sum(
        fillEggnog(temp := temp[1:], remainingCapacity - x, current + [x], capacities)
        for x in containers
    )


def minContainerCount(containers: list[int], remainingCapacity: int):
    capacities = {}
    fillEggnog(containers, remainingCapacity, None, capacities)
    return capacities[min(capacities)]


def main():
    containers = parseInput(file_reader.getInput().splitlines())

    print("part 1: ", fillEggnog(containers, 150))
    print("part 2: ", minContainerCount(containers, 150))


example = [
    20,
    15,
    10,
    5,
    5,
]

assert 4 == fillEggnog(example, 25)
assert 3 == minContainerCount(example, 25)

if __name__ == "__main__":
    main()
