import file_reader

WIRE: str = "->"
OPERATION_MAP: dict[str, callable] = {
    "AND": lambda x, y: x & y,
    "OR": lambda x, y: x | y,
    "LSHIFT": lambda x, y: x << y,
    "RSHIFT": lambda x, y: x >> y,
}


def parseValue(x, values: dict[str, int]):
    return int(x) if x not in values else values[x]


def parseLHS(values: dict[str, int], lhs: str) -> int:
    ops = [x.strip() for x in lhs.split() if x.strip()]
    if len(ops) == 1:
        return parseValue(ops[0], values)
    if len(ops) == 2:
        assert ops[0] == "NOT"
        return ~parseValue(ops[1], values)
    return OPERATION_MAP[ops[1]](parseValue(ops[0], values), parseValue(ops[2], values))


def circuitToDict(circuit: list[str]) -> dict[str, list[list[str]]]:
    circuitDict = {}
    for current in circuit:
        ops = [x.strip() for x in current.split(WIRE)[0].split() if x.strip()]
        if len(ops) == 1:
            circuitDict.setdefault(ops[0], []).append([current])
        elif len(ops) == 2:
            assert ops[0] == "NOT"
            circuitDict.setdefault(ops[1], []).append([current])
        elif len(ops) == 3:
            if not ops[0].isnumeric() and not ops[2].isnumeric():
                circuitDict.setdefault(ops[0], []).append([ops[2], current])
                circuitDict.setdefault(ops[2], []).append([ops[0], current])
            elif ops[0].isnumeric():
                circuitDict.setdefault(ops[2], []).append([current])
            else:
                circuitDict.setdefault(ops[0], []).append([current])
        else:
            assert False
    return circuitDict


def findInitialBooklets(circuitDict: dict[str, list[list[str]]]) -> list[str]:
    initial = []
    for key in circuitDict:
        if key.isnumeric():
            for current in circuitDict[key]:
                assert len(current) == 1
                if len(current) == 1:
                    initial.append(current[0])
    return initial


def getNewlyAvailableBooklets(
    key: str, circuitDict: dict[str, list[list[str]]], values: dict[str, int]
) -> list[str]:
    if key not in circuitDict:
        return []
    new = []
    for current in circuitDict[key]:
        if len(current) == 1:
            new.append(current[0])
        else:
            otherKey, booklet = current
            if otherKey in values or otherKey.isnumeric():
                new.append(booklet)
    return new


def processCircuit(circuit: list[str]) -> dict[str, int]:
    circuitDict = circuitToDict(circuit)
    booklets = findInitialBooklets(circuitDict)
    values = {}
    while booklets:
        current = booklets.pop(0)
        lhs, rhs = [x.strip() for x in current.split(WIRE)]
        values[rhs] = parseLHS(values, lhs) & 0xFFFF
        booklets += getNewlyAvailableBooklets(rhs, circuitDict, values)
    return values


def processCircuitOverwrite(circuit: list[str]) -> dict[str, int]:
    circuit = circuit.copy()
    overwriteValue = processCircuit(circuit)["a"]

    for current in circuit:
        if current.endswith("-> b"):
            circuit.remove(current)
            break
    circuit.append(str(overwriteValue) + " -> b")
    return processCircuit(circuit)


def main():
    input_str: str = file_reader.getInput().splitlines()

    print("part 1: ", processCircuit(input_str)["a"])
    print("part 2: ", processCircuitOverwrite(input_str)["a"])


example1 = """123 -> x
456 -> y
x AND y -> d
x OR y -> e
x LSHIFT 2 -> f
y RSHIFT 2 -> g
NOT x -> h
NOT y -> i""".split(
    "\n"
)

example1_result = processCircuit(example1)
assert 72 == example1_result["d"]
assert 507 == example1_result["e"]
assert 492 == example1_result["f"]
assert 114 == example1_result["g"]
assert 65412 == example1_result["h"]
assert 65079 == example1_result["i"]
assert 123 == example1_result["x"]
assert 456 == example1_result["y"]

if __name__ == "__main__":
    main()
