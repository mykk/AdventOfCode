import file_reader
import json


def addParsedNumbers(programInput, ignoreDict: callable) -> int:
    if isinstance(programInput, int):
        return programInput
    if isinstance(programInput, dict):
        if ignoreDict(programInput):
            return 0
        return sum(
            [addParsedNumbers(x, ignoreDict) for x in programInput.values()]
        ) + sum([addParsedNumbers(x, ignoreDict) for x in programInput])
    if isinstance(programInput, list):
        return sum([addParsedNumbers(x, ignoreDict) for x in programInput])
    return 0


def addNumbers(programInput: str, ignoreDict: callable) -> int:
    parsed = json.loads(programInput)
    return sum([addParsedNumbers(x, ignoreDict) for x in parsed.values()])


def main():
    input_str = file_reader.getInput()

    print("part 1: ", addNumbers(input_str, lambda x: False))
    print(
        "part 2: ",
        addNumbers(input_str, lambda x: "red" in x.keys() or "red" in x.values()),
    )


assert 123 == addNumbers('{"abd":123}', lambda x: False)
assert -123 == addNumbers('{"abd":-123}', lambda x: False)

if __name__ == "__main__":
    main()
