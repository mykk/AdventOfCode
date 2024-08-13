import file_reader


def lookAndSay(sequence: str, cycles: int) -> str:
    if cycles == 0:
        return sequence
    newSequence = ""
    count = 0
    previous = ""
    for c in sequence:
        if previous == "":
            previous = c
            count = 1
        elif previous == c:
            count += 1
        else:
            newSequence += str(count) + previous
            previous = c
            count = 1
    else:
        newSequence += str(count) + c
    return lookAndSay(newSequence, cycles - 1)


def main():
    input_str: str = file_reader.getInput()

    print("part 1: ", len(lookAndSay(input_str, 40)))
    print("part 2: ", len(lookAndSay(input_str, 50)))


assert "11" == lookAndSay("1", 1)
assert "21" == lookAndSay("11", 1)
assert "1211" == lookAndSay("21", 1)
assert "111221" == lookAndSay("1211", 1)
assert "312211" == lookAndSay("111221", 1)
assert "312211" == lookAndSay("1", 5)

if __name__ == "__main__":
    main()
