from collections import deque
import file_reader

example = """London to Dublin = 464
London to Belfast = 518
Dublin to Belfast = 141""".split('\n')

DEST_SEPARATOR = "to"
PRICE_SEPARATOR = "="

def appendDictIfNotFound(destinations: dict[str, dict[str, int]], location: str, destination: str, price: int):
   if (destination, price) not in destinations.setdefault(location, {}):
      destinations[location][destination] = price

def parseInput(programInput: list[str]) -> dict[str, dict[str, int]]:
   destinations = {}
   for line in programInput:
      location, destination = [x.strip() for x in line.split(DEST_SEPARATOR)]
      destination, price = [x.strip() for x in destination.split(PRICE_SEPARATOR)]
      appendDictIfNotFound(destinations, location, destination, int(price))
      appendDictIfNotFound(destinations, destination, location, int(price))
   return destinations

def findBestPath(destinations: dict[str, dict[str, int]], continueSearch: callable, betterCost: callable) -> int:
   bestCost = 0
   for startPoint in destinations:
      toVisit = deque([(startPoint, {startPoint}, 0)])
      while toVisit:
         location, visited, cost = toVisit.popleft()
         if not continueSearch(cost, bestCost): continue
         if len(visited) == len(destinations) and betterCost(cost, bestCost):
            bestCost = cost
            continue
         for dest in [x for x in destinations[location] if x not in visited]:
            toVisit.append((dest, visited|{dest}, cost + destinations[location][dest]))
   return bestCost

def findShortestPath(destinations: dict[str, dict[str, int]]) -> int:
   return findBestPath(destinations, lambda x, y: x < y or y == 0, lambda x, y: x < y or y == 0)

def findLongestPath(destinations: dict[str, dict[str, int]]) -> int:
   return findBestPath(destinations, lambda x, y: True, lambda x, y: x > y)

def main():
   destinations = parseInput(file_reader.getInput().splitlines())

   print("part 1: ", findShortestPath(destinations))
   print("part 2: ", findLongestPath(destinations))

assert(2 == len(parseInput(example)["London"]))
assert(2 == len(parseInput(example)["Dublin"]))
assert(2 == len(parseInput(example)["Belfast"]))
assert(parseInput(example)["Belfast"]["London"] == 518)
assert(parseInput(example)["London"]["Belfast"] == 518)
assert(parseInput(example)["London"]["Dublin"] == 464)
assert(parseInput(example)["Dublin"]["London"] == 464)
assert(parseInput(example)["Dublin"]["Belfast"] == 141)
assert(parseInput(example)["Belfast"]["Dublin"] == 141)

assert(605 == findShortestPath(parseInput(example)))
assert(982 == findLongestPath(parseInput(example)))

if __name__ == '__main__':
   main()
