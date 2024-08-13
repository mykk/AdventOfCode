import file_reader

TURN_ON = "turn on"
TOGGLE = "toggle"
TURN_OFF = "turn off"
THROUGH = "through"


def cleanUp(x):
    return [int(coord.strip()) for coord in x.split(",")]


def getCoords(line: str) -> tuple:
    start, end = line.split(THROUGH)
    return (cleanUp(start), cleanUp(end))


def assign(lighths: list, x: int, y: int, value: int) -> list:
    lighths[x][y] = value
    return lighths


def append(lighths: list, x: int, y: int, value: int) -> list:
    lighths[x][y] += value
    if lighths[x][y] < 0:
        lighths[x][y] = 0
    return lighths


def lightCount(programInput: list, switchFactory: dict) -> int:
    lighths = [[0] * 1000 for _ in range(1000)]
    for line in programInput:
        if TURN_ON in line:
            start, end = getCoords(line.removeprefix(TURN_ON))
            switch = switchFactory[TURN_ON]
        elif TURN_OFF in line:
            start, end = getCoords(line.removeprefix(TURN_OFF))
            switch = switchFactory[TURN_OFF]
        else:
            start, end = getCoords(line.removeprefix(TOGGLE))
            switch = switchFactory[TOGGLE]
        for x in range(start[0], end[0] + 1):
            for y in range(start[1], end[1] + 1):
                lighths = switch(x, y, lighths)
    return sum([sum(x) for x in lighths])


def lightCount1(programInput: list) -> int:
    return lightCount(
        programInput,
        {
            TURN_ON: lambda x, y, lighths: assign(lighths, x, y, 1),
            TURN_OFF: lambda x, y, lighths: assign(lighths, x, y, 0),
            TOGGLE: lambda x, y, lighths: assign(
                lighths, x, y, 1 if lighths[x][y] == 0 else 0
            ),
        },
    )


def lightCount2(programInput: list) -> int:
    return lightCount(
        programInput,
        {
            TURN_ON: lambda x, y, lighths: append(lighths, x, y, 1),
            TURN_OFF: lambda x, y, lighths: append(lighths, x, y, -1),
            TOGGLE: lambda x, y, lighths: append(lighths, x, y, 2),
        },
    )


def main():
    input_str: str = file_reader.getInput().splitlines()

    print("part 1: ", lightCount1(input_str))
    print("part 2: ", lightCount2(input_str))


example1 = ["turn on 0,0 through 999,999"]
example2 = ["toggle 0,0 through 999,0"]
assert 1000 * 1000 == lightCount1(example1)
assert 1000 == lightCount1(example2)
assert 1000 * 1000 - 1000 == lightCount1(example1 + example2)


if __name__ == "__main__":
    main()
