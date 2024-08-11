import file_reader

ARROW = '=>'

def parseReplacements(programInput: list[str]) -> dict[str, list[str]]:
    replacementMap = {}
    for line in programInput:
        x, y = line.split(ARROW)
        replacementMap.setdefault(x.strip(), []).append(y.strip())
    return replacementMap

def reverseReplacementMap(replacementMap: dict[str, str]) -> dict[str, list[str]]:
    reversedMap = {}
    for replacement, values in replacementMap.items():
        for value in values:
            reversedMap[value] = replacement
    return reversedMap

def parseInitialMolecule(programInput: list[str]) -> str:
    return programInput[-1]

def parseInput(programInput: list[str]) -> tuple[dict[str, str], str]:
    return parseReplacements(programInput[:-2]), parseInitialMolecule(programInput)

def findReplacements(replacementMap: dict[str, str], molecule: str) -> set[str]:
    molecules = set()
    for replacement, values in replacementMap.items():
        index = 0
        while (index := molecule.find(replacement, index)) != -1:
            for value in values:
                molecules.add(molecule[:index] + molecule[index:].replace(replacement, value, 1))
            index += len(replacement)
    return molecules


def findUnmatchStart(molecule, newMolecule) -> int:
    match = 0
    for j, c in enumerate(molecule):
        if j >= len(newMolecule) or c != newMolecule[j]: return match
        match += 1
    return 0
        
def findMakeSteps(replacementMap: dict[str, str], molecule: str, start: str) -> int:
    reversedMap = reverseReplacementMap(replacementMap)
    
    count = 0
    sortedReplacements = sorted([x for x in reversedMap], key = lambda x : (not 'Ca' in x, not 'Ti' in x, -len(x)))
    while molecule != start:
        for replacement in sortedReplacements:
            if replacement in molecule:
                while replacement in molecule:
                    molecule = molecule.replace(replacement, reversedMap[replacement], 1)
                    count += 1
                break
    return count

def main():
    programInput = file_reader.getInput().splitlines()
    print("part 1: ", len(findReplacements(*parseInput(programInput))))
    print("part 2: ", findMakeSteps(*parseInput(programInput), 'e'))

example ="""H => HO
H => OH
O => HH

HOH""".split('\n')

assert(4 == len(findReplacements(*parseInput(example))))

example = """e => H
e => O
H => HO
H => OH
O => HH""".split('\n')

assert(3 == findMakeSteps(parseReplacements(example), "HOH", 'e'))
assert(6 == findMakeSteps(parseReplacements(example), "HOHOHO", 'e'))

if __name__ == '__main__':
   main()


