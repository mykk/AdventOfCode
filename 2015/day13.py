import file_reader
ACTION_MAP = {'gain' : 1, 'lose' : -1}

def parseInput(programInput: str) -> dict[str, dict[str, int]]:
    happinessMap = {}
    for line in programInput:
        person1, _, action, amount, *_, person2 = [x.strip('.') for x in line.split()]
        happinessMap.setdefault(person1, {})[person2] = ACTION_MAP[action] * int(amount)
    return happinessMap

def extendInput(happinessMap: dict[str, dict[str, int]], newPerson: str, newValue: int):
    happinessMap = {key : items.copy() for key, items in happinessMap.items()}
    newEntry = {}
    for x in happinessMap:
        happinessMap[x][newPerson] = newValue
        newEntry[x] = newValue
    happinessMap[newPerson] = newEntry
    return happinessMap

def computeHappiness(happinessMap: dict[str, dict[str, int]], seated: list[str]) -> int:
    happiness = happinessMap[seated[0]][seated[-1]] + happinessMap[seated[-1]][seated[0]]  
    for i in range(len(seated) - 1):
        happiness += happinessMap[seated[i]][seated[i + 1]] + happinessMap[seated[i + 1]][seated[i]]
    return happiness

def findHappiness(happinessMap: dict[str, dict[str, int]], seated: list[str]) -> int:
    if len(happinessMap) == len(seated):
        return computeHappiness(happinessMap, seated)
    
    best = None
    for person1 in happinessMap:
        if person1 not in seated:
            happiness = findHappiness(happinessMap, seated + [person1])
            if best is None or happiness > best: best = happiness
    return best

def main():
   input_str = file_reader.getInput().splitlines()

   print("part 1: ", findHappiness(parseInput(input_str), []))
   print("part 2: ", findHappiness(extendInput(parseInput(input_str), 'mike', 0), []))

example = """Alice would gain 54 happiness units by sitting next to Bob.
Alice would lose 79 happiness units by sitting next to Carol.
Alice would lose 2 happiness units by sitting next to David.
Bob would gain 83 happiness units by sitting next to Alice.
Bob would lose 7 happiness units by sitting next to Carol.
Bob would lose 63 happiness units by sitting next to David.
Carol would lose 62 happiness units by sitting next to Alice.
Carol would gain 60 happiness units by sitting next to Bob.
Carol would gain 55 happiness units by sitting next to David.
David would gain 46 happiness units by sitting next to Alice.
David would lose 7 happiness units by sitting next to Bob.
David would gain 41 happiness units by sitting next to Carol."""

assert(330 == findHappiness(parseInput(example.splitlines()), []))

if __name__ == '__main__':
   main()
