import file_reader

REPORT = {"children": 3, "cats": 7, "samoyeds": 2, "pomeranians": 3, "akitas": 0, "vizslas": 0, "goldfish": 5, "trees": 3, "cars": 2, "perfumes": 1}
COMPARE_DICT = {
    "cats" : lambda actual, report : actual > report,
    "trees" : lambda actual, report : actual > report,
    "pomeranians" : lambda actual, report : actual < report,
    "goldfish" : lambda actual, report : actual < report
}

def parseInput(programInput: list[str]) -> dict[int, dict[str,int]]:
    aunts = {}
    for line in programInput:
        aunt, properties = line.split(':', 1)
        aunt = int(aunt.split()[1])
        aunts[aunt] = {}
        for property in properties.split(','):
            prop, amount = [x for x in property.split(':')]
            aunts[aunt][prop] = int(amount)
    return aunts

def findAunt(aunts: dict[int, dict[str,int]], compare: dict[str, callable]) -> int:
    for aunt, aunt_properties in aunts.items():
        for prop, value in aunt_properties.items():
            if prop in compare:
                if not compare[prop](value, REPORT[prop]): break
            elif value != REPORT[prop]: break
        else:
            return aunt
    return None

def main():
   input_str = file_reader.getInput().splitlines()

   print("part 1: ", findAunt(parseInput(input_str), {}))
   print("part 2: ", findAunt(parseInput(input_str), COMPARE_DICT))

if __name__ == '__main__':
   main()