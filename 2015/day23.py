from typing import Callable
import file_reader

INSTRUCTION_MAP: dict[str, Callable[[any], int]] = {
    "hlf": lambda r, reg_entry, pos: (
        reg_entry.update({r: reg_entry.setdefault(r, 0) / 2}),
        pos + 1,
    )[1],
    "tpl": lambda r, reg_entry, pos: (
        reg_entry.update({r: reg_entry.setdefault(r, 0) * 3}),
        pos + 1,
    )[1],
    "inc": lambda r, reg_entry, pos: (
        reg_entry.update({r: reg_entry.setdefault(r, 0) + 1}),
        pos + 1,
    )[1],
    "jmp": lambda offset, reg_entry, pos: pos + offset,
    "jie": lambda r, offset, reg_entry, pos: (
        pos + offset if reg_entry.setdefault(r, 0) % 2 == 0 else pos + 1
    ),
    "jio": lambda r, offset, reg_entry, pos: (
        pos + offset if reg_entry.setdefault(r, 0) == 1 else pos + 1
    ),
}

PARSER_MAP: dict[str, Callable[[str], any]] = {
    "hlf": lambda rhs: (rhs,),
    "tpl": lambda rhs: (rhs,),
    "inc": lambda rhs: (rhs,),
    "jmp": lambda rhs: (int(rhs),),
    "jie": lambda rhs: (lambda lhs, rhs: (lhs, int(rhs)))(*rhs.split(", ")),
    "jio": lambda rhs: (lambda lhs, rhs: (lhs, int(rhs)))(*rhs.split(", ")),
}


def executeProgram(program: list[str], initial_vals: dict[str, int]) -> dict[str, int]:
    reg_entry: dict[str, int] = initial_vals.copy()

    pos: int = 0
    while pos < len(program):
        command, rhs = program[pos].split(" ", 1)
        pos = INSTRUCTION_MAP[command](*PARSER_MAP[command](rhs), reg_entry, pos)
    return reg_entry


def main():
    programInput = file_reader.getInput().splitlines()

    print("part 1: ", executeProgram(programInput, {})["b"])
    print("part 2: ", executeProgram(programInput, {"a": 1})["b"])


example: list[str] = (
    """inc a
jio a, +2
tpl a
inc a""".split(
        "\n"
    )
)

example_result = executeProgram(example, {})

assert 2 == example_result["a"]

if __name__ == "__main__":
    main()
