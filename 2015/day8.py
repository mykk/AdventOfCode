import re
import file_reader


def transformImport(input: str) -> str:
    input = re.sub(r"\\\\", "A", input)
    input = re.sub(r'\\"', "A", input)
    input = re.sub(r"\\x[0-9a-fA-F]{2}", "A", input)
    return input


def countLetters(input: list[str]) -> int:
    return sum([len(x) - len(transformImport(x)) + 2 for x in input])


def makeItRaw(input: str) -> str:
    input = re.sub(r"\\\\", "AAAA", input)
    input = re.sub(r"\\x[0-9a-fA-F]{2}", r"AAAAA", input)
    input = re.sub(r'\\"', r"AAAA", input)
    return input


def countRaw(input: list[str]) -> int:
    return sum([len(makeItRaw(x)) - len(x) + 4 for x in input])


def main():
    input_str: str = file_reader.getInput().splitlines()

    print("part 1: ", countLetters(input_str))
    print("part 2: ", countRaw(input_str))


assert 2 == countLetters([r""])
assert 2 == countLetters([r"abc"])
assert 3 == countLetters([r"aaa\"aaa"])
assert 5 == countLetters([r"\x27"])
assert 12 == countLetters([r"\x27", r"", r"aaa\"aaa", r"abc"])

assert 4 == countRaw([r""])
assert 4 == countRaw([r"abc"])
assert 6 == countRaw([r"aaa\"aaa"])
assert 5 == countRaw([r"\x27"])
assert 7 == countRaw([r"\x27aaa\"aaa"])

assert 19 == countRaw([r"\x27", r"", r"aaa\"aaa", r"abc"])

if __name__ == "__main__":
    main()
