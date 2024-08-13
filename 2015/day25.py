import file_reader

STARTING_CODE: int = 20151125


def nextCode(current_code: int) -> int:
    MULTIPLIER: int = 252533
    DIVISOR: int = 33554393
    return current_code * MULTIPLIER % DIVISOR


def getIndex(row: int, column: int) -> int:
    counter = 0
    if row != 1:
        counter = column
        column = row + column - 2
    return column * (column + 1) // 2 + counter - 1


def findCode(row: int, column: int) -> int:
    index = getIndex(row, column)

    code = STARTING_CODE
    code_to_index = {STARTING_CODE: 0}
    index_to_code = {0: STARTING_CODE}
    for i in range(1, index + 1):
        code = nextCode(code)
        if code in code_to_index:
            duplicate_start_index = code_to_index[code]
            loop_size = i - duplicate_start_index
            index = index - duplicate_start_index
            loops = index // loop_size
            remainder = index - loops * loop_size
            return index_to_code[remainder]
        else:
            code_to_index[code] = i
            index_to_code[i] = code
    return code


def main():
    programInput = file_reader.getInput()
    row, column = [
        int(x)
        for x in programInput.replace(".", "").replace(",", "").split()
        if x.isnumeric()
    ]
    print("part 1: ", findCode(row, column))


assert 31916031 == nextCode(STARTING_CODE)
assert 18749137 == nextCode(nextCode(STARTING_CODE))
assert 16080970 == nextCode(nextCode(nextCode(STARTING_CODE)))

assert 31916031 == findCode(2, 1)
assert 18749137 == findCode(1, 2)

if __name__ == "__main__":
    main()
